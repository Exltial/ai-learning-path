# 测试运行指南

**项目：** AI 学习之路  
**更新日期：** 2026-03-10

---

## 📋 目录

1. [后端测试 (Go)](#1-后端测试-go)
2. [前端测试 (React)](#2-前端测试-react)
3. [测试覆盖率](#3-测试覆盖率)
4. [CI/CD 集成](#4-cicd-集成)
5. [常见问题](#5-常见问题)

---

## 1. 后端测试 (Go)

### 1.1 前置条件

确保已安装：
- Go 1.21+
- PostgreSQL (用于集成测试)
- Redis (用于缓存测试)

### 1.2 安装依赖

```bash
cd backend

# 下载所有依赖
go mod download

# 安装测试工具
go get github.com/stretchr/testify
go get github.com/DATA-DOG/go-sqlmock
go get github.com/testcontainers/testcontainers-go
```

### 1.3 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/handlers/...
go test ./internal/services/...
go test ./internal/repository/...

# 运行测试并显示详细输出
go test -v ./...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...

# 查看覆盖率报告（HTML 格式）
go tool cover -html=coverage.out

# 运行测试并显示覆盖率
go test -cover ./...

# 只运行匹配模式的测试
go test -run TestAuth ./...
go test -run TestRegister ./internal/handlers/...
```

### 1.4 测试数据库设置

对于集成测试，建议使用测试数据库：

```bash
# 创建测试数据库
createdb ai_learning_test

# 或者使用 Docker
docker run -d --name postgres-test \
  -e POSTGRES_PASSWORD=testpass \
  -e POSTGRES_DB=ai_learning_test \
  -p 5433:5432 \
  postgres:15

# 设置测试环境变量
export DATABASE_URL="postgres://postgres:testpass@localhost:5433/ai_learning_test?sslmode=disable"
export TEST_MODE=true
```

### 1.5 编写测试

测试文件命名约定：`*_test.go`

```go
package handlers

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
    // 测试代码
    assert.Equal(t, expected, actual)
}
```

---

## 2. 前端测试 (React)

### 2.1 前置条件

确保已安装：
- Node.js 18+
- npm 或 yarn

### 2.2 安装依赖

```bash
cd frontend

# 安装依赖
npm install

# 测试依赖已在 package.json 中配置
```

### 2.3 运行测试

```bash
# 运行所有测试（监视模式）
npm test

# 运行所有测试（单次运行）
npm run test:run

# 运行测试并生成覆盖率报告
npm run test:coverage

# 运行特定测试文件
npm test -- CodeEditor.test.tsx

# 运行匹配模式的测试
npm test -- --testNamePattern="应该正确渲染"

# 运行测试并显示详细输出
npm test -- --reporter=verbose
```

### 2.4 测试文件结构

```
frontend/src/
├── __tests__/
│   ├── setup.ts              # 测试配置文件
│   ├── CodeEditor.test.tsx   # 组件测试
│   ├── CourseCard.test.tsx
│   └── CoursesPage.test.tsx
├── components/
│   ├── CodeEditor.tsx
│   └── CourseCard.tsx
└── pages/
    └── CoursesPage.tsx
```

### 2.5 编写测试

```tsx
import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import MyComponent from './MyComponent'

describe('MyComponent', () => {
  it('应该正确渲染', () => {
    render(<MyComponent />)
    expect(screen.getByText('Hello')).toBeInTheDocument()
  })
})
```

### 2.6 测试配置

`vitest.config.ts`：

```typescript
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './src/__tests__/setup.ts',
    coverage: {
      reporter: ['text', 'json', 'html'],
    },
  },
})
```

---

## 3. 测试覆盖率

### 3.1 查看覆盖率

**后端：**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
# 在浏览器中打开 coverage.html
```

**前端：**
```bash
npm run test:coverage
# 打开 coverage/index.html
```

### 3.2 覆盖率目标

| 模块 | 目标覆盖率 |
|------|-----------|
| Handler 层 | 90% |
| Service 层 | 95% |
| Repository 层 | 70% |
| 前端组件 | 75% |
| 前端工具 | 90% |

### 3.3 生成合并覆盖率报告

```bash
# 后端
go test -coverprofile=backend-coverage.out ./...

# 前端
npm run test:coverage

# 可以使用 codecov 或 coveralls 等工具上传
```

---

## 4. CI/CD 集成

### 4.1 GitHub Actions 示例

`.github/workflows/test.yml`：

```yaml
name: Tests

on: [push, pull_request]

jobs:
  backend-test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: testpass
          POSTGRES_DB: ai_learning_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install dependencies
        run: |
          cd backend
          go mod download
      
      - name: Run tests
        run: |
          cd backend
          go test -v -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./backend/coverage.out
          flags: backend

  frontend-test:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Install dependencies
        run: |
          cd frontend
          npm ci
      
      - name: Run tests
        run: |
          cd frontend
          npm run test:run -- --coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./frontend/coverage/coverage-final.json
          flags: frontend
```

### 4.2 本地预提交钩子

`.git/hooks/pre-commit`：

```bash
#!/bin/bash

# 后端测试
echo "Running backend tests..."
cd backend
go test ./... || exit 1

# 前端测试
echo "Running frontend tests..."
cd ../frontend
npm run test:run || exit 1

echo "All tests passed!"
```

---

## 5. 常见问题

### Q1: 测试失败 "connection refused"

**原因：** 数据库或 Redis 未运行

**解决：**
```bash
# 启动 PostgreSQL
docker run -d --name postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=ai_learning \
  -p 5432:5432 \
  postgres:15

# 启动 Redis
docker run -d --name redis \
  -p 6379:6379 \
  redis:7
```

### Q2: 前端测试 "Cannot find module"

**原因：** 依赖未安装或路径别名配置问题

**解决：**
```bash
cd frontend
npm install
```

确保 `vitest.config.ts` 中有正确的路径别名配置。

### Q3: Go 测试 "import cycle not allowed"

**原因：** 包之间存在循环依赖

**解决：** 重构代码，使用接口解耦依赖。

### Q4: 测试运行很慢

**优化建议：**
- 使用 `-parallel` 参数并行运行测试
- 对慢的集成测试使用 `t.Parallel()`
- 使用测试数据库而非真实数据库
- Mock 外部依赖

### Q5: 如何 Mock HTTP 请求？

**后端：**
```go
import (
    "github.com/jarcoal/httpmock"
)

func TestWithMock(t *testing.T) {
    httpmock.Activate()
    defer httpmock.Deactivate()
    
    httpmock.RegisterResponder("GET", "https://api.example.com",
        httpmock.NewStringResponder(200, `{"data": "test"}`))
    
    // 执行测试
}
```

**前端：**
```tsx
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'

const server = setupServer(
  http.get('/api/courses', () => {
    return HttpResponse.json({ courses: [] })
  })
)

beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
afterAll(() => server.close())
```

---

## 📚 参考资源

- [Go Testing](https://pkg.go.dev/testing)
- [testify](https://github.com/stretchr/testify)
- [Vitest](https://vitest.dev/)
- [Testing Library](https://testing-library.com/)
- [Test Containers](https://www.testcontainers.org/)

---

**最后更新：** 2026-03-10
