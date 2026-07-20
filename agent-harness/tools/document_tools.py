try:
    from openai_agents import function_tool
except ImportError:
    from agents import function_tool

import os

def process_document(filepath: str) -> dict:
    try:
        import pymupdf4llm
    except ImportError:
        return {"error": "Error: pymupdf4llm is not installed."}

    # Ensure the filepath is resolved within the workspace context since this tool runs locally on the host
    full_path = os.path.join("/app/workspace", filepath.lstrip("/"))
    if not os.path.exists(full_path):
        return {"error": f"Error: Document not found at path {full_path}"}

    try:
        # Check if the file is an image or unsupported format, pymupdf handles pdfs and some image types,
        # but if it fails we just return the error.
        md_text = pymupdf4llm.to_markdown(full_path)
        return {
            "markdown": md_text
        }
    except Exception as e:
        return {"error": f"Error converting document: {str(e)}"}

@function_tool
def convert_document_to_markdown(filepath: str) -> str:
    """
    Ingests an unstructured document (PDF, DOCX, PPTX, XLSX, HTML, images) and converts it to Markdown.
    
    Args:
        filepath: The relative path to the document within the workspace (e.g., "docs/architecture.pdf").
        
    Returns:
        The markdown string representation of the document, preserving headings, reading order, and table structures.
    """
    res = process_document(filepath)
    if "error" in res:
        return res["error"]
    return res["markdown"]

def convert_document_to_markdown_impl(filepath: str) -> dict:
    """Used internally by the agent harness ingest endpoint."""
    return process_document(filepath)
