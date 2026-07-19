try:
    from openai_agents import function_tool
except ImportError:
    from agents import function_tool

import os

def convert_document_to_markdown_impl(filepath: str) -> str:
    try:
        from docling.document_converter import DocumentConverter
    except ImportError:
        return "Error: docling is not installed. Please install it with 'pip install docling'."

    # Ensure the filepath is resolved within the workspace context since this tool runs locally on the host
    full_path = os.path.join("/app/workspace", filepath.lstrip("/"))
    if not os.path.exists(full_path):
        return f"Error: Document not found at path {full_path}"

    try:
        from docling.datamodel.pipeline_options import PdfPipelineOptions, AcceleratorOptions, AcceleratorDevice
        from docling.document_converter import DocumentConverter, PdfFormatOption
        from docling.datamodel.base_models import InputFormat
        import torch

        # Check CUDA availability and log
        has_cuda = torch.cuda.is_available()
        print(f"CUDA Available for Docling: {has_cuda}")
        
        accel_options = AcceleratorOptions(num_threads=4, device=AcceleratorDevice.CUDA if has_cuda else AcceleratorDevice.AUTO)
        pipeline_options = PdfPipelineOptions(accelerator_options=accel_options)
        
        format_options = {
            InputFormat.PDF: PdfFormatOption(pipeline_options=pipeline_options)
        }
        
        converter = DocumentConverter(
            allowed_formats=[InputFormat.PDF, InputFormat.DOCX, InputFormat.PPTX, InputFormat.XLSX, InputFormat.HTML, InputFormat.IMAGE],
            format_options=format_options
        )
        
        try:
            result = converter.convert(full_path)
        except Exception as e:
            if "CUDA error" in str(e):
                print("CUDA failed with kernel error, falling back to CPU...")
                accel_options = AcceleratorOptions(num_threads=4, device=AcceleratorDevice.CPU)
                pipeline_options = PdfPipelineOptions(accelerator_options=accel_options)
                format_options = {
                    InputFormat.PDF: PdfFormatOption(pipeline_options=pipeline_options)
                }
                converter = DocumentConverter(
                    allowed_formats=[InputFormat.PDF, InputFormat.DOCX, InputFormat.PPTX, InputFormat.XLSX, InputFormat.HTML, InputFormat.IMAGE],
                    format_options=format_options
                )
                result = converter.convert(full_path)
            else:
                raise e
                
        return result.document.export_to_markdown()
    except Exception as e:
        return f"Error converting document: {str(e)}"

@function_tool
def convert_document_to_markdown(filepath: str) -> str:
    """
    Ingests an unstructured document (PDF, DOCX, PPTX, XLSX, HTML, images) and converts it to Markdown.
    
    Args:
        filepath: The relative path to the document within the workspace (e.g., "docs/architecture.pdf").
        
    Returns:
        The markdown string representation of the document, preserving headings, reading order, and table structures.
    """
    return convert_document_to_markdown_impl(filepath)
