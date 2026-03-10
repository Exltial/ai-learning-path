# AI Learning Platform - Backend Development Guide

## Phase 2 核心功能开发完成

### ✅ 已完成功能

#### 1. 用户认证系统（JWT）
- **用户注册 API** (`POST /api/v1/auth/register`)
  - 输入验证（用户名 3-50 字符，邮箱格式，密码 8+ 字符）
  - 用户名/邮箱唯一性检查
  - 密码 bcrypt 加密
  - 返回 JWT token

- **用户登录 API** (`POST /api/v1/auth/login`)
  - 邮箱/密码验证
  - JWT token 生成
  - 最后登录时间更新

- **Token 刷新 API** (`POST /api/v1/auth/refresh`)
  - 使用旧 token 换取新 token

- **JWT 中间件**
  - Token 验证
  - 用户信息注入到上下文
  - 角色权限检查

#### 2. 课程管理 API
- **课程列表** (`GET /api/v1/courses`)
  - 分页支持
  - 按分类/难度筛选
  - 只返回已发布课程

- **课程详情** (`GET /api/v1/courses/:id`)
  - 完整课程信息

- **创建课程** (`POST /api/v1/courses`)
  - 仅讲师/管理员
  - 输入验证

- **更新课程** (`PUT /api/v1/courses/:id`)
  - 部分更新支持

- **删除课程** (`DELETE /api/v1/courses/:id`)
  - 仅管理员

- **课程注册** (`POST /api/v1/courses/:id/enroll`)
  - 检查课程发布状态
  - 防止重复注册

- **课程评价** (`GET/POST /api/v1/courses/:id/reviews`)
  - 评分 1-5
  - 评论

#### 3. 作业提交 API
- **提交作业** (`POST /api/v1/exercises/:id/submit`)
  - 支持多种题型：
    - 选择题（自动判分）
    - 判断题（自动判分）
    - 填空题（自动判分）
    - 编程题（待实现代码执行）
    - 问答题（需人工评分）
  - 尝试次数限制
  - 自动反馈

- **查看提交** (`GET /api/v1/submissions/:id`)
  - 提交详情
  - 判分结果

- **提交历史** (`GET /api/v1/exercises/:id/submissions`)
  - 用户所有提交记录

- **人工评分** (`POST /api/v1/submissions/:id/grade`)
  - 仅讲师/管理员
  - 分数/反馈

#### 4. 数据库连接配置
- **配置文件** (`configs/config.go`)
  - 从环境变量加载配置
  - 数据库连接池配置
  - Redis 配置
  - JWT 配置

- **环境变量** (`.env.example`)
  - 完整的配置模板
  - 详细的注释说明

#### 5. API 文档（Swagger）
- **Swagger UI** (`/swagger/index.html`)
  - 交互式 API 文档
  - 可直接测试 API
  - 完整的请求/响应示例

### 📁 文件结构

```
backend/
├── cmd/
│   └── main.go                 # 应用入口，路由配置
├── configs/
│   └── config.go               # 配置管理
├── internal/
│   ├── handlers/
│   │   ├── auth_handler.go     # 认证处理器
│   │   ├── course_handler.go   # 课程处理器
│   │   ├── submission_handler.go # 提交处理器
│   │   └── ...
│   ├── services/
│   │   ├── auth_service.go     # 认证服务
│   │   ├── auth_service_test.go # 认证测试
│   │   ├── course_service.go   # 课程服务
│   │   ├── course_service_test.go # 课程测试
│   │   ├── submission_service.go # 提交服务
│   │   ├── submission_service_test.go # 提交测试
│   │   └── user_service.go     # 用户服务
│   ├── repository/
│   │   ├── user_repository.go
│   │   ├── course_repository.go
│   │   ├── submission_repository.go
│   │   └── ...
│   ├── middleware/
│   │   └── auth_middleware.go  # JWT 中间件
│   └── models/
│       └── models.go           # 数据模型
├── docs/
│   └── docs.go                 # Swagger 文档
├── .env.example                # 环境变量模板
├── go.mod                      # Go 模块定义
└── DEVELOPMENT.md              # 本文档
```

### 🚀 快速开始

#### 1. 安装依赖

```bash
cd /home/admin/.openclaw/workspace/projects/ai-learning-platform/backend

# 安装 Go 依赖
go mod tidy

# 安装 Swagger 工具
go install github.com/swaggo/swag/cmd/swag@latest
```

#### 2. 配置环境

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件，配置数据库和 Redis
# 或使用默认配置（本地 PostgreSQL + Redis）
```

#### 3. 启动数据库和 Redis

```bash
# 使用 Docker Compose
docker-compose up -d

# 或手动启动 PostgreSQL 和 Redis
```

#### 4. 运行应用

```bash
# 开发模式
go run cmd/main.go

# 或构建后运行
go build -o ai-learning-platform cmd/main.go
./ai-learning-platform
```

#### 5. 访问 Swagger 文档

打开浏览器访问：http://localhost:8080/swagger/index.html

### 🧪 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/services/...

# 带覆盖率
go test -cover ./...

# 详细输出
go test -v ./internal/services/...
```

### 📝 代码规范

- 所有公共函数都有注释
- 错误处理完整
- 输入验证严格
- 单元测试覆盖核心逻辑
- 遵循 Go 最佳实践

### 🔐 安全特性

- 密码 bcrypt 加密
- JWT token 认证
- 角色权限控制
- CORS 配置
- 输入验证
- SQL 注入防护（使用参数化查询）

### 📊 API 响应格式

所有 API 响应遵循统一格式：

**成功响应：**
```json
{
  "success": true,
  "message": "操作成功",
  "data": { ... }
}
```

**错误响应：**
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "错误描述"
  }
}
```

### 🔧 配置说明

| 配置项 | 环境变量 | 默认值 | 说明 |
|--------|---------|--------|------|
| 端口 | PORT | 8080 | HTTP 服务端口 |
| 数据库 URL | DATABASE_URL | postgres://... | PostgreSQL 连接字符串 |
| Redis 地址 | REDIS_ADDR | localhost:6379 | Redis 服务器地址 |
| JWT 密钥 | JWT_SECRET | your-secret-key... | JWT 签名密钥 |
| JWT 过期 | JWT_EXPIRATION_HOURS | 24 | Token 过期时间（小时） |

### 📚 下一步计划

- [ ] 实现代码执行沙箱（编程题自动判分）
- [ ] 实现学习进度追踪
- [ ] 实现用户成就系统
- [ ] 实现讨论区功能
- [ ] 添加邮件通知
- [ ] 性能优化和压力测试

### 🤝 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

---

**开发者**: AI 学习平台后端团队  
**版本**: 2.0.0  
**最后更新**: 2026-03-10
