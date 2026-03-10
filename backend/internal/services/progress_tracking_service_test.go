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

// Mock repositories for progress tracking tests
type MockProgressRepository struct {
	progresses []*models.Progress
}

func (m *MockProgressRepository) Create(ctx context.Context, progress *models.Progress) error {
	m.progresses = append(m.progresses, progress)
	return nil
}

func (m *MockProgressRepository) GetByUserAndLesson(ctx context.Context, userID, lessonID uuid.UUID) (*models.Progress, error) {
	for _, p := range m.progresses {
		if p.UserID == userID && p.LessonID == lessonID {
			return p, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *MockProgressRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Progress, error) {
	result := make([]*models.Progress, 0)
	for _, p := range m.progresses {
		if p.UserID == userID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *MockProgressRepository) MarkCompleted(ctx context.Context, userID, lessonID uuid.UUID) error {
	for _, p := range m.progresses {
		if p.UserID == userID && p.LessonID == lessonID {
			p.IsCompleted = true
			now := time.Now()
			p.CompletedAt = &now
			return nil
		}
	}
	return repository.ErrNotFound
}

func (m *MockProgressRepository) Update(ctx context.Context, progress *models.Progress) error {
	for i, p := range m.progresses {
		if p.ID == progress.ID {
			m.progresses[i] = progress
			return nil
		}
	}
	return repository.ErrNotFound
}

type MockEnrollmentRepository struct {
	enrollments []*models.Enrollment
}

func (m *MockEnrollmentRepository) GetByUserAndCourse(ctx context.Context, userID, courseID uuid.UUID) (*models.Enrollment, error) {
	for _, e := range m.enrollments {
		if e.UserID == userID && e.CourseID == courseID {
			return e, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *MockEnrollmentRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]*models.Enrollment, error) {
	result := make([]*models.Enrollment, 0)
	for _, e := range m.enrollments {
		if e.UserID == userID {
			result = append(result, e)
		}
	}
	return result, nil
}

func (m *MockEnrollmentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Enrollment, error) {
	return m.GetByUser(ctx, userID)
}

type MockLessonRepository struct {
	lessons []*models.Lesson
}

func (m *MockLessonRepository) GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]*models.Lesson, error) {
	result := make([]*models.Lesson, 0)
	for _, l := range m.lessons {
		if l.CourseID == courseID {
			result = append(result, l)
		}
	}
	return result, nil
}

func (m *MockLessonRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Lesson, error) {
	for _, l := range m.lessons {
		if l.ID == id {
			return l, nil
		}
	}
	return nil, repository.ErrNotFound
}

type MockCourseRepository struct {
	courses []*models.Course
}

func (m *MockCourseRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Course, error) {
	for _, c := range m.courses {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, repository.ErrNotFound
}

// ==================== Course Progress Tests ====================

func TestGetCourseProgress_BasicProgress(t *testing.T) {
	userID := uuid.New()
	courseID := uuid.New()
	lesson1ID := uuid.New()
	lesson2ID := uuid.New()

	mockProgressRepo := &MockProgressRepository{
		progresses: []*models.Progress{
			{
				ID:            uuid.New(),
				UserID:        userID,
				LessonID:      lesson1ID,
				IsCompleted:   true,
				VideoPosition: 300,
			},
			{
				ID:            uuid.New(),
				UserID:        userID,
				LessonID:      lesson2ID,
				IsCompleted:   false,
				VideoPosition: 150,
			},
		},
	}

	mockEnrollmentRepo := &MockEnrollmentRepository{
		enrollments: []*models.Enrollment{
			{
				ID:               uuid.New(),
				UserID:           userID,
				CourseID:         courseID,
				ProgressPercentage: 50,
				UpdatedAt:        time.Now(),
			},
		},
	}

	mockLessonRepo := &MockLessonRepository{
		lessons: []*models.Lesson{
			{ID: lesson1ID, CourseID: courseID, Title: "Lesson 1", VideoDuration: 300},
			{ID: lesson2ID, CourseID: courseID, Title: "Lesson 2", VideoDuration: 300},
		},
	}

	mockCourseRepo := &MockCourseRepository{}

	service := services.NewProgressTrackingService(
		mockProgressRepo,
		mockEnrollmentRepo,
		mockLessonRepo,
		mockCourseRepo,
	)

	progress, err := service.GetCourseProgress(context.Background(), userID, courseID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Logf("Course progress: %+v", progress)

	if progress["completed_lessons"].(int) != 1 {
		t.Errorf("Expected 1 completed lesson, got %d", progress["completed_lessons"])
	}

	if progress["total_lessons"].(int) != 2 {
		t.Errorf("Expected 2 total lessons, got %d", progress["total_lessons"])
	}

	progressPercent := progress["progress_percentage"].(float64)
	if progressPercent != 50.0 {
		t.Errorf("Expected 50%% progress, got %f", progressPercent)
	}
}

func TestUpdateVideoProgress_NewProgress(t *testing.T) {
	userID := uuid.New()
	lessonID := uuid.New()

	mockProgressRepo := &MockProgressRepository{}
	mockEnrollmentRepo := &MockEnrollmentRepository{}
	mockLessonRepo := &MockLessonRepository{}
	mockCourseRepo := &MockCourseRepository{}

	service := services.NewProgressTrackingService(
		mockProgressRepo,
		mockEnrollmentRepo,
		mockLessonRepo,
		mockCourseRepo,
	)

	err := service.UpdateVideoProgress(context.Background(), userID, lessonID, 120, 300)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(mockProgressRepo.progresses) != 1 {
		t.Errorf("Expected 1 progress record, got %d", len(mockProgressRepo.progresses))
	}

	progress := mockProgressRepo.progresses[0]
	if progress.VideoPosition != 120 {
		t.Errorf("Expected video position 120, got %d", progress.VideoPosition)
	}

	if !progress.IsWatching {
		t.Errorf("Expected IsWatching to be true")
	}
}

func TestUpdateVideoProgress_AutoComplete(t *testing.T) {
	userID := uuid.New()
	lessonID := uuid.New()

	mockProgressRepo := &MockProgressRepository{}
	mockEnrollmentRepo := &MockEnrollmentRepository{}
	mockLessonRepo := &MockLessonRepository{}
	mockCourseRepo := &MockCourseRepository{}

	service := services.NewProgressTrackingService(
		mockProgressRepo,
		mockEnrollmentRepo,
		mockLessonRepo,
		mockCourseRepo,
	)

	// Watch 95% of video (should auto-complete)
	err := service.UpdateVideoProgress(context.Background(), userID, lessonID, 285, 300)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	progress := mockProgressRepo.progresses[0]
	if !progress.IsCompleted {
		t.Errorf("Expected auto-completion at 95%% watch progress")
	}

	if progress.CompletedAt == nil {
		t.Errorf("Expected CompletedAt to be set")
	}
}

func TestUpdateVideoProgress_NoAutoComplete(t *testing.T) {
	userID := uuid.New()
	lessonID := uuid.New()

	mockProgressRepo := &MockProgressRepository{}
	mockEnrollmentRepo := &MockEnrollmentRepository{}
	mockLessonRepo := &MockLessonRepository{}
	mockCourseRepo := &MockCourseRepository{}

	service := services.NewProgressTrackingService(
		mockProgressRepo,
		mockEnrollmentRepo,
		mockLessonRepo,
		mockCourseRepo,
	)

	// Watch only 50% of video (should not auto-complete)
	err := service.UpdateVideoProgress(context.Background(), userID, lessonID, 150, 300)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	progress := mockProgressRepo.progresses[0]
	if progress.IsCompleted {
		t.Errorf("Expected no auto-completion at 50%% watch progress")
	}
}

func TestMarkLessonCompleted(t *testing.T) {
	userID := uuid.New()
	lessonID := uuid.New()

	mockProgressRepo := &MockProgressRepository{
		progresses: []*models.Progress{
			{
				ID:          uuid.New(),
				UserID:      userID,
				LessonID:    lessonID,
				IsCompleted: false,
			},
		},
	}

	mockEnrollmentRepo := &MockEnrollmentRepository{}
	mockLessonRepo := &MockLessonRepository{}
	mockCourseRepo := &MockCourseRepository{}

	service := services.NewProgressTrackingService(
		mockProgressRepo,
		mockEnrollmentRepo,
		mockLessonRepo,
		mockCourseRepo,
	)

	err := service.MarkLessonCompleted(context.Background(), userID, lessonID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !mockProgressRepo.progresses[0].IsCompleted {
		t.Errorf("Expected lesson to be marked as completed")
	}
}

// ==================== Learning Heatmap Tests ====================

func TestGetLearningHeatmapData(t *testing.T) {
	userID := uuid.New()
	now := time.Now()

	mockProgressRepo := &MockProgressRepository{
		progresses: []*models.Progress{
			{
				ID:             uuid.New(),
				UserID:         userID,
				LastAccessedAt: now,
				VideoPosition:  300,
				IsCompleted:    true,
			},
			{
				ID:             uuid.New(),
				UserID:         userID,
				LastAccessedAt: now.AddDate(0, 0, -1),
				VideoPosition:  600,
				IsCompleted:    true,
			},
			{
				ID:             uuid.New(),
				UserID:         userID,
				LastAccessedAt: now.AddDate(0, 0, -2),
				VideoPosition:  120,
				IsCompleted:    false,
			},
		},
	}

	mockEnrollmentRepo := &MockEnrollmentRepository{}
	mockLessonRepo := &MockLessonRepository{}
	mockCourseRepo := &MockCourseRepository{}

	service := services.NewProgressTrackingService(
		mockProgressRepo,
		mockEnrollmentRepo,
		mockLessonRepo,
		mockCourseRepo,
	)

	heatmapData, err := service.GetLearningHeatmapData(context.Background(), userID, 3)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Logf("Heatmap data points: %d", len(heatmapData))

	if len(heatmapData) == 0 {
		t.Errorf("Expected heatmap data, got none")
	}

	for _, data := range heatmapData {
		if data.Level < 0 || data.Level > 4 {
			t.Errorf("Invalid heatmap level: %d", data.Level)
		}
		t.Logf("Date: %s, Minutes: %d, Level: %d", data.Date, data.Count, data.Level)
	}
}

// ==================== Daily Stats Tests ====================

func TestGetDailyLearningStats(t *testing.T) {
	userID := uuid.New()
	now := time.Now()

	mockProgressRepo := &MockProgressRepository{
		progresses: []*models.Progress{
			{
				ID:             uuid.New(),
				UserID:         userID,
				LastAccessedAt: now,
				VideoPosition:  300,
				IsCompleted:    true,
			},
			{
				ID:             uuid.New(),
				UserID:         userID,
				LastAccessedAt: now.AddDate(0, 0, -1),
				VideoPosition:  600,
				IsCompleted:    true,
			},
		},
	}

	mockEnrollmentRepo := &MockEnrollmentRepository{}
	mockLessonRepo := &MockLessonRepository{}
	mockCourseRepo := &MockCourseRepository{}

	service := services.NewProgressTrackingService(
		mockProgressRepo,
		mockEnrollmentRepo,
		mockLessonRepo,
		mockCourseRepo,
	)

	stats, err := service.GetDailyLearningStats(context.Background(), userID, 7)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Logf("Daily stats entries: %d", len(stats))

	for _, stat := range stats {
		t.Logf("Date: %s, Seconds: %d, Lessons: %d", stat.Date, stat.TotalSeconds, stat.LessonsCompleted)
	}
}

// ==================== Weekly Report Tests ====================

func TestGenerateWeeklyReport(t *testing.T) {
	userID := uuid.New()
	courseID := uuid.New()

	mockProgressRepo := &MockProgressRepository{
		progresses: []*models.Progress{
			{
				ID:             uuid.New(),
				UserID:         userID,
				LastAccessedAt: time.Now(),
				VideoPosition:  300,
				IsCompleted:    true,
			},
		},
	}

	mockEnrollmentRepo := &MockEnrollmentRepository{
		enrollments: []*models.Enrollment{
			{
				ID:               uuid.New(),
				UserID:           userID,
				CourseID:         courseID,
				ProgressPercentage: 50,
				UpdatedAt:        time.Now(),
			},
		},
	}

	mockLessonRepo := &MockLessonRepository{
		lessons: []*models.Lesson{
			{ID: uuid.New(), CourseID: courseID, Title: "Test Lesson", VideoDuration: 300},
		},
	}

	mockCourseRepo := &MockCourseRepository{
		courses: []*models.Course{
			{ID: courseID, Title: "Test Course"},
		},
	}

	service := services.NewProgressTrackingService(
		mockProgressRepo,
		mockEnrollmentRepo,
		mockLessonRepo,
		mockCourseRepo,
	)

	report, err := service.GenerateWeeklyReport(context.Background(), userID, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Logf("Weekly Report:")
	t.Logf("  Week: %s to %s", report.WeekStart, report.WeekEnd)
	t.Logf("  Total Hours: %.2f", report.TotalHours)
	t.Logf("  Lessons Completed: %d", report.LessonsCompleted)
	t.Logf("  Avg Daily Minutes: %.2f", report.AvgDailyMinutes)

	if report.WeekStart == "" || report.WeekEnd == "" {
		t.Errorf("Expected week dates to be set")
	}
}

// ==================== Monthly Report Tests ====================

func TestGenerateMonthlyReport(t *testing.T) {
	userID := uuid.New()
	courseID := uuid.New()

	mockProgressRepo := &MockProgressRepository{
		progresses: []*models.Progress{
			{
				ID:             uuid.New(),
				UserID:         userID,
				LastAccessedAt: time.Now(),
				VideoPosition:  300,
				IsCompleted:    true,
			},
		},
	}

	mockEnrollmentRepo := &MockEnrollmentRepository{
		enrollments: []*models.Enrollment{
			{
				ID:               uuid.New(),
				UserID:           userID,
				CourseID:         courseID,
				Status:           "active",
				ProgressPercentage: 50,
				UpdatedAt:        time.Now(),
			},
		},
	}

	mockLessonRepo := &MockLessonRepository{
		lessons: []*models.Lesson{
			{ID: uuid.New(), CourseID: courseID, Title: "Test Lesson", VideoDuration: 300},
		},
	}

	mockCourseRepo := &MockCourseRepository{
		courses: []*models.Course{
			{ID: courseID, Title: "Test Course"},
		},
	}

	service := services.NewProgressTrackingService(
		mockProgressRepo,
		mockEnrollmentRepo,
		mockLessonRepo,
		mockCourseRepo,
	)

	report, err := service.GenerateMonthlyReport(context.Background(), userID, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Logf("Monthly Report:")
	t.Logf("  Month: %s %d", report.Month, report.Year)
	t.Logf("  Total Hours: %.2f", report.TotalHours)
	t.Logf("  Lessons Completed: %d", report.LessonsCompleted)
	t.Logf("  Best Day: %s (%.2f minutes)", report.BestDay, report.BestDayMinutes)

	if report.Month == "" {
		t.Errorf("Expected month to be set")
	}
}

// ==================== CSV Export Tests ====================

func TestExportReportToCSV(t *testing.T) {
	mockProgressRepo := &MockProgressRepository{}
	mockEnrollmentRepo := &MockEnrollmentRepository{}
	mockLessonRepo := &MockLessonRepository{}
	mockCourseRepo := &MockCourseRepository{}

	service := services.NewProgressTrackingService(
		mockProgressRepo,
		mockEnrollmentRepo,
		mockLessonRepo,
		mockCourseRepo,
	)

	report := &services.MonthlyReport{
		Month:            "March",
		Year:             2025,
		TotalHours:       10.5,
		LessonsCompleted: 15,
		CoursesCompleted: 2,
		AvgDailyMinutes:  45.5,
		BestDay:          "2025-03-15",
		BestDayMinutes:   120.0,
		DailyStats: []services.DailyLearningStats{
			{Date: "2025-03-01", TotalSeconds: 3600, LessonsCompleted: 2, CoursesAccessed: 1},
			{Date: "2025-03-02", TotalSeconds: 7200, LessonsCompleted: 3, CoursesAccessed: 1},
		},
		CoursesProgress: []services.CourseProgressSummary{
			{
				CourseTitle:      "Go Programming",
				ProgressPercent:  75.0,
				LessonsCompleted: 15,
				TotalLessons:     20,
				TimeSpentMinutes: 300,
			},
		},
	}

	csvData, err := service.ExportReportToCSV(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Logf("CSV Export successful (%d bytes)", len(csvData))
	t.Logf("First 500 chars:\n%s", csvData[:min(500, len(csvData))])
}

// ==================== Learning Time Stats Tests ====================

func TestGetLearningTimeStats(t *testing.T) {
	userID := uuid.New()

	mockProgressRepo := &MockProgressRepository{
		progresses: []*models.Progress{
			{
				ID:            uuid.New(),
				UserID:        userID,
				VideoPosition: 3600,
				IsCompleted:   true,
			},
			{
				ID:            uuid.New(),
				UserID:        userID,
				VideoPosition: 1800,
				IsCompleted:   false,
			},
			{
				ID:            uuid.New(),
				UserID:        userID,
				VideoPosition: 900,
				IsCompleted:   true,
			},
		},
	}

	mockEnrollmentRepo := &MockEnrollmentRepository{}
	mockLessonRepo := &MockLessonRepository{}
	mockCourseRepo := &MockCourseRepository{}

	service := services.NewProgressTrackingService(
		mockProgressRepo,
		mockEnrollmentRepo,
		mockLessonRepo,
		mockCourseRepo,
	)

	stats, err := service.GetLearningTimeStats(context.Background(), userID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Logf("Learning Time Stats: %+v", stats)

	totalHours := stats["total_learning_hours"].(float64)
	expectedHours := float64(3600+1800+900) / 3600.0
	if totalHours != expectedHours {
		t.Errorf("Expected %.2f hours, got %.2f", expectedHours, totalHours)
	}

	completionRate := stats["completion_rate"].(float64)
	expectedRate := 2.0 / 3.0 * 100.0
	if completionRate != expectedRate {
		t.Errorf("Expected %.2f%% completion rate, got %.2f", expectedRate, completionRate)
	}
}
