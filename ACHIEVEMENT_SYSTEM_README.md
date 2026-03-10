# 成就系统实现文档

## 概述

成就系统为 AI 学习平台提供完整的 gamification（游戏化）功能，包括成就徽章、等级系统、积分系统、学习 streak 追踪和排行榜。

## 文件结构

### 后端文件

```
backend/
├── internal/
│   ├── models/
│   │   └── achievement_models.go          # 数据模型定义
│   ├── services/
│   │   └── achievement_service.go         # 业务逻辑
│   ├── handlers/
│   │   └── achievement_handler.go         # HTTP 处理器
│   └── repository/
│       ├── achievement_repository.go      # 成就数据访问
│       ├── user_achievement_repository.go # 用户成就数据访问
│       ├── user_level_repository.go       # 用户等级数据访问
│       ├── points_transaction_repository.go # 积分交易数据访问
│       ├── streak_repository.go           # 学习 streak 数据访问
│       └── leaderboard_repository.go      # 排行榜数据访问
└── migrations/
    └── 003_achievement_system.sql         # 数据库迁移
```

### 前端文件

```
frontend/
└── src/
    ├── pages/
    │   └── AchievementsPage.tsx           # 成就展示页面
    └── components/
        └── Leaderboard.tsx                # 排行榜组件
```

## 核心功能

### 1. 成就体系

#### 成就类型
- **general**: 综合成就
- **course**: 课程相关成就
- **streak**: 连续学习成就
- **exercise**: 练习相关成就
- **social**: 社交成就
- **milestone**: 里程碑成就

#### 成就等级（稀有度）
- **bronze** (青铜): 普通成就
- **silver** (白银): 罕见成就
- **gold** (黄金): 稀有成就
- **platinum** (铂金): 史诗成就
- **diamond** (钻石): 传奇成就

#### 默认成就示例
| 成就名称 | 描述 | 积分 | 类型 | 等级 |
|---------|------|------|------|------|
| 初学者 | 完成第一节课 | 10 | course | bronze |
| 持之以恒 | 连续学习 7 天 | 50 | streak | silver |
| 学霸 | 完成所有课程 | 200 | course | gold |
| 刷题达人 | 完成 100 道练习题 | 150 | exercise | gold |
| 坚持不懈 | 连续学习 30 天 | 300 | streak | platinum |
| 知识渊博 | 获得 1000 积分 | 100 | milestone | gold |

### 2. 等级系统

#### 等级称号
1. 初学者
2. 学习者
3. 进阶者
4. 熟练者
5. 专家
6. 高手
7. 大师
8. 宗师
9. 传奇
10. 神话

#### 升级机制
- 通过获得积分提升经验值
- 每级所需经验：`next_level_exp = current_exp * 1.5`
- 升级时自动更新称号

### 3. 积分系统

#### 积分来源
- 解锁成就
- 完成课程
- 完成练习
- 每日奖励
- 特殊活动

#### 积分记录
所有积分变动都会记录在 `points_transactions` 表中，包括：
- 变动金额
- 变动后余额
- 来源类型
- 来源 ID
- 描述
- 时间戳

### 4. 学习 Streak

- 追踪用户连续学习天数
- 每天首次学习时自动更新
- 中断后重新计数
- 记录历史最高 streak

### 5. 排行榜

#### 排行榜类型
- **周榜**: 每周重置
- **月榜**: 每月重置
- **总榜**: 历史累计
- **好友榜**: 仅显示好友

## API 接口

### 成就相关

```
GET /api/achievements
- 获取当前用户的所有成就（含进度）
- 返回：AchievementWithProgress[]

GET /api/achievements/summary
- 获取用户成就摘要
- 返回：UserAchievementSummary

POST /api/achievements/check
- 手动触发成就检查
- Body: { event_type, event_data }
- 返回：解锁的成就列表

POST /api/achievements/activity
- 更新用户活动（用于 streak 追踪）
- 返回：是否解锁新成就
```

### 等级和积分

```
GET /api/achievements/level
- 获取用户等级信息
- 返回：UserLevel

GET /api/achievements/streak
- 获取用户 streak 信息
- 返回：UserStreak

GET /api/achievements/points/history
- 获取积分历史记录
- 返回：UserPointsTransaction[]

POST /api/achievements/points
- 奖励积分（管理员或特定事件）
- Body: { amount, source_type, source_id, description }
```

### 排行榜

```
GET /api/leaderboard?type=all_time&limit=100
- 获取排行榜
- 参数：type (weekly|monthly|all_time|friends), limit
- 返回：LeaderboardEntry[]
```

## 使用示例

### 后端集成

```go
// 在服务层初始化成就服务
achievementService := services.NewAchievementService(
    achievementRepo,
    userAchievementRepo,
    userLevelRepo,
    pointsTransactionRepo,
    streakRepo,
    leaderboardRepo,
    progressRepo,
    enrollmentRepo,
    exerciseRepo,
    submissionRepo,
    userRepo,
)

// 初始化默认成就
err := achievementService.InitializeDefaultAchievements(ctx)

// 用户完成课程时检查成就
eventData := map[string]interface{}{
    "course_id": courseID,
}
unlocked, err := achievementService.CheckAndUnlockAchievements(
    ctx, 
    userID, 
    "course_complete", 
    eventData,
)

// 更新用户活动（streak）
err = achievementService.UpdateStreak(ctx, userID)

// 奖励积分
err = achievementService.AwardPoints(
    ctx, 
    userID, 
    50, 
    "course_complete", 
    &courseID, 
    "完成课程：Go 语言入门",
)
```

### 前端集成

```tsx
// 在路由中添加成就页面
<Route path="/achievements" element={<AchievementsPage />} />

// 在页面中使用排行榜组件
<Leaderboard type="all_time" limit={10} showTitle={true} />

// 获取成就数据
const response = await fetch('/api/achievements', {
  headers: {
    'Authorization': `Bearer ${token}`,
  },
})
const data = await response.json()
```

## 数据库表结构

### 主要表

1. **achievements**: 成就定义
2. **user_achievements**: 用户已解锁成就
3. **user_levels**: 用户等级和积分
4. **points_transactions**: 积分交易记录
5. **user_streaks**: 学习 streak
6. **level_titles**: 等级称号
7. **daily_challenges**: 每日挑战
8. **user_daily_challenge_progress**: 每日挑战进度

### 运行迁移

```bash
# PostgreSQL
psql -d ai_learning_platform -f backend/migrations/003_achievement_system.sql
```

## 扩展建议

### 短期扩展
1. 添加成就通知系统（WebSocket/邮件/推送）
2. 实现每日挑战系统
3. 添加成就分享功能（社交媒体）
4. 实现成就进度实时追踪

### 长期扩展
1. 添加成就赛季系统
2. 实现公会/团队成就
3. 添加成就交易系统
4. 实现特殊成就（隐藏成就）
5. 添加成就展示墙（个人主页）

## 注意事项

1. **性能优化**: 排行榜查询已添加索引，但在用户量大时考虑使用缓存（Redis）
2. **并发控制**: 积分更新需要考虑并发情况，建议使用数据库事务
3. **防作弊**: 成就解锁应验证前置条件，防止刷分
4. **数据安全**: 积分变动需要审计日志
5. **国际化**: 成就名称和描述支持多语言

## 测试建议

1. 单元测试：测试成就解锁逻辑
2. 集成测试：测试 API 端点
3. 压力测试：测试排行榜性能
4. 边界测试：测试 streak 计算（时区、闰年等）

## 维护

- 定期检查成就解锁率，调整难度
- 根据用户反馈添加新成就
- 监控积分通胀，调整获取速度
- 定期清理异常数据

---

**版本**: 1.0.0  
**创建日期**: 2026-03-10  
**作者**: AI Achievement Engineer
