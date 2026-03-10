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
│   │   └── LessonCard.tsx
│   ├── pages/           # 页面组件
│   │   ├── HomePage.tsx
│   │   ├── CoursesPage.tsx
│   │   ├── CourseDetailPage.tsx
│   │   ├── LearningPathPage.tsx
│   │   └── ProfilePage.tsx
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
└── postcss.config.js
```

## 功能特性

### 页面路由

- **首页 (`/`)** - 展示精选课程、学习路径、平台特色
- **课程列表 (`/courses`)** - 浏览所有课程，支持搜索和筛选
- **课程详情 (`/courses/:id`)** - 课程介绍、大纲、代码练习
- **学习路径 (`/learning-path`)** - 系统化学习路线
- **个人中心 (`/profile`)** - 学习进度、成就徽章、账户设置

### 核心组件

- **Header** - 响应式导航栏，支持移动端
- **CourseCard** - 课程卡片展示
- **CodeEditor** - 基于 Monaco 的在线代码编辑器
- **ProgressBar** - 进度条组件
- **LessonCard** - 课程章节卡片

## 快速开始

### 安装依赖

```bash
npm install
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

## 下一步计划

- [ ] 集成后端 API
- [ ] 添加用户认证
- [ ] 实现代码提交和评测
- [ ] 添加暗色模式切换
- [ ] 优化移动端体验
- [ ] 添加单元测试
- [ ] 性能优化

## 开发规范

- 使用 TypeScript 严格模式
- 组件采用函数式写法
- 使用 TailwindCSS 进行样式开发
- 遵循 React Hooks 最佳实践
- 保持组件单一职责

---

**AI 学习之路** - 让 AI 陪你一起成长 🚀
