"""
AI Learning Platform - Code Sandbox Service
Provides secure, isolated Python code execution with resource limits.
"""

import asyncio
import docker
import io
import json
import tarfile
import tempfile
import uuid
from datetime import datetime
from typing import Optional
from contextlib import asynccontextmanager

from fastapi import FastAPI, HTTPException, BackgroundTasks
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, Field


# ============== Data Models ==============

class CodeExecutionRequest(BaseModel):
    """Request model for code execution"""
    code: str = Field(..., description="Python code to execute", min_length=1, max_length=50000)
    timeout: int = Field(default=30, description="Execution timeout in seconds", ge=1, le=300)
    memory_limit: str = Field(default="512m", description="Memory limit (e.g., '512m', '1g')")
    cpu_limit: float = Field(default=1.0, description="CPU limit (number of cores)", ge=0.1, le=4.0)


class CodeExecutionResponse(BaseModel):
    """Response model for code execution"""
    success: bool
    stdout: str
    stderr: str
    return_value: Optional[str] = None
    execution_time: float
    memory_used: Optional[str] = None
    error_type: Optional[str] = None
    container_id: Optional[str] = None


class HealthResponse(BaseModel):
    """Health check response"""
    status: str
    timestamp: str
    docker_available: bool
    version: str


# ============== Global State ==============

docker_client = None
active_containers = {}


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan manager"""
    global docker_client
    # Startup
    try:
        docker_client = docker.from_env()
        print("✓ Docker client initialized")
    except Exception as e:
        print(f"✗ Docker initialization failed: {e}")
        docker_client = None
    yield
    # Shutdown
    await cleanup_containers()
    if docker_client:
        docker_client.close()


app = FastAPI(
    title="AI Learning Platform - Code Sandbox",
    description="Secure Python code execution service with Docker isolation",
    version="1.0.0",
    lifespan=lifespan
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


# ============== Helper Functions ==============

async def cleanup_containers():
    """Clean up all active containers"""
    global active_containers
    for container_id, container_info in list(active_containers.items()):
        try:
            container = container_info.get("container")
            if container:
                container.stop(timeout=5)
                container.remove(force=True)
                print(f"✓ Container {container_id} cleaned up")
        except Exception as e:
            print(f"✗ Error cleaning up container {container_id}: {e}")
    active_containers.clear()


def create_sandbox_script(code: str) -> str:
    """Create execution script that captures output and return value"""
    return f'''
import sys
import json
from io import StringIO

# Capture stdout
old_stdout = sys.stdout
sys.stdout = StringIO()

# Capture stderr
old_stderr = sys.stderr
sys.stderr = StringIO()

result = None
error = None

try:
    # Execute user code
    exec_globals = {{}}
    exec({code!r}, exec_globals)
    
    # Try to get return value from 'result' variable if exists
    if 'result' in exec_globals:
        result = exec_globals['result']
    elif '_' in exec_globals:
        result = exec_globals['_']
        
except Exception as e:
    error = {{
        "type": type(e).__name__,
        "message": str(e)
    }}

# Get outputs
stdout_output = sys.stdout.getvalue()
stderr_output = sys.stderr.getvalue()

# Restore streams
sys.stdout = old_stdout
sys.stderr = old_stderr

# Prepare result
output = {{
    "stdout": stdout_output,
    "stderr": stderr_output,
    "result": result,
    "error": error
}}

print(json.dumps(output))
'''


async def execute_in_container(
    code: str,
    timeout: int = 30,
    memory_limit: str = "512m",
    cpu_limit: float = 1.0
) -> CodeExecutionResponse:
    """Execute code in isolated Docker container"""
    global docker_client
    
    if not docker_client:
        raise HTTPException(status_code=503, detail="Docker service unavailable")
    
    container_id = f"sandbox-{uuid.uuid4().hex[:8]}"
    start_time = datetime.now()
    
    try:
        # Create execution script
        script = create_sandbox_script(code)
        
        # Create container with resource limits
        container = docker_client.containers.run(
            "ai-learning-sandbox:latest",
            command=["python", "-c", script],
            name=container_id,
            detach=True,
            mem_limit=memory_limit,
            nano_cpus=int(cpu_limit * 1e9),  # Convert to nano CPUs
            network_disabled=True,  # No network access for security
            read_only=True,  # Read-only filesystem
            tmpfs={"/tmp": "rw,noexec,nosuid,size=100m"},  # Writable temp directory
            security_opt=["no-new-privileges:true"],  # Prevent privilege escalation
            cap_drop=["ALL"],  # Drop all capabilities
            pids_limit=100,  # Limit process count
        )
        
        # Track active container
        active_containers[container_id] = {
            "container": container,
            "created_at": start_time
        }
        
        # Wait for completion with timeout
        try:
            result = container.wait(timeout=timeout)
            exit_code = result.get("StatusCode", -1) if isinstance(result, dict) else -1
        except Exception as e:
            # Timeout or other error
            container.stop(timeout=2)
            raise TimeoutError(f"Execution exceeded {timeout}s timeout") from e
        
        # Get logs
        logs = container.logs().decode('utf-8', errors='replace')
        
        # Parse output
        stdout = ""
        stderr = ""
        return_value = None
        error_type = None
        
        try:
            output = json.loads(logs.strip())
            stdout = output.get("stdout", "")
            stderr = output.get("stderr", "")
            return_value = str(output.get("result")) if output.get("result") is not None else None
            if output.get("error"):
                error_type = output["error"].get("type", "UnknownError")
                if not stderr:
                    stderr = output["error"].get("message", "")
        except json.JSONDecodeError:
            # If not JSON, treat as raw output
            stdout = logs
            if exit_code != 0:
                stderr = "Execution failed"
                error_type = "ExecutionError"
        
        # Get container stats
        memory_used = None
        try:
            stats = container.stats(stream=False)
            if "memory_stats" in stats:
                usage = stats["memory_stats"].get("usage", 0)
                if usage:
                    memory_used = f"{usage / (1024 * 1024):.2f}MB"
        except Exception:
            pass
        
        # Calculate execution time
        execution_time = (datetime.now() - start_time).total_seconds()
        
        # Cleanup container
        container.remove(force=True)
        if container_id in active_containers:
            del active_containers[container_id]
        
        return CodeExecutionResponse(
            success=exit_code == 0 and not error_type,
            stdout=stdout,
            stderr=stderr,
            return_value=return_value,
            execution_time=execution_time,
            memory_used=memory_used,
            error_type=error_type,
            container_id=container_id
        )
        
    except docker.errors.ImageNotFound:
        raise HTTPException(
            status_code=503,
            detail="Sandbox image not found. Please build the Docker image first."
        )
    except docker.errors.DockerException as e:
        raise HTTPException(status_code=500, detail=f"Docker error: {str(e)}")
    except TimeoutError as e:
        return CodeExecutionResponse(
            success=False,
            stdout="",
            stderr=str(e),
            execution_time=timeout,
            error_type="TimeoutError",
            container_id=container_id
        )
    except Exception as e:
        # Cleanup on error
        try:
            container = docker_client.containers.get(container_id)
            container.stop(timeout=2)
            container.remove(force=True)
        except Exception:
            pass
        if container_id in active_containers:
            del active_containers[container_id]
        raise


# ============== API Endpoints ==============

@app.get("/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint"""
    docker_available = False
    try:
        if docker_client:
            docker_client.ping()
            docker_available = True
    except Exception:
        pass
    
    return HealthResponse(
        status="healthy" if docker_available else "degraded",
        timestamp=datetime.now().isoformat(),
        docker_available=docker_available,
        version="1.0.0"
    )


@app.post("/execute", response_model=CodeExecutionResponse)
async def execute_code(request: CodeExecutionRequest):
    """
    Execute Python code in an isolated sandbox environment.
    
    - **code**: Python code to execute (max 50KB)
    - **timeout**: Execution timeout in seconds (1-300)
    - **memory_limit**: Memory limit (e.g., '512m', '1g')
    - **cpu_limit**: CPU limit in cores (0.1-4.0)
    """
    return await execute_in_container(
        code=request.code,
        timeout=request.timeout,
        memory_limit=request.memory_limit,
        cpu_limit=request.cpu_limit
    )


@app.get("/containers")
async def list_active_containers():
    """List currently active sandbox containers"""
    containers = []
    for cid, info in active_containers.items():
        containers.append({
            "container_id": cid,
            "created_at": info["created_at"].isoformat(),
            "age_seconds": (datetime.now() - info["created_at"]).total_seconds()
        })
    return {"active_containers": containers, "count": len(containers)}


@app.delete("/containers/{container_id}")
async def stop_container(container_id: str):
    """Stop and remove a specific sandbox container"""
    if container_id not in active_containers:
        raise HTTPException(status_code=404, detail="Container not found")
    
    try:
        container_info = active_containers[container_id]
        container = container_info.get("container")
        if container:
            container.stop(timeout=5)
            container.remove(force=True)
        del active_containers[container_id]
        return {"message": f"Container {container_id} stopped and removed"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/")
async def root():
    """Root endpoint with API information"""
    return {
        "name": "AI Learning Platform - Code Sandbox",
        "version": "1.0.0",
        "description": "Secure Python code execution service",
        "endpoints": {
            "health": "GET /health",
            "execute": "POST /execute",
            "containers": "GET /containers",
            "docs": "GET /docs"
        }
    }


# ============== Background Cleanup ==============

async def periodic_cleanup():
    """Periodically clean up old containers"""
    while True:
        await asyncio.sleep(300)  # Every 5 minutes
        now = datetime.now()
        to_remove = []
        
        for cid, info in active_containers.items():
            age = (now - info["created_at"]).total_seconds()
            if age > 600:  # Remove containers older than 10 minutes
                to_remove.append(cid)
        
        for cid in to_remove:
            try:
                container_info = active_containers[cid]
                container = container_info.get("container")
                if container:
                    container.stop(timeout=2)
                    container.remove(force=True)
                del active_containers[cid]
                print(f"✓ Cleaned up old container {cid}")
            except Exception as e:
                print(f"✗ Error cleaning up {cid}: {e}")


# Start background cleanup task
@app.on_event("startup")
async def start_cleanup_task():
    asyncio.create_task(periodic_cleanup())


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
