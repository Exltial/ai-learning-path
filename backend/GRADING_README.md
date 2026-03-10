# 作业自动批改系统 (Auto-Grading System)

## 概述

作业自动批改系统为 AI 学习平台提供全面的自动评分功能，支持多种题型，包括选择题、判断题、填空题、编程题和问答题。

## 功能特性

### 1. 选择题自动判分
- ✅ 支持单选和多选
- ✅ 部分得分支持（选对部分正确选项）
- ✅ 错误选项扣分机制
- ✅ 详细反馈（正确/错误/缺失的选项）

### 2. 判断题自动判分
- ✅ 支持多种答案格式（true/false, T/F, 1/0, Yes/No）
- ✅ 大小写不敏感
- ✅ 即时反馈正确答案

### 3. 填空题自动判分
- ✅ 支持多个正确答案
- ✅ 模糊匹配（Levenshtein 距离算法）
- ✅ 关键词部分得分
- ✅ 大小写和标点符号容错

### 4. 编程题自动判分
- ✅ 单元测试执行
- ✅ 输出对比
- ✅ 多测试用例支持
- ✅ 执行时间跟踪
- ✅ 详细测试报告（通过的/失败的测试）
- ✅ 通过率计算

### 5. 问答题评分框架
- ✅ 人工评分支持
- ✅ 基础结构验证（段落数量、字数）
- ✅ 评分历史记录
- ✅ 评分反馈

## 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Layer                              │
│                   (grading_handler.go)                       │
│  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │
│  │ Auto-Grade  │ Manual-Grade│  History    │   Analytics │  │
│  │  Endpoint   │  Endpoint   │  Endpoint   │  Endpoint   │  │
│  └─────────────┴─────────────┴─────────────┴─────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                   Business Logic Layer                       │
│                  (grading_service.go)                        │
│  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │
│  │ Multiple    │  True/      │   Fill      │   Coding    │  │
│  │  Choice     │  False      │   Blank     │   Grading   │  │
│  └─────────────┴─────────────┴─────────────┴─────────────┘  │
│  ┌─────────────┬─────────────┬─────────────────────────────┐ │
│  │   Essay     │  Partial    │   Grading                   │ │
│  │  Grading    │  Credit     │   History                   │ │
│  └─────────────┴─────────────┴─────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                    Data Access Layer                         │
│                (grading_repository.go)                       │
│  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │
│  │   Create    │    Get      │   Update    │   Delete    │  │
│  │  History    │   History   │   Stats     │   History   │  │
│  └─────────────┴─────────────┴─────────────┴─────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      Database                                │
│              (grading_history table)                         │
└─────────────────────────────────────────────────────────────┘
```

## API 端点

### 1. 自动评分
```http
POST /api/v1/submissions/{submission_id}/auto-grade
Content-Type: application/json
Authorization: Bearer <token>

{
  "force": false  // 强制重新评分
}
```

**响应示例:**
```json
{
  "success": true,
  "message": "Submission graded successfully",
  "data": {
    "submission_id": "uuid",
    "is_correct": true,
    "score": 85.5,
    "max_score": 100,
    "percentage": 85.5,
    "feedback": "Correct! All test cases passed.",
    "detailed_feedback": {
      "correct_parts": ["Option A", "Option C"],
      "incorrect_parts": [],
      "missing_parts": []
    },
    "pass_rate": 100,
    "test_results": [
      {
        "test_name": "Test Case 1",
        "passed": true,
        "expected": "42",
        "actual": "42",
        "execution_time_ms": 10
      }
    ]
  }
}
```

### 2. 人工评分
```http
POST /api/v1/submissions/{submission_id}/manual-grade
Content-Type: application/json
Authorization: Bearer <token>

{
  "score": 85.5,
  "feedback": "Good work! Consider improving the structure.",
  "is_correct": true,
  "reason": "Re-grade after student appeal"
}
```

### 3. 获取评分历史
```http
GET /api/v1/submissions/{submission_id}/grading-history
Authorization: Bearer <token>
```

### 4. 获取我的评分历史
```http
GET /api/v1/grading/my-history?limit=20
Authorization: Bearer <token>
```

### 5. 获取练习评分统计
```http
GET /api/v1/exercises/{exercise_id}/grading-stats
Authorization: Bearer <token>
```

### 6. 批量评分
```http
POST /api/v1/grading/batch-grade
Content-Type: application/json
Authorization: Bearer <token>

{
  "submission_ids": ["uuid1", "uuid2"],
  "score": 80,
  "feedback": "Good effort"
}
```

### 7. 重新评分
```http
POST /api/v1/submissions/{submission_id}/regrade
Authorization: Bearer <token>
```

### 8. 评分分析
```http
GET /api/v1/grading/analytics?days=30&exercise_id=uuid
Authorization: Bearer <token>
```

## 评分算法

### 选择题
```
Score = MaxScore × (CorrectRatio - Penalty)
- CorrectRatio = CorrectSelections / TotalCorrectOptions
- Penalty = IncorrectSelections × 0.25
- Minimum score: 0
```

### 填空题
```
1. 精确匹配：100% 得分
2. 模糊匹配（编辑距离 ≤ 2）：90% 得分
3. 关键词匹配（≥50% 关键词）：最高 50% 部分得分
```

### 编程题
```
Score = MaxScore × PassRate
- PassRate = PassedTests / TotalTests
- 通过率 ≥ 80% 视为正确
```

### 问答题
```
需要人工评分，系统仅进行基础验证：
- 最小字数检查
- 段落结构检查
- 提交确认
```

## 数据库表结构

### grading_history

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | UUID | 主键 |
| submission_id | UUID | 提交 ID（外键） |
| exercise_id | UUID | 练习 ID（外键） |
| user_id | UUID | 用户 ID（外键） |
| grading_type | VARCHAR(20) | 评分类型（auto/manual/semi_auto） |
| previous_score | DECIMAL(5,2) | 之前的分数 |
| new_score | DECIMAL(5,2) | 新的分数 |
| score_change | DECIMAL(5,2) | 分数变化 |
| reason | TEXT | 评分原因 |
| graded_by | UUID | 评分人 ID（人工评分时） |
| graded_at | TIMESTAMP | 评分时间 |
| metadata | JSONB | 额外元数据 |

## 使用示例

### Go 代码示例

```go
// 初始化服务
gradingRepo := repository.NewGradingRepository(db)
submissionRepo := repository.NewSubmissionRepository(db)
exerciseRepo := repository.NewExerciseRepository(db)

gradingService := services.NewGradingService(gradingRepo, submissionRepo, exerciseRepo)
gradingHandler := handlers.NewGradingHandler(gradingService, submissionService, exerciseRepo)

// 自动评分
submission, _ := submissionRepo.GetByID(ctx, submissionID)
exercise, _ := exerciseRepo.GetByID(ctx, exerciseID)

result, err := gradingService.GradeSubmission(ctx, submission, exercise)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Score: %.2f/%.2f (%.2f%%)\n", 
    result.Score, result.MaxScore, result.Percentage)
fmt.Printf("Feedback: %s\n", result.Feedback)
```

## 配置

### 部分得分配置

在 `grading_service.go` 中可以调整以下参数：

```go
// 选择题错误选项扣分比例
penalty := float64(incorrectCount) * 0.25  // 25%

// 填空题模糊匹配容错距离
if distance <= 2 {  // 允许 2 个字符差异
    return true
}

// 编程题通过阈值
if passRate >= 0.8 {  // 80% 通过率视为正确
    isCorrect = true
}

// 问答题及格线
if score >= maxScore*0.6 {  // 60% 及格
    isCorrect = true
}
```

## 扩展性

### 添加新的题型

1. 在 `models.Exercise` 中添加新的 `ExerciseType`
2. 在 `grading_service.go` 中实现新的评分方法
3. 在 `GradeSubmission` 的 switch 语句中添加新的 case
4. 更新验证逻辑 `validateSubmission`

### 自定义评分规则

可以通过继承 `GradingService` 或注入自定义评分器来实现特定课程的评分规则。

## 测试

运行测试：
```bash
cd backend
go test ./internal/services/grading_service_test.go -v
go test ./internal/handlers/grading_handler_test.go -v
```

## 注意事项

1. **编程题执行安全**: 生产环境中必须使用沙箱（Docker 容器）执行代码
2. **评分公平性**: 所有自动评分规则应提前向学生公开
3. **申诉机制**: 提供重新评分和人工复审的渠道
4. **数据隐私**: 评分历史包含学生表现数据，需妥善保管
5. **性能优化**: 批量评分时注意数据库连接池和事务管理

## 未来改进

- [ ] AI 辅助问答题评分
- [ ] 代码相似性检测
- [ ] 评分规则可视化配置
- [ ] 实时评分进度追踪
- [ ] 评分质量分析（评分者间一致性）
- [ ] 自适应难度调整建议

## 贡献者

- AI 学习平台开发团队
- 作业自动批改工程师

## 许可证

与主项目许可证相同
