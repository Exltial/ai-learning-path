# AI 学习之路 - AI Interactive Learning Platform

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18-blue.svg)](https://reactjs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)](https://www.postgresql.org/)

🎓 **从零开始学习 AI 的交互式平台** - 由浅到深，边学边练，包含课后作业与自动批改系统

🎓 **Interactive Platform for Learning AI from Scratch** - Progressive learning with hands-on exercises and automated grading

---

## 🌟 项目特色 / Features

### 中文
- 📚 **系统化课程** - 5 个 Level，从 Python 基础到深度学习实战
- 💻 **交互式学习** - 在线代码编辑器，实时运行 Python 代码
- ✅ **自动批改** - 课后作业自动评分，即时反馈
- 📊 **进度追踪** - 可视化学习进度，成就系统
- 🎯 **零基础友好** - 无需编程基础，循序渐进

### English
- 📚 **Systematic Curriculum** - 5 Levels, from Python basics to deep learning practice
- 💻 **Interactive Learning** - Online code editor with real-time Python execution
- ✅ **Automated Grading** - Instant feedback on exercises with auto-scoring
- 📊 **Progress Tracking** - Visual learning progress and achievement system
- 🎯 **Beginner Friendly** - No programming background required, step-by-step learning

---

## 📖 课程大纲 / Curriculum

| Level | 主题 / Theme | 课程数 / Courses | 核心内容 / Core Content |
|-------|------|--------|----------|
| Level 1 | Python 编程基础 / Python Basics | 5 课 | 变量、流程控制、循环、函数、模块 / Variables, Control Flow, Loops, Functions, Modules |
| Level 2 | 数学基础 / Math Fundamentals | 5 课 | 向量、矩阵、导数、概率、统计 / Vectors, Matrices, Derivatives, Probability, Statistics |
| Level 3 | 机器学习入门 / ML Introduction | 7 课 | 数据预处理、线性回归、逻辑回归、决策树、SVM、聚类 / Data Preprocessing, Linear/Logistic Regression, Decision Trees, SVM, Clustering |
| Level 4 | 深度学习 / Deep Learning | 7 课 | 神经网络、CNN、RNN、Transformer、生成式 AI、PyTorch 实战 / Neural Networks, CNN, RNN, Transformer, Generative AI, PyTorch |
| Level 5 | 实战项目 / Capstone Projects | 5 个项目 | 房价预测、图像分类、情感分析、推荐系统、综合毕业设计 / Price Prediction, Image Classification, Sentiment Analysis, Recommendation System, Capstone |

---

## 🛠️ 技术栈 / Tech Stack

### 前端 / Frontend
- **React 18** + **TypeScript** - 现代化 UI 框架 / Modern UI Framework
- **TailwindCSS** - 原子化 CSS / Utility-first CSS
- **Monaco Editor** - VSCode 同款代码编辑器 / VSCode Code Editor
- **Vite** - 快速构建工具 / Fast Build Tool

### 后端 / Backend
- **Go 1.21+** - 高性能后端语言 / High-performance Backend Language
- **Gin** - Web 框架 / Web Framework
- **PostgreSQL 15** - 关系型数据库 / Relational Database
- **Redis 7** - 缓存与会话管理 / Cache & Session Management
- **GORM** - ORM 框架 / ORM Framework

### 代码沙箱 / Code Sandbox
- **Docker** - 容器化隔离 / Containerized Isolation
- **Python 3.11** - 运行环境 / Runtime Environment
- **预装库 / Pre-installed Libraries**: numpy, pandas, matplotlib, scikit-learn, torch

### 部署 / Deployment
- **Docker Compose** - 容器编排 / Container Orchestration
- **Nginx** - 反向代理 / Reverse Proxy

---

## 🚀 快速开始 / Quick Start

### 环境要求 / Prerequisites
- Docker & Docker Compose
- Node.js 18+
- Go 1.21+

### 1. 克隆项目 / Clone Repository
```bash
git clone https://github.com/exltial/ai-learning-path.git
cd ai-learning-path
```

### 2. 启动后端 / Start Backend
```bash
cd backend
docker-compose up -d
cp .env.example .env
# 编辑 .env 配置数据库连接 / Edit .env for database connection
go run cmd/main.go
```

### 3. 启动前端 / Start Frontend
```bash
cd frontend
npm install
npm run dev
```

### 4. 访问 / Access
- 前端 / Frontend: http://localhost:5173
- 后端 API / Backend API: http://localhost:8080
- API 文档 / API Docs: http://localhost:8080/swagger

📖 **详细部署指南请查看** / **For detailed deployment guide, see:**
- [部署指南 / Deployment Guide](docs/DEPLOYMENT.md)
- [快速开始 / Quick Start Guide](QUICKSTART.md)

---

## 📁 项目结构 / Project Structure

```
ai-learning-path/
├── frontend/              # React 前端 / React Frontend
│   ├── src/
│   │   ├── components/    # 可复用组件 / Reusable Components
│   │   ├── pages/         # 页面组件 / Page Components
│   │   ├── hooks/         # 自定义 hooks / Custom Hooks
│   │   └── utils/         # 工具函数 / Utility Functions
│   └── package.json
├── backend/               # Go 后端 / Go Backend
│   ├── cmd/               # 应用入口 / Application Entry
│   ├── internal/          # 业务逻辑 / Business Logic
│   │   ├── handlers/      # HTTP 处理器 / HTTP Handlers
│   │   ├── services/      # 服务层 / Service Layer
│   │   ├── repository/    # 数据访问 / Data Access
│   │   └── models/        # 数据模型 / Data Models
│   └── go.mod
├── sandbox/               # 代码沙箱服务 / Code Sandbox Service
├── docs/                  # 文档 / Documentation
│   ├── DEPLOYMENT.md      # 部署指南 / Deployment Guide
│   ├── API.md             # API 文档 / API Documentation
│   ├── CONTRIBUTING.md    # 贡献指南 / Contributing Guide
│   └── curriculum.md      # 课程大纲 / Curriculum
├── docker-compose.yml
├── README.md
└── QUICKSTART.md          # 快速开始 / Quick Start
```

---

## 📋 开发计划 / Development Roadmap

### Phase 1: MVP ✅ (已完成 / Completed)
- [x] 项目架构设计 / Project Architecture
- [x] 数据库 Schema / Database Schema
- [x] API 接口设计 / API Design
- [x] 前端项目初始化 / Frontend Initialization
- [x] 课程大纲设计 / Curriculum Design

### Phase 2: 核心功能 (进行中 / In Progress)
- [ ] 用户注册/登录系统 / User Registration/Login
- [ ] 课程列表与详情页 / Course List & Detail Pages
- [ ] 代码编辑器集成 / Code Editor Integration
- [ ] 作业提交功能 / Assignment Submission

### Phase 3: 交互增强 / Enhanced Interaction
- [ ] 代码沙箱服务 / Code Sandbox Service
- [ ] 作业自动批改 / Automated Grading
- [ ] 进度追踪系统 / Progress Tracking
- [ ] 成就系统 / Achievement System

### Phase 4: 内容丰富 / Content Expansion
- [ ] 完整课程内容填充 / Complete Course Content
- [ ] 视频讲解 / Video Lectures
- [ ] 社区讨论功能 / Community Discussion
- [ ] 移动端适配 / Mobile Adaptation

---

## 📚 文档 / Documentation

| 文档 / Document | 描述 / Description |
|----------------|-------------------|
| [README.md](README.md) | 项目概述 / Project Overview |
| [QUICKSTART.md](QUICKSTART.md) | 快速开始指南 / Quick Start Guide |
| [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) | 部署指南 / Deployment Guide |
| [docs/API.md](docs/API.md) | API 参考文档 / API Reference |
| [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) | 贡献指南 / Contributing Guide |
| [CHANGELOG.md](CHANGELOG.md) | 更新日志 / Changelog |
| [docs/curriculum.md](docs/curriculum.md) | 课程大纲 / Curriculum |

---

## 👥 团队 / Team

- **项目经理 / Project Manager**: AI Assistant
- **架构师 / Architect**: AI Architect
- **前端开发 / Frontend Developer**: AI Frontend Dev
- **内容策划 / Content Designer**: AI Content Designer

---

## 🤝 贡献 / Contributing

我们欢迎各种形式的贡献！请查看 [贡献指南](docs/CONTRIBUTING.md) 了解如何参与。

We welcome contributions of all kinds! Please see our [Contributing Guide](docs/CONTRIBUTING.md) for details.

---

## 📄 开源协议 / License

MIT License - 详见 [LICENSE](LICENSE) 文件 / See [LICENSE](LICENSE) file

---

## 📬 联系方式 / Contact

- GitHub: [@exltial](https://github.com/exltial)
- 邮箱 / Email: exltial@163.com

---

**🎯 让 AI 学习变得简单有趣！ / Make AI Learning Simple and Fun!**

---

*最后更新 / Last Updated: 2026-03-10*
