# AI Learning Platform - Code Sandbox Service

安全、隔离的 Python 代码执行服务，使用 Docker 容器进行资源限制和环境隔离。

## 🎯 功能特性

- **Docker 隔离**: 每个代码执行请求在独立容器中运行
- **资源限制**: 可配置 CPU、内存、执行时间限制
- **输出捕获**: 自动捕获 stdout、stderr 和返回值
- **预装库**: 内置常用数据科学库（numpy, pandas, matplotlib, scikit-learn）
- **安全加固**: 禁用网络、只读文件系统、权限限制
- **健康检查**: 内置健康检查和自动清理机制

## 📁 项目结构

```
sandbox/
├── Dockerfile              # Docker 镜像构建文件
├── main.py                 # FastAPI 沙箱服务
├── docker-compose.yml      # Docker Compose 配置
├── README.md               # 使用说明
└── nginx.conf              # Nginx 配置（生产环境可选）
```

## 🚀 快速开始

### 1. 构建 Docker 镜像

```bash
cd /home/admin/.openclaw/workspace/projects/ai-learning-platform/sandbox
docker-compose build
```

### 2. 启动服务

```bash
# 仅启动沙箱服务
docker-compose up -d sandbox

# 或启动完整栈（包括 Nginx）
docker-compose --profile production up -d
```

### 3. 验证服务

```bash
curl http://localhost:8000/health
```

预期响应:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00",
  "docker_available": true,
  "version": "1.0.0"
}
```

## 📡 API 使用

### 执行代码

**请求:**
```bash
curl -X POST http://localhost:8000/execute \
  -H "Content-Type: application/json" \
  -d '{
    "code": "import numpy as np\nresult = np.array([1, 2, 3]) * 2\nprint(result)",
    "timeout": 30,
    "memory_limit": "512m",
    "cpu_limit": 1.0
  }'
```

**响应:**
```json
{
  "success": true,
  "stdout": "[2 4 6]\n",
  "stderr": "",
  "return_value": "[2 4 6]",
  "execution_time": 0.523,
  "memory_used": "45.32MB",
  "error_type": null,
  "container_id": "sandbox-a1b2c3d4"
}
```

### 使用 pandas 示例

```bash
curl -X POST http://localhost:8000/execute \
  -H "Content-Type: application/json" \
  -d '{
    "code": "import pandas as pd\nimport numpy as np\ndf = pd.DataFrame({\"A\": [1, 2, 3], \"B\": [4, 5, 6]})\nprint(df.describe())\nresult = df[\"A\"].mean()",
    "timeout": 30
  }'
```

### 使用 matplotlib 示例

```bash
curl -X POST http://localhost:8000/execute \
  -H "Content-Type: application/json" \
  -d '{
    "code": "import matplotlib\nmatplotlib.use(\"Agg\")\nimport matplotlib.pyplot as plt\nimport numpy as np\nx = np.linspace(0, 10, 100)\ny = np.sin(x)\nplt.plot(x, y)\nplt.savefig(\"/tmp/plot.png\")\nprint(\"Plot saved to /tmp/plot.png\")",
    "timeout": 30
  }'
```

## 🔧 API 端点

| 端点 | 方法 | 描述 |
|------|------|------|
| `/` | GET | API 信息 |
| `/health` | GET | 健康检查 |
| `/execute` | POST | 执行代码 |
| `/containers` | GET | 列出活跃容器 |
| `/containers/{id}` | DELETE | 停止指定容器 |
| `/docs` | GET | Swagger API 文档 |

## 📋 执行请求参数

| 参数 | 类型 | 默认值 | 范围 | 描述 |
|------|------|--------|------|------|
| `code` | string | - | 1-50000 字符 | Python 代码 |
| `timeout` | integer | 30 | 1-300 秒 | 执行超时 |
| `memory_limit` | string | "512m" | - | 内存限制 (如 "512m", "1g") |
| `cpu_limit` | float | 1.0 | 0.1-4.0 | CPU 核心数限制 |

## 🔒 安全特性

### 容器安全配置

- **网络禁用**: 容器无法访问外部网络
- **只读文件系统**: 防止文件系统修改
- **权限限制**: 禁止提权 (`no-new-privileges`)
- **能力丢弃**: 丢弃所有 Linux capabilities
- **进程限制**: 最多 100 个进程
- **临时目录**: 仅 `/tmp` 可写，大小限制 100MB

### 资源限制

- CPU 使用率可配置 (0.1-4.0 核心)
- 内存使用可配置 (默认 512MB)
- 执行时间可配置 (1-300 秒)
- 自动清理超时容器

## 🛠️ 开发指南

### 本地测试

```bash
# 构建镜像
docker build -t ai-learning-sandbox:latest .

# 运行容器
docker run -d \
  -p 8000:8000 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --name test-sandbox \
  ai-learning-sandbox:latest

# 查看日志
docker logs -f test-sandbox

# 停止并清理
docker stop test-sandbox && docker rm test-sandbox
```

### Python 客户端示例

```python
import requests

SANDBOX_URL = "http://localhost:8000"

def execute_code(code: str, timeout: int = 30):
    """Execute Python code in sandbox"""
    response = requests.post(
        f"{SANDBOX_URL}/execute",
        json={
            "code": code,
            "timeout": timeout,
            "memory_limit": "512m",
            "cpu_limit": 1.0
        }
    )
    return response.json()

# 示例
result = execute_code("print('Hello from sandbox!')")
print(result)
```

### 错误处理

```python
result = execute_code("1 / 0")
if not result["success"]:
    print(f"Error: {result['error_type']}")
    print(f"Stderr: {result['stderr']}")
```

## 📊 监控与日志

### 查看服务日志

```bash
docker-compose logs -f sandbox
```

### 查看活跃容器

```bash
curl http://localhost:8000/containers
```

### 查看 Docker 容器

```bash
docker ps | grep ai-learning
```

## 🧹 维护

### 停止服务

```bash
docker-compose down
```

### 清理所有容器

```bash
docker-compose down -v
docker rmi ai-learning-sandbox:latest
```

### 更新服务

```bash
docker-compose pull
docker-compose build
docker-compose up -d --force-recreate
```

## ⚠️ 注意事项

1. **Docker Socket**: 服务需要访问 `/var/run/docker.sock` 以创建子容器
2. **资源消耗**: 每个代码执行请求会创建一个新容器，注意资源使用
3. **代码大小**: 代码限制为 50KB，适合执行脚本而非大型程序
4. **超时设置**: 建议设置合理的超时时间，避免资源占用
5. **生产环境**: 生产环境建议添加认证、速率限制和审计日志

## 📝 返回值说明

代码执行后，如果定义了 `result` 变量，其值会被捕获并返回：

```python
# 示例 1: 显式定义 result
result = [1, 2, 3, 4, 5]

# 示例 2: 计算结果
import numpy as np
data = np.random.randn(100)
result = {"mean": data.mean(), "std": data.std()}

# 示例 3: 无返回值
print("Hello")  # return_value 为 null
```

## 🎓 使用场景

- **在线编程练习**: 安全执行用户提交的代码
- **数据科学教学**: 运行 numpy/pandas/matplotlib 示例
- **代码评测**: 自动化测试和评分
- **实验环境**: 快速验证代码片段
- **API 后端**: 为学习平台提供代码执行能力

## 📄 License

MIT License

---

**AI Learning Platform** - 让代码学习更安全、更简单 🚀
