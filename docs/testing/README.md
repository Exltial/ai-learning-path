# 测试文档目录

本目录包含 AI 学习之路项目的完整测试文档和测试用例。

## 📁 文件结构

```
docs/testing/
├── README.md                 # 本文件
├── test-report.md           # 代码审查与测试报告
├── test-cases.md            # 详细测试用例文档
├── TESTING_GUIDE.md         # 测试运行指南
└── ... (更多测试文档)
```

## 📄 文档说明

### 1. test-report.md - 代码审查与测试报告

包含：
- 代码审查发现的 15 个问题（4 个严重、7 个中等、4 个轻微）
- 每个问题的详细描述、影响和修复建议
- 测试覆盖率目标
- 修复优先级和行动计划

**适用场景：**
- 了解项目当前代码质量状况
- 确定需要优先修复的问题
- 制定开发计划

### 2. test-cases.md - 测试用例文档

包含：
- 认证模块测试用例（注册、登录、Token 刷新）
- 用户模块测试用例（获取信息、更新、修改密码）
- 课程模块测试用例（CRUD、enrollment、评论）
- 学习进度模块测试用例
- 练习与提交模块测试用例
- 前端组件测试用例

**适用场景：**
- 编写新测试时参考
- 确保测试覆盖所有功能点
- QA 团队执行手动测试

### 3. TESTING_GUIDE.md - 测试运行指南

包含：
- 后端测试运行方法（Go）
- 前端测试运行方法（React）
- 测试覆盖率生成和查看
- CI/CD 集成示例
- 常见问题解答

**适用场景：**
- 新成员上手测试
- 本地运行测试
- 配置 CI/CD 流水线

## 🧪 测试文件位置

### 后端测试
```
backend/
└── internal/
    ├── handlers/
    │   ├── auth_handler_test.go
    │   ├── user_handler_test.go
    │   └── course_handler_test.go
    └── services/
        └── auth_service_test.go
```

### 前端测试
```
frontend/
├── vitest.config.ts         # Vitest 配置
└── src/
    └── __tests__/
        ├── setup.ts         # 测试配置文件
        ├── CodeEditor.test.tsx
        ├── CourseCard.test.tsx
        └── CoursesPage.test.tsx
```

## 🚀 快速开始

### 运行后端测试

```bash
cd backend

# 运行所有测试
go test ./...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 运行前端测试

```bash
cd frontend

# 安装测试依赖
npm install

# 运行所有测试
npm run test:run

# 运行测试并生成覆盖率
npm run test:coverage
```

## 📊 测试覆盖率目标

| 模块 | 当前 | 目标 | 状态 |
|------|------|------|------|
| Auth Handler | ~20% | 90% | 🟡 进行中 |
| User Handler | ~20% | 85% | 🟡 进行中 |
| Course Handler | ~20% | 85% | 🟡 进行中 |
| Auth Service | ~80% | 95% | 🟢 良好 |
| 前端组件 | ~60% | 75% | 🟡 进行中 |

## 📝 测试编写规范

### 后端 (Go)

```go
func TestFunction_Scenario_ExpectedResult(t *testing.T) {
    // 1. 准备测试数据
    // 2. 执行测试
    // 3. 断言结果
}
```

### 前端 (React)

```tsx
describe('ComponentName', () => {
  describe('渲染', () => {
    it('应该正确渲染', () => {
      // 测试代码
    })
  })
  
  describe('交互', () => {
    it('应该响应点击事件', () => {
      // 测试代码
    })
  })
})
```

## 🔧 测试工具

### 后端
- **testing** - Go 标准测试库
- **testify** - 断言和 mock 库
- **testcontainers-go** - 集成测试容器

### 前端
- **Vitest** - 快速单元测试框架
- **Testing Library** - React 组件测试
- **jsdom** - 浏览器环境模拟

## 📚 相关资源

- [Go Testing 最佳实践](https://go.dev/doc/tutorial/add-a-test)
- [Vitest 文档](https://vitest.dev/)
- [Testing Library 文档](https://testing-library.com/)

## 🤝 贡献指南

添加新测试时：

1. 参考 `test-cases.md` 中的测试用例
2. 遵循现有测试文件的命名和结构
3. 确保测试覆盖率不低于目标值
4. 更新本文档中的覆盖率表格

## 📅 更新日志

- **2026-03-10** - 初始版本，包含代码审查报告和基础测试用例

---

**维护者：** AI 测试工程师  
**最后更新：** 2026-03-10
