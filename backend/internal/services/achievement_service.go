package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"

	"github.com/google/uuid"
)

// AchievementService handles achievement-related business logic
type AchievementService struct {
	achievementRepo    *repository.AchievementRepository
	userAchievementRepo *repository.UserAchievementRepository
	userLevelRepo      *repository.UserLevelRepository
	pointsTransactionRepo *repository.PointsTransactionRepository
	streakRepo         *repository.StreakRepository
	leaderboardRepo    *repository.LeaderboardRepository
	progressRepo       *repository.ProgressRepository
	enrollmentRepo     *repository.EnrollmentRepository
	exerciseRepo       *repository.ExerciseRepository
	submissionRepo     *repository.SubmissionRepository
	userRepo           *repository.UserRepository
}

// NewAchievementService creates a new AchievementService
func NewAchievementService(
	achievementRepo *repository.AchievementRepository,
	userAchievementRepo *repository.UserAchievementRepository,
	userLevelRepo *repository.UserLevelRepository,
	pointsTransactionRepo *repository.PointsTransactionRepository,
	streakRepo *repository.StreakRepository,
	leaderboardRepo *repository.LeaderboardRepository,
	progressRepo *repository.ProgressRepository,
	enrollmentRepo *repository.EnrollmentRepository,
	exerciseRepo *repository.ExerciseRepository,
	submissionRepo *repository.SubmissionRepository,
	userRepo *repository.UserRepository,
) *AchievementService {
	return &AchievementService{
		achievementRepo:     achievementRepo,
		userAchievementRepo: userAchievementRepo,
		userLevelRepo:       userLevelRepo,
		pointsTransactionRepo: pointsTransactionRepo,
		streakRepo:          streakRepo,
		leaderboardRepo:     leaderboardRepo,
		progressRepo:        progressRepo,
		enrollmentRepo:      enrollmentRepo,
		exerciseRepo:        exerciseRepo,
		submissionRepo:      submissionRepo,
		userRepo:            userRepo,
	}
}

// Error definitions
var (
	ErrAchievementNotFound  = errors.New("achievement not found")
	ErrAchievementAlreadyUnlocked = errors.New("achievement already unlocked")
	ErrInvalidCriteria      = errors.New("invalid achievement criteria")
)

// CheckAndUnlockAchievements checks if user qualifies for any achievements and unlocks them
func (s *AchievementService) CheckAndUnlockAchievements(ctx context.Context, userID uuid.UUID, eventType string, eventData map[string]interface{}) ([]*models.Achievement, error) {
	unlockedAchievements := []*models.Achievement{}

	// Get all enabled achievements
	achievements, err := s.achievementRepo.GetAllEnabled(ctx)
	if err != nil {
		return nil, err
	}

	// Get user's unlocked achievements to skip them
	unlockedIDs, err := s.userAchievementRepo.GetUserAchievementIDs(ctx, userID)
	if err != nil {
		unlockedIDs = []uuid.UUID{}
	}

	// Check each achievement
	for _, achievement := range achievements {
		// Skip if already unlocked
		if containsUUID(unlockedIDs, achievement.ID) {
			continue
		}

		// Check if criteria is met
		isUnlocked, err := s.checkAchievementCriteria(ctx, userID, achievement, eventType, eventData)
		if err != nil {
			continue
		}

		if isUnlocked {
			// Unlock the achievement
			err := s.UnlockAchievement(ctx, userID, achievement.ID)
			if err != nil {
				continue
			}
			unlockedAchievements = append(unlockedAchievements, achievement)
		}
	}

	return unlockedAchievements, nil
}

// checkAchievementCriteria checks if user meets the achievement criteria
func (s *AchievementService) checkAchievementCriteria(ctx context.Context, userID uuid.UUID, achievement *models.Achievement, eventType string, eventData map[string]interface{}) (bool, error) {
	criteria := achievement.Criteria

	switch criteria.Type {
	case "complete_lesson":
		return s.checkLessonCompletion(ctx, userID, criteria, eventData)
	case "complete_course":
		return s.checkCourseCompletion(ctx, userID, criteria, eventData)
	case "streak_days":
		return s.checkStreakDays(ctx, userID, criteria)
	case "complete_exercises":
		return s.checkExerciseCompletion(ctx, userID, criteria)
	case "total_points":
		return s.checkTotalPoints(ctx, userID, criteria)
	case "daily_activity":
		return s.checkDailyActivity(ctx, userID, criteria, eventData)
	default:
		return false, ErrInvalidCriteria
	}
}

// checkLessonCompletion checks if user has completed required lessons
func (s *AchievementService) checkLessonCompletion(ctx context.Context, userID uuid.UUID, criteria models.AchievementCriteria, eventData map[string]interface{}) (bool, error) {
	progresses, err := s.progressRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	completedCount := 0
	for _, progress := range progresses {
		if progress.IsCompleted {
			if criteria.TargetID == nil || progress.LessonID == *criteria.TargetID {
				completedCount++
			}
		}
	}

	return completedCount >= criteria.Threshold, nil
}

// checkCourseCompletion checks if user has completed required courses
func (s *AchievementService) checkCourseCompletion(ctx context.Context, userID uuid.UUID, criteria models.AchievementCriteria, eventData map[string]interface{}) (bool, error) {
	enrollments, err := s.enrollmentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	completedCount := 0
	for _, enrollment := range enrollments {
		if enrollment.Status == "completed" {
			if criteria.TargetID == nil || enrollment.CourseID == *criteria.TargetID {
				completedCount++
			}
		}
	}

	return completedCount >= criteria.Threshold, nil
}

// checkStreakDays checks if user has achieved required streak
func (s *AchievementService) checkStreakDays(ctx context.Context, userID uuid.UUID, criteria models.AchievementCriteria) (bool, error) {
	streak, err := s.streakRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, nil
	}

	return streak.CurrentStreak >= criteria.Threshold, nil
}

// checkExerciseCompletion checks if user has completed required exercises
func (s *AchievementService) checkExerciseCompletion(ctx context.Context, userID uuid.UUID, criteria models.AchievementCriteria) (bool, error) {
	submissions, err := s.submissionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	correctCount := 0
	for _, submission := range submissions {
		if submission.IsCorrect != nil && *submission.IsCorrect {
			correctCount++
		}
	}

	return correctCount >= criteria.Threshold, nil
}

// checkTotalPoints checks if user has accumulated required points
func (s *AchievementService) checkTotalPoints(ctx context.Context, userID uuid.UUID, criteria models.AchievementCriteria) (bool, error) {
	level, err := s.userLevelRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, nil
	}

	return level.TotalPoints >= criteria.Threshold, nil
}

// checkDailyActivity checks if user has been active today
func (s *AchievementService) checkDailyActivity(ctx context.Context, userID uuid.UUID, criteria models.AchievementCriteria, eventData map[string]interface{}) (bool, error) {
	streak, err := s.streakRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, nil
	}

	// Check if last activity was today
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	lastActivity := time.Date(streak.LastActivityAt.Year(), streak.LastActivityAt.Month(), streak.LastActivityAt.Day(), 0, 0, 0, 0, streak.LastActivityAt.Location())

	return !lastActivity.Before(today), nil
}

// UnlockAchievement unlocks an achievement for a user
func (s *AchievementService) UnlockAchievement(ctx context.Context, userID uuid.UUID, achievementID uuid.UUID) error {
	// Check if already unlocked
	existing, err := s.userAchievementRepo.GetByUserAndAchievement(ctx, userID, achievementID)
	if err == nil && existing != nil {
		return ErrAchievementAlreadyUnlocked
	}

	// Get achievement to get points
	achievement, err := s.achievementRepo.GetByID(ctx, achievementID)
	if err != nil {
		return ErrAchievementNotFound
	}

	// Create user achievement
	userAchievement := &models.UserAchievement{
		ID:            uuid.New(),
		UserID:        userID,
		AchievementID: achievementID,
		EarnedAt:      time.Now(),
		IsNotified:    false,
	}

	err = s.userAchievementRepo.Create(ctx, userAchievement)
	if err != nil {
		return err
	}

	// Award points
	if achievement.Points > 0 {
		err = s.AwardPoints(ctx, userID, achievement.Points, "achievement", &achievementID, fmt.Sprintf("解锁成就：%s", achievement.Name))
		if err != nil {
			return err
		}
	}

	return nil
}

// AwardPoints awards points to a user
func (s *AchievementService) AwardPoints(ctx context.Context, userID uuid.UUID, amount int, sourceType string, sourceID *uuid.UUID, description string) error {
	// Get or create user level
	level, err := s.userLevelRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Create new level entry
		level = &models.UserLevel{
			ID:            uuid.New(),
			UserID:        userID,
			Level:         1,
			CurrentPoints: 0,
			TotalPoints:   0,
			Experience:    0,
			NextLevelExp:  100,
			Title:         "初学者",
			UpdatedAt:     time.Now(),
		}
		err = s.userLevelRepo.Create(ctx, level)
		if err != nil {
			return err
		}
	}

	// Update points
	oldBalance := level.CurrentPoints
	newBalance := oldBalance + amount
	level.CurrentPoints = newBalance
	level.TotalPoints += amount
	level.Experience += amount
	level.UpdatedAt = time.Now()

	// Check for level up
	for level.Experience >= level.NextLevelExp {
		level.Level++
		level.Experience -= level.NextLevelExp
		level.NextLevelExp = int(float64(level.NextLevelExp) * 1.5) // Increase difficulty
		
		// Update title based on level
		level.Title = s.getTitleForLevel(level.Level)
	}

	err = s.userLevelRepo.Update(ctx, level)
	if err != nil {
		return err
	}

	// Record transaction
	transaction := &models.UserPointsTransaction{
		ID:           uuid.New(),
		UserID:       userID,
		Amount:       amount,
		BalanceAfter: newBalance,
		SourceType:   sourceType,
		SourceID:     sourceID,
		Description:  description,
		CreatedAt:    time.Now(),
	}

	return s.pointsTransactionRepo.Create(ctx, transaction)
}

// getTitleForLevel returns the title for a given level
func (s *AchievementService) getTitleForLevel(level int) string {
	titles := []string{
		"初学者",      // 1
		"学习者",      // 2
		"进阶者",      // 3
		"熟练者",      // 4
		"专家",        // 5
		"高手",        // 6
		"大师",        // 7
		"宗师",        // 8
		"传奇",        // 9
		"神话",        // 10
	}

	if level <= len(titles) {
		return titles[level-1]
	}
	return "传奇大师"
}

// UpdateStreak updates user's learning streak
func (s *AchievementService) UpdateStreak(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	streak, err := s.streakRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Create new streak
		streak = &models.UserStreak{
			ID:             uuid.New(),
			UserID:         userID,
			CurrentStreak:  1,
			LongestStreak:  1,
			LastActivityAt: now,
			UpdatedAt:      now,
		}
		return s.streakRepo.Create(ctx, streak)
	}

	lastActivity := time.Date(streak.LastActivityAt.Year(), streak.LastActivityAt.Month(), streak.LastActivityAt.Day(), 0, 0, 0, 0, streak.LastActivityAt.Location())

	// If already active today, don't update
	if !lastActivity.Before(today) {
		return nil
	}

	// Check if it's consecutive day
	yesterday := today.AddDate(0, 0, -1)
	if lastActivity.Before(yesterday) {
		// Streak broken
		streak.CurrentStreak = 1
	} else {
		// Consecutive day
		streak.CurrentStreak++
		if streak.CurrentStreak > streak.LongestStreak {
			streak.LongestStreak = streak.CurrentStreak
		}
	}

	streak.LastActivityAt = now
	streak.UpdatedAt = now

	return s.streakRepo.Update(ctx, streak)
}

// GetUserAchievements gets all achievements for a user
func (s *AchievementService) GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]*models.UserAchievement, error) {
	return s.userAchievementRepo.GetByUserID(ctx, userID)
}

// GetUserAchievementsWithProgress gets achievements with progress info
func (s *AchievementService) GetUserAchievementsWithProgress(ctx context.Context, userID uuid.UUID) ([]*models.AchievementWithProgress, error) {
	achievements, err := s.achievementRepo.GetAllEnabled(ctx)
	if err != nil {
		return nil, err
	}

	unlockedIDs, err := s.userAchievementRepo.GetUserAchievementIDs(ctx, userID)
	if err != nil {
		unlockedIDs = []uuid.UUID{}
	}

	result := []*models.AchievementWithProgress{}

	for _, achievement := range achievements {
		progress := &models.AchievementWithProgress{
			Achievement: *achievement,
		}

		// Check if unlocked
		if containsUUID(unlockedIDs, achievement.ID) {
			progress.IsUnlocked = true
			ua, err := s.userAchievementRepo.GetByUserAndAchievement(ctx, userID, achievement.ID)
			if err == nil && ua != nil {
				progress.EarnedAt = &ua.EarnedAt
			}
			progress.ProgressPercent = 100
			progress.CurrentCount = achievement.Criteria.Threshold
		} else {
			// Calculate progress
			count, err := s.getProgressCount(ctx, userID, achievement.Criteria)
			if err == nil {
				progress.CurrentCount = count
				if achievement.Criteria.Threshold > 0 {
					progress.ProgressPercent = float64(count) / float64(achievement.Criteria.Threshold) * 100
					if progress.ProgressPercent > 100 {
						progress.ProgressPercent = 100
					}
				}
			}
		}

		result = append(result, progress)
	}

	return result, nil
}

// getProgressCount gets current progress count for a criteria type
func (s *AchievementService) getProgressCount(ctx context.Context, userID uuid.UUID, criteria models.AchievementCriteria) (int, error) {
	switch criteria.Type {
	case "complete_lesson":
		progresses, err := s.progressRepo.GetByUserID(ctx, userID)
		if err != nil {
			return 0, err
		}
		count := 0
		for _, p := range progresses {
			if p.IsCompleted {
				count++
			}
		}
		return count, nil
	case "complete_course":
		enrollments, err := s.enrollmentRepo.GetByUserID(ctx, userID)
		if err != nil {
			return 0, err
		}
		count := 0
		for _, e := range enrollments {
			if e.Status == "completed" {
				count++
			}
		}
		return count, nil
	case "streak_days":
		streak, err := s.streakRepo.GetByUserID(ctx, userID)
		if err != nil {
			return 0, err
		}
		return streak.CurrentStreak, nil
	case "complete_exercises":
		submissions, err := s.submissionRepo.GetByUserID(ctx, userID)
		if err != nil {
			return 0, err
		}
		count := 0
		for _, sub := range submissions {
			if sub.IsCorrect != nil && *sub.IsCorrect {
				count++
			}
		}
		return count, nil
	case "total_points":
		level, err := s.userLevelRepo.GetByUserID(ctx, userID)
		if err != nil {
			return 0, err
		}
		return level.TotalPoints, nil
	default:
		return 0, nil
	}
}

// GetUserLevel gets user's level info
func (s *AchievementService) GetUserLevel(ctx context.Context, userID uuid.UUID) (*models.UserLevel, error) {
	return s.userLevelRepo.GetByUserID(ctx, userID)
}

// GetUserStreak gets user's streak info
func (s *AchievementService) GetUserStreak(ctx context.Context, userID uuid.UUID) (*models.UserStreak, error) {
	return s.streakRepo.GetByUserID(ctx, userID)
}

// GetLeaderboard gets leaderboard entries
func (s *AchievementService) GetLeaderboard(ctx context.Context, leaderboardType models.LeaderboardType, limit int, userID *uuid.UUID) ([]*models.LeaderboardEntry, error) {
	switch leaderboardType {
	case models.LeaderboardTypeWeekly:
		return s.leaderboardRepo.GetWeekly(ctx, limit)
	case models.LeaderboardTypeMonthly:
		return s.leaderboardRepo.GetMonthly(ctx, limit)
	case models.LeaderboardTypeAllTime:
		return s.leaderboardRepo.GetAllTime(ctx, limit)
	case models.LeaderboardTypeFriends:
		if userID == nil {
			return []*models.LeaderboardEntry{}, nil
		}
		return s.leaderboardRepo.GetFriends(ctx, *userID, limit)
	default:
		return s.leaderboardRepo.GetAllTime(ctx, limit)
	}
}

// GetUserAchievementSummary gets comprehensive achievement summary for a user
func (s *AchievementService) GetUserAchievementSummary(ctx context.Context, userID uuid.UUID) (*models.UserAchievementSummary, error) {
	// Get unlocked achievements count
	unlockedCount, err := s.userAchievementRepo.GetCountByUserID(ctx, userID)
	if err != nil {
		unlockedCount = 0
	}

	// Get total achievements
	totalCount, err := s.achievementRepo.GetCount(ctx)
	if err != nil {
		totalCount = 0
	}

	// Get user level
	level, err := s.userLevelRepo.GetByUserID(ctx, userID)
	if err != nil {
		level = &models.UserLevel{
			Level:       1,
			TotalPoints: 0,
			Title:       "初学者",
		}
	}

	// Get streak
	streak, err := s.streakRepo.GetByUserID(ctx, userID)
	if err != nil {
		streak = &models.UserStreak{
			CurrentStreak: 0,
			LongestStreak: 0,
		}
	}

	// Get rank
	rank, err := s.leaderboardRepo.GetUserRank(ctx, userID)
	if err != nil {
		rank = 0
	}

	return &models.UserAchievementSummary{
		UserID:              userID,
		TotalAchievements:   totalCount,
		UnlockedAchievements: unlockedCount,
		TotalPoints:         level.TotalPoints,
		CurrentLevel:        level.Level,
		CurrentTitle:        level.Title,
		CurrentStreak:       streak.CurrentStreak,
		LongestStreak:       streak.LongestStreak,
		Rank:                rank,
	}, nil
}

// GetPointsHistory gets user's points transaction history
func (s *AchievementService) GetPointsHistory(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]*models.UserPointsTransaction, error) {
	return s.pointsTransactionRepo.GetByUserID(ctx, userID, limit, offset)
}

// containsUUID checks if a slice contains a UUID
func containsUUID(slice []uuid.UUID, id uuid.UUID) bool {
	for _, item := range slice {
		if item == id {
			return true
		}
	}
	return false
}

// InitializeDefaultAchievements creates default achievements in the system
func (s *AchievementService) InitializeDefaultAchievements(ctx context.Context) error {
	defaultAchievements := []struct {
		Name        string
		Description string
		Points      int
		Type        models.AchievementType
		Tier        models.AchievementTier
		Criteria    models.AchievementCriteria
	}{
		{
			Name:        "初学者",
			Description: "完成第一节课",
			Points:      10,
			Type:        models.AchievementTypeCourse,
			Tier:        models.AchievementTierBronze,
			Criteria:    models.AchievementCriteria{Type: "complete_lesson", Threshold: 1},
		},
		{
			Name:        "持之以恒",
			Description: "连续学习 7 天",
			Points:      50,
			Type:        models.AchievementTypeStreak,
			Tier:        models.AchievementTierSilver,
			Criteria:    models.AchievementCriteria{Type: "streak_days", Threshold: 7},
		},
		{
			Name:        "学霸",
			Description: "完成所有课程",
			Points:      200,
			Type:        models.AchievementTypeCourse,
			Tier:        models.AchievementTierGold,
			Criteria:    models.AchievementCriteria{Type: "complete_course", Threshold: 10},
		},
		{
			Name:        "刷题达人",
			Description: "完成 100 道练习题",
			Points:      150,
			Type:        models.AchievementTypeExercise,
			Tier:        models.AchievementTierGold,
			Criteria:    models.AchievementCriteria{Type: "complete_exercises", Threshold: 100},
		},
		{
			Name:        "坚持不懈",
			Description: "连续学习 30 天",
			Points:      300,
			Type:        models.AchievementTypeStreak,
			Tier:        models.AchievementTierPlatinum,
			Criteria:    models.AchievementCriteria{Type: "streak_days", Threshold: 30},
		},
		{
			Name:        "知识渊博",
			Description: "获得 1000 积分",
			Points:      100,
			Type:        models.AchievementTypeMilestone,
			Tier:        models.AchievementTierGold,
			Criteria:    models.AchievementCriteria{Type: "total_points", Threshold: 1000},
		},
		{
			Name:        "每日精进",
			Description: "连续学习 3 天",
			Points:      20,
			Type:        models.AchievementTypeStreak,
			Tier:        models.AchievementTierBronze,
			Criteria:    models.AchievementCriteria{Type: "streak_days", Threshold: 3},
		},
		{
			Name:        "入门弟子",
			Description: "完成 5 节课",
			Points:      30,
			Type:        models.AchievementTypeCourse,
			Tier:        models.AchievementTierBronze,
			Criteria:    models.AchievementCriteria{Type: "complete_lesson", Threshold: 5},
		},
		{
			Name:        "登堂入室",
			Description: "完成 20 节课",
			Points:      80,
			Type:        models.AchievementTypeCourse,
			Tier:        models.AchievementTierSilver,
			Criteria:    models.AchievementCriteria{Type: "complete_lesson", Threshold: 20},
		},
		{
			Name:        "练习生",
			Description: "完成 10 道练习题",
			Points:      25,
			Type:        models.AchievementTypeExercise,
			Tier:        models.AchievementTierBronze,
			Criteria:    models.AchievementCriteria{Type: "complete_exercises", Threshold: 10},
		},
	}

	for _, ach := range defaultAchievements {
		existing, err := s.achievementRepo.GetByName(ctx, ach.Name)
		if err != nil || existing == nil {
			newAchievement := &models.Achievement{
				ID:              uuid.New(),
				Name:            ach.Name,
				Description:     ach.Description,
				Points:          ach.Points,
				AchievementType: ach.Type,
				Tier:            ach.Tier,
				Criteria:        ach.Criteria,
				IsEnabled:       true,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}
			err = s.achievementRepo.Create(ctx, newAchievement)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
