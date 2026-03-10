# 贡献指南 / Contributing Guide

我们欢迎各种形式的贡献！本指南将帮助你了解如何参与项目开发。

We welcome contributions of all kinds! This guide will help you understand how to participate in the project development.

---

## 📋 目录 / Table of Contents

1. [行为准则 / Code of Conduct](#行为准则--code-of-conduct)
2. [贡献方式 / Ways to Contribute](#贡献方式--ways-to-contribute)
3. [开发环境设置 / Development Setup](#开发环境设置--development-setup)
4. [开发流程 / Development Workflow](#开发流程--development-workflow)
5. [代码规范 / Coding Standards](#代码规范--coding-standards)
6. [提交规范 / Commit Guidelines](#提交规范--commit-guidelines)
7. [Pull Request 流程 / PR Process](#pull-request-流程--pr-process)
8. [测试指南 / Testing Guide](#测试指南--testing-guide)
9. [文档贡献 / Documentation](#文档贡献--documentation)
10. [常见问题 / FAQ](#常见问题--faq)

---

## 行为准则 / Code of Conduct

### 我们的承诺

为了营造一个开放和友好的环境，我们承诺：

- 使用友好和包容的语言
- 尊重不同的观点和经验
- 优雅地接受建设性批评
- 关注对社区最有利的事情
- 对其他社区成员表示同理心

### 不可接受的行为

- 使用性化的语言或图像
- 人身攻击或侮辱性评论
- 公开或私下骚扰
- 未经许可发布他人信息
- 其他不道德或不专业的行为

### 执行

不当行为可通过 [exltial@163.com](mailto:exltial@163.com) 报告。所有投诉都将得到审查和调查。

---

## 贡献方式 / Ways to Contribute

### 1. 报告 Bug / Reporting Bugs

发现 Bug？请创建 Issue 并包含：

- 清晰描述性的标题
- 详细的问题描述
- 复现步骤
- 预期行为 vs 实际行为
- 环境信息（操作系统、浏览器、版本等）
- 截图或录屏（如适用）

**Bug 报告模板:**

```markdown
### 问题描述


### 复现步骤

1. 
2. 
3. 

### 预期行为


### 实际行为


### 环境信息

- OS: 
- Browser: 
- Version: 

### 附加信息

```

### 2. 功能建议 / Feature Requests

有新想法？欢迎提出！

- 描述功能需求
- 说明使用场景
- 提供可能的实现方案
- 讨论替代方案

### 3. 提交代码 / Submitting Code

- 修复 Bug
- 实现新功能
- 性能优化
- 文档改进
- 测试用例

### 4. 其他贡献 / Other Contributions

- 文档翻译
- 设计 UI/UX
- 内容创作
- 社区帮助
- 宣传推广

---

## 开发环境设置 / Development Setup

### 前置要求 / Prerequisites

```bash
# 必需 / Required
Git
Node.js 18+
Go 1.21+
PostgreSQL 15+
Redis 7+
Docker & Docker Compose (推荐)

# 可选 / Optional
Make
```

### 克隆项目 / Clone Repository

```bash
# Fork 项目
# 访问 https://github.com/exltial/ai-learning-path 并点击 Fork

# 克隆你的 Fork
git clone https://github.com/YOUR_USERNAME/ai-learning-path.git
cd ai-learning-path

# 添加上游仓库
git remote add upstream https://github.com/exltial/ai-learning-path.git

# 验证远程仓库
git remote -v
```

### 安装依赖 / Install Dependencies

```bash
# 后端 / Backend
cd backend
go mod download

# 前端 / Frontend
cd frontend
npm install
```

### 配置环境变量 / Configure Environment

```bash
# 后端
cd backend
cp .env.example .env
# 编辑 .env 文件

# 前端
cd frontend
cp .env.example .env
# 编辑 .env 文件
```

### 启动开发服务 / Start Development Servers

```bash
# 方式一：使用 Docker Compose（推荐）
cd backend
docker-compose up -d postgres redis
cd ..
make dev  # 或手动启动前后端

# 方式二：手动启动
# 终端 1 - 后端
cd backend
go run cmd/main.go

# 终端 2 - 前端
cd frontend
npm run dev
```

---

## 开发流程 / Development Workflow

### 1. 创建分支 / Create Branch

```bash
# 确保基于最新的主分支
git checkout main
git pull upstream main

# 创建功能分支
git checkout -b feature/your-feature-name
# 或修复 Bug
git checkout -b fix/bug-description
```

### 分支命名规范 / Branch Naming Convention

```
feature/xxx      # 新功能
fix/xxx          # Bug 修复
docs/xxx         # 文档更新
refactor/xxx     # 代码重构
test/xxx         # 测试相关
chore/xxx        # 构建/工具相关
perf/xxx         # 性能优化
style/xxx        # 代码格式
```

### 2. 开发与提交 / Develop and Commit

```bash
# 开发你的功能
# ...

# 添加更改
git add .

# 提交（遵循提交规范）
git commit -m "feat: add user authentication system"

# 推送到你的 Fork
git push origin feature/your-feature-name
```

### 3. 保持同步 / Stay Synced

```bash
# 定期同步上游仓库
git fetch upstream
git rebase upstream/main

# 解决冲突后
git push origin feature/your-feature-name --force-with-lease
```

---

## 代码规范 / Coding Standards

### Go 代码规范 / Go Coding Standards

```go
// ✅ 好的做法

// 使用有意义的变量名
func CalculateUserScore(userID string) (float64, error) {
    // ...
}

// 错误处理
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// 注释使用英文或中文，保持一致性
// GetUser retrieves a user by ID
func GetUser(id string) (*User, error) {
    // ...
}

// 遵循 Go 格式化
gofmt -w .
```

```go
// ❌ 不好的做法

// 变量名无意义
func calc(uid string) (float64, error) {
    // ...
}

// 忽略错误
result, _ := doSomething()

// 混合语言注释
// GetUser 获取用户 by ID
```

### TypeScript/React 代码规范

```typescript
// ✅ 好的做法

// 使用 TypeScript 严格模式
interface User {
  id: string;
  username: string;
  email: string;
}

// 函数组件
const UserProfile: React.FC<UserProps> = ({ user }) => {
  return <div>{user.username}</div>;
};

// 使用有意义的类型
type LoadingState = 'idle' | 'loading' | 'success' | 'error';

// 错误边界
class ErrorBoundary extends React.Component {
  // ...
}
```

```typescript
// ❌ 不好的做法

// 使用 any 类型
const data: any = fetchData();

// 混合默认导出和命名导出
export default function Component() {}
export const util = () => {};

// 过大的组件（>500 行）
```

### 通用规范 / General Standards

1. **代码格式 / Code Formatting**
   - Go: 使用 `gofmt`
   - TypeScript: 使用 Prettier
   - 保持一致的缩进（2 空格或 4 空格）

2. **命名规范 / Naming**
   - Go: 驼峰命名（CamelCase）
   - TypeScript: 驼峰命名（camelCase）
   - 文件名：小写 + 连字符（kebab-case）

3. **注释 / Comments**
   - 公共函数必须有文档注释
   - 复杂逻辑需要解释
   - 避免无意义的注释

4. **函数长度 / Function Length**
   - 函数不超过 50 行
   - 单一职责原则

---

## 提交规范 / Commit Guidelines

### Commit Message 格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type 类型

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建/工具相关
- `perf`: 性能优化
- `ci`: CI/CD 配置

### 示例 / Examples

```bash
# 新功能
feat(auth): add user registration system

- Implement registration API endpoint
- Add email validation
- Add password strength checker

Closes #123

# Bug 修复
fix(api): resolve database connection timeout

The connection pool was not properly configured,
causing timeouts under high load.

Fixes #456

# 文档更新
docs(readme): add installation instructions

Add detailed installation steps for different
operating systems.
```

### 提交前检查 / Pre-commit Checklist

```bash
# Go
cd backend
go fmt ./...
go vet ./...
go test ./...

# Frontend
cd frontend
npm run lint
npm run test
npm run build
```

---

## Pull Request 流程 / PR Process

### 1. 创建 Pull Request

```bash
# 推送到你的 Fork
git push origin feature/your-feature-name

# 访问 GitHub，点击 "New Pull Request"
# 选择 base: main <- compare: your-branch
```

### 2. PR 模板 / PR Template

```markdown
## 描述 / Description

简要描述此 PR 的目的。

## 类型 / Type

- [ ] 🚀 新功能 (feat)
- [ ] 🐛 Bug 修复 (fix)
- [ ] 📝 文档 (docs)
- [ ] ♻️ 重构 (refactor)
- [ ] 🎨 样式 (style)
- [ ] ✅ 测试 (test)
- [ ] ⚡ 性能 (perf)
- [ ] 🔧 配置 (chore)

## 测试 / Testing

- [ ] 已添加测试用例
- [ ] 所有测试通过
- [ ] 已手动测试

## 截图 / Screenshots

（如适用）

## 相关问题 / Related Issues

Closes #123
```

### 3. Code Review

- 至少需要 1 个维护者批准
- 解决所有评论和建议
- 确保 CI 检查通过

### 4. 合并 / Merge

- 使用 Squash and Merge（推荐）
- 或删除合并后删除分支

---

## 测试指南 / Testing Guide

### Go 测试 / Go Testing

```go
// backend/internal/handlers/auth_handler_test.go
package handlers

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
    // 设置测试环境
    setupTestDB()
    defer teardownTestDB()

    // 创建测试用户
    createUser(testUser)

    // 测试登录
    t.Run("successful login", func(t *testing.T) {
        req := httptest.NewRequest("POST", "/api/v1/auth/login", 
            strings.NewReader(`{"email":"test@example.com","password":"password123"}`))
        w := httptest.NewRecorder()
        
        LoginHandler(w, req)
        
        assert.Equal(t, 200, w.Code)
        assert.Contains(t, w.Body.String(), "token")
    })

    t.Run("invalid credentials", func(t *testing.T) {
        // ...
    })
}
```

运行测试：

```bash
cd backend
go test ./... -v
go test ./... -race -cover
```

### Frontend 测试 / Frontend Testing

```typescript
// frontend/src/components/__tests__/LoginForm.test.tsx
import { render, screen, fireEvent } from '@testing-library/react';
import LoginForm from '../LoginForm';

describe('LoginForm', () => {
  it('renders correctly', () => {
    render(<LoginForm />);
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
  });

  it('validates email format', async () => {
    render(<LoginForm />);
    
    fireEvent.change(screen.getByLabelText(/email/i), {
      target: { value: 'invalid-email' }
    });
    
    expect(await screen.findByText(/invalid email/i)).toBeInTheDocument();
  });

  it('submits form successfully', async () => {
    // ...
  });
});
```

运行测试：

```bash
cd frontend
npm test
npm run test:coverage
```

---

## 文档贡献 / Documentation

### 文档结构 / Documentation Structure

```
docs/
├── README.md              # 项目概述
├── DEPLOYMENT.md          # 部署指南
├── API.md                 # API 文档
├── CONTRIBUTING.md        # 贡献指南
├── curriculum.md          # 课程大纲
└── ...
```

### 文档规范 / Documentation Standards

1. **双语支持 / Bilingual Support**
   - 标题使用中英双语
   - 代码示例保持英文
   - 注释可中英混合

2. **格式要求 / Formatting**
   - 使用 Markdown
   - 代码块指定语言
   - 使用表格对比

3. **更新文档 / Updating Docs**
   - 代码变更时同步更新文档
   - 添加变更日志
   - 标注最后更新时间

---

## 常见问题 / FAQ

### Q: 我应该如何开始？

A: 从简单的任务开始！查看 GitHub Issues 中标记为 `good first issue` 的问题。

### Q: 我的 PR 多久会被 review？

A: 通常在 1-3 个工作日内。如果超过一周没有回复，可以礼貌地 @ 维护者。

### Q: 我可以提交多个功能的 PR 吗？

A: 建议一个 PR 只做一件事。多个功能请分开提交。

### Q: 如何测试我的代码？

A: 参考 [测试指南](#测试指南--testing-guide)，确保单元测试和集成测试都通过。

### Q: 代码风格有问题怎么办？

A: CI 会自动检查代码风格。如果失败，根据错误信息修复即可。

### Q: 我可以在文档中使用中文吗？

A: 可以！我们鼓励中英双语文档。

### Q: 如何成为维护者？

A: 持续贡献高质量的代码和 PR，积极参与社区讨论。

---

## 开发工具推荐 / Recommended Tools

### 编辑器 / Editors

- **VS Code** (推荐)
  - Go 扩展
  - ESLint
  - Prettier
  - GitLens

- **GoLand**
  - 完整的 Go IDE

### 浏览器扩展 / Browser Extensions

- React Developer Tools
- Redux DevTools

### CLI 工具 / CLI Tools

```bash
# Go
go install golang.org/x/tools/cmd/goimports@latest

# Frontend
npm install -g prettier eslint
```

---

## 认可与奖励 / Recognition

贡献者将获得：

- GitHub 贡献者列表展示
- 项目 README 中的感谢
- 社区认可
- 优秀贡献者奖励（TBD）

---

## 联系方式 / Contact

- GitHub Issues: [项目 Issues](https://github.com/exltial/ai-learning-path/issues)
- 邮箱：exltial@163.com
- 讨论区：项目 Discussions

---

**感谢你的贡献！🎉**

**Thank you for contributing! 🎉**

---

*最后更新 / Last Updated: 2026-03-10*
