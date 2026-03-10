# Phase 2 功能测试报告

**项目名称**: AI 学习平台  
**测试阶段**: Phase 2 核心功能  
**测试日期**: 2026-03-10  
**测试工程师**: AI 测试工程师  
**报告版本**: v1.0

---

## 📋 执行摘要

### 测试环境状态

| 服务 | 预期端口 | 状态 | 说明 |
|------|----------|------|------|
| 后端 API | 8080 | ❌ 未运行 | 端口被 SearXNG 占用，Go 环境未安装 |
| 前端 | 5173 | ❌ 未运行 | 需要启动 Vite 开发服务器 |
| 沙箱服务 | 8000 | ❌ 未运行 | 需要启动 Docker 容器 |

### 测试结论

由于服务未启动，本次测试主要完成：
- ✅ 测试用例设计
- ✅ API 测试脚本编写
- ✅ 测试环境配置文档
- ⏳ 功能测试（待服务启动后执行）

---

## 🎯 测试范围

### 1. 后端 API 测试

#### 1.1 认证模块

| 测试项 | API 端点 | 测试状态 | 优先级 |
|--------|----------|----------|--------|
| 用户注册 - 正常流程 | POST /api/v1/auth/register | ⏳ 待测试 | P0 |
| 用户注册 - 用户名重复 | POST /api/v1/auth/register | ⏳ 待测试 | P0 |
| 用户注册 - 邮箱重复 | POST /api/v1/auth/register | ⏳ 待测试 | P0 |
| 用户注册 - 密码强度 | POST /api/v1/auth/register | ⏳ 待测试 | P1 |
| 用户注册 - 参数验证 | POST /api/v1/auth/register | ⏳ 待测试 | P1 |
| 用户登录 - 正常流程 | POST /api/v1/auth/login | ⏳ 待测试 | P0 |
| 用户登录 - 错误凭证 | POST /api/v1/auth/login | ⏳ 待测试 | P0 |
| 用户登录 - 用户不存在 | POST /api/v1/auth/login | ⏳ 待测试 | P1 |
| Token 刷新 | POST /api/v1/auth/refresh | ⏳ 待测试 | P2 |
| 用户登出 | POST /api/v1/auth/logout | ⏳ 待测试 | P2 |

#### 1.2 课程模块

| 测试项 | API 端点 | 测试状态 | 优先级 |
|--------|----------|----------|--------|
| 课程列表 - 正常查询 | GET /api/v1/courses | ⏳ 待测试 | P0 |
| 课程列表 - 分页 | GET /api/v1/courses?page=1&limit=10 | ⏳ 待测试 | P1 |
| 课程列表 - 分类筛选 | GET /api/v1/courses?category=Programming | ⏳ 待测试 | P1 |
| 课程列表 - 难度筛选 | GET /api/v1/courses?difficulty=beginner | ⏳ 待测试 | P1 |
| 课程详情 - 正常查询 | GET /api/v1/courses/:id | ⏳ 待测试 | P0 |
| 课程详情 - 不存在 | GET /api/v1/courses/invalid-id | ⏳ 待测试 | P1 |
| 创建课程 - 讲师权限 | POST /api/v1/courses | ⏳ 待测试 | P2 |
| 注册课程 | POST /api/v1/courses/:id/enroll | ⏳ 待测试 | P1 |

#### 1.3 作业提交模块

| 测试项 | API 端点 | 测试状态 | 优先级 |
|--------|----------|----------|--------|
| 提交选择题 | POST /api/v1/exercises/:id/submit | ⏳ 待测试 | P0 |
| 提交判断题 | POST /api/v1/exercises/:id/submit | ⏳ 待测试 | P0 |
| 提交填空题 | POST /api/v1/exercises/:id/submit | ⏳ 待测试 | P0 |
| 提交编程题 | POST /api/v1/exercises/:id/submit | ⏳ 待测试 | P1 |
| 提交问答题 | POST /api/v1/exercises/:id/submit | ⏳ 待测试 | P1 |
| 达到最大尝试次数 | POST /api/v1/exercises/:id/submit | ⏳ 待测试 | P1 |
| 获取提交历史 | GET /api/v1/exercises/:id/submissions | ⏳ 待测试 | P2 |
| 人工评分 | POST /api/v1/submissions/:id/grade | ⏳ 待测试 | P2 |

### 2. 前端功能测试

| 测试页面 | 测试项 | 测试状态 | 优先级 |
|----------|--------|----------|--------|
| 登录页 | 表单渲染 | ⏳ 待测试 | P0 |
| 登录页 | 表单验证 | ⏳ 待测试 | P0 |
| 登录页 | 错误提示 | ⏳ 待测试 | P0 |
| 注册页 | 表单验证 | ⏳ 待测试 | P0 |
| 注册页 | 密码强度提示 | ⏳ 待测试 | P1 |
| 课程列表页 | 课程卡片渲染 | ⏳ 待测试 | P0 |
| 课程列表页 | 搜索功能 | ⏳ 待测试 | P1 |
| 课程列表页 | 筛选功能 | ⏳ 待测试 | P1 |
| 课程列表页 | 分页功能 | ⏳ 待测试 | P1 |
| 课程详情页 | 课程信息展示 | ⏳ 待测试 | P0 |
| 课程详情页 | 章节大纲展示 | ⏳ 待测试 | P0 |
| 课程详情页 | 进度显示 | ⏳ 待测试 | P1 |
| 全局 | 响应式布局 | ⏳ 待测试 | P2 |
| 全局 | 浏览器兼容性 | ⏳ 待测试 | P2 |

### 3. 沙箱服务测试

| 测试项 | 测试状态 | 优先级 |
|--------|----------|--------|
| Python 基础代码执行 | ⏳ 待测试 | P0 |
| numpy 库使用 | ⏳ 待测试 | P1 |
| pandas 库使用 | ⏳ 待测试 | P1 |
| 超时限制测试 | ⏳ 待测试 | P1 |
| 内存限制测试 | ⏳ 待测试 | P1 |
| 错误处理 | ⏳ 待测试 | P0 |

### 4. 集成测试

| 测试场景 | 测试状态 | 优先级 |
|----------|----------|--------|
| 用户注册 → 登录 → 浏览课程 | ⏳ 待测试 | P0 |
| 用户注册 → 登录 → 开始学习 → 提交作业 | ⏳ 待测试 | P0 |
| 完整学习流程 | ⏳ 待测试 | P1 |

---

## 🔧 测试环境配置

### 启动服务命令

#### 1. 启动后端服务

```bash
# 方案 A: 使用 Docker（推荐）
cd /home/admin/.openclaw/workspace/projects/ai-learning-platform/backend
docker-compose up -d

# 方案 B: 本地运行（需要安装 Go 1.21+）
cd /home/admin/.openclaw/workspace/projects/ai-learning-platform/backend
cp .env.example .env
# 编辑 .env 配置数据库连接
go mod tidy
go run cmd/main.go
```

#### 2. 启动前端服务

```bash
cd /home/admin/.openclaw/workspace/projects/ai-learning-platform/frontend
npm install
npm run dev
```

#### 3. 启动沙箱服务

```bash
cd /home/admin/.openclaw/workspace/projects/ai-learning-platform/sandbox
docker-compose up -d
```

### 测试工具

- API 测试：curl, Postman, 或提供的 test-api.sh 脚本
- 前端测试：Chrome DevTools, Playwright
- 性能测试：wrk, ab, JMeter

---

## 📊 测试用例详情

### 认证模块测试用例

#### TC-AUTH-001: 用户注册 - 正常流程

**前置条件**: 数据库干净，无重复用户

**测试步骤**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Password123"
  }'
```

**预期结果**:
- 状态码：201 Created
- 响应包含 user 和 token 字段
- 用户信息正确

#### TC-AUTH-002: 用户注册 - 用户名重复

**前置条件**: 已存在用户名 "testuser"

**测试步骤**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test2@example.com",
    "password": "Password123"
  }'
```

**预期结果**:
- 状态码：409 Conflict
- 错误码：ErrUsernameExists

#### TC-AUTH-003: 用户登录 - 正常流程

**前置条件**: 用户已注册

**测试步骤**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123"
  }'
```

**预期结果**:
- 状态码：200 OK
- 响应包含 user 和 token 字段
- token 有效期 24 小时

### 课程模块测试用例

#### TC-COURSE-001: 获取课程列表

**测试步骤**:
```bash
curl -X GET http://localhost:8080/api/v1/courses
```

**预期结果**:
- 状态码：200 OK
- 返回课程列表和分页信息

#### TC-COURSE-002: 获取课程详情

**测试步骤**:
```bash
curl -X GET http://localhost:8080/api/v1/courses/{course_id}
```

**预期结果**:
- 状态码：200 OK
- 返回完整课程信息

### 作业提交测试用例

#### TC-SUBMIT-001: 提交选择题 - 正确答案

**测试步骤**:
```bash
curl -X POST http://localhost:8080/api/v1/exercises/{id}/submit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{
    "selected_options": ["Option A"]
  }'
```

**预期结果**:
- 状态码：201 Created
- is_correct: true
- score: 满分

---

## 🐛 已知问题

| 问题 ID | 描述 | 严重程度 | 状态 |
|---------|------|----------|------|
| BUG-001 | 8080 端口被 SearXNG 占用 | 高 | 待解决 |
| BUG-002 | Go 环境未安装 | 高 | 待解决 |
| BUG-003 | 前端服务未启动 | 中 | 待解决 |
| BUG-004 | 沙箱服务未部署 | 中 | 待解决 |

---

## 📝 测试建议

1. **环境准备**: 优先解决端口冲突问题，建议：
   - 停止 SearXNG 容器，或
   - 修改后端服务端口为 8081

2. **自动化测试**: 使用提供的 `scripts/test-api.sh` 脚本进行回归测试

3. **性能测试**: 服务稳定后，进行并发测试：
   ```bash
   wrk -t12 -c400 -d30s http://localhost:8080/api/v1/courses
   ```

4. **安全测试**: 建议进行：
   - SQL 注入测试
   - XSS 测试
   - CSRF 测试
   - 权限绕过测试

---

## 📈 测试进度

| 测试类型 | 总数 | 已通过 | 失败 | 待测试 | 通过率 |
|----------|------|--------|------|--------|--------|
| API 测试 | 25 | 0 | 0 | 25 | 0% |
| 前端测试 | 14 | 0 | 0 | 14 | 0% |
| 沙箱测试 | 6 | 0 | 0 | 6 | 0% |
| 集成测试 | 3 | 0 | 0 | 3 | 0% |
| **总计** | **48** | **0** | **0** | **48** | **0%** |

---

## 📅 下一步计划

1. [ ] 解决端口冲突，启动后端服务
2. [ ] 启动前端开发服务器
3. [ ] 部署沙箱服务
4. [ ] 执行 API 自动化测试
5. [ ] 执行前端功能测试
6. [ ] 执行集成测试
7. [ ] 输出完整测试报告

---

**报告生成时间**: 2026-03-10 13:35 GMT+8  
**下次更新**: 服务启动后执行实际测试
