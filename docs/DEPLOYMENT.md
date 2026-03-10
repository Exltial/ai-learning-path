# 部署指南 / Deployment Guide

本指南详细介绍如何将 AI 学习平台部署到生产环境。

This guide provides detailed instructions for deploying the AI Learning Platform to production environments.

---

## 📋 目录 / Table of Contents

1. [系统要求 / System Requirements](#系统要求--system-requirements)
2. [环境配置 / Environment Configuration](#环境配置--environment-configuration)
3. [本地开发环境 / Local Development](#本地开发环境--local-development)
4. [Docker 部署 / Docker Deployment](#docker-部署--docker-deployment)
5. [生产环境部署 / Production Deployment](#生产环境部署--production-deployment)
6. [Nginx 配置 / Nginx Configuration](#nginx-配置--nginx-configuration)
7. [SSL/HTTPS 配置 / SSL/HTTPS Configuration](#sslhttps-配置--sslhttps-configuration)
8. [数据库迁移 / Database Migration](#数据库迁移--database-migration)
9. [监控与日志 / Monitoring & Logging](#监控与日志--monitoring--logging)
10. [常见问题 / Troubleshooting](#常见问题--troubleshooting)

---

## 系统要求 / System Requirements

### 最低配置 / Minimum Requirements

| 组件 / Component | 要求 / Requirement |
|-----------------|-------------------|
| CPU | 2 核心 / 2 Cores |
| 内存 / RAM | 4 GB |
| 存储 / Storage | 20 GB SSD |
| 操作系统 / OS | Linux (Ubuntu 20.04+ / CentOS 7+) |

### 推荐配置 / Recommended Requirements

| 组件 / Component | 要求 / Requirement |
|-----------------|-------------------|
| CPU | 4 核心 / 4 Cores |
| 内存 / RAM | 8 GB |
| 存储 / Storage | 50 GB SSD |
| 操作系统 / OS | Ubuntu 22.04 LTS |

### 软件依赖 / Software Dependencies

```bash
# Docker & Docker Compose
Docker: 24.0+
Docker Compose: 2.20+

# 开发环境 / Development Environment
Node.js: 18+
Go: 1.21+
PostgreSQL: 15+
Redis: 7+
```

---

## 环境配置 / Environment Configuration

### 后端环境变量 / Backend Environment Variables

创建 `.env` 文件在 `backend/` 目录：

```bash
# ===========================================
# AI Learning Platform - Production Environment
# ===========================================

# -------------------------------------------
# Server Configuration
# -------------------------------------------
PORT=8080
GIN_MODE=release  # MUST be 'release' in production
SERVER_READ_TIMEOUT=30
SERVER_WRITE_TIMEOUT=30

# -------------------------------------------
# Database Configuration (PostgreSQL)
# -------------------------------------------
DATABASE_URL=postgres://username:password@host:5432/database_name?sslmode=require
DATABASE_MAX_CONNECTIONS=50
DATABASE_MAX_IDLE_CONNECTIONS=10
DATABASE_CONNECTION_TIMEOUT=30
DATABASE_MAX_LIFETIME=60

# -------------------------------------------
# Redis Configuration
# -------------------------------------------
REDIS_URL=redis://:password@host:6379/0
REDIS_ADDR=host:6379
REDIS_PASSWORD=your_redis_password

# -------------------------------------------
# JWT Configuration
# CRITICAL: Use strong random key in production!
# Generate with: openssl rand -base64 32
# -------------------------------------------
JWT_SECRET=your-super-secret-jwt-key-min-32-chars
JWT_EXPIRATION_HOURS=24

# -------------------------------------------
# CORS Configuration
# -------------------------------------------
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# -------------------------------------------
# Rate Limiting
# -------------------------------------------
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60

# -------------------------------------------
# File Upload Configuration
# -------------------------------------------
MAX_UPLOAD_SIZE=10485760  # 10MB
UPLOAD_DIR=/app/uploads

# -------------------------------------------
# Logging Configuration
# -------------------------------------------
LOG_LEVEL=info
LOG_FORMAT=json

# -------------------------------------------
# Email Configuration (Optional)
# -------------------------------------------
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password
SMTP_FROM=noreply@yourdomain.com

# -------------------------------------------
# Code Sandbox Configuration
# -------------------------------------------
CODE_EXECUTOR_URL=http://sandbox:8081
CODE_EXECUTOR_TIMEOUT=30
```

### 前端环境变量 / Frontend Environment Variables

创建 `.env` 文件在 `frontend/` 目录：

```bash
# API Configuration
VITE_API_URL=https://api.yourdomain.com/api

# App Configuration
VITE_APP_NAME=AI 学习之路
VITE_APP_VERSION=1.0.0

# Feature Flags
VITE_ENABLE_SANDBOX=true
VITE_ENABLE_ANALYTICS=true

# Analytics (Optional)
VITE_ANALYTICS_ID=UA-XXXXXXXXX-X
```

---

## 本地开发环境 / Local Development

### 1. 克隆项目 / Clone Repository

```bash
git clone https://github.com/exltial/ai-learning-path.git
cd ai-learning-path
```

### 2. 启动数据库和缓存 / Start Database & Cache

```bash
cd backend
docker-compose up -d postgres redis
```

### 3. 初始化数据库 / Initialize Database

```bash
# 连接到 PostgreSQL
docker-compose exec postgres psql -U postgres -c "CREATE DATABASE ai_learning;"

# 导入 Schema
docker-compose exec postgres psql -U postgres -d ai_learning < schema.sql
```

### 4. 启动后端服务 / Start Backend

```bash
cd backend
cp .env.example .env
# 编辑 .env 文件配置
go mod download
go run cmd/main.go
```

### 5. 启动前端服务 / Start Frontend

```bash
cd frontend
npm install
npm run dev
```

### 6. 验证 / Verify

- 前端 / Frontend: http://localhost:5173
- 后端 / Backend: http://localhost:8080
- API 文档 / API Docs: http://localhost:8080/swagger

---

## Docker 部署 / Docker Deployment

### 完整的 Docker Compose 配置

创建 `docker-compose.prod.yml` 在项目根目录：

```yaml
version: '3.8'

services:
  # PostgreSQL Database
  postgres:
    image: postgres:15-alpine
    container_name: ai_learning_db
    environment:
      POSTGRES_USER: ${DB_USER:-ai_learning}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-change_me_in_production}
      POSTGRES_DB: ai_learning
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backend/schema.sql:/docker-entrypoint-initdb.d/schema.sql
    ports:
      - "5432:5432"
    networks:
      - ai_learning_network
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-ai_learning}"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis Cache
  redis:
    image: redis:7-alpine
    container_name: ai_learning_redis
    command: redis-server --requirepass ${REDIS_PASSWORD:-change_me}
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    networks:
      - ai_learning_network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Backend API
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: ai_learning_backend
    environment:
      - PORT=8080
      - GIN_MODE=release
      - DATABASE_URL=postgres://${DB_USER:-ai_learning}:${DB_PASSWORD:-change_me_in_production}@postgres:5432/ai_learning?sslmode=disable
      - REDIS_URL=redis://:${REDIS_PASSWORD:-change_me}@redis:6379/0
      - JWT_SECRET=${JWT_SECRET:-change_me_in_production}
      - CORS_ALLOWED_ORIGINS=${CORS_ALLOWED_ORIGINS:-http://localhost:5173}
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - ai_learning_network
    restart: unless-stopped
    volumes:
      - ./backend/uploads:/app/uploads
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Frontend
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
      args:
        - VITE_API_URL=${VITE_API_URL:-http://localhost:8080/api}
    container_name: ai_learning_frontend
    ports:
      - "80:80"
    depends_on:
      - backend
    networks:
      - ai_learning_network
    restart: unless-stopped

  # Code Sandbox (Optional)
  sandbox:
    build:
      context: ./sandbox
      dockerfile: Dockerfile
    container_name: ai_learning_sandbox
    ports:
      - "8081:8081"
    networks:
      - ai_learning_network
    restart: unless-stopped
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock

volumes:
  postgres_data:
    driver: local
  redis_data:
    driver: local

networks:
  ai_learning_network:
    driver: bridge
```

### 构建并启动 / Build and Start

```bash
# 设置环境变量
export DB_USER=ai_learning
export DB_PASSWORD=$(openssl rand -base64 32)
export REDIS_PASSWORD=$(openssl rand -base64 32)
export JWT_SECRET=$(openssl rand -base64 32)
export CORS_ALLOWED_ORIGINS=https://yourdomain.com
export VITE_API_URL=https://api.yourdomain.com/api

# 构建并启动所有服务
docker-compose -f docker-compose.prod.yml up -d --build

# 查看服务状态
docker-compose -f docker-compose.prod.yml ps

# 查看日志
docker-compose -f docker-compose.prod.yml logs -f
```

### 停止服务 / Stop Services

```bash
docker-compose -f docker-compose.prod.yml down

# 删除数据卷 (谨慎使用!)
docker-compose -f docker-compose.prod.yml down -v
```

---

## 生产环境部署 / Production Deployment

### 方案一：单机部署 / Option 1: Single Server Deployment

适用于小型项目或测试环境。

```bash
# 1. 安装 Docker
curl -fsSL https://get.docker.com | sh

# 2. 克隆项目
git clone https://github.com/exltial/ai-learning-path.git
cd ai-learning-path

# 3. 配置环境变量
cp .env.example .env
# 编辑 .env 文件

# 4. 启动服务
docker-compose -f docker-compose.prod.yml up -d

# 5. 配置 Nginx (见下一节)
```

### 方案二：Kubernetes 部署 / Option 2: Kubernetes Deployment

创建 `k8s/` 目录和相关配置文件：

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ai-learning-backend
  namespace: ai-learning
spec:
  replicas: 3
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend
        image: your-registry/ai-learning-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: jwt-secret
              key: secret
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

### 方案三：云服务器部署 / Option 3: Cloud Deployment

#### AWS 部署

```bash
# 使用 ECS + RDS + ElastiCache
# 1. 创建 RDS PostgreSQL 实例
# 2. 创建 ElastiCache Redis 实例
# 3. 创建 ECS 集群
# 4. 配置 Task Definition
# 5. 部署服务
```

#### 阿里云部署

```bash
# 使用 ACK + RDS + Redis
# 1. 创建 RDS PostgreSQL 实例
# 2. 创建 Redis 实例
# 3. 创建 ACK 集群
# 4. 配置 Deployment 和 Service
# 5. 部署应用
```

---

## Nginx 配置 / Nginx Configuration

### 安装 Nginx

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install nginx

# CentOS/RHEL
sudo yum install nginx
```

### Nginx 配置文件

创建 `/etc/nginx/sites-available/ai-learning`：

```nginx
server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;
    
    # 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com www.yourdomain.com;
    
    # SSL 证书
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    
    # 安全头
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    
    # 前端静态文件
    location / {
        root /var/www/ai-learning;
        try_files $uri $uri/ /index.html;
        
        # 缓存静态资源
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }
    
    # 后端 API 代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # 超时配置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
    
    # Swagger 文档
    location /swagger/ {
        proxy_pass http://localhost:8080/swagger/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    
    # 限制请求频率
    location /api/auth/ {
        proxy_pass http://localhost:8080;
        limit_req zone=auth_limit burst=10 nodelay;
    }
    
    # 日志
    access_log /var/log/nginx/ai-learning.access.log;
    error_log /var/log/nginx/ai-learning.error.log;
}

# 限制请求频率
limit_req_zone $binary_remote_addr zone=auth_limit:10m rate=10r/s;
```

### 启用配置

```bash
# 创建软链接
sudo ln -s /etc/nginx/sites-available/ai-learning /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重启 Nginx
sudo systemctl restart nginx
sudo systemctl enable nginx
```

---

## SSL/HTTPS 配置 / SSL/HTTPS Configuration

### 使用 Let's Encrypt

```bash
# 安装 Certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# 自动续期 (已自动配置 cron)
sudo certbot renew --dry-run
```

### 配置自动续期

```bash
# 编辑 crontab
sudo crontab -e

# 添加以下行
0 3 * * * certbot renew --quiet --post-hook "systemctl reload nginx"
```

---

## 数据库迁移 / Database Migration

### 自动迁移

```bash
# 在容器内运行
docker-compose exec backend ./bin/migrate up

# 或直接运行
cd backend
go run cmd/migrate.go up
```

### 手动迁移

```bash
# 连接到数据库
psql -h localhost -U postgres -d ai_learning

# 运行迁移脚本
\i migrations/20260310120000_create_users_table.sql
```

### 备份数据库

```bash
# 备份
pg_dump -h localhost -U postgres ai_learning > backup_$(date +%Y%m%d).sql

# 恢复
psql -h localhost -U postgres ai_learning < backup_20260310.sql
```

### 定时备份

创建 `/etc/cron.daily/backup-ai-learning`：

```bash
#!/bin/bash
BACKUP_DIR="/backups/ai-learning"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="ai_learning"
DB_USER="postgres"

mkdir -p $BACKUP_DIR
pg_dump -h localhost -U $DB_USER $DB_NAME | gzip > $BACKUP_DIR/backup_$DATE.sql.gz

# 保留最近 7 天的备份
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete
```

---

## 监控与日志 / Monitoring & Logging

### 应用日志

```bash
# 查看后端日志
docker-compose logs -f backend

# 查看前端日志
docker-compose logs -f frontend

# 查看 Nginx 日志
tail -f /var/log/nginx/ai-learning.access.log
tail -f /var/log/nginx/ai-learning.error.log
```

### 性能监控

#### 使用 Prometheus + Grafana

```yaml
# 添加 Prometheus 配置到 docker-compose
prometheus:
  image: prom/prometheus:latest
  volumes:
    - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
  ports:
    - "9090:9090"
  networks:
    - ai_learning_network

grafana:
  image: grafana/grafana:latest
  volumes:
    - grafana_data:/var/lib/grafana
  ports:
    - "3000:3000"
  networks:
    - ai_learning_network
```

### 健康检查

```bash
# 后端健康检查
curl http://localhost:8080/health

# 数据库连接检查
curl http://localhost:8080/health/db

# Redis 连接检查
curl http://localhost:8080/health/redis
```

---

## 常见问题 / Troubleshooting

### 1. 后端无法启动 / Backend Won't Start

**问题 / Issue:**
```
panic: could not connect to database
```

**解决方案 / Solution:**
```bash
# 检查数据库是否运行
docker-compose ps postgres

# 检查数据库连接
docker-compose exec postgres psql -U postgres -c "SELECT 1"

# 查看后端日志
docker-compose logs backend

# 检查环境变量
docker-compose exec backend env | grep DATABASE
```

### 2. 前端无法连接后端 / Frontend Can't Connect to Backend

**问题 / Issue:**
CORS 错误或网络错误

**解决方案 / Solution:**
```bash
# 检查 CORS 配置
grep CORS_ALLOWED_ORIGINS backend/.env

# 确保地址正确
# 应该是：https://yourdomain.com 而不是 http://localhost:8080

# 重启后端
docker-compose restart backend
```

### 3. 数据库迁移失败 / Database Migration Failed

**问题 / Issue:**
```
error: relation "users" already exists
```

**解决方案 / Solution:**
```bash
# 检查迁移状态
docker-compose exec backend ./bin/migrate status

# 回滚迁移
docker-compose exec backend ./bin/migrate down 1

# 重新迁移
docker-compose exec backend ./bin/migrate up
```

### 4. 内存不足 / Out of Memory

**问题 / Issue:**
容器被 OOM Killer 杀死

**解决方案 / Solution:**
```bash
# 增加 Docker 内存限制
# 编辑 docker-compose.prod.yml，添加：
services:
  backend:
    deploy:
      resources:
        limits:
          memory: 1G
        reservations:
          memory: 512M

# 或增加系统 Swap
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

### 5. SSL 证书过期 / SSL Certificate Expired

**解决方案 / Solution:**
```bash
# 手动续期
sudo certbot renew

# 检查自动续期
sudo systemctl status certbot.timer

# 查看续期日志
sudo cat /var/log/letsencrypt/letsencrypt.log
```

### 6. 性能问题 / Performance Issues

**问题 / Issue:**
API 响应慢

**解决方案 / Solution:**
```bash
# 检查数据库慢查询
docker-compose exec postgres psql -U postgres -d ai_learning -c \
  "SELECT query, mean_exec_time FROM pg_stat_statements ORDER BY mean_exec_time DESC LIMIT 10;"

# 检查 Redis 性能
docker-compose exec redis redis-cli INFO stats

# 启用查询日志
# 在 backend/.env 中设置 LOG_LEVEL=debug
```

### 7. 文件上传失败 / File Upload Failed

**问题 / Issue:**
```
http: request body too large
```

**解决方案 / Solution:**
```bash
# 检查 Nginx 配置
# 在 server 块中添加：
client_max_body_size 20M;

# 检查后端配置
# 在 backend/.env 中设置：
MAX_UPLOAD_SIZE=20971520  # 20MB
```

---

## 安全检查清单 / Security Checklist

- [ ] 修改所有默认密码 / Change all default passwords
- [ ] 使用强 JWT Secret / Use strong JWT secret
- [ ] 启用 HTTPS / Enable HTTPS
- [ ] 配置防火墙 / Configure firewall
- [ ] 定期更新依赖 / Update dependencies regularly
- [ ] 启用数据库备份 / Enable database backups
- [ ] 配置日志监控 / Configure log monitoring
- [ ] 限制 API 访问频率 / Rate limit API access
- [ ] 禁用调试模式 / Disable debug mode
- [ ] 审查 CORS 配置 / Review CORS configuration

---

## 性能优化建议 / Performance Tips

1. **启用缓存 / Enable Caching**
   - Redis 缓存热点数据
   - CDN 加速静态资源

2. **数据库优化 / Database Optimization**
   - 添加合适的索引
   - 使用连接池
   - 定期 VACUUM

3. **前端优化 / Frontend Optimization**
   - 代码分割
   - 图片压缩
   - 启用 Gzip/Brotli

4. **使用负载均衡 / Use Load Balancing**
   - 多实例部署后端
   - Nginx 负载均衡

---

## 联系支持 / Support

如遇到问题，请查看：
- [GitHub Issues](https://github.com/exltial/ai-learning-path/issues)
- [API 文档](./API.md)
- [贡献指南](./CONTRIBUTING.md)

---

*最后更新 / Last Updated: 2026-03-10*
