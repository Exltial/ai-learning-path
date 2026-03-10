#!/bin/bash

# AI 学习平台 - 路由修复脚本
# 用于修复 Phase 4 中缺失的 API 路由配置

set -e

echo "======================================"
echo "AI 学习平台 - 路由修复脚本"
echo "======================================"
echo ""

BACKEND_DIR="/home/admin/.openclaw/workspace/projects/ai-learning-platform/backend"
MAIN_GO="$BACKEND_DIR/cmd/main.go"

# 备份原文件
echo "1. 备份 main.go..."
cp "$MAIN_GO" "$MAIN_GO.backup"
echo "   ✓ 备份完成：$MAIN_GO.backup"

# 检查是否已安装 swag
echo ""
echo "2. 检查 swag 工具..."
if ! command -v swag &> /dev/null; then
    echo "   ⚠ swag 未安装，正在安装..."
    go install github.com/swaggo/swag/cmd/swag@latest
    echo "   ✓ swag 安装完成"
else
    echo "   ✓ swag 已安装"
fi

# 更新 main.go 中的路由配置
echo ""
echo "3. 更新 API 路由配置..."

# 读取文件内容
CONTENT=$(cat "$MAIN_GO")

# 检查是否已包含讨论路由
if grep -q "discussionHandler" "$MAIN_GO"; then
    echo "   ⚠ 讨论路由已存在，跳过"
else
    echo "   - 添加讨论处理器初始化..."
    # 在 progressHandler 初始化后添加讨论处理器
    sed -i '/progressHandler := handlers.NewProgressHandler/a\
\
\t// Initialize discussion handler\
\tdiscussionRepo := repository.NewDiscussionRepository(dbPool)\
\tenrollmentRepoForDiscussion := repository.NewEnrollmentRepository(dbPool)\
\tdiscussionService := services.NewDiscussionService(discussionRepo, enrollmentRepoForDiscussion)\
\tdiscussionHandler := handlers.NewDiscussionHandler(discussionService)' "$MAIN_GO"
    echo "   ✓ 讨论处理器初始化已添加"
fi

# 检查是否已包含成就处理器
if grep -q "achievementHandler" "$MAIN_GO"; then
    echo "   ⚠ 成就处理器已存在，跳过"
else
    echo "   - 添加成就处理器初始化..."
    # 在 discussionHandler 初始化后添加成就处理器
    sed -i '/discussionHandler := handlers.NewDiscussionHandler/a\
\
\t// Initialize achievement handler\
\tachievementRepo := repository.NewAchievementRepository(dbPool)\
\tuserAchievementRepo := repository.NewUserAchievementRepository(dbPool)\
\tpointsTransactionRepo := repository.NewPointsTransactionRepository(dbPool)\
\tleaderboardRepo := repository.NewLeaderboardRepository(dbPool)\
\tuserLevelRepo := repository.NewUserLevelRepository(dbPool)\
\tstreakRepo := repository.NewStreakRepository(dbPool)\
\tachievementService := services.NewAchievementService(achievementRepo, userAchievementRepo, userLevelRepo, pointsTransactionRepo, streakRepo, leaderboardRepo, progressRepo, enrollmentRepo, exerciseRepo, submissionRepo, userRepo)\
\tachievementHandler := handlers.NewAchievementHandler(achievementService)' "$MAIN_GO"
    echo "   ✓ 成就处理器初始化已添加"
fi

# 添加讨论路由
echo ""
echo "4. 添加讨论系统路由..."
if grep -q 'discussions.GET("", discussionHandler.GetDiscussions)' "$MAIN_GO"; then
    echo "   ⚠ 讨论路由已存在，跳过"
else
    # 替换讨论路由占位符
    sed -i 's|_ = protected.Group("/discussions") // TODO: Implement discussion handlers|discussions := protected.Group("/discussions")\
\t\t{\
\t\t\tdiscussions.GET("", discussionHandler.ListDiscussions)\
\t\t\tdiscussions.GET("/hot", discussionHandler.GetHotDiscussions)\
\t\t\tdiscussions.GET("/tags", discussionHandler.GetTags)\
\t\t\tdiscussions.GET("/favorites", discussionHandler.GetUserFavorites)\
\t\t\tdiscussions.POST("", discussionHandler.CreateDiscussion)\
\t\t\tdiscussions.GET("/:discussion_id", discussionHandler.GetDiscussion)\
\t\t\tdiscussions.PUT("/:discussion_id", discussionHandler.UpdateDiscussion)\
\t\t\tdiscussions.DELETE("/:discussion_id", discussionHandler.DeleteDiscussion)\
\t\t\tdiscussions.POST("/:discussion_id/like", discussionHandler.ToggleLike)\
\t\t\tdiscussions.POST("/:discussion_id/favorite", discussionHandler.ToggleFavorite)\
\t\t\tdiscussions.POST("/:discussion_id/resolve", discussionHandler.ResolveDiscussion)\
\t\t}|g' "$MAIN_GO"
    echo "   ✓ 讨论路由已添加"
fi

# 添加成就和排行榜路由
echo ""
echo "5. 添加成就和排行榜路由..."
if grep -q 'achievements.GET("", achievementHandler.ListAchievements)' "$MAIN_GO"; then
    echo "   ⚠ 成就路由已存在，跳过"
else
    # 在通知路由前添加成就路由
    sed -i '/\/\/ Notification routes/i\
\t\t\t// Achievement and Leaderboard routes\
\t\t\tachievements := protected.Group("/achievements")\
\t\t\t{\
\t\t\t\tachievements.GET("", achievementHandler.ListAchievements)\
\t\t\t\tachievements.GET("/user", achievementHandler.GetUserAchievements)\
\t\t\t\tachievements.POST("/check", achievementHandler.TriggerAchievementCheck)\
\t\t\t}\
\t\t\t\
\t\t\t// Leaderboard routes\
\t\t\tleaderboard := protected.Group("/leaderboard")\
\t\t\t{\
\t\t\t\tleaderboard.GET("", achievementHandler.GetLeaderboard)\
\t\t\t}\
\t\t\t' "$MAIN_GO"
    echo "   ✓ 成就和排行榜路由已添加"
fi

# 添加进度快捷路由
echo ""
echo "6. 添加进度追踪快捷路由..."
if grep -q 'progress.GET("", progressHandler.GetUserProgress)' "$MAIN_GO"; then
    echo "   ⚠ 进度路由已存在，跳过"
else
    # 在通知路由前添加进度路由
    sed -i '/\/\/ Notification routes/i\
\t\t\t// Progress routes (shortcut)\
\t\t\tprogress := protected.Group("/progress")\
\t\t\t{\
\t\t\t\tprogress.GET("", progressHandler.GetUserProgress)\
\t\t\t}\
\t\t\t' "$MAIN_GO"
    echo "   ✓ 进度路由已添加"
fi

# 重新生成 Swagger 文档
echo ""
echo "7. 重新生成 Swagger 文档..."
cd "$BACKEND_DIR"
swag init --parseDependency --parseInternal
echo "   ✓ Swagger 文档已更新"

# 编译检查
echo ""
echo "8. 编译检查..."
if go build -o bin/server_test cmd/main.go 2>&1; then
    echo "   ✓ 编译成功"
    rm -f bin/server_test
else
    echo "   ✗ 编译失败，请检查错误信息"
    echo ""
    echo "正在恢复备份..."
    cp "$MAIN_GO.backup" "$MAIN_GO"
    echo "✓ 已恢复到原始状态"
    exit 1
fi

# 清理
echo ""
echo "9. 清理..."
rm -f "$MAIN_GO.backup"
echo "   ✓ 备份已删除"

echo ""
echo "======================================"
echo "✓ 路由修复完成！"
echo "======================================"
echo ""
echo "下一步："
echo "1. 重启后端服务：cd backend && ./bin/server"
echo "2. 运行测试：bash scripts/test-all.sh"
echo "3. 检查 Swagger 文档：http://localhost:8080/swagger"
echo ""
