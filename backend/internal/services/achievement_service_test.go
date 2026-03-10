package services_test

import (
	"context"
	"testing"
	"time"

	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"
	"ai-learning-platform/internal/services"

	"github.com/google/uuid"
)

// Mock repositories for achievement tests
type MockAchievementRepository struct {
	achievements []*models.Achievement
}

func (m *MockAchievementRepository) Create(ctx context.Context, achievement *models.Achievement) error {
	m.achievements = append(m.achievements, achievement)
	return nil
}

func (m *MockAchievementRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Achievement, error) {
	for _, a := range m.achievements {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *MockAchievementRepository) GetByName(ctx context.Context, name string) (*models.Achievement, error) {
	for _, a := range m.achievements {
		if a.Name == name {
			return a, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *MockAchievementRepository) GetAllEnabled(ctx context.Context) ([]*models.Achievement, error) {
	result := make([]*models.Achievement, 0)
	for _, a := range m.achievements {
		if a.IsEnabled {
			result = append(result, a)
		}
	}
	return result, nil
}

func (m *MockAchievementRepository) GetCount(ctx context.Context) (int, error) {
	count := 0
	for _, a := range m.achievements {
		if a.IsEnabled {
			count++
		}
	}
	return count, nil
}

func (m *MockAchievementRepository) Update(ctx context.Context, achievement *models.Achievement) error {
	for i, a := range m.achievements {
		if a.ID == achievement.ID {
			m.achievements[i] = achievement
			return nil
		}
	}
	return repository.ErrNotFound
}

func (m *MockAchievementRepository) Delete(ctx context.Context, id uuid.UUID) error {
	for i, a := range m.achievements {
		if a.ID == id {
			m.achievements = append(m.achievements[:i], m.achievements[i+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

type MockUserAchievementRepository struct {
	userAchievements []*models.UserAchievement
}

func (m *MockUserAchievementRepository) Create(ctx context.Context, ua *models.UserAchievement) error {
	m.userAchievements = append(m.userAchievements, ua)
	return nil
}

func (m *MockUserAchievementRepository) GetByUserAndAchievement(ctx context.Context, userID, achievementID uuid.UUID) (*models.UserAchievement, error) {
	for _, ua := range m.userAchievements {
		if ua.UserID == userID && ua.AchievementID == achievementID {
			return ua, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *MockUserAchievementRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.UserAchievement, error) {
	result := make([]*models.UserAchievement, 0)
	for _, ua := range m.userAchievements {
		if ua.UserID == userID {
			result = append(result, ua)
		}
	}
	return result, nil
}

func (m *MockUserAchievementRepository) GetUserAchievementIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0)
	for _, ua := range m.userAchievements {
		if ua.UserID == userID {
			ids = append(ids, ua.AchievementID)
		}
	}
	return ids, nil
}

func (m *MockUserAchievementRepository) GetCountByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	count := 0
	for _, ua := range m.userAchievements {
		if ua.UserID == userID {
			count++
		}
	}
	return count, nil
}

type MockUserLevelRepository struct {
	levels map[uuid.UUID]*models.UserLevel
}

func NewMockUserLevelRepository() *MockUserLevelRepository {
	return &MockUserLevelRepository{levels: make(map[uuid.UUID]*models.UserLevel)}
}

func (m *MockUserLevelRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserLevel, error) {
	if level, exists := m.levels[userID]; exists {
		return level, nil
	}
	return nil, repository.ErrNotFound
}

func (m *MockUserLevelRepository) Create(ctx context.Context, level *models.UserLevel) error {
	m.levels[level.UserID] = level
	return nil
}

func (m *MockUserLevelRepository) Update(ctx context.Context, level *models.UserLevel) error {
	m.levels[level.UserID] = level
	return nil
}

type MockPointsTransactionRepository struct {
	transactions []*models.UserPointsTransaction
}

func (m *MockPointsTransactionRepository) Create(ctx context.Context, tx *models.UserPointsTransaction) error {
	m.transactions = append(m.transactions, tx)
	return nil
}

func (m *MockPointsTransactionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.UserPointsTransaction, error) {
	result := make([]*models.UserPointsTransaction, 0)
	for _, tx := range m.transactions {
		if tx.UserID == userID {
			result = append(result, tx)
		}
	}
	return result, nil
}

type MockStreakRepository struct {
	streaks map[uuid.UUID]*models.UserStreak
}

func NewMockStreakRepository() *MockStreakRepository {
	return &MockStreakRepository{streaks: make(map[uuid.UUID]*models.UserStreak)}
}

func (m *MockStreakRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserStreak, error) {
	if streak, exists := m.streaks[userID]; exists {
		return streak, nil
	}
	return nil, repository.ErrNotFound
}

func (m *MockStreakRepository) Create(ctx context.Context, streak *models.UserStreak) error {
	m.streaks[streak.UserID] = streak
	return nil
}

func (m *MockStreakRepository) Update(ctx context.Context, streak *models.UserStreak) error {
	m.streaks[streak.UserID] = streak
	return nil
}

type MockLeaderboardRepository struct{}

func (m *MockLeaderboardRepository) GetWeekly(ctx context.Context, limit int) ([]*models.LeaderboardEntry, error) {
	return []*models.LeaderboardEntry{}, nil
}

func (m *MockLeaderboardRepository) GetMonthly(ctx context.Context, limit int) ([]*models.LeaderboardEntry, error) {
	return []*models.LeaderboardEntry{}, nil
}

func (m *MockLeaderboardRepository) GetAllTime(ctx context.Context, limit int) ([]*models.LeaderboardEntry, error) {
	return []*models.LeaderboardEntry{}, nil
}

func (m *MockLeaderboardRepository) GetFriends(ctx context.Context, userID uuid.UUID, limit int) ([]*models.LeaderboardEntry, error) {
	return []*models.LeaderboardEntry{}, nil
}

func (m *MockLeaderboardRepository) GetUserRank(ctx context.Context, userID *uuid.UUID) (int, error) {
	return 0, nil
}

type MockProgressRepositoryForAchievement struct {
	progresses []*models.Progress
}

func (m *MockProgressRepositoryForAchievement) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Progress, error) {
	result := make([]*models.Progress, 0)
	for _, p := range m.progresses {
		if p.UserID == userID {
			result = append(result, p)
		}
	}
	return result, nil
}

type MockEnrollmentRepositoryForAchievement struct {
	enrollments []*models.Enrollment
}

func (m *MockEnrollmentRepositoryForAchievement) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Enrollment, error) {
	result := make([]*models.Enrollment, 0)
	for _, e := range m.enrollments {
		if e.UserID == userID {
			result = append(result, e)
		}
	}
	return result, nil
}

type MockExerciseRepositoryForAchievement struct{}

func (m *MockExerciseRepositoryForAchievement) GetByID(ctx context.Context, id uuid.UUID) (*models.Exercise, error) {
	return &models.Exercise{}, nil
}

func (m *MockExerciseRepositoryForAchievement) GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*models.Exercise, error) {
	return []*models.Exercise{}, nil
}

type MockSubmissionRepositoryForAchievement struct {
	submissions []*models.Submission
}

func (m *MockSubmissionRepositoryForAchievement) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Submission, error) {
	result := make([]*models.Submission, 0)
	for _, s := range m.submissions {
		if s.UserID == userID {
			result = append(result, s)
		}
	}
	return result, nil
}

type MockUserRepository struct{}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return &models.User{ID: id, Username: "testuser"}, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return &models.User{Email: email}, nil
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	return nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	return nil
}

// ==================== Achievement Unlock Tests ====================

func TestUnlockAchievement_Basic(t *testing.T) {
	userID := uuid.New()
	achievementID := uuid.New()

	mockAchievementRepo := &MockAchievementRepository{
		achievements: []*models.Achievement{
			{
				ID:              achievementID,
				Name:            "First Lesson",
				Description:     "Complete your first lesson",
				Points:          10,
				AchievementType: models.AchievementTypeCourse,
				Tier:            models.AchievementTierBronze,
				IsEnabled:       true,
			},
		},
	}

	mockUserAchievementRepo := &MockUserAchievementRepository{}
	mockUserLevelRepo := NewMockUserLevelRepository()
	mockPointsTxRepo := &MockPointsTransactionRepository{}
	mockStreakRepo := NewMockStreakRepository()
	mockLeaderboardRepo := &MockLeaderboardRepository{}
	mockProgressRepo := &MockProgressRepositoryForAchievement{}
	mockEnrollmentRepo := &MockEnrollmentRepositoryForAchievement{}
	mockExerciseRepo := &MockExerciseRepositoryForAchievement{}
	mockSubmissionRepo := &MockSubmissionRepositoryForAchievement{}
	mockUserRepo := &MockUserRepository{}

	service := services.NewAchievementService(
		mockAchievementRepo,
		mockUserAchievementRepo,
		mockUserLevelRepo,
		mockPointsTxRepo,
		mockStreakRepo,
		mockLeaderboardRepo,
		mockProgressRepo,
		mockEnrollmentRepo,
		mockExerciseRepo,
		mockSubmissionRepo,
		mockUserRepo,
	)

	err := service.UnlockAchievement(context.Background(), userID, achievementID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(mockUserAchievementRepo.userAchievements) != 1 {
		t.Errorf("Expected 1 user achievement, got %d", len(mockUserAchievementRepo.userAchievements))
	}

	if len(mockPointsTxRepo.transactions) != 1 {
		t.Errorf("Expected 1 points transaction, got %d", len(mockPointsTxRepo.transactions))
	}

	tx := mockPointsTxRepo.transactions[0]
	if tx.Amount != 10 {
		t.Errorf("Expected 10 points, got %d", tx.Amount)
	}
}

func TestUnlockAchievement_AlreadyUnlocked(t *testing.T) {
	userID := uuid.New()
	achievementID := uuid.New()

	mockAchievementRepo := &MockAchievementRepository{
		achievements: []*models.Achievement{
			{
				ID:              achievementID,
				Name:            "First Lesson",
				Points:          10,
				AchievementType: models.AchievementTypeCourse,
				Tier:            models.AchievementTierBronze,
				IsEnabled:       true,
			},
		},
	}

	mockUserAchievementRepo := &MockUserAchievementRepository{
		userAchievements: []*models.UserAchievement{
			{
				ID:            uuid.New(),
				UserID:        userID,
				AchievementID: achievementID,
				EarnedAt:      time.Now(),
			},
		},
	}

	mockUserLevelRepo := NewMockUserLevelRepository()
	mockPointsTxRepo := &MockPointsTransactionRepository{}
	mockStreakRepo := NewMockStreakRepository()
	mockLeaderboardRepo := &MockLeaderboardRepository{}
	mockProgressRepo := &MockProgressRepositoryForAchievement{}
	mockEnrollmentRepo := &MockEnrollmentRepositoryForAchievement{}
	mockExerciseRepo := &MockExerciseRepositoryForAchievement{}
	mockSubmissionRepo := &MockSubmissionRepositoryForAchievement{}
	mockUserRepo := &MockUserRepository{}

	service := services.NewAchievementService(
		mockAchievementRepo,
		mockUserAchievementRepo,
		mockUserLevelRepo,
		mockPointsTxRepo,
		mockStreakRepo,
		mockLeaderboardRepo,
		mockProgressRepo,
		mockEnrollmentRepo,
		mockExerciseRepo,
		mockSubmissionRepo,
		mockUserRepo,
	)

	err := service.UnlockAchievement(context.Background(), userID, achievementID)
	if err == nil {
		t.Errorf("Expected error for already unlocked achievement")
	}

	t.Logf("Expected error received: %v", err)
}

// ==================== Points System Tests ====================

func TestAwardPoints_NewUser(t *testing.T) {
	userID := uuid.New()

	mockAchievementRepo := &MockAchievementRepository{}
	mockUserAchievementRepo := &MockUserAchievementRepository{}
	mockUserLevelRepo := NewMockUserLevelRepository()
	mockPointsTxRepo := &MockPointsTransactionRepository{}
	mockStreakRepo := NewMockStreakRepository()
	mockLeaderboardRepo := &MockLeaderboardRepository{}
	mockProgressRepo := &MockProgressRepositoryForAchievement{}
	mockEnrollmentRepo := &MockEnrollmentRepositoryForAchievement{}
	mockExerciseRepo := &MockExerciseRepositoryForAchievement{}
	mockSubmissionRepo := &MockSubmissionRepositoryForAchievement{}
	mockUserRepo := &MockUserRepository{}

	service := services.NewAchievementService(
		mockAchievementRepo,
		mockUserAchievementRepo,
		mockUserLevelRepo,
		mockPointsTxRepo,
		mockStreakRepo,
		mockLeaderboardRepo,
		mockProgressRepo,
		mockEnrollmentRepo,
		mockExerciseRepo,
		mockSubmissionRepo,
		mockUserRepo,
	)

	err := service.AwardPoints(context.Background(), userID, 50, "achievement", nil, "Test achievement")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	level, err := mockUserLevelRepo.GetByUserID(context.Background(), userID)
	if err != nil {
		t.Fatalf("Failed to get user level: %v", err)
	}

	if level.CurrentPoints != 50 {
		t.Errorf("Expected 50 current points, got %d", level.CurrentPoints)
	}

	if level.TotalPoints != 50 {
		t.Errorf("Expected 50 total points, got %d", level.TotalPoints)
	}

	if level.Experience != 50 {
		t.Errorf("Expected 50 experience, got %d", level.Experience)
	}

	if level.Level != 1 {
		t.Errorf("Expected level 1, got %d", level.Level)
	}
}

func TestAwardPoints_LevelUp(t *testing.T) {
	userID := uuid.New()

	mockAchievementRepo := &MockAchievementRepository{}
	mockUserAchievementRepo := &MockUserAchievementRepository{}
	mockUserLevelRepo := NewMockUserLevelRepository()
	mockPointsTxRepo := &MockPointsTransactionRepository{}
	mockStreakRepo := NewMockStreakRepository()
	mockLeaderboardRepo := &MockLeaderboardRepository{}
	mockProgressRepo := &MockProgressRepositoryForAchievement{}
	mockEnrollmentRepo := &MockEnrollmentRepositoryForAchievement{}
	mockExerciseRepo := &MockExerciseRepositoryForAchievement{}
	mockSubmissionRepo := &MockSubmissionRepositoryForAchievement{}
	mockUserRepo := &MockUserRepository{}

	service := services.NewAchievementService(
		mockAchievementRepo,
		mockUserAchievementRepo,
		mockUserLevelRepo,
		mockPointsTxRepo,
		mockStreakRepo,
		mockLeaderboardRepo,
		mockProgressRepo,
		mockEnrollmentRepo,
		mockExerciseRepo,
		mockSubmissionRepo,
		mockUserRepo,
	)

	// Award enough points to level up (100 XP needed for level 2)
	err := service.AwardPoints(context.Background(), userID, 150, "achievement", nil, "Big achievement")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	level, err := mockUserLevelRepo.GetByUserID(context.Background(), userID)
	if err != nil {
		t.Fatalf("Failed to get user level: %v", err)
	}

	t.Logf("Level: %d, XP: %d, Next Level XP: %d, Title: %s", level.Level, level.Experience, level.NextLevelExp, level.Title)

	if level.Level < 2 {
		t.Errorf("Expected level up to at least 2, got %d", level.Level)
	}
}

func TestAwardPoints_MultipleAwards(t *testing.T) {
	userID := uuid.New()

	mockAchievementRepo := &MockAchievementRepository{}
	mockUserAchievementRepo := &MockUserAchievementRepository{}
	mockUserLevelRepo := NewMockUserLevelRepository()
	mockPointsTxRepo := &MockPointsTransactionRepository{}
	mockStreakRepo := NewMockStreakRepository()
	mockLeaderboardRepo := &MockLeaderboardRepository{}
	mockProgressRepo := &MockProgressRepositoryForAchievement{}
	mockEnrollmentRepo := &MockEnrollmentRepositoryForAchievement{}
	mockExerciseRepo := &MockExerciseRepositoryForAchievement{}
	mockSubmissionRepo := &MockSubmissionRepositoryForAchievement{}
	mockUserRepo := &MockUserRepository{}

	service := services.NewAchievementService(
		mockAchievementRepo,
		mockUserAchievementRepo,
		mockUserLevelRepo,
		mockPointsTxRepo,
		mockStreakRepo,
		mockLeaderboardRepo,
		mockProgressRepo,
		mockEnrollmentRepo,
		mockExerciseRepo,
		mockSubmissionRepo,
		mockUserRepo,
	)

	// Award points multiple times
	for i := 0; i < 5; i++ {
		err := service.AwardPoints(context.Background(), userID, 20, "exercise", nil, "Exercise completed")
		if err != nil {
			t.Fatalf("Unexpected error on award %d: %v", i, err)
		}
	}

	level, err := mockUserLevelRepo.GetByUserID(context.Background(), userID)
	if err != nil {
		t.Fatalf("Failed to get user level: %v", err)
	}

	if level.TotalPoints != 100 {
		t.Errorf("Expected 100 total points, got %d", level.TotalPoints)
	}

	if len(mockPointsTxRepo.transactions) != 5 {
		t.Errorf("Expected 5 transactions, got %d", len(mockPointsTxRepo.transactions))
	}
}

// ==================== Streak Tests ====================

func TestUpdateStreak_NewUser(t *testing.T) {
	userID := uuid.New()

	mockAchievementRepo := &MockAchievementRepository{}
	mockUserAchievementRepo := &MockUserAchievementRepository{}
	mockUserLevelRepo := NewMockUserLevelRepository()
	mockPointsTxRepo := &MockPointsTransactionRepository{}
	mockStreakRepo := NewMockStreakRepository()
	mockLeaderboardRepo := &MockLeaderboardRepository{}
	mockProgressRepo := &MockProgressRepositoryForAchievement{}
	mockEnrollmentRepo := &MockEnrollmentRepositoryForAchievement{}
	mockExerciseRepo := &MockExerciseRepositoryForAchievement{}
	mockSubmissionRepo := &MockSubmissionRepositoryForAchievement{}
	mockUserRepo := &MockUserRepository{}

	service := services.NewAchievementService(
		mockAchievementRepo,
		mockUserAchievementRepo,
		mockUserLevelRepo,
		mockPointsTxRepo,
		mockStreakRepo,
		mockLeaderboardRepo,
		mockProgressRepo,
		mockEnrollmentRepo,
		mockExerciseRepo,
		mockSubmissionRepo,
		mockUserRepo,
	)

	err := service.UpdateStreak(context.Background(), userID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	streak, err := mockStreakRepo.GetByUserID(context.Background(), userID)
	if err != nil {
		t.Fatalf("Failed to get streak: %v", err)
	}

	if streak.CurrentStreak != 1 {
		t.Errorf("Expected streak of 1, got %d", streak.CurrentStreak)
	}

	if streak.LongestStreak != 1 {
		t.Errorf("Expected longest streak of 1, got %d", streak.LongestStreak)
	}
}

func TestUpdateStreak_ConsecutiveDays(t *testing.T) {
	userID := uuid.New()
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	mockStreakRepo := &MockStreakRepository{
		streaks: map[uuid.UUID]*models.UserStreak{
			userID: {
				ID:             uuid.New(),
				UserID:         userID,
				CurrentStreak:  5,
				LongestStreak:  5,
				LastActivityAt: yesterday,
			},
		},
	}

	mockAchievementRepo := &MockAchievementRepository{}
	mockUserAchievementRepo := &MockUserAchievementRepository{}
	mockUserLevelRepo := NewMockUserLevelRepository()
	mockPointsTxRepo := &MockPointsTransactionRepository{}
	mockLeaderboardRepo := &MockLeaderboardRepository{}
	mockProgressRepo := &MockProgressRepositoryForAchievement{}
	mockEnrollmentRepo := &MockEnrollmentRepositoryForAchievement{}
	mockExerciseRepo := &MockExerciseRepositoryForAchievement{}
	mockSubmissionRepo := &MockSubmissionRepositoryForAchievement{}
	mockUserRepo := &MockUserRepository{}

	service := services.NewAchievementService(
		mockAchievementRepo,
		mockUserAchievementRepo,
		mockUserLevelRepo,
		mockPointsTxRepo,
		mockStreakRepo,
		mockLeaderboardRepo,
		mockProgressRepo,
		mockEnrollmentRepo,
		mockExerciseRepo,
		mockSubmissionRepo,
		mockUserRepo,
	)

	err := service.UpdateStreak(context.Background(), userID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	streak, err := mockStreakRepo.GetByUserID(context.Background(), userID)
	if err != nil {
		t.Fatalf("Failed to get streak: %v", err)
	}

	if streak.CurrentStreak != 6 {
		t.Errorf("Expected streak of 6, got %d", streak.CurrentStreak)
	}

	if streak.LongestStreak != 6 {
		t.Errorf("Expected longest streak of 6, got %d", streak.LongestStreak)
	}
}

func TestUpdateStreak_SameDay(t *testing.T) {
	userID := uuid.New()
	now := time.Now()

	mockStreakRepo := &MockStreakRepository{
		streaks: map[uuid.UUID]*models.UserStreak{
			userID: {
				ID:             uuid.New(),
				UserID:         userID,
				CurrentStreak:  5,
				LongestStreak:  10,
				LastActivityAt: now.Add(-2 * time.Hour), // Earlier today
			},
		},
	}

	mockAchievementRepo := &MockAchievementRepository{}
	mockUserAchievementRepo := &MockUserAchievementRepository{}
	mockUserLevelRepo := NewMockUserLevelRepository()
	mockPointsTxRepo := &MockPointsTransactionRepository{}
	mockLeaderboardRepo := &MockLeaderboardRepository{}
	mockProgressRepo := &MockProgressRepositoryForAchievement{}
	mockEnrollmentRepo := &MockEnrollmentRepositoryForAchievement{}
	mockExerciseRepo := &MockExerciseRepositoryForAchievement{}
	mockSubmissionRepo := &MockSubmissionRepositoryForAchievement{}
	mockUserRepo := &MockUserRepository{}

	service := services.NewAchievementService(
		mockAchievementRepo,
		mockUserAchievementRepo,
		mockUserLevelRepo,
		mockPointsTxRepo,
		mockStreakRepo,
		mockLeaderboardRepo,
		mockProgressRepo,
		mockEnrollmentRepo,
		mockExerciseRepo,
		mockSubmissionRepo,
		mockUserRepo,
	)

	err := service.UpdateStreak(context.Background(), userID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	streak, err := mockStreakRepo.GetByUserID(context.Background(), userID)
	if err != nil {
		t.Fatalf("Failed to get streak: %v", err)
	}

	if streak.CurrentStreak != 5 {
		t.Errorf("Expected streak to remain 5 (same day), got %d", streak.CurrentStreak)
	}
}

func TestUpdateStreak_Broken(t *testing.T) {
	userID := uuid.New()
	now := time.Now()
	twoDaysAgo := now.AddDate(0, 0, -2)

	mockStreakRepo := &MockStreakRepository{
		streaks: map[uuid.UUID]*models.UserStreak{
			userID: {
				ID:             uuid.New(),
				UserID:         userID,
				CurrentStreak:  5,
				LongestStreak:  10,
				LastActivityAt: twoDaysAgo,
			},
		},
	}

	mockAchievementRepo := &MockAchievementRepository{}
	mockUserAchievementRepo := &MockUserAchievementRepository{}
	mockUserLevelRepo := NewMockUserLevelRepository()
	mockPointsTxRepo := &MockPointsTransactionRepository{}
	mockLeaderboardRepo := &MockLeaderboardRepository{}
	mockProgressRepo := &MockProgressRepositoryForAchievement{}
	mockEnrollmentRepo := &MockEnrollmentRepositoryForAchievement{}
	mockExerciseRepo := &MockExerciseRepositoryForAchievement{}
	mockSubmissionRepo := &MockSubmissionRepositoryForAchievement{}
	mockUserRepo := &MockUserRepository{}

	service := services.NewAchievementService(
		mockAchievementRepo,
		mockUserAchievementRepo,
		mockUserLevelRepo,
		mockPointsTxRepo,
		mockStreakRepo,
		mockLeaderboardRepo,
		mockProgressRepo,
		mockEnrollmentRepo,
		mockExerciseRepo,
		mockSubmissionRepo,
		mockUserRepo,
	)

	err := service.UpdateStreak(context.Background(), userID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	streak, err := mockStreakRepo.GetByUserID(context.Background(), userID)
	if err != nil {
		t.Fatalf("Failed to get streak: %v", err)
	}

	if streak.CurrentStreak != 1 {
		t.Errorf("Expected streak to reset to 1, got %d", streak.CurrentStreak)
	}

	if streak.LongestStreak != 10 {
		t.Errorf("Expected longest streak to remain 10, got %d", streak.LongestStreak)
	}
}

// ==================== Achievement Criteria Tests ====================

func TestCheckAndUnlockAchievements_LessonCompletion(t *testing.T) {
	userID := uuid.New()
	lessonID := uuid.New()
	achievementID := uuid.New()

	mockAchievementRepo := &MockAchievementRepository{
		achievements: []*models.Achievement{
			{
				ID:              achievementID,
				Name:            "First Lesson",
				Description:     "Complete your first lesson",
				Points:          10,
				AchievementType: models.AchievementTypeCourse,
				Tier:            models.AchievementTierBronze,
				Criteria:        models.AchievementCriteria{Type: "complete_lesson", Threshold: 1},
				IsEnabled:       true,
			},
		},
	}

	mockUserAchievementRepo := &MockUserAchievementRepository{}
	mockUserLevelRepo := NewMockUserLevelRepository()
	mockPointsTxRepo := &MockPointsTransactionRepository{}
	mockStreakRepo := NewMockStreakRepository()
	mockLeaderboardRepo := &MockLeaderboardRepository{}
	mockProgressRepo := &MockProgressRepositoryForAchievement{
		progresses: []*models.Progress{
			{
				ID:          uuid.New(),
				UserID:      userID,
				LessonID:    lessonID,
				IsCompleted: true,
			},
		},
	}
	mockEnrollmentRepo := &MockEnrollmentRepositoryForAchievement{}
	mockExerciseRepo := &MockExerciseRepositoryForAchievement{}
	mockSubmissionRepo := &MockSubmissionRepositoryForAchievement{}
	mockUserRepo := &MockUserRepository{}

	service := services.NewAchievementService(
		mockAchievementRepo,
		mockUserAchievementRepo,
		mockUserLevelRepo,
		mockPointsTxRepo,
		mockStreakRepo,
		mockLeaderboardRepo,
		mockProgressRepo,
		mockEnrollmentRepo,
		mockExerciseRepo,
		mockSubmissionRepo,
		mockUserRepo,
	)

	unlocked, err := service.CheckAndUnlockAchievements(context.Background(), userID, "lesson_complete", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(unlocked) != 1 {
		t.Errorf("Expected 1 unlocked achievement, got %d", len(unlocked))
	} else {
		t.Logf("Unlocked achievement: %s", unlocked[0].Name)
	}
}

// ==================== User Achievement Summary Tests ====================

func TestGetUserAchievementSummary(t *testing.T) {
	userID := uuid.New()

	mockAchievementRepo := &MockAchievementRepository{
		achievements: []*models.Achievement{
			{ID: uuid.New(), Name: "Achievement 1", IsEnabled: true},
			{ID: uuid.New(), Name: "Achievement 2", IsEnabled: true},
			{ID: uuid.New(), Name: "Achievement 3", IsEnabled: true},
		},
	}

	mockUserAchievementRepo := &MockUserAchievementRepository{
		userAchievements: []*models.UserAchievement{
			{ID: uuid.New(), UserID: userID, AchievementID: mockAchievementRepo.achievements[0].ID},
			{ID: uuid.New(), UserID: userID, AchievementID: mockAchievementRepo.achievements[1].ID},
		},
	}

	mockUserLevelRepo := NewMockUserLevelRepository()
	mockUserLevelRepo.levels[userID] = &models.UserLevel{
		ID:            uuid.New(),
		UserID:        userID,
		Level:         5,
		CurrentPoints: 250,
		TotalPoints:   500,
		Title:         "专家",
	}

	mockStreakRepo := &MockStreakRepository{
		streaks: map[uuid.UUID]*models.UserStreak{
			userID: {
				ID:            uuid.New(),
				UserID:        userID,
				CurrentStreak: 7,
				LongestStreak: 15,
			},
		},
	}

	mockAchievementRepo2 := &MockAchievementRepository{}
	mockPointsTxRepo := &MockPointsTransactionRepository{}
	mockLeaderboardRepo := &MockLeaderboardRepository{}
	mockProgressRepo := &MockProgressRepositoryForAchievement{}
	mockEnrollmentRepo := &MockEnrollmentRepositoryForAchievement{}
	mockExerciseRepo := &MockExerciseRepositoryForAchievement{}
	mockSubmissionRepo := &MockSubmissionRepositoryForAchievement{}
	mockUserRepo := &MockUserRepository{}

	service := services.NewAchievementService(
		mockAchievementRepo,
		mockUserAchievementRepo,
		mockUserLevelRepo,
		mockPointsTxRepo,
		mockStreakRepo,
		mockLeaderboardRepo,
		mockProgressRepo,
		mockEnrollmentRepo,
		mockExerciseRepo,
		mockSubmissionRepo,
		mockUserRepo,
	)

	summary, err := service.GetUserAchievementSummary(context.Background(), userID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Logf("User Achievement Summary:")
	t.Logf("  Total Achievements: %d", summary.TotalAchievements)
	t.Logf("  Unlocked: %d", summary.UnlockedAchievements)
	t.Logf("  Total Points: %d", summary.TotalPoints)
	t.Logf("  Level: %d (%s)", summary.CurrentLevel, summary.CurrentTitle)
	t.Logf("  Current Streak: %d", summary.CurrentStreak)
	t.Logf("  Longest Streak: %d", summary.LongestStreak)

	if summary.UnlockedAchievements != 2 {
		t.Errorf("Expected 2 unlocked achievements, got %d", summary.UnlockedAchievements)
	}

	if summary.CurrentLevel != 5 {
		t.Errorf("Expected level 5, got %d", summary.CurrentLevel)
	}
}

// ==================== Default Achievements Initialization ====================

func TestInitializeDefaultAchievements(t *testing.T) {
	mockAchievementRepo := &MockAchievementRepository{}
	mockUserAchievementRepo := &MockUserAchievementRepository{}
	mockUserLevelRepo := NewMockUserLevelRepository()
	mockPointsTxRepo := &MockPointsTransactionRepository{}
	mockStreakRepo := NewMockStreakRepository()
	mockLeaderboardRepo := &MockLeaderboardRepository{}
	mockProgressRepo := &MockProgressRepositoryForAchievement{}
	mockEnrollmentRepo := &MockEnrollmentRepositoryForAchievement{}
	mockExerciseRepo := &MockExerciseRepositoryForAchievement{}
	mockSubmissionRepo := &MockSubmissionRepositoryForAchievement{}
	mockUserRepo := &MockUserRepository{}

	service := services.NewAchievementService(
		mockAchievementRepo,
		mockUserAchievementRepo,
		mockUserLevelRepo,
		mockPointsTxRepo,
		mockStreakRepo,
		mockLeaderboardRepo,
		mockProgressRepo,
		mockEnrollmentRepo,
		mockExerciseRepo,
		mockSubmissionRepo,
		mockUserRepo,
	)

	err := service.InitializeDefaultAchievements(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Logf("Created %d default achievements", len(mockAchievementRepo.achievements))

	if len(mockAchievementRepo.achievements) == 0 {
		t.Errorf("Expected some default achievements to be created")
	}

	// Verify specific achievements exist
	achievementNames := make(map[string]bool)
	for _, a := range mockAchievementRepo.achievements {
		achievementNames[a.Name] = true
	}

	expectedAchievements := []string{
		"初学者",
		"持之以恒",
		"学霸",
		"刷题达人",
		"坚持不懈",
	}

	for _, expected := range expectedAchievements {
		if !achievementNames[expected] {
			t.Errorf("Expected achievement '%s' not found", expected)
		}
	}
}

// ==================== Points History Tests ====================

func TestGetPointsHistory(t *testing.T) {
	userID := uuid.New()

	mockAchievementRepo := &MockAchievementRepository{}
	mockUserAchievementRepo := &MockUserAchievementRepository{}
	mockUserLevelRepo := NewMockUserLevelRepository()
	mockPointsTxRepo := &MockPointsTransactionRepository{
		transactions: []*models.UserPointsTransaction{
			{
				ID:           uuid.New(),
				UserID:       userID,
				Amount:       10,
				BalanceAfter: 10,
				SourceType:   "achievement",
				Description:  "First achievement",
				CreatedAt:    time.Now(),
			},
			{
				ID:           uuid.New(),
				UserID:       userID,
				Amount:       20,
				BalanceAfter: 30,
				SourceType:   "exercise",
				Description:  "Exercise completed",
				CreatedAt:    time.Now(),
			},
			{
				ID:           uuid.New(),
				UserID:       userID,
				Amount:       50,
				BalanceAfter: 80,
				SourceType:   "achievement",
				Description:  "Big achievement",
				CreatedAt:    time.Now(),
			},
		},
	}
	mockStreakRepo := NewMockStreakRepository()
	mockLeaderboardRepo := &MockLeaderboardRepository{}
	mockProgressRepo := &MockProgressRepositoryForAchievement{}
	mockEnrollmentRepo := &MockEnrollmentRepositoryForAchievement{}
	mockExerciseRepo := &MockExerciseRepositoryForAchievement{}
	mockSubmissionRepo := &MockSubmissionRepositoryForAchievement{}
	mockUserRepo := &MockUserRepository{}

	service := services.NewAchievementService(
		mockAchievementRepo,
		mockUserAchievementRepo,
		mockUserLevelRepo,
		mockPointsTxRepo,
		mockStreakRepo,
		mockLeaderboardRepo,
		mockProgressRepo,
		mockEnrollmentRepo,
		mockExerciseRepo,
		mockSubmissionRepo,
		mockUserRepo,
	)

	history, err := service.GetPointsHistory(context.Background(), userID, 10, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Logf("Points history entries: %d", len(history))

	if len(history) != 3 {
		t.Errorf("Expected 3 transactions, got %d", len(history))
	}
}
