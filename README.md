# AI 学习之路 - AI Interactive Learning Platform

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18-blue.svg)](https://reactjs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)](https://www.postgresql.org/)

🎓 **从零开始学习 AI 的交互式平台** - 由浅到深，边学边练，包含课后作业与自动批改系统

---

## 🌟 项目特色

- 📚 **系统化课程** - 5 个 Level，从 Python 基础到深度学习实战
- 💻 **交互式学习** - 在线代码编辑器，实时运行 Python 代码
- ✅ **自动批改** - 课后作业自动评分，即时反馈
- 📊 **进度追踪** - 可视化学习进度，成就系统
- 🎯 **零基础友好** - 无需编程基础，循序渐进

---

## 📖 课程大纲

| Level | 主题 | 课程数 | 核心内容 |
|-------|------|--------|----------|
| Level 1 | Python 编程基础 | 5 课 | 变量、流程控制、循环、函数、模块 |
| Level 2 | 数学基础 | 5 课 | 向量、矩阵、导数、概率、统计 |
| Level 3 | 机器学习入门 | 7 课 | 数据预处理、线性回归、逻辑回归、决策树、SVM、聚类 |
| Level 4 | 深度学习 | 7 课 | 神经网络、CNN、RNN、Transformer、生成式 AI、PyTorch 实战 |
| Level 5 | 实战项目 | 5 个项目 | 房价预测、图像分类、情感分析、推荐系统、综合毕业设计 |

---

## 🛠️ 技术栈

### 前端
- **React 18** + **TypeScript** - 现代化 UI 框架
- **TailwindCSS** - 原子化 CSS
- **Monaco Editor** - VSCode 同款代码编辑器
- **Vite** - 快速构建工具

### 后端
- **Go 1.21+** - 高性能后端语言
- **Gin** - Web 框架
- **PostgreSQL 15** - 关系型数据库
- **Redis 7** - 缓存与会话管理
- **GORM** - ORM 框架

### 代码沙箱
- **Docker** - 容器化隔离
- **Python 3.11** - 运行环境
- **预装库**: numpy, pandas, matplotlib, scikit-learn, torch

### 部署
- **Docker Compose** - 容器编排
- **Nginx** - 反向代理

---

## 🚀 快速开始

### 环境要求
- Docker & Docker Compose
- Node.js 18+
- Go 1.21+

### 1. 克隆项目
```bash
git clone https://github.com/exltial/ai-learning-path.git
cd ai-learning-path
```

### 2. 启动后端
```bash
cd backend
docker-compose up -d
cp .env.example .env
# 编辑 .env 配置数据库连接
go run cmd/main.go
```

### 3. 启动前端
```bash
cd frontend
npm install
npm run dev
```

### 4. 访问
- 前端：http://localhost:5173
- 后端 API：http://localhost:8080
- API 文档：http://localhost:8080/swagger

---

## 📁 项目结构

```
ai-learning-path/
├── frontend/              # React 前端
│   ├── src/
│   │   ├── components/    # 可复用组件
│   │   ├── pages/         # 页面组件
│   │   ├── hooks/         # 自定义 hooks
│   │   └── utils/         # 工具函数
│   └── package.json
├── backend/               # Go 后端
│   ├── cmd/               # 应用入口
│   ├── internal/          # 业务逻辑
│   │   ├── handlers/      # HTTP 处理器
│   │   ├── services/      # 服务层
│   │   ├── repository/    # 数据访问
│   │   └── models/        # 数据模型
│   └── go.mod
├── sandbox/               # 代码沙箱服务
├── docs/                  # 文档
│   └── curriculum.md      # 课程大纲
├── docker-compose.yml
└── README.md
```

---

## 📋 开发计划

### Phase 1: MVP ✅ (已完成)
- [x] 项目架构设计
- [x] 数据库 Schema
- [x] API 接口设计
- [x] 前端项目初始化
- [x] 课程大纲设计

### Phase 2: 核心功能 (进行中)
- [ ] 用户注册/登录系统
- [ ] 课程列表与详情页
- [ ] 代码编辑器集成
- [ ] 作业提交功能

### Phase 3: 交互增强
- [ ] 代码沙箱服务
- [ ] 作业自动批改
- [ ] 进度追踪系统
- [ ] 成就系统

### Phase 4: 内容丰富
- [ ] 完整课程内容填充
- [ ] 视频讲解
- [ ] 社区讨论功能
- [ ] 移动端适配

---

## 👥 团队

- **项目经理**: AI Assistant
- **架构师**: AI Architect
- **前端开发**: AI Frontend Dev
- **内容策划**: AI Content Designer

---

## 📄 开源协议

MIT License - 详见 [LICENSE](LICENSE) 文件

---

## 📬 联系方式

- GitHub: [@exltial](https://github.com/exltial)
- 邮箱：exltial@163.com

---

**🎯 让 AI 学习变得简单有趣！**
