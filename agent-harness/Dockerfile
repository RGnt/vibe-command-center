FROM python:3.11-slim

WORKDIR /app

# Install system dependencies required by opencv/docling
RUN apt-get update && apt-get install -y --no-install-recommends \
    libgl1 \
    libglib2.0-0 \
    libxcb1 \
    libxext6 \
    libxrender1 \
    && rm -rf /var/lib/apt/lists/*

# Install PyTorch with CUDA 12.4 support first (required for RTX 50 series)
RUN pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu124

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy source code
COPY . .

# Expose port
EXPOSE 8000

# Run FastAPI app
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
