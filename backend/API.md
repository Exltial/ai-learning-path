# AI 学习之路 - RESTful API 接口文档

**版本:** v1.0.0  
**基础路径:** `/api/v1`  
**认证方式:** JWT Bearer Token

---

## 📋 目录

1. [认证与授权](#认证与授权)
2. [用户管理](#用户管理)
3. [课程管理](#课程管理)
4. [章节管理](#章节管理)
5. [练习管理](#练习管理)
6. [提交与评分](#提交与评分)
7. [学习进度](#学习进度)
8. [测验管理](#测验管理)
9. [成就系统](#成就系统)
10. [讨论区](#讨论区)
11. [通知系统](#通知系统)
12. [课程评价](#课程评价)

---

## 🔐 认证与授权

### 1.1 用户注册
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "string (required, 3-50 chars)",
  "email": "string (required, valid email)",
  "password": "string (required, min 8 chars)",
  "avatar_url": "string (optional)"
}
```

**响应:** `201 Created`
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "username": "string",
      "email": "string",
      "role": "student",
      "created_at": "timestamp"
    },
    "token": "jwt_token_string"
  }
}
```

### 1.2 用户登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "string (required)",
  "password": "string (required)"
}
```

**响应:** `200 OK`
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "username": "string",
      "email": "string",
      "role": "student",
      "avatar_url": "string"
    },
    "token": "jwt_token_string",
    "expires_in": 86400
  }
}
```

### 1.3 刷新 Token
```http
POST /api/v1/auth/refresh
Content-Type: application/json
Authorization: Bearer {refresh_token}
```

**响应:** `200 OK`
```json
{
  "success": true,
  "data": {
    "token": "new_jwt_token",
    "expires_in": 86400
  }
}
```

### 1.4 登出
```http
POST /api/v1/auth/logout
Authorization: Bearer {token}
```

**响应:** `200 OK`
```json
{
  "success": true,
  "message": "Successfully logged out"
}
```

---

## 👤 用户管理

### 2.1 获取当前用户信息
```http
GET /api/v1/users/me
Authorization: Bearer {token}
```

**响应:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "username": "string",
    "email": "string",
    "avatar_url": "string",
    "role": "student",
    "created_at": "timestamp",
    "last_login_at": "timestamp"
  }
}
```

### 2.2 更新用户信息
```http
PUT /api/v1/users/me
Authorization: Bearer {token}
Content-Type: application/json

{
  "username": "string (optional)",
  "avatar_url": "string (optional)"
}
```

**响应:** `200 OK`

### 2.3 修改密码
```http
PUT /api/v1/users/me/password
Authorization: Bearer {token}
Content-Type: application/json

{
  "current_password": "string (required)",
  "new_password": "string (required, min 8 chars)"
}
```

**响应:** `200 OK`

### 2.4 获取用户统计
```http
GET /api/v1/users/me/stats
Authorization: Bearer {token}
```

**响应:** `200 OK`
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
    "achievements_count": 3
  }
}
```

---

## 📚 课程管理

### 3.1 获取课程列表
```http
GET /api/v1/courses
Authorization: Bearer {token}
```

**查询参数:**
- `category` - 分类筛选
- `difficulty` - 难度筛选 (beginner/intermediate/advanced)
- `page` - 页码 (default: 1)
- `limit` - 每页数量 (default: 20, max: 100)
- `sort` - 排序字段 (created_at/rating/enrollment_count)
- `order` - 排序方向 (asc/desc)

**响应:** `200 OK`
```json
{
  "success": true,
  "data": {
    "courses": [
      {
        "id": "uuid",
        "title": "string",
        "description": "string",
        "thumbnail_url": "string",
        "instructor": {
          "id": "uuid",
          "username": "string",
          "avatar_url": "string"
        },
        "category": "string",
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

### 3.2 获取课程详情
```http
GET /api/v1/courses/{course_id}
Authorization: Bearer {token}
```

**响应:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "title": "string",
    "description": "string",
    "thumbnail_url": "string",
    "instructor": {
      "id": "uuid",
      "username": "string",
      "avatar_url": "string"
    },
    "category": "string",
    "difficulty_level": "beginner",
    "estimated_hours": 20,
    "price": 99.00,
    "rating": 4.5,
    "enrollment_count": 1200,
    "is_published": true,
    "lessons_count": 15,
    "exercises_count": 45,
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

### 3.3 创建课程 (Instructor/Admin)
```http
POST /api/v1/courses
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "string (required)",
  "description": "string",
  "category": "string",
  "difficulty_level": "beginner",
  "estimated_hours": 20,
  "price": 99.00
}
```

**响应:** `201 Created`

### 3.4 更新课程 (Instructor/Admin)
```http
PUT /api/v1/courses/{course_id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "string",
  "description": "string",
  "category": "string",
  "difficulty_level": "beginner",
  "estimated_hours": 20,
  "price": 99.00,
  "is_published": true
}
```

**响应:** `200 OK`

### 3.5 删除课程 (Instructor/Admin)
```http
DELETE /api/v1/courses/{course_id}
Authorization: Bearer {token}
```

**响应:** `204 No Content`

### 3.6 注册课程
```http
POST /api/v1/courses/{course_id}/enroll
Authorization: Bearer {token}
```

**响应:** `201 Created`
```json
{
  "success": true,
  "data": {
    "enrollment_id": "uuid",
    "course_id": "uuid",
    "enrolled_at": "timestamp"
  }
}
```

### 3.7 获取课程章节列表
```http
GET /api/v1/courses/{course_id}/lessons
Authorization: Bearer {token}
```

**响应:** `200 OK`
```json
{
  "success": true,
  "data": {
    "lessons": [
      {
        "id": "uuid",
        "title": "string",
        "description": "string",
        "order_index": 1,
        "is_free_preview": true,
        "is_completed": false,
        "video_duration": 600
      }
    ]
  }
}
```

---

## 📖 章节管理

### 4.1 获取章节详情
```http
GET /api/v1/lessons/{lesson_id}
Authorization: Bearer {token}
```

**响应:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "course_id": "uuid",
    "title": "string",
    "description": "string",
    "content": "string",
    "video_url": "string",
    "video_duration": 600,
    "order_index": 1,
    "is_free_preview": true,
    "exercises": [
      {
        "id": "uuid",
        "title": "string",
        "exercise_type": "multiple_choice",
        "difficulty": "easy",
        "points": 10
      }
    ]
  }
}
```

### 4.2 创建章节 (Instructor/Admin)
```http
POST /api/v1/courses/{course_id}/lessons
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "string (required)",
  "description": "string",
  "content": "string",
  "video_url": "string",
  "video_duration": 600,
  "is_free_preview": false
}
```

**响应:** `201 Created`

### 4.3 更新章节 (Instructor/Admin)
```http
PUT /api/v1/lessons/{lesson_id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "string",
  "description": "string",
  "content": "string",
  "video_url": "string",
  "video_duration": 600,
  "is_free_preview": false
}
```

**响应:** `200 OK`

### 4.4 删除章节 (Instructor/Admin)
```http
DELETE /api/v1/lessons/{lesson_id}
Authorization: Bearer {token}
```

**响应:** `204 No Content`

---

## ✍️ 练习管理

### 5.1 获取练习详情
```http
GET /api/v1/exercises/{exercise_id}
Authorization: Bearer {token}
```

**响应:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "lesson_id": "uuid",
    "title": "string",
    "description": "string",
    "exercise_type": "multiple_choice",
    "difficulty": "easy",
    "points": 10,
    "max_attempts": 3,
    "time_limit": 300,
    "starter_code": "string",
    "options": [
      {"text": "选项 A", "is_correct": false},
      {"text": "选项 B", "is_correct": true}
    ]
  }
}
```

### 5.2 创建练习 (Instructor/Admin)
```http
POST /api/v1/lessons/{lesson_id}/exercises
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "string (required)",
  "description": "string",
  "exercise_type": "multiple_choice (required)",
  "difficulty": "easy",
  "points": 10,
  "max_attempts": 3,
  "time_limit": 300,
  "starter_code": "string",
  "test_cases": {},
  "expected_answer": {},
  "options": []
}
```

**响应:** `201 Created`

### 5.3 更新练习 (Instructor/Admin)
```http
PUT /api/v1/exercises/{exercise_id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "string",
  "description": "string",
  "exercise_type": "multiple_choice",
  "difficulty": "easy",
  "points": 10,
  "max_attempts": 3,
  "time_limit": 300
}
```

**响应:** `200 OK`

### 5.4 删除练习 (Instructor/Admin)
```http
DELETE /api/v1/exercises/{exercise_id}
Authorization: Bearer {token}
```

**响应:** `204 No Content`

---

## 📤 提交与评分

### 6.1 提交练习答案
```http
POST /api/v1/exercises/{exercise_id}/submit
Authorization: Bearer {token}
Content-Type: application/json

{
  "answer": "string (for text answers)",
  "code": "string (for coding exercises)",
  "selected_options": ["option_id_1", "option_id_2"] (for multiple choice)
}
```

**响应:** `201 Created`
```json
{
  "success": true,
  "data": {
    "submission_id": "uuid",
    "is_correct": true,
    "score": 10.0,
    "feedback": "Great job!",
    "attempt_number": 1,
    "remaining_attempts": 2
  }
}
```

### 6.2 获取提交历史
```http
GET /api/v1/exercises/{exercise_id}/submissions
Authorization: Bearer {token}
```

**查询参数:**
- `page` - 页码
- `limit` - 每页数量

**响应:** `200 OK`
```json
{
  "success": true,
  "data": {
    "submissions": [
      {
        "id": "uuid",
        "answer": "string",
        "code": "string",
        "is_correct": true,
        "score": 10.0,
        "feedback": "string",
        "attempt_number": 1,
        "submitted_at": "timestamp"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 3
    }
  }
}
```

### 6.3 获取提交详情
```http
GET /api/v1/submissions/{submission_id}
Authorization: Bearer {token}
```

**响应:** `200 OK`

### 6.4 手动评分 (Instructor/Admin, for essay type)
```http
POST /api/v1/submissions/{submission_id}/grade
Authorization: Bearer {token}
Content-Type: application/json

{
  "score": 8.5,
  "feedback": "string",
  "is_correct": true
}
```

**响应:** `200 OK`

---

## 📊 学习进度

### 7.1 获取学习进度
```http
GET /api/v1/courses/{course_id}/progress
Authorization: Bearer {token}
```

**响应:** `200 OK`
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
    "lessons": [
      {
        "lesson_id": "uuid",
        "title": "string",
        "is_completed": true,
        "completed_at": "timestamp"
      }
    ]
  }
}
```

### 7.2 更新章节进度
```http
PUT /api/v1/lessons/{lesson_id}/progress
Authorization: Bearer {token}
Content-Type: application/json

{
  "is_completed": true,
  "video_position": 120
}
```

**响应:** `200 OK`

### 7.3 获取用户所有进度
```http
GET /api/v1/users/me/progress
Authorization: Bearer {token}
```

**响应:** `200 OK`
```json
{
  "success": true,
  "data": {
    "active_courses": [
      {
        "course_id": "uuid",
        "title": "string",
        "progress_percentage": 45.5,
        "last_accessed_at": "timestamp"
      }
    ],
    "completed_courses": [
      {
        "course_id": "uuid",
        "title": "string",
        "completed_at": "timestamp"
      }
    ]
  }
}
```

---

## 📝 测验管理

### 8.1 开始测验
```http
POST /api/v1/courses/{course_id}/quiz/start
Authorization: Bearer {token}
```

**响应:** `201 Created`
```json
{
  "success": true,
  "data": {
    "quiz_id": "uuid",
    "total_questions": 20,
    "time_limit": 1800,
    "started_at": "timestamp"
  }
}
```

### 8.2 提交测验答案
```http
POST /api/v1/quiz/{quiz_id}/submit
Authorization: Bearer {token}
Content-Type: application/json

{
  "answers": [
    {"question_id": "uuid", "answer": "string"}
  ]
}
```

**响应:** `200 OK`
```json
{
  "success": true,
  "data": {
    "total_questions": 20,
    "correct_answers": 16,
    "score_percentage": 80.0,
    "time_taken": 1200,
    "passed": true
  }
}
```

### 8.3 获取测验历史
```http
GET /api/v1/users/me/quiz-results
Authorization: Bearer {token}
```

**响应:** `200 OK`

---

## 🏆 成就系统

### 9.1 获取用户成就
```http
GET /api/v1/users/me/achievements
Authorization: Bearer {token}
```

**响应:** `200 OK`
```json
{
  "success": true,
  "data": {
    "achievements": [
      {
        "id": "uuid",
        "name": "First Steps",
        "description": "完成第一个练习",
        "icon_url": "string",
        "points": 10,
        "earned_at": "timestamp"
      }
    ],
    "total_points": 450
  }
}
```

### 9.2 获取所有成就
```http
GET /api/v1/achievements
Authorization: Bearer {token}
```

**响应:** `200 OK`

---

## 💬 讨论区

### 10.1 获取课程讨论列表
```http
GET /api/v1/courses/{course_id}/discussions
Authorization: Bearer {token}
```

**查询参数:**
- `page` - 页码
- `limit` - 每页数量
- `sort` - 排序 (created_at/upvotes)

**响应:** `200 OK`
```json
{
  "success": true,
  "data": {
    "discussions": [
      {
        "id": "uuid",
        "title": "string",
        "content": "string",
        "author": {
          "id": "uuid",
          "username": "string"
        },
        "upvotes": 15,
        "replies_count": 5,
        "is_resolved": false,
        "created_at": "timestamp"
      }
    ],
    "pagination": {}
  }
}
```

### 10.2 创建讨论
```http
POST /api/v1/courses/{course_id}/discussions
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "string",
  "content": "string (required)",
  "lesson_id": "uuid (optional)"
}
```

**响应:** `201 Created`

### 10.3 回复讨论
```http
POST /api/v1/discussions/{discussion_id}/replies
Authorization: Bearer {token}
Content-Type: application/json

{
  "content": "string (required)"
}
```

**响应:** `201 Created`

### 10.4 点赞讨论
```http
POST /api/v1/discussions/{discussion_id}/upvote
Authorization: Bearer {token}
```

**响应:** `200 OK`

### 10.5 标记为已解决 (Instructor/Admin 或作者)
```http
PATCH /api/v1/discussions/{discussion_id}/resolve
Authorization: Bearer {token}
```

**响应:** `200 OK`

---

## 🔔 通知系统

### 11.1 获取通知列表
```http
GET /api/v1/users/me/notifications
Authorization: Bearer {token}
```

**查询参数:**
- `is_read` - 筛选已读/未读
- `page` - 页码
- `limit` - 每页数量

**响应:** `200 OK`
```json
{
  "success": true,
  "data": {
    "notifications": [
      {
        "id": "uuid",
        "title": "string",
        "message": "string",
        "notification_type": "info",
        "is_read": false,
        "action_url": "string",
        "created_at": "timestamp"
      }
    ],
    "unread_count": 5
  }
}
```

### 11.2 标记通知为已读
```http
PATCH /api/v1/notifications/{notification_id}/read
Authorization: Bearer {token}
```

**响应:** `200 OK`

### 11.3 标记所有通知为已读
```http
PATCH /api/v1/users/me/notifications/read-all
Authorization: Bearer {token}
```

**响应:** `200 OK`

---

## ⭐ 课程评价

### 12.1 创建课程评价
```http
POST /api/v1/courses/{course_id}/reviews
Authorization: Bearer {token}
Content-Type: application/json

{
  "rating": 5 (required, 1-5),
  "comment": "string"
}
```

**响应:** `201 Created`

### 12.2 更新课程评价
```http
PUT /api/v1/courses/{course_id}/reviews
Authorization: Bearer {token}
Content-Type: application/json

{
  "rating": 4,
  "comment": "string"
}
```

**响应:** `200 OK`

### 12.3 获取课程评价列表
```http
GET /api/v1/courses/{course_id}/reviews
Authorization: Bearer {token}
```

**响应:** `200 OK`
```json
{
  "success": true,
  "data": {
    "reviews": [
      {
        "id": "uuid",
        "user": {
          "id": "uuid",
          "username": "string",
          "avatar_url": "string"
        },
        "rating": 5,
        "comment": "string",
        "is_verified": true,
        "created_at": "timestamp"
      }
    ],
    "average_rating": 4.5,
    "total_reviews": 120
  }
}
```

---

## ❌ 错误响应格式

所有错误响应遵循统一格式:

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

| 错误码 | HTTP 状态码 | 说明 |
|--------|------------|------|
| `UNAUTHORIZED` | 401 | 未授权或 Token 过期 |
| `FORBIDDEN` | 403 | 权限不足 |
| `NOT_FOUND` | 404 | 资源不存在 |
| `BAD_REQUEST` | 400 | 请求参数错误 |
| `CONFLICT` | 409 | 资源冲突 (如重复注册) |
| `RATE_LIMITED` | 429 | 请求频率超限 |
| `INTERNAL_ERROR` | 500 | 服务器内部错误 |

---

## 🔒 权限说明

| 角色 | 说明 |
|------|------|
| `student` | 普通学员，可学习课程、提交练习、参与讨论 |
| `instructor` | 讲师，可创建/管理课程、评分作业 |
| `admin` | 管理员，拥有全部权限 |

---

## 📌 注意事项

1. **所有需要认证的接口**必须在 Header 中包含 `Authorization: Bearer {token}`
2. **时间格式**统一使用 ISO 8601 格式 (e.g., `2026-03-10T12:00:00Z`)
3. **分页参数**默认 page=1, limit=20
4. **文件上传**使用 multipart/form-data
5. **Rate Limiting**: 每个用户每分钟最多 100 次请求

---

*文档版本: v1.0.0 | 最后更新：2026-03-10*
