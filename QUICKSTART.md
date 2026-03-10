# 快速开始指南 / Quick Start Guide

5 分钟内快速启动 AI 学习平台开发环境。

Get the AI Learning Platform development environment up and running in 5 minutes.

---

## 📋 目录 / Table of Contents

1. [系统要求 / Prerequisites](#系统要求--prerequisites)
2. [快速安装 / Quick Install](#快速安装--quick-install)
3. [开发模式 / Development Mode](#开发模式--development-mode)
4. [生产部署 / Production Deployment](#生产部署--production-deployment)
5. [验证安装 / Verify Installation](#验证安装--verify-installation)
6. [下一步 / Next Steps](#下一步--next-steps)
7. [常见问题 / Troubleshooting](#常见问题--troubleshooting)

---

## 系统要求 / Prerequisites

### 必需软件 / Required Software

```bash
# 检查已安装的软件
node --version      # 需要 v18+
npm --version       # 需要 v9+
go version          # 需要 go1.21+
docker --version    # 需要 v24+
docker-compose --version  # 需要 v2.20+
git --version
```

### 安装指南 / Installation Guide

#### macOS

```bash
# 使用 Homebrew
brew install node
brew install go
brew install --cask docker
```

#### Ubuntu/Debian

```bash
# Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
```

#### Windows

```bash
# 使用 Chocolatey
choco install nodejs-lts
choco install golang
choco install docker-desktop
```

---

## 快速安装 / Quick Install

### 方式一：Docker Compose（推荐）/ Option 1: Docker Compose (Recommended)

**最快！5 分钟内完成 / Fastest! Complete in 5 minutes**

```bash
# 1. 克隆项目 / Clone repository
git clone https://github.com/exltial/ai-learning-path.git
cd ai-learning-path

# 2. 配置环境变量 / Configure environment
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env

# 3. 启动所有服务 / Start all services
cd backend
docker-compose up -d

# 4. 查看服务状态 / Check service status
docker-compose ps

# 5. 查看日志 / View logs
docker-compose logs -f
```

**完成！访问以下地址：**
- 前端 / Frontend: http://localhost:5173
- 后端 / Backend: http://localhost:8080
- API 文档 / API Docs: http://localhost:8080/swagger

### 方式二：手动安装 / Option 2: Manual Installation

**适合需要自定义配置的场景 / For custom configurations**

#### 步骤 1: 克隆项目 / Clone Repository

```bash
git clone https://github.com/exltial/ai-learning-path.git
cd ai-learning-path
```

#### 步骤 2: 启动数据库和缓存 / Start Database & Cache

```bash
cd backend
docker-compose up -d postgres redis

# 等待服务就绪 / Wait for services
sleep 10

# 验证服务 / Verify services
docker-compose ps
```

#### 步骤 3: 初始化数据库 / Initialize Database

```bash
# 创建数据库
docker-compose exec postgres psql -U postgres -c "CREATE DATABASE ai_learning;"

# 导入 Schema
docker-compose exec postgres psql -U postgres -d ai_learning -f /docker-entrypoint-initdb.d/schema.sql

# 或使用本地文件
psql -h localhost -U postgres -d ai_learning -f schema.sql
```

#### 步骤 4: 配置后端 / Configure Backend

```bash
cd backend

# 复制环境配置
cp .env.example .env

# 编辑 .env 文件（可选）
# 默认配置已可用，如需修改请编辑：
# - DATABASE_URL
# - REDIS_URL
# - JWT_SECRET
```

#### 步骤 5: 启动后端 / Start Backend

```bash
cd backend

# 下载依赖
go mod download

# 运行服务
go run cmd/main.go

# 或构建后运行
go build -o bin/server cmd/main.go
./bin/server
```

#### 步骤 6: 配置前端 / Configure Frontend

```bash
cd frontend

# 复制环境配置
cp .env.example .env

# 编辑 .env 文件（可选）
# 默认配置：
# VITE_API_URL=http://localhost:8080/api
```

#### 步骤 7: 启动前端 / Start Frontend

```bash
cd frontend

# 安装依赖（首次需要）
npm install

# 开发模式
npm run dev

# 访问 http://localhost:5173
```

---

## 开发模式 / Development Mode

### 启动开发环境 / Start Development Environment

```bash
# 终端 1 - 后端 / Terminal 1 - Backend
cd backend
go run cmd/main.go

# 终端 2 - 前端 / Terminal 2 - Frontend
cd frontend
npm run dev
```

### 热重载 / Hot Reload

- **前端**: 保存文件自动刷新
- **后端**: 使用 air 实现热重载

```bash
# 安装 air
go install github.com/cosmtrek/air@latest

# 热重载运行
cd backend
air
```

### 开发工具 / Development Tools

```bash
# 后端工具
cd backend
go test ./...           # 运行测试
go fmt ./...            # 格式化代码
go vet ./...            # 代码检查

# 前端工具
cd frontend
npm test                # 运行测试
npm run lint            # 代码检查
npm run build           # 构建生产版本
```

---

## 生产部署 / Production Deployment

### Docker 部署 / Docker Deployment

```bash
# 1. 配置生产环境变量
export DB_USER=ai_learning
export DB_PASSWORD=$(openssl rand -base64 32)
export REDIS_PASSWORD=$(openssl rand -base64 32)
export JWT_SECRET=$(openssl rand -base64 32)
export CORS_ALLOWED_ORIGINS=https://yourdomain.com

# 2. 构建并启动
docker-compose -f docker-compose.prod.yml up -d --build

# 3. 查看状态
docker-compose -f docker-compose.prod.yml ps

# 4. 查看日志
docker-compose -f docker-compose.prod.yml logs -f
```

### 详细部署指南 / Detailed Deployment Guide

请查看 [部署指南](docs/DEPLOYMENT.md) 获取完整的生产部署说明。

---

## 验证安装 / Verify Installation

### 健康检查 / Health Checks

```bash
# 检查后端
curl http://localhost:8080/health

# 预期响应 / Expected response:
# {"status":"ok","timestamp":"2026-03-10T12:00:00Z"}

# 检查数据库连接
curl http://localhost:8080/health/db

# 检查 Redis 连接
curl http://localhost:8080/health/redis
```

### 测试 API / Test API

```bash
# 用户注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# 用户登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# 获取课程列表（需要 Token）
curl http://localhost:8080/api/v1/courses \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 访问前端 / Access Frontend

打开浏览器访问：
- http://localhost:5173

应该看到首页，包含：
- 导航栏
- 精选课程
- 平台特色介绍

---

## 下一步 / Next Steps

### 1. 浏览文档 / Browse Documentation

- [项目概述](README.md) - 了解项目特色和功能
- [API 文档](docs/API.md) - 查看完整的 API 接口
- [部署指南](docs/DEPLOYMENT.md) - 生产环境部署
- [贡献指南](docs/CONTRIBUTING.md) - 参与项目开发

### 2. 开始开发 / Start Development

```bash
# 创建功能分支
git checkout -b feature/my-feature

# 开发你的功能
# ...

# 提交代码
git add .
git commit -m "feat: add my awesome feature"
git push origin feature/my-feature
```

### 3. 运行测试 / Run Tests

```bash
# 后端测试
cd backend
go test ./... -v

# 前端测试
cd frontend
npm test
```

### 4. 探索代码库 / Explore Codebase

```bash
# 后端结构
tree backend -L 2

# 前端结构
tree frontend -L 2
```

### 5. 加入社区 / Join Community

- 关注 GitHub 项目
- 报告 Bug 或提出功能建议
- 参与讨论

---

## 常见问题 / Troubleshooting

### 问题 1: 端口被占用 / Port Already in Use

**错误 / Error:**
```
bind: address already in use
```

**解决方案 / Solution:**

```bash
# 查找占用端口的进程
lsof -i :8080  # 后端端口
lsof -i :5173  # 前端端口
lsof -i :5432  # PostgreSQL 端口
lsof -i :6379  # Redis 端口

# 杀死进程
kill -9 <PID>

# 或修改配置文件中的端口
```

### 问题 2: Docker 容器无法启动 / Docker Container Won't Start

**错误 / Error:**
```
Cannot start service backend: driver failed programming external connectivity
```

**解决方案 / Solution:**

```bash
# 重启 Docker
sudo systemctl restart docker

# 或
docker-compose down
docker-compose up -d

# 检查 Docker 日志
docker-compose logs backend
```

### 问题 3: 数据库连接失败 / Database Connection Failed

**错误 / Error:**
```
could not connect to database: connection refused
```

**解决方案 / Solution:**

```bash
# 检查 PostgreSQL 是否运行
docker-compose ps postgres

# 查看 PostgreSQL 日志
docker-compose logs postgres

# 等待数据库就绪（首次启动需要时间）
sleep 30

# 测试连接
docker-compose exec postgres psql -U postgres -c "SELECT 1"
```

### 问题 4: npm install 失败 / npm install Fails

**错误 / Error:**
```
npm ERR! code ENOENT
npm ERR! syscall open
```

**解决方案 / Solution:**

```bash
# 清理缓存
npm cache clean --force

# 删除 node_modules 和 package-lock.json
rm -rf node_modules package-lock.json

# 重新安装
npm install

# 或使用国内镜像
npm config set registry https://registry.npmmirror.com
npm install
```

### 问题 5: Go 依赖下载失败 / Go Dependencies Download Failed

**错误 / Error:**
```
dial tcp: lookup proxy.golang.org: no such host
```

**解决方案 / Solution:**

```bash
# 配置 Go 代理
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.org

# 重新下载依赖
go mod download
```

### 问题 6: 前端页面空白 / Frontend Page is Blank

**问题 / Issue:**
打开 http://localhost:5173 显示空白页面

**解决方案 / Solution:**

```bash
# 检查控制台错误
# 打开浏览器开发者工具查看错误

# 检查 API 连接
# 确认后端服务正在运行

# 清除缓存
rm -rf frontend/node_modules/.vite
npm run dev

# 检查环境变量
cat frontend/.env
# 确保 VITE_API_URL 正确
```

### 问题 7: CORS 错误 / CORS Error

**错误 / Error:**
```
Access to fetch at 'http://localhost:8080' from origin 'http://localhost:5173' has been blocked by CORS policy
```

**解决方案 / Solution:**

```bash
# 检查后端 CORS 配置
# 在 backend/.env 中:
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000

# 重启后端
docker-compose restart backend
# 或
go run cmd/main.go
```

---

## 命令速查 / Command Cheat Sheet

### 后端命令 / Backend Commands

```bash
# 开发模式
go run cmd/main.go

# 构建
go build -o bin/server cmd/main.go

# 测试
go test ./... -v

# 格式化
gofmt -w .

# 下载依赖
go mod download

# 更新依赖
go get -u ./...
```

### 前端命令 / Frontend Commands

```bash
# 开发模式
npm run dev

# 构建
npm run build

# 预览生产构建
npm run preview

# 测试
npm test

# 代码检查
npm run lint

# 格式化
npm run format
```

### Docker 命令 / Docker Commands

```bash
# 启动所有服务
docker-compose up -d

# 停止所有服务
docker-compose down

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 重启服务
docker-compose restart

# 重建服务
docker-compose up -d --build

# 进入容器
docker-compose exec backend bash
docker-compose exec postgres bash
```

---

## 获取帮助 / Get Help

如遇到问题：

1. 查看 [常见问题](#常见问题--troubleshooting)
2. 查看 [部署指南](docs/DEPLOYMENT.md)
3. 查看 [API 文档](docs/API.md)
4. 提交 GitHub Issue
5. 联系维护者：exltial@163.com

---

**🎉 恭喜！你已成功设置开发环境！**

**🎉 Congratulations! You've successfully set up the development environment!**

---

*最后更新 / Last Updated: 2026-03-10*
