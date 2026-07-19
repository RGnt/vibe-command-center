FROM nvidia/cuda:13.2.1-cudnn-runtime-ubuntu24.04

WORKDIR /app

ENV DEBIAN_FRONTEND=noninteractive

ENV PIP_BREAK_SYSTEM_PACKAGES=1

# Install python and system dependencies required by opencv/docling
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3 \
    python3-venv \
    python3-dev \
    python3-pip \
    libgl1 \
    libglib2.0-0 \
    libxcb1 \
    libxext6 \
    libxrender1 \
    && rm -rf /var/lib/apt/lists/*

RUN update-alternatives --install /usr/bin/python python /usr/bin/python3 1

# Install PyTorch nightly with CUDA 13.2 support (required for RTX 5080 sm_120 Blackwell architecture)
RUN python -m pip install --pre torch torchvision --index-url https://download.pytorch.org/whl/cu132

# Install dependencies
COPY requirements.txt .
RUN python -m pip install --no-cache-dir -r requirements.txt

# Copy source code
COPY . .

# Expose port
EXPOSE 8000

# Run FastAPI app
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
