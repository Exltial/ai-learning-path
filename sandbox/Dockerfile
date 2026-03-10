# Python Code Sandbox Dockerfile
# Provides isolated Python execution environment with pre-installed data science libraries

FROM python:3.11-slim

# Set working directory
WORKDIR /sandbox

# Install system dependencies required for Python packages
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    g++ \
    libblas-dev \
    liblapack-dev \
    libatlas-base-dev \
    gfortran \
    && rm -rf /var/lib/apt/lists/*

# Install Python dependencies
RUN pip install --no-cache-dir \
    numpy==1.26.2 \
    pandas==2.1.3 \
    matplotlib==3.8.2 \
    scikit-learn==1.3.2 \
    fastapi==0.109.0 \
    uvicorn==0.27.0 \
    docker==7.0.0 \
    pydantic==2.5.3

# Create non-root user for security
RUN useradd -m -u 1000 sandbox_user

# Set ownership
RUN chown -R sandbox_user:sandbox_user /sandbox

# Switch to non-root user
USER sandbox_user

# Expose port
EXPOSE 8000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD python -c "import requests; requests.get('http://localhost:8000/health', timeout=5)" || exit 1

# Default command
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
