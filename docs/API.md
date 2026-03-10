# API 文档 / API Reference

完整的 RESTful API 接口文档，包含 Swagger/OpenAPI 规范。

Complete RESTful API documentation with Swagger/OpenAPI specification.

---

## 📋 目录 / Table of Contents

1. [API 概览 / API Overview](#api-概览--api-overview)
2. [认证与授权 / Authentication](#认证与授权--authentication)
3. [接口详情 / API Endpoints](#接口详情--api-endpoints)
4. [Swagger/OpenAPI 规范 / Swagger Specification](#swaggeropenapi-规范--swagger-specification)
5. [错误码 / Error Codes](#错误码--error-codes)
6. [速率限制 / Rate Limiting](#速率限制--rate-limiting)
7. [版本控制 / Versioning](#版本控制--versioning)
8. [最佳实践 / Best Practices](#最佳实践--best-practices)

---

## API 概览 / API Overview

### 基础信息 / Base Information

| 项目 / Item | 值 / Value |
|------------|-----------|
| 基础 URL / Base URL | `https://api.yourdomain.com/api/v1` |
| 认证方式 / Auth Method | JWT Bearer Token |
| 数据格式 / Format | JSON |
| 字符编码 / Encoding | UTF-8 |
| 时间格式 / Time Format | ISO 8601 (UTC) |

### Swagger UI

访问 Swagger UI 查看交互式 API 文档：
```
http://localhost:8080/swagger
https://api.yourdomain.com/swagger
```

### OpenAPI 规范文件

下载 OpenAPI 3.0 规范文件：
```
http://localhost:8080/swagger/doc.json
https://api.yourdomain.com/swagger/doc.json
```

---

## 认证与授权 / Authentication

### JWT Token 认证

所有需要认证的接口必须在请求头中包含：

```http
Authorization: Bearer <your_jwt_token>
```

### 获取 Token

#### 用户注册 / Register

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "string (3-50 chars, required)",
  "email": "string (valid email, required)",
  "password": "string (min 8 chars, required)",
  "avatar_url": "string (optional)"
}
```

**响应 / Response:** `201 Created`

```json
{
  "success": true,
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "email": "john@example.com",
      "role": "student",
      "created_at": "2026-03-10T12:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4..."
  }
}
```

#### 用户登录 / Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "string (required)",
  "password": "string (required)"
}
```

**响应 / Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "email": "john@example.com",
      "avatar_url": "https://example.com/avatar.jpg",
      "role": "student"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4...",
    "expires_in": 86400
  }
}
```

### Token 刷新 / Refresh Token

```http
POST /api/v1/auth/refresh
Content-Type: application/json
Authorization: Bearer {refresh_token}
```

**响应 / Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "token": "new_jwt_token...",
    "refresh_token": "new_refresh_token...",
    "expires_in": 86400
  }
}
```

### 登出 / Logout

```http
POST /api/v1/auth/logout
Authorization: Bearer {token}
```

**响应 / Response:** `200 OK`

```json
{
  "success": true,
  "message": "Successfully logged out"
}
```

---

## 接口详情 / API Endpoints

### 用户管理 / User Management

#### 获取当前用户信息 / Get Current User

```http
GET /api/v1/users/me
Authorization: Bearer {token}
```

**响应 / Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "username": "john_doe",
    "email": "john@example.com",
    "avatar_url": "https://example.com/avatar.jpg",
    "role": "student",
    "created_at": "2026-03-10T12:00:00Z",
    "last_login_at": "2026-03-10T15:30:00Z"
  }
}
```

#### 更新用户信息 / Update User

```http
PUT /api/v1/users/me
Authorization: Bearer {token}
Content-Type: application/json

{
  "username": "new_username",
  "avatar_url": "https://example.com/new-avatar.jpg"
}
```

#### 修改密码 / Change Password

```http
PUT /api/v1/users/me/password
Authorization: Bearer {token}
Content-Type: application/json

{
  "current_password": "old_password",
  "new_password": "new_secure_password"
}
```

#### 获取用户统计 / Get User Stats

```http
GET /api/v1/users/me/stats
Authorization: Bearer {token}
```

**响应 / Response:**

```json
{
  "success": true,
  "data": {
    "total_courses": 5,
    "completed_courses": 2,
    "total_exercises": 50,
    "completed_exercises": 35,
    "total_points": 450,
    "learning_streak": 7,
    "achievements_count": 3,
    "total_learning_time_minutes": 1200
  }
}
```

### 课程管理 / Course Management

#### 获取课程列表 / List Courses

```http
GET /api/v1/courses
Authorization: Bearer {token}
```

**查询参数 / Query Parameters:**

| 参数 / Param | 类型 / Type | 默认 / Default | 说明 / Description |
|-------------|------------|---------------|-------------------|
| category | string | - | 分类筛选 / Category filter |
| difficulty | string | - | 难度 / beginner, intermediate, advanced |
| page | integer | 1 | 页码 / Page number |
| limit | integer | 20 | 每页数量 / Items per page (max: 100) |
| sort | string | created_at | 排序字段 / Sort field |
| order | string | desc | 排序方向 / asc, desc |

**响应 / Response:**

```json
{
  "success": true,
  "data": {
    "courses": [
      {
        "id": "uuid",
        "title": "Python 编程基础",
        "title_en": "Python Programming Basics",
        "description": "从零开始学习 Python...",
        "thumbnail_url": "https://example.com/thumb.jpg",
        "instructor": {
          "id": "uuid",
          "username": "instructor_name",
          "avatar_url": "https://example.com/avatar.jpg"
        },
        "category": "Programming",
        "difficulty_level": "beginner",
        "estimated_hours": 20,
        "price": 99.00,
        "rating": 4.5,
        "enrollment_count": 1200,
        "is_enrolled": true,
        "progress_percentage": 45.5
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "total_pages": 8
    }
  }
}
```

#### 获取课程详情 / Get Course Detail

```http
GET /api/v1/courses/{course_id}
Authorization: Bearer {token}
```

#### 注册课程 / Enroll Course

```http
POST /api/v1/courses/{course_id}/enroll
Authorization: Bearer {token}
```

#### 获取课程章节 / Get Course Lessons

```http
GET /api/v1/courses/{course_id}/lessons
Authorization: Bearer {token}
```

### 章节管理 / Lesson Management

#### 获取章节详情 / Get Lesson Detail

```http
GET /api/v1/lessons/{lesson_id}
Authorization: Bearer {token}
```

**响应 / Response:**

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "course_id": "uuid",
    "title": "变量与数据类型",
    "title_en": "Variables and Data Types",
    "description": "学习 Python 中的变量和基础数据类型",
    "content": "# 变量与数据类型\n\n在 Python 中...",
    "video_url": "https://example.com/video.mp4",
    "video_duration": 600,
    "order_index": 1,
    "is_free_preview": true,
    "exercises": [
      {
        "id": "uuid",
        "title": "变量定义练习",
        "exercise_type": "coding",
        "difficulty": "easy",
        "points": 10
      }
    ]
  }
}
```

### 练习管理 / Exercise Management

#### 获取练习详情 / Get Exercise Detail

```http
GET /api/v1/exercises/{exercise_id}
Authorization: Bearer {token}
```

**响应 / Response:**

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "lesson_id": "uuid",
    "title": "变量定义练习",
    "title_en": "Variable Definition Exercise",
    "description": "练习定义不同类型的变量",
    "exercise_type": "coding",
    "difficulty": "easy",
    "points": 10,
    "max_attempts": 3,
    "time_limit": 300,
    "starter_code": "# 请定义一个变量 name，值为你的名字\n",
    "test_cases": [
      {
        "input": "",
        "expected_output": "name = 'John'"
      }
    ],
    "options": null
  }
}
```

### 提交与评分 / Submission & Grading

#### 提交练习答案 / Submit Exercise

```http
POST /api/v1/exercises/{exercise_id}/submit
Authorization: Bearer {token}
Content-Type: application/json

{
  "answer": "string (for text answers)",
  "code": "print('Hello, World!') (for coding exercises)",
  "selected_options": ["option_id_1", "option_id_2"] (for multiple choice)
}
```

**响应 / Response:** `201 Created`

```json
{
  "success": true,
  "data": {
    "submission_id": "uuid",
    "is_correct": true,
    "score": 10.0,
    "max_score": 10.0,
    "feedback": "Great job! Your code is correct.",
    "feedback_zh": "做得好！你的代码是正确的。",
    "attempt_number": 1,
    "remaining_attempts": 2,
    "execution_result": {
      "stdout": "Hello, World!\n",
      "stderr": "",
      "exit_code": 0,
      "execution_time_ms": 45
    }
  }
}
```

#### 获取提交历史 / Get Submission History

```http
GET /api/v1/exercises/{exercise_id}/submissions
Authorization: Bearer {token}
```

**查询参数 / Query Parameters:**

| 参数 / Param | 类型 / Type | 默认 / Default |
|-------------|------------|---------------|
| page | integer | 1 |
| limit | integer | 10 |

#### 获取提交详情 / Get Submission Detail

```http
GET /api/v1/submissions/{submission_id}
Authorization: Bearer {token}
```

### 学习进度 / Learning Progress

#### 获取课程进度 / Get Course Progress

```http
GET /api/v1/courses/{course_id}/progress
Authorization: Bearer {token}
```

**响应 / Response:**

```json
{
  "success": true,
  "data": {
    "course_id": "uuid",
    "enrollment_id": "uuid",
    "progress_percentage": 45.5,
    "completed_lessons": 7,
    "total_lessons": 15,
    "completed_exercises": 20,
    "total_exercises": 45,
    "last_accessed_at": "2026-03-10T15:30:00Z",
    "lessons": [
      {
        "lesson_id": "uuid",
        "title": "变量与数据类型",
        "is_completed": true,
        "completed_at": "2026-03-09T10:00:00Z",
        "exercise_completion": {
          "total": 3,
          "completed": 3
        }
      }
    ]
  }
}
```

#### 更新章节进度 / Update Lesson Progress

```http
PUT /api/v1/lessons/{lesson_id}/progress
Authorization: Bearer {token}
Content-Type: application/json

{
  "is_completed": true,
  "video_position": 120
}
```

#### 获取用户所有进度 / Get All User Progress

```http
GET /api/v1/users/me/progress
Authorization: Bearer {token}
```

### 成就系统 / Achievement System

#### 获取用户成就 / Get User Achievements

```http
GET /api/v1/users/me/achievements
Authorization: Bearer {token}
```

**响应 / Response:**

```json
{
  "success": true,
  "data": {
    "achievements": [
      {
        "id": "uuid",
        "name": "第一步 / First Steps",
        "description": "完成第一个练习 / Complete your first exercise",
        "icon_url": "https://example.com/icons/first-steps.png",
        "category": "learning",
        "points": 10,
        "rarity": "common",
        "earned_at": "2026-03-10T12:30:00Z"
      }
    ],
    "total_points": 450,
    "next_available": [
      {
        "id": "uuid",
        "name": "连续学习 7 天",
        "progress": 5,
        "required": 7
      }
    ]
  }
}
```

### 讨论区 / Discussion Forum

#### 获取讨论列表 / List Discussions

```http
GET /api/v1/courses/{course_id}/discussions
Authorization: Bearer {token}
```

**查询参数 / Query Parameters:**

| 参数 / Param | 类型 / Type | 默认 / Default |
|-------------|------------|---------------|
| page | integer | 1 |
| limit | integer | 20 |
| sort | string | created_at |
| lesson_id | uuid | - |

#### 创建讨论 / Create Discussion

```http
POST /api/v1/courses/{course_id}/discussions
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "string (required)",
  "content": "string (required)",
  "lesson_id": "uuid (optional)",
  "is_anonymous": false
}
```

#### 回复讨论 / Reply to Discussion

```http
POST /api/v1/discussions/{discussion_id}/replies
Authorization: Bearer {token}
Content-Type: application/json

{
  "content": "string (required)",
  "parent_reply_id": "uuid (optional, for nested replies)"
}
```

#### 点赞讨论 / Upvote Discussion

```http
POST /api/v1/discussions/{discussion_id}/upvote
Authorization: Bearer {token}
```

### 通知系统 / Notification System

#### 获取通知列表 / List Notifications

```http
GET /api/v1/users/me/notifications
Authorization: Bearer {token}
```

**查询参数 / Query Parameters:**

| 参数 / Param | 类型 / Type | 默认 / Default |
|-------------|------------|---------------|
| is_read | boolean | - |
| page | integer | 1 |
| limit | integer | 20 |

**响应 / Response:**

```json
{
  "success": true,
  "data": {
    "notifications": [
      {
        "id": "uuid",
        "title": "作业已批改",
        "message": "你的 Python 练习已获得 10 分",
        "notification_type": "grading",
        "is_read": false,
        "action_url": "/exercises/uuid/submissions",
        "created_at": "2026-03-10T16:00:00Z"
      }
    ],
    "unread_count": 5
  }
}
```

#### 标记通知为已读 / Mark Notification as Read

```http
PATCH /api/v1/notifications/{notification_id}/read
Authorization: Bearer {token}
```

#### 标记所有通知为已读 / Mark All as Read

```http
PATCH /api/v1/users/me/notifications/read-all
Authorization: Bearer {token}
```

---

## Swagger/OpenAPI 规范 / Swagger Specification

### OpenAPI 3.0 定义文件

```yaml
openapi: 3.0.3
info:
  title: AI Learning Platform API
  description: |
    RESTful API for AI Interactive Learning Platform
    
    ## 认证
    所有需要认证的接口需要在 Header 中包含:
    ```
    Authorization: Bearer <token>
    ```
  version: 1.0.0
  contact:
    name: API Support
    email: exltial@163.com
    url: https://github.com/exltial/ai-learning-path

servers:
  - url: https://api.yourdomain.com/api/v1
    description: Production server
  - url: http://localhost:8080/api/v1
    description: Development server

tags:
  - name: Authentication
    description: 认证与授权接口
  - name: Users
    description: 用户管理接口
  - name: Courses
    description: 课程管理接口
  - name: Lessons
    description: 章节管理接口
  - name: Exercises
    description: 练习管理接口
  - name: Submissions
    description: 提交与评分接口
  - name: Progress
    description: 学习进度接口
  - name: Achievements
    description: 成就系统接口
  - name: Discussions
    description: 讨论区接口
  - name: Notifications
    description: 通知系统接口

paths:
  /auth/register:
    post:
      tags:
        - Authentication
      summary: 用户注册
      operationId: registerUser
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/RegisterRequest'
      responses:
        '201':
          description: 注册成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AuthResponse'
        '400':
          description: 请求参数错误
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /auth/login:
    post:
      tags:
        - Authentication
      summary: 用户登录
      operationId: loginUser
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/LoginRequest'
      responses:
        '200':
          description: 登录成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AuthResponse'

components:
  schemas:
    RegisterRequest:
      type: object
      required:
        - username
        - email
        - password
      properties:
        username:
          type: string
          minLength: 3
          maxLength: 50
          example: john_doe
        email:
          type: string
          format: email
          example: john@example.com
        password:
          type: string
          minLength: 8
          example: SecurePass123
        avatar_url:
          type: string
          format: uri
          example: https://example.com/avatar.jpg

    LoginRequest:
      type: object
      required:
        - email
        - password
      properties:
        email:
          type: string
          format: email
        password:
          type: string

    AuthResponse:
      type: object
      properties:
        success:
          type: boolean
        data:
          type: object
          properties:
            user:
              $ref: '#/components/schemas/User'
            token:
              type: string
            refresh_token:
              type: string
            expires_in:
              type: integer

    User:
      type: object
      properties:
        id:
          type: string
          format: uuid
        username:
          type: string
        email:
          type: string
          format: email
        role:
          type: string
          enum: [student, instructor, admin]
        avatar_url:
          type: string
        created_at:
          type: string
          format: date-time

    Error:
      type: object
      properties:
        success:
          type: boolean
          example: false
        error:
          type: object
          properties:
            code:
              type: string
            message:
              type: string
            details:
              type: object

  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

security:
  - bearerAuth: []
```

### 在 Go 中集成 Swagger

```go
// cmd/main.go
import (
    "github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
)

func main() {
    r := gin.Default()
    
    // Swagger UI
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    // ...
}

//go:generate swag init -g cmd/main.go -o ./docs/swagger
```

---

## 错误码 / Error Codes

### 统一错误响应格式

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": {}
  }
}
```

### 常见错误码

| 错误码 / Code | HTTP 状态码 | 说明 / Description |
|--------------|------------|-------------------|
| `UNAUTHORIZED` | 401 | 未授权或 Token 过期 / Unauthorized or token expired |
| `FORBIDDEN` | 403 | 权限不足 / Insufficient permissions |
| `NOT_FOUND` | 404 | 资源不存在 / Resource not found |
| `BAD_REQUEST` | 400 | 请求参数错误 / Invalid request parameters |
| `VALIDATION_ERROR` | 400 | 数据验证失败 / Validation failed |
| `CONFLICT` | 409 | 资源冲突 / Resource conflict |
| `RATE_LIMITED` | 429 | 请求频率超限 / Rate limit exceeded |
| `INTERNAL_ERROR` | 500 | 服务器内部错误 / Internal server error |
| `SERVICE_UNAVAILABLE` | 503 | 服务暂时不可用 / Service temporarily unavailable |
| `TIMEOUT` | 408 | 请求超时 / Request timeout |

### 验证错误示例

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "fields": {
        "email": ["Invalid email format"],
        "password": ["Password must be at least 8 characters"]
      }
    }
  }
}
```

---

## 速率限制 / Rate Limiting

### 限制策略

| 接口类型 / Endpoint Type | 限制 / Limit | 时间窗口 / Window |
|-------------------------|-------------|------------------|
| 认证接口 / Auth | 10 次 / requests | 1 分钟 / minute |
| 普通 API / General API | 100 次 / requests | 1 分钟 / minute |
| 文件上传 / File Upload | 10 次 / requests | 1 分钟 / minute |
| 代码执行 / Code Execution | 5 次 / requests | 1 分钟 / minute |

### 速率限制响应

```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1678453200

{
  "success": false,
  "error": {
    "code": "RATE_LIMITED",
    "message": "Too many requests. Please try again in 60 seconds.",
    "details": {
      "retry_after": 60
    }
  }
}
```

---

## 版本控制 / Versioning

### URL 版本控制

当前 API 版本：**v1**

```
/api/v1/users
/api/v1/courses
```

### 版本弃用策略

1. 新版本发布后，旧版本继续支持 6 个月
2. 提前 3 个月发送弃用通知
3. 弃用后返回警告头：
   ```http
   Deprecation: true
   Sunset: Sat, 10 Sep 2026 00:00:00 GMT
   ```

---

## 最佳实践 / Best Practices

### 1. 请求优化 / Request Optimization

```javascript
// ✅ 好的做法 - 批量请求
GET /api/v1/courses?ids=1,2,3,4,5

// ❌ 不好的做法 - 多次单独请求
GET /api/v1/courses/1
GET /api/v1/courses/2
GET /api/v1/courses/3
```

### 2. 分页处理 / Pagination

```javascript
// 使用分页参数
GET /api/v1/courses?page=1&limit=20

// 处理分页响应
const loadCourses = async (page) => {
  const response = await fetch(`/api/v1/courses?page=${page}&limit=20`);
  const data = await response.json();
  return {
    courses: data.data.courses,
    hasMore: page < data.data.pagination.total_pages
  };
};
```

### 3. 错误处理 / Error Handling

```javascript
// 统一的错误处理
const handleApiError = (error) => {
  if (error.response) {
    switch (error.response.status) {
      case 401:
        // Token 过期，刷新或重新登录
        refreshToken();
        break;
      case 429:
        // 速率限制，等待后重试
        const retryAfter = error.response.headers['x-ratelimit-reset'];
        setTimeout(() => retry(), retryAfter * 1000);
        break;
      default:
        console.error('API Error:', error.response.data.error);
    }
  }
};
```

### 4. 缓存策略 / Caching Strategy

```javascript
// 使用 ETag 进行条件请求
const fetchCourse = async (id) => {
  const cachedEtag = localStorage.getItem(`course_${id}_etag`);
  const headers = cachedEtag ? { 'If-None-Match': cachedEtag } : {};
  
  const response = await fetch(`/api/v1/courses/${id}`, { headers });
  
  if (response.status === 304) {
    // 使用缓存数据
    return JSON.parse(localStorage.getItem(`course_${id}_data`));
  }
  
  const data = await response.json();
  localStorage.setItem(`course_${id}_data`, JSON.stringify(data));
  localStorage.setItem(`course_${id}_etag`, response.headers.get('etag'));
  return data;
};
```

---

## 代码示例 / Code Examples

### JavaScript/TypeScript

```typescript
// api/client.ts
class ApiClient {
  private baseURL: string;
  private token: string | null = null;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  setToken(token: string) {
    this.token = token;
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...(this.token && { Authorization: `Bearer ${this.token}` }),
      ...options.headers,
    };

    const response = await fetch(url, { ...options, headers });

    if (!response.ok) {
      const error = await response.json();
      throw new ApiError(error);
    }

    return response.json();
  }

  async login(email: string, password: string) {
    const data = await this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    this.setToken(data.data.token);
    return data;
  }

  async getCourses(params?: CourseParams) {
    const query = new URLSearchParams(params as any).toString();
    return this.request(`/courses${query ? '?' + query : ''}`);
  }
}
```

### Go

```go
// pkg/api/client.go
package api

type Client struct {
    baseURL string
    token   string
    client  *http.Client
}

func NewClient(baseURL string) *Client {
    return &Client{
        baseURL: baseURL,
        client:  &http.Client{},
    }
}

func (c *Client) SetToken(token string) {
    c.token = token
}

func (c *Client) request(method, endpoint string, body interface{}) (*Response, error) {
    var reqBody io.Reader
    if body != nil {
        jsonData, _ := json.Marshal(body)
        reqBody = bytes.NewBuffer(jsonData)
    }

    req, _ := http.NewRequest(method, c.baseURL+endpoint, reqBody)
    req.Header.Set("Content-Type", "application/json")
    if c.token != "" {
        req.Header.Set("Authorization", "Bearer "+c.token)
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Handle response...
}
```

### Python

```python
# api/client.py
import requests
from typing import Optional, Dict, Any

class APIClient:
    def __init__(self, base_url: str):
        self.base_url = base_url
        self.session = requests.Session()
        self.token: Optional[str] = None
    
    def set_token(self, token: str):
        self.token = token
        self.session.headers.update({'Authorization': f'Bearer {token}'})
    
    def request(self, method: str, endpoint: str, **kwargs) -> Dict[str, Any]:
        url = f'{self.base_url}{endpoint}'
        if 'headers' not in kwargs:
            kwargs['headers'] = {}
        kwargs['headers']['Content-Type'] = 'application/json'
        
        response = self.session.request(method, url, **kwargs)
        response.raise_for_status()
        return response.json()
    
    def login(self, email: str, password: str) -> Dict[str, Any]:
        data = self.request('POST', '/auth/login', 
                          json={'email': email, 'password': password})
        self.set_token(data['data']['token'])
        return data
    
    def get_courses(self, **params) -> Dict[str, Any]:
        return self.request('GET', '/courses', params=params)
```

---

*最后更新 / Last Updated: 2026-03-10*
