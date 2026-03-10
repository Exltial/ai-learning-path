# 更新日志 / Changelog

本项目的所有重要更改都将记录在此文件中。

All notable changes to this project will be documented in this file.

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

This format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [未发布] / [Unreleased]

### 计划中 / Planned

#### 新增 / Added
- 用户注册/登录系统
- 课程列表与详情页
- 代码编辑器集成
- 作业提交功能
- 学习进度追踪
- 成就系统

#### 改进 / Improved
- 前端性能优化
- API 响应速度提升
- 数据库查询优化

#### 修复 / Fixed
- 待报告的 Bug

---

## [0.2.0] - 2026-03-10

### 新增 / Added

#### 前端 / Frontend
- ✨ 用户认证系统
  - 登录页面（邮箱 + 密码）
  - 注册页面（用户名、邮箱、密码）
  - 表单验证（前端验证 + 错误提示）
  - 密码强度检测
  - 记住我功能
  - 密码显示/隐藏切换
  - JWT Token 管理
  - 自动登录状态保持

- ✨ 课程列表增强
  - 实时搜索（课程名、描述、标签）
  - 难度筛选（入门/进阶/高级）
  - 标签筛选（多选）
  - 排序功能（最受欢迎/评分最高/最新发布）
  - 视图切换（网格/列表）
  - 活跃筛选器展示
  - 一键清除筛选
  - 加载状态骨架屏

- ✨ 课程详情页增强
  - 课程大纲展示（可折叠章节）
  - 开始学习/继续学习按钮
  - 课程进度展示
  - 课时完成状态标记
  - 锁定课时显示
  - 代码练习标签页
  - 课程概述（学习目标、要求、讲师）
  - 多标签切换

- ✨ Monaco 代码编辑器集成
  - 多语言支持（Python、JavaScript、TypeScript 等）
  - 主题切换（亮色/暗色）
  - 代码高亮
  - 自动布局
  - 运行/重置代码按钮
  - 占位符提示

- ✨ 学习进度展示
  - 学习统计概览（课程数、完成数、学习时长）
  - 课程进度列表
  - 进度条可视化
  - 最后学习时间
  - 成就徽章系统
  - 紧凑型和完整型两种展示模式

- ✨ API 服务封装
  - 统一的 API 客户端
  - JWT Token 自动注入
  - 错误处理
  - 类型安全的请求/响应

#### 后端 / Backend
- ✨ 完整的 RESTful API
  - 认证与授权（注册、登录、刷新 Token、登出）
  - 用户管理（CRUD、统计信息）
  - 课程管理（CRUD、注册、章节列表）
  - 章节管理（CRUD、内容管理）
  - 练习管理（CRUD、多种题型支持）
  - 提交与评分（自动批改、反馈）
  - 学习进度（追踪、更新）
  - 成就系统（徽章、积分）
  - 讨论区（创建、回复、点赞）
  - 通知系统（推送、标记已读）
  - 课程评价（评分、评论）

- ✨ 数据库设计
  - PostgreSQL Schema 设计
  - 用户表、课程表、章节表、练习表
  - 提交表、进度表、成就表
  - 讨论表、通知表、评价表
  - 索引优化

- ✨ 中间件
  - JWT 认证中间件
  - CORS 配置
  - 速率限制
  - 错误处理
  - 日志记录

- ✨ 代码沙箱集成
  - Docker 容器隔离
  - Python 代码执行
  - 超时控制
  - 资源限制

#### 文档 / Documentation
- ✨ 项目文档
  - README.md（中英双语）
  - DEPLOYMENT.md（部署指南）
  - API.md（API 参考文档）
  - CONTRIBUTING.md（贡献指南）
  - CHANGELOG.md（更新日志）
  - QUICKSTART.md（快速开始）
  - curriculum.md（课程大纲）

- ✨ API 文档
  - Swagger/OpenAPI 规范
  - 交互式 API 文档
  - 代码示例（JavaScript、Go、Python）

### 改进 / Improved

#### 前端 / Frontend
- 🚀 性能优化
  - 组件懒加载
  - 代码分割
  - 图片优化
  - 缓存策略

- 🎨 UI/UX 改进
  - 响应式设计优化
  - 移动端适配
  - 加载状态优化
  - 错误提示优化

#### 后端 / Backend
- 🚀 性能优化
  - 数据库连接池
  - Redis 缓存
  - 查询优化
  - 索引优化

- 🔒 安全增强
  - 密码 bcrypt 加密
  - JWT Token 安全配置
  - CORS 严格模式
  - 输入验证

### 修复 / Fixed

- 🐛 修复前端路由问题
- 🐛 修复 API 跨域问题
- 🐛 修复数据库连接问题
- 🐛 修复代码编辑器渲染问题

### 技术债务 / Technical Debt

- ⚠️ 需要添加单元测试
- ⚠️ 需要添加集成测试
- ⚠️ 需要完善错误处理
- ⚠️ 需要优化数据库迁移

---

## [0.1.0] - 2026-03-05

### 新增 / Added

#### 项目初始化 / Project Initialization
- ✨ 项目架构设计
  - 前后端分离架构
  - Docker 容器化部署
  - 微服务设计

- ✨ 技术栈选型
  - 前端：React 18 + TypeScript + Vite
  - 后端：Go 1.21+ + Gin
  - 数据库：PostgreSQL 15+
  - 缓存：Redis 7+
  - 代码沙箱：Docker

- ✨ 课程大纲设计
  - Level 1: Python 编程基础（5 课）
  - Level 2: 数学基础（5 课）
  - Level 3: 机器学习入门（7 课）
  - Level 4: 深度学习（7 课）
  - Level 5: 实战项目（5 个项目）

#### 基础设施 / Infrastructure
- ✨ Docker Compose 配置
  - PostgreSQL 服务
  - Redis 服务
  - 后端服务
  - 前端服务
  - 代码沙箱服务

- ✨ 环境变量配置
  - 开发环境配置
  - 生产环境配置
  - 敏感信息管理

- ✨ Git 配置
  - .gitignore
  - 分支策略
  - 提交规范

#### 后端基础 / Backend Foundation
- ✨ 项目结构
  - cmd/ - 应用入口
  - internal/ - 业务逻辑
  - pkg/ - 公共库
  - configs/ - 配置文件
  - migrations/ - 数据库迁移

- ✨ 数据库 Schema
  - 用户表设计
  - 课程表设计
  - 章节表设计
  - 练习表设计
  - 基础索引

- ✨ API 设计
  - RESTful 规范
  - 错误处理规范
  - 响应格式规范

#### 前端基础 / Frontend Foundation
- ✨ 项目结构
  - src/components/ - 组件
  - src/pages/ - 页面
  - src/hooks/ - Hooks
  - src/services/ - API 服务
  - src/types/ - 类型定义

- ✨ 基础组件
  - Header - 导航栏
  - Footer - 页脚
  - Layout - 布局
  - CourseCard - 课程卡片
  - ProgressBar - 进度条

- ✨ 页面框架
  - HomePage - 首页
  - CoursesPage - 课程列表
  - CourseDetailPage - 课程详情
  - LearningPathPage - 学习路径
  - ProfilePage - 个人中心

- ✨ 样式系统
  - TailwindCSS 配置
  - 设计令牌
  - 响应式断点

### 技术栈 / Tech Stack

#### 前端 / Frontend
- React 18
- TypeScript 5
- Vite 5
- TailwindCSS 3
- React Router v6
- Monaco Editor
- Lucide React

#### 后端 / Backend
- Go 1.21+
- Gin
- GORM
- PostgreSQL 15
- Redis 7
- JWT
- Swagger

#### 开发工具 / Development Tools
- Docker & Docker Compose
- Git
- Make
- ESLint
- Prettier

---

## 版本说明 / Version Notes

### 语义化版本 / Semantic Versioning

本项目遵循语义化版本 2.0.0：

- **主版本号（Major）**：不兼容的 API 变更
- **次版本号（Minor）**：向后兼容的功能性新增
- **修订号（Patch）**：向后兼容的问题修正

### 发布周期 / Release Cycle

- **主要版本**：每季度发布
- **次要版本**：每月发布
- **补丁版本**：按需发布

### 支持政策 / Support Policy

- **当前版本**：完全支持
- **上一个版本**：关键 Bug 修复
- **更早版本**：不再支持

---

## 贡献者 / Contributors

感谢所有为这个项目做出贡献的人！

Thank you to all the people who have contributed to this project!

---

## 相关链接 / Related Links

- [GitHub Repository](https://github.com/exltial/ai-learning-path)
- [Issue Tracker](https://github.com/exltial/ai-learning-path/issues)
- [Project Board](https://github.com/exltial/ai-learning-path/projects)
- [Documentation](./docs/)

---

*最后更新 / Last Updated: 2026-03-10*
