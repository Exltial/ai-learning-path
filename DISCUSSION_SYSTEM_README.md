# 社区讨论系统实现文档

## 概述

已成功实现 AI 学习平台的社区讨论系统，支持课程讨论区、帖子 CRUD、嵌套回复、点赞/收藏、热门排序等功能。

## 实现文件清单

### 后端文件 (Go)

1. **backend/migrations/007_create_discussions.up.sql** (6.3KB)
   - 创建 discussions 主表（支持嵌套回复）
   - 创建 discussion_likes 表（点赞/点踩）
   - 创建 discussion_favorites 表（收藏）
   - 创建 discussion_mentions 表（@提及用户）
   - 创建 discussion_tags 和 discussion_tag_mapping 表（标签系统）
   - 创建必要的索引优化查询性能

2. **backend/migrations/007_create_discussions.down.sql** (1.6KB)
   - 回滚迁移脚本

3. **backend/internal/models/discussion_models.go** (6.6KB)
   - Discussion - 讨论主题/回复模型
   - DiscussionLike - 点赞模型
   - DiscussionFavorite - 收藏模型
   - DiscussionMention - 提及模型
   - DiscussionTag - 标签模型
   - 请求/响应模型（CreateDiscussionRequest, UpdateDiscussionRequest 等）

4. **backend/internal/repository/discussion_repository.go** (22KB)
   - 完整的数据库操作层
   - 支持 CRUD 操作
   - 支持嵌套回复查询
   - 支持热门讨论计算
   - 支持标签管理
   - 支持提及用户功能

5. **backend/internal/services/discussion_service.go** (17KB)
   - 业务逻辑层
   - Markdown 渲染（使用 goldmark）
   - @提及用户处理
   - 点赞/收藏逻辑
   - 热门讨论排序算法
   - 权限验证

6. **backend/internal/handlers/discussion_handler.go** (26KB)
   - HTTP 处理器
   - RESTful API 端点
   - 请求验证
   - 错误处理
   - Swagger 文档注释

### 前端文件 (TypeScript/React)

1. **frontend/src/types/discussion.ts** (2.3KB)
   - TypeScript 类型定义
   - 接口定义（Discussion, DiscussionTag, User 等）

2. **frontend/src/components/DiscussionThread.tsx** (17KB)
   - 可复用的讨论线程组件
   - 支持嵌套回复（最多 10 层）
   - Markdown 编辑器（带工具栏）
   - 点赞/收藏功能
   - 编辑/删除功能
   - @提及用户支持

3. **frontend/src/pages/DiscussionPage.tsx** (25KB)
   - 课程讨论区主页面
   - 讨论列表与筛选
   - 创建新讨论对话框
   - 标签筛选
   - 排序功能（最新、热门、点赞、回复）
   - 分页支持
   - 统计信息展示

4. **frontend/src/services/api.ts** (已更新)
   - 添加讨论相关 API 方法
   - getDiscussions, getDiscussion
   - createDiscussion, updateDiscussion, deleteDiscussion
   - toggleLike, toggleFavorite, resolveDiscussion
   - getDiscussionTags, getHotDiscussions

5. **frontend/src/pages/CourseDetailPage.tsx** (已更新)
   - 添加"课程讨论"按钮
   - 链接到课程讨论区

6. **frontend/src/App.tsx** (已更新)
   - 添加讨论区路由：/courses/:courseId/discussions

## 功能特性

### 1. 课程讨论区
- ✅ 每个课程独立讨论区
- ✅ 支持按课程 ID 筛选
- ✅ 支持按课时筛选（可选）

### 2. 帖子 CRUD
- ✅ 发布新讨论（支持标题、内容、标签、匿名）
- ✅ 编辑讨论（仅作者）
- ✅ 删除讨论（软删除，仅作者/管理员）
- ✅ 查看讨论详情

### 3. 回复功能
- ✅ 嵌套回复（支持最多 10 层）
- ✅ 回复计数
- ✅ @提及用户
- ✅ 回复通知（通过 discussion_mentions 表）

### 4. 点赞/收藏
- ✅ 点赞/点踩功能
- ✅ 收藏/取消收藏
- ✅ 实时计数更新
- ✅ 用户收藏列表

### 5. 热门讨论排序
- ✅ 热门算法：`(upvotes - downvotes) * 2 + reply_count * 3 + hours_since_creation`
- ✅ 多种排序方式：最新、热门、最多点赞、最多回复
- ✅ 置顶讨论支持

### 6. Markdown 支持
- ✅ 完整 Markdown 渲染（使用 goldmark/rehype-highlight）
- ✅ 代码块高亮（支持多种语言）
- ✅ Markdown 编辑器工具栏
- ✅ 实时预览

### 7. @提及用户
- ✅ 自动识别 @username
- ✅ 创建提及记录
- ✅ 未读提及查询

### 8. 标签系统
- ✅ 预定义标签（提问、分享、讨论、建议、已解决、精华）
- ✅ 标签筛选
- ✅ 标签使用计数
- ✅ 自定义标签颜色

### 9. 管理功能
- ✅ 标记为已解决（作者）
- ✅ 锁定/解锁讨论（管理员）
- ✅ 置顶/取消置顶（管理员）
- ✅ 举报功能（UI 已实现）

## API 端点

### 讨论管理
```
POST   /api/v1/discussions              # 创建讨论
GET    /api/v1/discussions              # 获取讨论列表
GET    /api/v1/discussions/:id          # 获取讨论详情
PUT    /api/v1/discussions/:id          # 更新讨论
DELETE /api/v1/discussions/:id          # 删除讨论
```

### 互动功能
```
POST   /api/v1/discussions/:id/like     # 点赞/点踩
POST   /api/v1/discussions/:id/favorite # 收藏
POST   /api/v1/discussions/:id/resolve  # 标记为已解决
```

### 其他
```
GET    /api/v1/discussions/hot          # 热门讨论
GET    /api/v1/discussions/tags         # 所有标签
GET    /api/v1/discussions/favorites    # 用户收藏列表
```

## 数据库设计亮点

1. **自引用嵌套结构**：使用 parent_id 和 depth 字段实现高效嵌套查询
2. **软删除**：使用 deleted_at 字段，支持数据恢复
3. **全文搜索优化**：对 title 和 content 建立 ILIKE 索引
4. **热门计算**：使用数据库计算热门分数，减少应用层负担
5. **唯一约束**：防止重复点赞和收藏

## 使用说明

### 运行迁移
```bash
cd backend
migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname" up
```

### 访问讨论区
1. 进入课程详情页
2. 点击"课程讨论"按钮
3. 或直接访问 `/courses/:courseId/discussions`

## 后续优化建议

1. **性能优化**
   - 添加 Redis 缓存热门讨论
   - 实现讨论列表分页游标
   - 异步处理提及通知

2. **功能增强**
   - 实现实时通知（WebSocket）
   - 添加富文本编辑器选项
   - 支持图片上传
   - 实现搜索高亮

3. **安全加固**
   - 添加内容审核机制
   - 实现反垃圾评论
   - 添加速率限制

## 技术栈

- **后端**: Go 1.21+, Gin, pgx, goldmark
- **前端**: React 18, TypeScript, TailwindCSS
- **数据库**: PostgreSQL 14+
- **Markdown**: goldmark (后端), react-markdown + rehype-highlight (前端)

## 完成时间

2026-03-10

---

*社区讨论系统已完整实现，所有要求的功能均已交付。*
