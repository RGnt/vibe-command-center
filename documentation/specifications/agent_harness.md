# Agent Harness Specification

This document describes the architecture of the Agent Harness, a FastAPI Python service that executes autonomous AI agents against the KanbanX workspace.

## Architecture Overview

The Agent Harness runs as a standalone Python FastAPI server, serving a REST endpoint intended to be consumed by the Frontend (via Nginx reverse proxy). 

**Stack:**
- Language: Python 3.11+
- Framework: FastAPI + Uvicorn
- AI Framework: `openai-agents` (v0.18.3)
- Models: Uses local models (e.g. Gemma4, Llama3) running on the host machine via Ollama compatible endpoints.

## HTTP Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check returning a standard OK status. |
| `/api/v1/agents/run` | POST | Main entrypoint for agent execution. Accepts a JSON payload containing the context and task. |

## RunRequest Schema

The `POST /api/v1/agents/run` endpoint expects the following JSON payload:

```json
{
  "wiki_prompt": "# Agents\nCore instructions from a user's wiki.",
  "task_description": "Task: Title\nDescription: Detailed description from a todo item.",
  "additional_prompt": "Optional extra instructions specific to this run."
}
```

## Execution Flow

1. **Initialization**: The server initializes an `Agent` class instance from the `openai-agents` library, passing the `wiki_prompt` as the system instructions.
2. **Tool Injection**: The agent is provided with python tools wrapped in `function_tool()`.
    - `read_file(filepath: str) -> str`
    - `write_file(filepath: str, content: str) -> str`
3. **Task Launch**: `Runner.run(agent, task_prompt)` is invoked, triggering an autonomous loop where the agent can call tools, read files, and write code iteratively until it decides the task is finished.
4. **Diff Generation**: The harness intercepts all `write_file` calls. At the end of the run, it generates standard unified diff patches comparing the file's original state to its newly written state. It does not overwrite files permanently without tracking the diff.
5. **Rollback**: To keep the environment safe during testing, the harness currently rolls back the files to their original state after generating the diff. The diffs are shipped to the client to be displayed or applied.

## Response Format

A successful agent execution returns the following JSON:

```json
{
  "status": "success",
  "agent_response": "The text response from the agent describing what it did.",
  "changes": [
    {
      "file": "frontend/src/App.tsx",
      "diff": "--- a/frontend/src/App.tsx\n+++ b/frontend/src/App.tsx\n@@ -1,3 +1,4 @@\n..."
    }
  ]
}
```

## Security & Context Limits

- **File Pathing**: All read and write operations are strictly sandboxed and prefixed to the `/app/workspace` directory via `os.path.join`.
- **Environment**: Sensitive environment variables for local LLM routing (`LLM_BASE_URL` and `LLM_API_KEY`) are dynamically loaded, preventing exposure of host networking details to the client.
