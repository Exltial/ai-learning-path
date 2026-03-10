# AI 学习之路 - 后端服务

AI Interactive Learning Platform Backend Service

## 技术栈

- **语言**: Go 1.21+
- **框架**: Gin
- **数据库**: PostgreSQL 15+
- **缓存**: Redis 7+
- **ORM**: GORM / pgx
- **认证**: JWT

## 项目结构

```
backend/
├── cmd/
│   └── main.go              # 应用入口
├── internal/
│   ├── handlers/            # HTTP 处理器
│   │   └── auth_handler.go
│   ├── services/            # 业务逻辑层
│   │   └── auth_service.go
│   ├── repository/          # 数据访问层
│   │   ├── user_repository.go
│   │   └── course_repository.go
│   ├── models/              # 数据模型
│   │   └── models.go
│   └── middleware/          # 中间件
│       └── auth_middleware.go
├── pkg/                     # 公共库
├── configs/                 # 配置文件
├── migrations/              # 数据库迁移
├── schema.sql               # 数据库建表语句
├── API.md                   # API 接口文档
├── go.mod                   # Go 模块定义
└── README.md                # 项目说明
```

## 快速开始

### 环境要求

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (可选)

### 安装依赖

```bash
go mod download
```

### 配置环境变量

```bash
export PORT=8080
export DATABASE_URL=postgres://postgres:password@localhost:5432/ai_learning?sslmode=disable
export REDIS_URL=redis://localhost:6379
export JWT_SECRET=your-secret-key-change-in-production
```

### 初始化数据库

```bash
psql -U postgres -d ai_learning -f schema.sql
```

### 运行服务

```bash
go run cmd/main.go
```

服务将在 `http://localhost:8080` 启动

### 使用 Docker (推荐)

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

## API 文档

完整的 API 文档请参阅 [API.md](./API.md)

### 健康检查

```bash
curl http://localhost:8080/health
```

### 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

## 开发指南

### 代码规范

遵循 Go 官方代码规范，使用 `gofmt` 格式化代码:

```bash
gofmt -w .
```

### 运行测试

```bash
go test ./...
```

### 构建生产版本

```bash
go build -o bin/server cmd/main.go
```

## 数据库迁移

数据库迁移文件放在 `migrations/` 目录下，按时间戳命名:

```
migrations/
├── 20260310120000_create_users_table.sql
├── 20260310120001_create_courses_table.sql
└── ...
```

## 安全注意事项

1. **JWT Secret**: 生产环境必须使用强随机密钥
2. **密码存储**: 使用 bcrypt 加密，cost factor >= 10
3. **HTTPS**: 生产环境必须启用 HTTPS
4. **CORS**: 配置合适的 CORS 策略
5. **Rate Limiting**: 实现 API 限流防止滥用

## 许可证

MIT License
