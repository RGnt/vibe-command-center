try:
    from openai_agents import function_tool
except ImportError:
    from agents import function_tool

import os

@function_tool
def generate_knowledge_graph(directory: str) -> str:
    """
    Runs Graphify to parse unstructured code, documentation, and media into a persistent, queryable Knowledge Graph.
    This creates/updates the graph.json and GRAPH_REPORT.md in the workspace.
    
    Args:
        directory: The relative path to the directory to ingest (e.g., "." for the entire workspace, or "backend/").
        
    Returns:
        A summary of the generated knowledge graph, typically reading from GRAPH_REPORT.md.
    """
    try:
        from graphifyy import Graphify
    except ImportError:
        return "Error: graphifyy is not installed. Please install it with 'pip install graphifyy'."

    # Ensure the directory is resolved within the workspace context
    full_path = os.path.join("/app/workspace", directory.lstrip("/"))
    if not os.path.exists(full_path):
        return f"Error: Directory not found at path {full_path}"

    try:
        # Initialize Graphify (it uses OPENAI_API_KEY from environment by default for semantic extraction)
        # Note: Depending on the size of the codebase, Pass 1 (AST) is fast, but Pass 2 (Semantic) might take time.
        graphify = Graphify(input_dir=full_path, output_dir=full_path)
        graphify.run()
        
        report_path = os.path.join(full_path, "GRAPH_REPORT.md")
        if os.path.exists(report_path):
            with open(report_path, "r", encoding="utf-8") as f:
                report = f.read()
            return f"Knowledge graph successfully generated. Report summary:\n\n{report[:2000]}...\n(truncated)"
        return "Knowledge graph generated, but GRAPH_REPORT.md was not found."
    except Exception as e:
        return f"Error generating knowledge graph: {str(e)}"
