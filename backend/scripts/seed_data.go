// ============================================================
// AI 学习平台 - 数据种子脚本
// 一键导入所有测试数据
// ============================================================
// 使用方法:
//   go run scripts/seed_data.go
//
// 环境变量:
//   DATABASE_URL - 数据库连接字符串 (默认：postgres://localhost:5432/ai_learning?sslmode=disable)
// ============================================================

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/google/uuid"
)

// 数据库配置
var (
	dbURL = getEnv("DATABASE_URL", "postgres://localhost:5432/ai_learning?sslmode=disable")
	db    *sql.DB
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// 初始化数据库连接
func initDB() error {
	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("打开数据库连接失败：%w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("数据库连接失败：%w", err)
	}

	fmt.Println("✓ 数据库连接成功")
	return nil
}

// 执行 SQL 文件
func executeSQLFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败 [%s]: %w", filePath, err)
	}

	// 分割 SQL 语句（按分号分割）
	statements := strings.Split(string(content), ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		_, err := db.Exec(stmt)
		if err != nil {
			// 忽略已存在数据的错误
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			return fmt.Errorf("执行 SQL 失败 [%s]: %w", filePath, err)
		}
	}

	return nil
}

// 执行 SQL 文件（支持多语句）
func executeSQLFileBatch(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败 [%s]: %w", filePath, err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		// 忽略已存在数据的错误
		if strings.Contains(err.Error(), "duplicate key") {
			fmt.Printf("⚠ 文件 [%s] 部分数据已存在\n", filePath)
			return nil
		}
		return fmt.Errorf("执行 SQL 失败 [%s]: %w", filePath, err)
	}

	return nil
}

// 创建额外成就数据
func seedAchievements() error {
	fmt.Println("\n📍 创建成就数据...")

	achievements := []struct {
		name        string
		description string
		points      int
		achType     string
		criteria    string
	}{
		{"初学者", "完成第一个课程", 100, "course", `{"type": "complete_course", "count": 1}`},
		{"学习达人", "完成 5 个课程", 500, "course", `{"type": "complete_course", "count": 5}`},
		{"学霸", "完成 10 个课程", 1000, "course", `{"type": "complete_course", "count": 10}`},
		{"持之以恒", "连续学习 7 天", 200, "streak", `{"type": "learning_streak", "days": 7}`},
		{"坚持不懈", "连续学习 30 天", 1000, "streak", `{"type": "learning_streak", "days": 30}`},
		{"练习之王", "完成 100 道练习题", 500, "exercise", `{"type": "complete_exercise", "count": 100}`},
		{"完美主义者", "练习题全部正确", 300, "exercise", `{"type": "perfect_score", "count": 50}`},
		{"社交达人", "发布 10 个讨论帖", 150, "social", `{"type": "create_discussion", "count": 10}`},
		{"助人者", "获得 50 个赞", 400, "social", `{"type": "receive_upvotes", "count": 50}`},
		{"早起鸟", "连续 7 天早上学习", 250, "streak", `{"type": "morning_streak", "days": 7}`},
		{"夜猫子", "连续 7 天晚上学习", 250, "streak", `{"type": "night_streak", "days": 7}`},
		{"全栈开发者", "完成所有编程课程", 2000, "milestone", `{"type": "complete_all_programming", "count": 1}`},
		{"AI 专家", "完成所有 AI 课程", 2000, "milestone", `{"type": "complete_all_ai", "count": 1}`},
		{"终身学习", "学习时长达到 100 小时", 3000, "milestone", `{"type": "learning_hours", "hours": 100}`},
		{"平台贡献者", "帮助 100 位同学", 1500, "social", `{"type": "help_others", "count": 100}`},
	}

	for _, ach := range achievements {
		id := uuid.New()
		_, err := db.Exec(`
			INSERT INTO achievements (id, name, description, points, achievement_type, criteria, created_at)
			VALUES ($1, $2, $3, $4, $5, $6::jsonb, NOW())
			ON CONFLICT (id) DO NOTHING
		`, id, ach.name, ach.description, ach.points, achType(ach.achType), ach.criteria)
		
		if err != nil && !strings.Contains(err.Error(), "duplicate") {
			return fmt.Errorf("插入成就失败 [%s]: %w", ach.name, err)
		}
	}

	fmt.Printf("✓ 创建 %d 个成就\n", len(achievements))
	return nil
}

func achType(t string) string {
	switch t {
	case "course":
		return "course"
	case "streak":
		return "streak"
	case "exercise":
		return "exercise"
	case "social":
		return "social"
	case "milestone":
		return "milestone"
	default:
		return "general"
	}
}

// 创建用户等级标题
func seedLevelTitles() error {
	fmt.Println("\n📍 创建等级标题...")

	titles := []struct {
		level       int
		title       string
		description string
	}{
		{1, "初学者", "刚刚踏上学习之旅"},
		{2, "探索者", "开始探索知识的海洋"},
		{3, "学习者", "持续学习新知识"},
		{4, "进阶者", "技能不断提升"},
		{5, "熟练者", "掌握核心技能"},
		{6, "专家", "在领域内有深入理解"},
		{7, "大师", "精通多个领域"},
		{8, "导师", "能够指导他人学习"},
		{9, "传奇", "创造学习传奇"},
		{10, "终身学习者", "学习永不止息"},
	}

	for _, t := range titles {
		id := uuid.New()
		_, err := db.Exec(`
			INSERT INTO level_titles (id, level, title, description, created_at)
			VALUES ($1, $2, $3, $4, NOW())
			ON CONFLICT (id) DO NOTHING
		`, id, t.level, t.title, t.description)
		
		if err != nil && !strings.Contains(err.Error(), "duplicate") {
			return fmt.Errorf("插入等级标题失败 [%s]: %w", t.title, err)
		}
	}

	fmt.Printf("✓ 创建 %d 个等级标题\n", len(titles))
	return nil
}

// 学员 ID 列表（全局常量）
var studentIDs = []string{
	"20000000-0000-0000-0000-000000000001",
	"20000000-0000-0000-0000-000000000002",
	"20000000-0000-0000-0000-000000000003",
}

// 创建示例注册和进度数据
func seedEnrollmentsAndProgress() error {
	fmt.Println("\n📍 创建课程注册和进度数据...")

	// 课程 ID 列表
	courseIDs := []string{
		"30000000-0000-0000-0000-000000000001",
		"30000000-0000-0000-0000-000000000002",
		"30000000-0000-0000-0000-000000000003",
		"30000000-0000-0000-0000-000000000004",
		"30000000-0000-0000-0000-000000000005",
	}

	// 为每个学员注册部分课程
	enrollments := []struct {
		studentID string
		courseID  string
		status    string
		progress  float64
	}{
		{studentIDs[0], courseIDs[0], "active", 75.5},
		{studentIDs[0], courseIDs[1], "active", 45.0},
		{studentIDs[0], courseIDs[2], "completed", 100.0},
		{studentIDs[1], courseIDs[0], "active", 60.0},
		{studentIDs[1], courseIDs[3], "active", 30.0},
		{studentIDs[1], courseIDs[4], "active", 15.0},
		{studentIDs[2], courseIDs[1], "completed", 100.0},
		{studentIDs[2], courseIDs[2], "active", 80.0},
		{studentIDs[2], courseIDs[4], "active", 25.0},
	}

	for _, e := range enrollments {
		enrollID := uuid.New()
		enrolledAt := time.Now().AddDate(0, 0, -int(e.progress/3))
		
		var completedAt *time.Time
		if e.status == "completed" {
			t := time.Now().AddDate(0, 0, -5)
			completedAt = &t
		}

		_, err := db.Exec(`
			INSERT INTO enrollments (id, user_id, course_id, enrolled_at, completed_at, status, progress_percentage, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
			ON CONFLICT (user_id, course_id) DO UPDATE SET 
				status = EXCLUDED.status,
				progress_percentage = EXCLUDED.progress_percentage,
				updated_at = NOW()
		`, enrollID, e.studentID, e.courseID, enrolledAt, completedAt, e.status, e.progress)
		
		if err != nil {
			return fmt.Errorf("插入注册记录失败：%w", err)
		}
	}

	fmt.Printf("✓ 创建 %d 条课程注册记录\n", len(enrollments))
	return nil
}

// 创建排行榜数据
func seedLeaderboard() error {
	fmt.Println("\n📍 创建排行榜数据...")

	// 为学员创建用户等级和积分
	userLevels := []struct {
		userID        string
		level         int
		currentPoints int
		totalPoints   int
		experience    int
		title         string
	}{
		{studentIDs[0], 5, 2500, 5800, 5800, "熟练者"},
		{studentIDs[1], 4, 1800, 4200, 4200, "进阶者"},
		{studentIDs[2], 6, 3200, 7500, 7500, "专家"},
	}

	for _, ul := range userLevels {
		id := uuid.New()
		_, err := db.Exec(`
			INSERT INTO user_levels (id, user_id, level, current_points, total_points, experience, title, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
			ON CONFLICT (user_id) DO UPDATE SET
				level = EXCLUDED.level,
				current_points = EXCLUDED.current_points,
				total_points = EXCLUDED.total_points,
				experience = EXCLUDED.experience,
				title = EXCLUDED.title,
				updated_at = NOW()
		`, id, ul.userID, ul.level, ul.currentPoints, ul.totalPoints, ul.experience, ul.title)
		
		if err != nil {
			return fmt.Errorf("插入用户等级失败：%w", err)
		}

		// 创建积分交易记录
		transactionID := uuid.New()
		_, err = db.Exec(`
			INSERT INTO user_points_transactions (id, user_id, amount, balance_after, source_type, description, created_at)
			VALUES ($1, $2, $3, $4, 'initial', '初始积分', NOW())
		`, transactionID, ul.userID, ul.currentPoints, ul.currentPoints)
		
		if err != nil {
			return fmt.Errorf("插入积分交易失败：%w", err)
		}
	}

	// 创建学习连续记录
	streaks := []struct {
		userID       string
		currentStreak int
		longestStreak int
	}{
		{studentIDs[0], 15, 30},
		{studentIDs[1], 7, 21},
		{studentIDs[2], 25, 45},
	}

	for _, s := range streaks {
		id := uuid.New()
		_, err := db.Exec(`
			INSERT INTO user_streaks (id, user_id, current_streak, longest_streak, last_activity_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW())
			ON CONFLICT (user_id) DO UPDATE SET
				current_streak = EXCLUDED.current_streak,
				longest_streak = EXCLUDED.longest_streak,
				updated_at = NOW()
		`, id, s.userID, s.currentStreak, s.longestStreak)
		
		if err != nil {
			return fmt.Errorf("插入学习连续记录失败：%w", err)
		}
	}

	fmt.Println("✓ 创建排行榜数据完成")
	return nil
}

// 创建用户成就关联
func seedUserAchievements() error {
	fmt.Println("\n📍 创建用户成就关联...")

	// 为学员分配一些成就
	userAchievements := []struct {
		userID        string
		achievementID int // 索引，从 0 开始
	}{
		{studentIDs[0], 0}, // 初学者
		{studentIDs[0], 3}, // 持之以恒
		{studentIDs[0], 5}, // 练习之王
		{studentIDs[1], 0}, // 初学者
		{studentIDs[1], 2}, // 学霸
		{studentIDs[2], 0}, // 初学者
		{studentIDs[2], 1}, // 学习达人
		{studentIDs[2], 3}, // 持之以恒
		{studentIDs[2], 4}, // 坚持不懈
		{studentIDs[2], 11}, // 全栈开发者
	}

	for _, ua := range userAchievements {
		// 查询成就 ID
		var achievementID string
		err := db.QueryRow(`
			SELECT id FROM achievements ORDER BY created_at LIMIT 1 OFFSET $1
		`, ua.achievementID).Scan(&achievementID)
		
		if err != nil {
			continue // 跳过不存在的成就
		}

		id := uuid.New()
		_, err = db.Exec(`
			INSERT INTO user_achievements (id, user_id, achievement_id, earned_at)
			VALUES ($1, $2, $3, NOW())
			ON CONFLICT (user_id, achievement_id) DO NOTHING
		`, id, ua.userID, achievementID)
		
		if err != nil && !strings.Contains(err.Error(), "duplicate") {
			return fmt.Errorf("插入用户成就失败：%w", err)
		}
	}

	fmt.Printf("✓ 创建 %d 条用户成就记录\n", len(userAchievements))
	return nil
}

// 主函数
func main() {
	fmt.Println("============================================")
	fmt.Println("  AI 学习平台 - 数据种子脚本")
	fmt.Println("============================================")
	fmt.Println()

	// 初始化数据库
	if err := initDB(); err != nil {
		log.Fatalf("❌ 数据库初始化失败：%v", err)
	}
	defer db.Close()

	// 获取脚本所在目录
	scriptDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		scriptDir = "."
	}

	// SQL 文件列表
	sqlFiles := []string{
		"seed_users.sql",
		"seed_courses.sql",
		"seed_exercises.sql",
		"seed_discussions.sql",
	}

	// 执行 SQL 文件
	for _, file := range sqlFiles {
		filePath := filepath.Join(scriptDir, file)
		fmt.Printf("\n📍 执行 %s...\n", file)
		
		if err := executeSQLFileBatch(filePath); err != nil {
			log.Printf("⚠ 警告：%v", err)
			continue
		}
		
		fmt.Printf("✓ %s 执行完成\n", file)
	}

	// 创建额外数据
	if err := seedAchievements(); err != nil {
		log.Printf("⚠ 成就数据创建失败：%v", err)
	}

	if err := seedLevelTitles(); err != nil {
		log.Printf("⚠ 等级标题创建失败：%v", err)
	}

	if err := seedEnrollmentsAndProgress(); err != nil {
		log.Printf("⚠ 注册和进度数据创建失败：%v", err)
	}

	if err := seedLeaderboard(); err != nil {
		log.Printf("⚠ 排行榜数据创建失败：%v", err)
	}

	if err := seedUserAchievements(); err != nil {
		log.Printf("⚠ 用户成就数据创建失败：%v", err)
	}

	fmt.Println()
	fmt.Println("============================================")
	fmt.Println("  ✅ 数据种子导入完成！")
	fmt.Println("============================================")
	fmt.Println()
	fmt.Println("创建的数据包括:")
	fmt.Println("  - 9 个用户账户 (3 管理员 + 3 讲师 + 3 学员)")
	fmt.Println("  - 5 个课程 (共 50 个章节)")
	fmt.Println("  - 60+ 道练习题 (5 种题型)")
	fmt.Println("  - 30+ 个讨论帖子")
	fmt.Println("  - 15 个成就")
	fmt.Println("  - 排行榜和等级数据")
	fmt.Println("  - 课程注册和学习进度")
	fmt.Println()
	fmt.Println("默认密码：Password123!")
	fmt.Println()
}
