from fastapi import FastAPI, HTTPException
from fastapi.responses import StreamingResponse
from pydantic import BaseModel
import asyncio
import json
import httpx
import os
import traceback
import difflib

app = FastAPI(title="Agent Harness API")

BACKEND_URL = os.getenv("BACKEND_URL", "http://gemma4-backend:8080/api")
LLM_BASE_URL = os.getenv("LLM_BASE_URL", "http://host.docker.internal:8091/v1")
LLM_API_KEY = os.getenv("LLM_API_KEY", "ollama")
LLM_MODEL = os.getenv("LLM_MODEL", "Gemma4")
MAX_TURNS = int(os.getenv("MAX_TURNS", "25"))

os.environ["OPENAI_API_KEY"] = LLM_API_KEY
os.environ["OPENAI_BASE_URL"] = LLM_BASE_URL

Agent = None
function_tool = None
try:
    from agents import Agent, function_tool
except ImportError:
    try:
        from openai_agents import Agent, function_tool
    except ImportError:
        pass

class RunRequest(BaseModel):
    wiki_prompt: str
    task_description: str
    additional_prompt: str

# Track file changes
original_files = {}

def read_file(filepath: str) -> str:
    """Read a file from the workspace."""
    full_path = os.path.join("/app/workspace", filepath.lstrip("/"))
    try:
        with open(full_path, "r", encoding="utf-8") as f:
            return f.read()
    except Exception as e:
        return f"Error reading file: {str(e)}"

def write_file(filepath: str, content: str) -> str:
    """Write content to a file in the workspace."""
    full_path = os.path.join("/app/workspace", filepath.lstrip("/"))
    try:
        # Save original content for diff generation
        if full_path not in original_files:
            if os.path.exists(full_path):
                with open(full_path, "r", encoding="utf-8") as f:
                    original_files[full_path] = f.read()
            else:
                original_files[full_path] = None # Indicates new file

        os.makedirs(os.path.dirname(full_path), exist_ok=True)
        with open(full_path, "w", encoding="utf-8") as f:
            f.write(content)
        return "File successfully written."
    except Exception as e:
        return f"Error writing file: {str(e)}"

@app.post("/api/v1/agents/run")
async def run_agent_task(req: RunRequest):
    task_prompt = f"Task: {req.task_description}\nAdditional Info: {req.additional_prompt}"
    
    try:
        from agents.sandbox.sandboxes import DockerSandboxClient, DockerSandboxClientOptions
        from agents.sandbox import SandboxAgent, SandboxRunConfig, Manifest, LocalSnapshotSpec, Permissions, User, FileMode
        from agents.sandbox.entries import Dir
        import docker
        from pathlib import Path
        import os
        from agents import Runner
        
        client = DockerSandboxClient(docker_client=docker.from_env())
        session = await client.create(
            manifest=Manifest(root="/workspace"),
            options=DockerSandboxClientOptions(image="ubuntu:latest")
        )
        await session.start()
        agent_response = ""
        diffs = []
        
        try:
            # await session.hydrate_workspace() # Fails with openai-agents 0.18.3
            container_id = getattr(session.state, 'container_id', None)
            if container_id:
                docker_client = docker.from_env()
                container = docker_client.containers.get(container_id)
                import io, tarfile
                stream = io.BytesIO()
                def tar_filter(tarinfo):
                    if any(part in ["node_modules", ".git", "__pycache__", "dist", "build", ".venv", "venv"] for part in tarinfo.name.split('/')):
                        return None
                    return tarinfo
                    
                with tarfile.open(fileobj=stream, mode='w') as tar:
                    for item in os.listdir("/app/workspace"):
                        if item in ["node_modules", ".git", "__pycache__", "dist", "build", ".venv", "venv"]:
                            continue
                        tar.add(os.path.join("/app/workspace", item), arcname=item, filter=tar_filter)
                container.put_archive('/workspace', stream.getvalue())
            import tools
            agent = SandboxAgent(
                name="TaskExecutor",
                instructions=req.wiki_prompt,
                model=LLM_MODEL,
                tools=[tools.convert_document_to_markdown, tools.generate_knowledge_graph]
            )
            
            from agents import RunConfig
            run_config = RunConfig(sandbox=SandboxRunConfig(session=session))
            result = await Runner.run(agent, task_prompt, run_config=run_config, max_turns=MAX_TURNS)
            agent_response = str(result.final_output)
            
            # Generate Diffs - Only read files modified within the last 5 minutes (i.e. touched by the agent)
            exec_res = await session.exec(["find", "/workspace", "-type", "f", "-mmin", "-5", "-not", "-path", "*/node_modules/*", "-not", "-path", "*/.git/*"])
            sandbox_files = exec_res.stdout.split()
            
            for sfile in sandbox_files:
                sfile = sfile.strip()
                if not sfile: continue
                rel_path = os.path.relpath(sfile, "/workspace")
                local_path = os.path.join("/app/workspace", rel_path)
                
                try:
                    new_content = await session.read(sfile)
                except Exception:
                    continue
                    
                if os.path.exists(local_path):
                    try:
                        with open(local_path, "r", encoding="utf-8") as f:
                            old_content = f.read()
                    except UnicodeDecodeError:
                        continue # Skip binary files
                else:
                    old_content = ""
                    
                if new_content != old_content:
                    orig_lines = old_content.splitlines(keepends=True)
                    new_lines = new_content.splitlines(keepends=True)
                    diff = "".join(difflib.unified_diff(
                        orig_lines, new_lines,
                        fromfile=f"a/{rel_path}",
                        tofile=f"b/{rel_path}"
                    ))
                    if diff:
                        diffs.append({"file": rel_path, "diff": diff})

            # Check for deleted files
            for root_dir, dirs, files in os.walk("/app/workspace"):
                for excluded in ["node_modules", ".git", "dist", "build", ".venv", "venv", "__pycache__"]:
                    if excluded in dirs: dirs.remove(excluded)
                for file in files:
                    rel_path = os.path.relpath(os.path.join(root_dir, file), "/app/workspace")
                    sfile = f"/workspace/{rel_path}"
                    if sfile not in sandbox_files:
                        try:
                            with open(os.path.join(root_dir, file), "r", encoding="utf-8") as f:
                                old_content = f.read()
                        except UnicodeDecodeError:
                            continue # Skip binary files
                        orig_lines = old_content.splitlines(keepends=True)
                        diff = "".join(difflib.unified_diff(
                            orig_lines, [],
                            fromfile=f"a/{rel_path}",
                            tofile=f"b/{rel_path}"
                        ))
                        if diff:
                            diffs.append({"file": rel_path, "diff": diff})

        finally:
            try:
                await session.stop()
            except Exception:
                pass
            try:
                await client.delete(session)
            except Exception:
                pass

        return {
            "status": "success",
            "agent_response": agent_response,
            "changes": diffs
        }
    except Exception as e:
        traceback.print_exc()
        raise HTTPException(status_code=500, detail=f"Agent execution failed: {str(e)}")

@app.get("/health")
def health_check():
    return {"status": "ok", "message": "Agent harness v2 is running."}

class ConvertRequest(BaseModel):
    filepath: str
    original_filename: str = ""

def extract_knowledge_graph(markdown_text: str) -> dict:
    import openai
    client = openai.OpenAI(
        base_url=LLM_BASE_URL,
        api_key=LLM_API_KEY
    )
    
    prompt = f"""
    Extract a Knowledge Graph from the following text.
    Return ONLY a valid JSON object matching this schema, nothing else:
    {{
      "nodes": [{{"id": "string", "label": "Person|Organization|Concept|Location|Event|Entity", "properties": {{"name": "string"}} }}],
      "edges": [{{"source": "node_id", "target": "node_id", "label": "RELATED_TO|PART_OF|DEPENDS_ON|AFFECTS", "properties": {{"description": "string"}} }}]
    }}
    
    Text:
    {markdown_text[:12000]}
    """
    
    try:
        response = client.chat.completions.create(
            model=LLM_MODEL,
            messages=[{"role": "user", "content": prompt}],
            response_format={"type": "json_object"},
            temperature=0.1
        )
        content = response.choices[0].message.content
        return json.loads(content)
    except Exception as e:
        traceback.print_exc()
        return {"nodes": [], "edges": []}

@app.post("/api/v1/convert")
async def convert_document(req: ConvertRequest):
    async def sse_generator():
        import tools
        # Run the blocking conversion in a background thread
        loop = asyncio.get_running_loop()
        task = loop.run_in_executor(None, tools.convert_document_to_markdown_impl, req.filepath)
        
        # Yield heartbeats every 5 seconds while waiting
        while not task.done():
            yield f"data: {json.dumps({'status': 'processing'})}\n\n"
            try:
                await asyncio.wait_for(asyncio.shield(task), timeout=5.0)
            except asyncio.TimeoutError:
                pass
                
        # Task is complete, yield the final result
        try:
            result = task.result()
            if "error" in result:
                yield f"data: {json.dumps({'status': 'error', 'detail': result['error']})}\n\n"
            else:
                markdown = result["markdown"]
                docling_doc = result["document"]
                yield f"data: {json.dumps({'status': 'extracting_graph'})}\n\n"
                
                graph_task = loop.run_in_executor(None, extract_knowledge_graph, markdown)
                
                # Extract Structural Graph
                import hashlib
                from docling_core.transforms.chunker.hierarchical_chunker import HierarchicalChunker
                
                chunker = HierarchicalChunker()
                chunks = list(chunker.chunk(docling_doc))
                
                doc_name = req.original_filename if req.original_filename else os.path.basename(req.filepath)
                doc_id = hashlib.md5(doc_name.encode()).hexdigest()
                structural_nodes = [{"id": doc_id, "label": "Document", "properties": {"name": doc_name}}]
                structural_edges = []
                
                for c in chunks:
                    chunk_id = hashlib.md5(c.text.encode()).hexdigest()
                    structural_nodes.append({"id": chunk_id, "label": "Segment", "properties": {"text": c.text[:200]}})
                    
                    headings = c.meta.headings if hasattr(c.meta, 'headings') and c.meta.headings else []
                    parent_id = doc_id
                    
                    if headings:
                        chapter_name = headings[-1]
                        chapter_id = hashlib.md5((doc_name + chapter_name).encode()).hexdigest()
                        
                        if not any(n["id"] == chapter_id for n in structural_nodes):
                            structural_nodes.append({"id": chapter_id, "label": "Chapter", "properties": {"name": chapter_name}})
                            structural_edges.append({"source": chapter_id, "target": doc_id, "label": "PART_OF"})
                            
                        parent_id = chapter_id
                        
                    structural_edges.append({"source": chunk_id, "target": parent_id, "label": "PART_OF"})

                while not graph_task.done():
                    yield f"data: {json.dumps({'status': 'extracting_graph'})}\n\n"
                    try:
                        await asyncio.wait_for(asyncio.shield(graph_task), timeout=5.0)
                    except asyncio.TimeoutError:
                        pass
                        
                graph_data = graph_task.result()
                graph_data.setdefault("nodes", []).extend(structural_nodes)
                graph_data.setdefault("edges", []).extend(structural_edges)
                
                # To avoid breaking the SSE format, escape newlines in JSON or rely on json.dumps doing it correctly
                yield f"data: {json.dumps({'status': 'complete', 'markdown': markdown, 'graph': graph_data})}\n\n"
        except Exception as e:
            traceback.print_exc()
            yield f"data: {json.dumps({'status': 'error', 'detail': str(e)})}\n\n"

    return StreamingResponse(sse_generator(), media_type="text/event-stream")
