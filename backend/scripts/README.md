# 数据种子脚本使用说明

本目录包含 AI 学习平台的测试数据种子脚本，用于快速创建演示和测试环境。

## 📁 文件说明

| 文件 | 说明 |
|------|------|
| `seed_data.go` | Go 主脚本，一键导入所有数据 |
| `seed_users.sql` | 用户账户数据（管理员/讲师/学员） |
| `seed_courses.sql` | 课程和章节数据 |
| `seed_exercises.sql` | 练习题数据 |
| `seed_discussions.sql` | 讨论区数据 |

## 🚀 快速开始

### 方法一：使用 Go 脚本（推荐）

```bash
# 进入后端目录
cd backend

# 运行种子脚本
go run scripts/seed_data.go
```

### 方法二：手动执行 SQL 文件

```bash
# 使用 psql 执行
psql -h localhost -U postgres -d ai_learning -f scripts/seed_users.sql
psql -h localhost -U postgres -d ai_learning -f scripts/seed_courses.sql
psql -h localhost -U postgres -d ai_learning -f scripts/seed_exercises.sql
psql -h localhost -U postgres -d ai_learning -f scripts/seed_discussions.sql
```

### 方法三：使用 Docker

```bash
# 如果使用 Docker 部署
docker-compose exec backend go run scripts/seed_data.go
```

## 📊 创建的数据

### 用户账户（共 9 个）

**管理员账户:**
- `admin_zhang` / `zhang.admin@aiplatform.com`
- `admin_wang` / `wang.admin@aiplatform.com`
- `admin_li` / `li.admin@aiplatform.com`

**讲师账户:**
- `prof_chen` / `chen.professor@aiplatform.com` - 负责 Python 和机器学习课程
- `prof_liu` / `liu.professor@aiplatform.com` - 负责 Web 开发课程
- `prof_huang` / `huang.professor@aiplatform.com` - 负责数据结构和深度学习课程

**学员账户:**
- `student_zhao` / `zhao.student@aiplatform.com`
- `student_wu` / `wu.student@aiplatform.com`
- `student_zhou` / `zhou.student@aiplatform.com`

**默认密码：** `Password123!`（所有账户）

### 课程（共 5 个）

| 课程 | 讲师 | 难度 | 章节数 | 预计学时 |
|------|------|------|--------|----------|
| Python 编程入门 | prof_chen | 初级 | 10 | 40 小时 |
| 机器学习基础 | prof_chen | 中级 | 10 | 60 小时 |
| Web 开发实战 | prof_liu | 初级 | 10 | 50 小时 |
| 数据结构与算法 | prof_huang | 中级 | 10 | 45 小时 |
| 深度学习入门 | prof_huang | 高级 | 10 | 70 小时 |

### 练习题（共 60+ 道）

每门课程包含 10+ 道练习题，涵盖 5 种题型：

1. **multiple_choice** - 选择题
2. **coding** - 编程题
3. **fill_blank** - 填空题
4. **true_false** - 判断题
5. **essay** - 简答题

### 讨论区（共 30+ 个帖子）

每门课程包含 5+ 个讨论主题，包括：
- 学习问题求助
- 概念理解讨论
- 最佳实践分享
- 项目部署问题
- 课程公告

### 成就系统（共 15 个）

- 学习成就（完成课程）
- 连续学习成就（学习天数）
- 练习成就（完成题目）
- 社交成就（讨论互动）
- 里程碑成就（重大突破）

### 排行榜数据

- 用户等级和积分
- 学习连续记录
- 成就解锁记录
- 积分交易记录

## ⚙️ 配置

### 数据库连接

通过环境变量配置数据库连接：

```bash
export DATABASE_URL="postgres://username:password@localhost:5432/ai_learning?sslmode=disable"
```

或在 `.env` 文件中配置：

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/ai_learning?sslmode=disable
```

## 🔄 重新导入数据

如果需要重新导入数据（清空现有数据）：

```bash
# 方法 1：重置数据库
psql -h localhost -U postgres -d ai_learning -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
psql -h localhost -U postgres -d ai_learning -f schema.sql
go run scripts/seed_data.go

# 方法 2：使用迁移工具（如果有）
go run cmd/migrate/main.go reset
go run scripts/seed_data.go
```

## 📝 自定义数据

可以根据需要修改 SQL 文件来自定义数据：

1. 修改 `seed_users.sql` 添加更多用户
2. 修改 `seed_courses.sql` 添加新课程
3. 修改 `seed_exercises.sql` 添加新题目
4. 修改 `seed_discussions.sql` 添加新讨论

修改后重新运行 `go run scripts/seed_data.go` 即可。

## 🐛 故障排除

### 问题：数据库连接失败

```
❌ 数据库初始化失败：dial tcp [::1]:5432: connect: connection refused
```

**解决方案：**
1. 确保 PostgreSQL 服务正在运行
2. 检查数据库连接字符串是否正确
3. 确认数据库 `ai_learning` 已创建

### 问题：数据已存在

```
⚠ 警告：执行 SQL 失败：duplicate key value violates unique constraint
```

**解决方案：**
- 这是正常提示，脚本会跳过已存在的数据
- 如需重新导入，先清空数据（见上方"重新导入数据"）

### 问题：缺少依赖

```
cannot find package "github.com/lib/pq"
```

**解决方案：**
```bash
cd backend
go mod download
go mod tidy
```

## 📞 支持

如有问题，请在项目 Issues 中反馈。
