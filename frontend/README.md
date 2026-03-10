# AI 学习之路 - 前端项目

基于 React 18 + TypeScript + TailwindCSS + Monaco Editor 的现代化在线学习平台前端。

## 技术栈

- **React 18** - 前端框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **TailwindCSS** - 样式框架
- **React Router v6** - 路由管理
- **Monaco Editor** - 代码编辑器
- **Lucide React** - 图标库

## 项目结构

```
frontend/
├── public/              # 静态资源
├── src/
│   ├── components/      # 可复用组件
│   │   ├── Header.tsx
│   │   ├── Footer.tsx
│   │   ├── Layout.tsx
│   │   ├── CourseCard.tsx
│   │   ├── CodeEditor.tsx
│   │   ├── ProgressBar.tsx
│   │   ├── LessonCard.tsx
│   │   └── LearningProgress.tsx
│   ├── pages/           # 页面组件
│   │   ├── HomePage.tsx
│   │   ├── CoursesPage.tsx
│   │   ├── CourseDetailPage.tsx
│   │   ├── LearningPathPage.tsx
│   │   ├── ProfilePage.tsx
│   │   ├── LoginPage.tsx
│   │   └── RegisterPage.tsx
│   ├── contexts/        # React Context
│   │   └── AuthContext.tsx
│   ├── services/        # API 服务
│   │   └── api.ts
│   ├── types/           # TypeScript 类型定义
│   │   └── index.ts
│   ├── hooks/           # 自定义 Hooks
│   ├── utils/           # 工具函数
│   ├── assets/          # 资源文件
│   ├── App.tsx          # 应用入口
│   ├── main.tsx         # React 入口
│   └── index.css        # 全局样式
├── package.json
├── tsconfig.json
├── vite.config.ts
├── tailwind.config.js
├── postcss.config.js
└── .env.example
```

## 功能特性

### 页面路由

- **登录/注册 (`/login`, `/register`)** - 用户认证页面，支持表单验证
- **首页 (`/`)** - 展示精选课程、学习路径、平台特色
- **课程列表 (`/courses`)** - 浏览所有课程，支持搜索、筛选、排序
- **课程详情 (`/courses/:id`)** - 课程介绍、大纲、代码练习
- **学习路径 (`/learning-path`)** - 系统化学习路线
- **个人中心 (`/profile`)** - 学习进度、成就徽章、账户设置

### 核心组件

- **Header** - 响应式导航栏，支持移动端，集成用户认证状态
- **CourseCard** - 课程卡片展示
- **CodeEditor** - 基于 Monaco 的在线代码编辑器，支持多语言
- **ProgressBar** - 进度条组件，支持多种尺寸和颜色
- **LessonCard** - 课程章节卡片
- **LearningProgress** - 学习进度展示组件

### Phase 2 新增功能

#### 1. 用户认证系统
- ✅ 登录页面（邮箱 + 密码）
- ✅ 注册页面（用户名、邮箱、密码）
- ✅ 表单验证（前端验证 + 错误提示）
- ✅ 密码强度检测
- ✅ 记住我功能
- ✅ 密码显示/隐藏切换
- ✅ JWT Token 管理
- ✅ 自动登录状态保持

#### 2. 课程列表增强
- ✅ 实时搜索（课程名、描述、标签）
- ✅ 难度筛选（入门/进阶/高级）
- ✅ 标签筛选（多选）
- ✅ 排序功能（最受欢迎/评分最高/最新发布）
- ✅ 视图切换（网格/列表）
- ✅ 活跃筛选器展示
- ✅ 一键清除筛选
- ✅ 加载状态骨架屏

#### 3. 课程详情页增强
- ✅ 课程大纲展示（可折叠章节）
- ✅ 开始学习/继续学习按钮
- ✅ 课程进度展示
- ✅ 课时完成状态标记
- ✅ 锁定课时显示
- ✅ 代码练习标签页
- ✅ 课程概述（学习目标、要求、讲师）
- ✅ 多标签切换

#### 4. Monaco 代码编辑器集成
- ✅ 多语言支持（Python、JavaScript、TypeScript 等）
- ✅ 主题切换（亮色/暗色）
- ✅ 代码高亮
- ✅ 自动布局
- ✅ 运行/重置代码按钮
- ✅ 占位符提示

#### 5. 学习进度展示
- ✅ 学习统计概览（课程数、完成数、学习时长）
- ✅ 课程进度列表
- ✅ 进度条可视化
- ✅ 最后学习时间
- ✅ 成就徽章系统
- ✅ 紧凑型和完整型两种展示模式

### API 服务封装

- ✅ 统一的 API 客户端
- ✅ JWT Token 自动注入
- ✅ 错误处理
- ✅ 类型安全的请求/响应
- ✅ 支持以下接口：
  - 用户注册/登录/登出
  - 获取当前用户信息
  - 获取课程列表（带筛选）
  - 获取课程详情
  - 获取课程章节
  - 获取/更新学习进度

## 快速开始

### 安装依赖

```bash
npm install
```

### 环境配置

复制环境变量文件并修改配置：

```bash
cp .env.example .env
```

编辑 `.env` 文件：

```env
# API Configuration
VITE_API_URL=http://localhost:3000/api

# App Configuration
VITE_APP_NAME=AI 学习之路
VITE_APP_VERSION=1.0.0
```

### 开发模式

```bash
npm run dev
```

访问 http://localhost:3000

### 构建生产版本

```bash
npm run build
```

### 预览生产构建

```bash
npm run preview
```

### 代码检查

```bash
npm run lint
```

## 设计系统

### 颜色

- **Primary** - 主色调 (蓝色系)
- **Secondary** - 次色调 (灰色系)
- **Success** - 成功状态 (绿色)
- **Warning** - 警告状态 (黄色)
- **Danger** - 危险状态 (红色)

### 组件样式

使用 TailwindCSS 工具类，配合自定义组件类：

- `.btn-primary` - 主按钮
- `.btn-secondary` - 次按钮
- `.card` - 卡片容器

## 响应式设计

- ✅ 移动端优先
- ✅ 断点：sm (640px), md (768px), lg (1024px), xl (1280px)
- ✅ 移动端导航菜单
- ✅ 自适应网格布局
- ✅ 触摸友好的按钮尺寸

## 表单验证

### 登录表单
- 邮箱格式验证
- 密码长度验证（最少 6 位）
- 必填项验证
- 实时错误提示

### 注册表单
- 用户名长度验证（3-20 位）
- 邮箱格式验证
- 密码强度检测（最少 8 位，包含字母和数字）
- 密码确认验证
- 密码强度指示器
- 服务条款同意

## API 接口约定

### 认证接口
```
POST /api/auth/register - 用户注册
POST /api/auth/login - 用户登录
POST /api/auth/logout - 用户登出
GET  /api/auth/me - 获取当前用户
```

### 课程接口
```
GET  /api/courses - 获取课程列表（支持筛选参数）
GET  /api/courses/:id - 获取课程详情
GET  /api/courses/:id/lessons - 获取课程章节
```

### 进度接口
```
GET  /api/progress - 获取学习进度
POST /api/progress - 更新学习进度
```

## 下一步计划

- [ ] 集成真实后端 API
- [ ] 添加代码提交和评测功能
- [ ] 实现视频播放器
- [ ] 添加暗色模式切换
- [ ] 添加单元测试
- [ ] 性能优化（代码分割、懒加载）
- [ ] PWA 支持
- [ ] 添加学习提醒功能

## 开发规范

- 使用 TypeScript 严格模式
- 组件采用函数式写法
- 使用 TailwindCSS 进行样式开发
- 遵循 React Hooks 最佳实践
- 保持组件单一职责
- 使用语义化的 HTML 标签
- 添加适当的 ARIA 属性

## 浏览器支持

- Chrome (最新)
- Firefox (最新)
- Safari (最新)
- Edge (最新)
- 移动端浏览器

---

**AI 学习之路** - 让 AI 陪你一起成长 🚀
