package services

import (
	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ProgressTrackingService handles comprehensive progress tracking
type ProgressTrackingService struct {
	progressRepo   *repository.ProgressRepository
	enrollmentRepo *repository.EnrollmentRepository
	lessonRepo     *repository.LessonRepository
	courseRepo     *repository.CourseRepository
}

// LearningSession represents a learning session for time tracking
type LearningSession struct {
	UserID    uuid.UUID `json:"user_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  int       `json:"duration"` // in seconds
	CourseID  uuid.UUID `json:"course_id"`
	LessonID  uuid.UUID `json:"lesson_id"`
}

// DailyLearningStats represents daily learning statistics
type DailyLearningStats struct {
	Date             string `json:"date"`
	TotalSeconds     int    `json:"total_seconds"`
	LessonsCompleted int    `json:"lessons_completed"`
	CoursesAccessed  int    `json:"courses_accessed"`
}

// WeeklyReport represents a weekly learning report
type WeeklyReport struct {
	WeekStart        string                  `json:"week_start"`
	WeekEnd          string                  `json:"week_end"`
	TotalHours       float64                 `json:"total_hours"`
	LessonsCompleted int                     `json:"lessons_completed"`
	CoursesProgress  []CourseProgressSummary `json:"courses_progress"`
	DailyStats       []DailyLearningStats    `json:"daily_stats"`
	AvgDailyMinutes  float64                 `json:"avg_daily_minutes"`
}

// MonthlyReport represents a monthly learning report
type MonthlyReport struct {
	Month            string                  `json:"month"`
	Year             int                     `json:"year"`
	TotalHours       float64                 `json:"total_hours"`
	LessonsCompleted int                     `json:"lessons_completed"`
	CoursesCompleted int                     `json:"courses_completed"`
	CoursesProgress  []CourseProgressSummary `json:"courses_progress"`
	DailyStats       []DailyLearningStats    `json:"daily_stats"`
	AvgDailyMinutes  float64                 `json:"avg_daily_minutes"`
	BestDay          string                  `json:"best_day"`
	BestDayMinutes   float64                 `json:"best_day_minutes"`
}

// CourseProgressSummary summarizes course progress
type CourseProgressSummary struct {
	CourseID         uuid.UUID `json:"course_id"`
	CourseTitle      string    `json:"course_title"`
	ProgressPercent  float64   `json:"progress_percent"`
	LessonsCompleted int       `json:"lessons_completed"`
	TotalLessons     int       `json:"total_lessons"`
	TimeSpentMinutes int       `json:"time_spent_minutes"`
}

// HeatmapData represents data for learning heatmap visualization
type HeatmapData struct {
	Date  string `json:"date"`
	Count int    `json:"count"` // minutes spent
	Level int    `json:"level"` // 0-4 intensity level
}

// NewProgressTrackingService creates a new ProgressTrackingService
func NewProgressTrackingService(
	progressRepo *repository.ProgressRepository,
	enrollmentRepo *repository.EnrollmentRepository,
	lessonRepo *repository.LessonRepository,
	courseRepo *repository.CourseRepository,
) *ProgressTrackingService {
	return &ProgressTrackingService{
		progressRepo:   progressRepo,
		enrollmentRepo: enrollmentRepo,
		lessonRepo:     lessonRepo,
		courseRepo:     courseRepo,
	}
}

// GetCourseProgress retrieves detailed progress for a user in a course
func (s *ProgressTrackingService) GetCourseProgress(ctx context.Context, userID, courseID uuid.UUID) (map[string]interface{}, error) {
	// Get enrollment
	enrollment, err := s.enrollmentRepo.GetByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}

	// Get all lessons for the course
	lessons, err := s.lessonRepo.GetByCourseID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lessons: %w", err)
	}

	// Get progress for each lesson
	progressList := make([]map[string]interface{}, 0)
	completedLessons := 0
	totalVideoPosition := 0
	totalVideoDuration := 0

	for _, lesson := range lessons {
		progress, err := s.progressRepo.GetByUserAndLesson(ctx, userID, lesson.ID)
		lessonProgress := map[string]interface{}{
			"lesson_id":      lesson.ID,
			"title":          lesson.Title,
			"is_completed":   false,
			"completed_at":   nil,
			"video_position": 0,
			"video_duration": lesson.VideoDuration,
		}

		if err == nil && progress != nil {
			if progress.IsCompleted {
				completedLessons++
				lessonProgress["is_completed"] = true
				lessonProgress["completed_at"] = progress.CompletedAt
			}
			lessonProgress["video_position"] = progress.VideoPosition
			totalVideoPosition += progress.VideoPosition
		}

		totalVideoDuration += lesson.VideoDuration
		progressList = append(progressList, lessonProgress)
	}

	// Calculate progress percentage
	totalLessons := len(lessons)
	progressPercentage := 0.0
	if totalLessons > 0 {
		progressPercentage = float64(completedLessons) / float64(totalLessons) * 100
	}

	// Calculate video progress percentage
	videoProgressPercentage := 0.0
	if totalVideoDuration > 0 {
		videoProgressPercentage = float64(totalVideoPosition) / float64(totalVideoDuration) * 100
	}

	return map[string]interface{}{
		"course_id":                 courseID,
		"enrollment_id":             enrollment.ID,
		"progress_percentage":       progressPercentage,
		"video_progress_percentage": videoProgressPercentage,
		"completed_lessons":         completedLessons,
		"total_lessons":             totalLessons,
		"lessons":                   progressList,
		"last_accessed_at":          enrollment.UpdatedAt,
	}, nil
}

// UpdateVideoProgress updates video playback position in real-time
func (s *ProgressTrackingService) UpdateVideoProgress(ctx context.Context, userID, lessonID uuid.UUID, position int, duration int) error {
	// Get or create progress record
	progress, err := s.progressRepo.GetByUserAndLesson(ctx, userID, lessonID)
	now := time.Now()

	if err != nil {
		// Create new progress record
		progress = &models.Progress{
			ID:             uuid.New(),
			UserID:         userID,
			LessonID:       lessonID,
			EnrollmentID:   uuid.Nil, // TODO: Get actual enrollment ID
			IsCompleted:    false,
			IsWatching:     true,
			VideoPosition:  position,
			LastAccessedAt: now,
		}
		return s.progressRepo.Create(ctx, progress)
	}

	// Update existing progress
	progress.VideoPosition = position
	progress.IsWatching = true
	progress.LastAccessedAt = now

	// Auto-complete if video is watched to 90% or more
	if duration > 0 && position >= int(float64(duration)*0.9) {
		progress.IsCompleted = true
		if progress.CompletedAt == nil {
			progress.CompletedAt = &now
		}
	}

	return s.progressRepo.Update(ctx, progress)
}

// MarkLessonCompleted marks a lesson as completed
func (s *ProgressTrackingService) MarkLessonCompleted(ctx context.Context, userID, lessonID uuid.UUID) error {
	return s.progressRepo.MarkCompleted(ctx, userID, lessonID)
}

// GetLearningHeatmapData retrieves learning data for heatmap visualization
func (s *ProgressTrackingService) GetLearningHeatmapData(ctx context.Context, userID uuid.UUID, months int) ([]HeatmapData, error) {
	// Get all progress records for the user
	progresses, err := s.progressRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get progress records: %w", err)
	}

	// Aggregate learning activity by date
	dateMap := make(map[string]int) // date -> minutes spent
	cutoffDate := time.Now().AddDate(0, -months, 0)

	for _, progress := range progresses {
		if progress.LastAccessedAt.Before(cutoffDate) {
			continue
		}

		dateStr := progress.LastAccessedAt.Format("2006-01-02")

		// Estimate time spent based on video position and lesson completion
		minutesSpent := 5 // Default minimum
		if progress.VideoPosition > 0 {
			minutesSpent = progress.VideoPosition / 60
		}
		if progress.IsCompleted {
			minutesSpent += 10 // Bonus for completion
		}

		dateMap[dateStr] += minutesSpent
	}

	// Convert to heatmap data
	heatmapData := make([]HeatmapData, 0)
	maxMinutes := 0

	for dateStr, minutes := range dateMap {
		if minutes > maxMinutes {
			maxMinutes = minutes
		}
		heatmapData = append(heatmapData, HeatmapData{
			Date:  dateStr,
			Count: minutes,
			Level: 0, // Will be calculated below
		})
	}

	// Calculate intensity levels (0-4)
	for i := range heatmapData {
		if maxMinutes > 0 {
			ratio := float64(heatmapData[i].Count) / float64(maxMinutes)
			if ratio >= 0.8 {
				heatmapData[i].Level = 4
			} else if ratio >= 0.6 {
				heatmapData[i].Level = 3
			} else if ratio >= 0.4 {
				heatmapData[i].Level = 2
			} else if ratio >= 0.2 {
				heatmapData[i].Level = 1
			}
		}
	}

	return heatmapData, nil
}

// GetDailyLearningStats retrieves daily learning statistics
func (s *ProgressTrackingService) GetDailyLearningStats(ctx context.Context, userID uuid.UUID, days int) ([]DailyLearningStats, error) {
	progresses, err := s.progressRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get progress records: %w", err)
	}

	// Aggregate by date
	dateStats := make(map[string]*DailyLearningStats)
	cutoffDate := time.Now().AddDate(0, 0, -days)

	for _, progress := range progresses {
		if progress.LastAccessedAt.Before(cutoffDate) {
			continue
		}

		dateStr := progress.LastAccessedAt.Format("2006-01-02")

		if _, exists := dateStats[dateStr]; !exists {
			dateStats[dateStr] = &DailyLearningStats{
				Date: dateStr,
			}
		}

		// Calculate time spent
		minutesSpent := 5
		if progress.VideoPosition > 0 {
			minutesSpent = progress.VideoPosition / 60
		}
		if progress.IsCompleted {
			minutesSpent += 10
			dateStats[dateStr].LessonsCompleted++
		}

		dateStats[dateStr].TotalSeconds += minutesSpent * 60
		dateStats[dateStr].CoursesAccessed++
	}

	// Convert to slice
	stats := make([]DailyLearningStats, 0, len(dateStats))
	for _, stat := range dateStats {
		stats = append(stats, *stat)
	}

	return stats, nil
}

// GenerateWeeklyReport generates a weekly learning report
func (s *ProgressTrackingService) GenerateWeeklyReport(ctx context.Context, userID uuid.UUID, weekOffset int) (*WeeklyReport, error) {
	now := time.Now()
	// Calculate week start (Monday)
	weekStart := now.AddDate(0, 0, -int(now.Weekday()-1)-weekOffset*7)
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, now.Location())
	weekEnd := weekStart.AddDate(0, 0, 6)

	// Get daily stats for the week
	dailyStats, err := s.GetDailyLearningStats(ctx, userID, 7+weekOffset*7)
	if err != nil {
		return nil, err
	}

	// Filter to only this week
	weekDailyStats := make([]DailyLearningStats, 0)
	totalSeconds := 0

	for _, stat := range dailyStats {
		statDate, _ := time.Parse("2006-01-02", stat.Date)
		if !statDate.Before(weekStart) && !statDate.After(weekEnd) {
			weekDailyStats = append(weekDailyStats, stat)
			totalSeconds += stat.TotalSeconds
		}
	}

	// Get enrollments to calculate course progress
	enrollments, err := s.enrollmentRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	coursesProgress := make([]CourseProgressSummary, 0)
	totalLessonsCompleted := 0

	for _, enrollment := range enrollments {
		course, err := s.courseRepo.GetByID(ctx, enrollment.CourseID)
		if err != nil {
			continue
		}

		lessons, err := s.lessonRepo.GetByCourseID(ctx, enrollment.CourseID)
		if err != nil {
			continue
		}

		completedLessons := 0
		for _, lesson := range lessons {
			progress, err := s.progressRepo.GetByUserAndLesson(ctx, userID, lesson.ID)
			if err == nil && progress != nil && progress.IsCompleted {
				if progress.CompletedAt != nil &&
					!progress.CompletedAt.Before(weekStart) &&
					!progress.CompletedAt.After(weekEnd) {
					completedLessons++
				}
			}
		}

		totalLessonsCompleted += completedLessons

		coursesProgress = append(coursesProgress, CourseProgressSummary{
			CourseID:         enrollment.CourseID,
			CourseTitle:      course.Title,
			ProgressPercent:  enrollment.ProgressPercentage,
			LessonsCompleted: completedLessons,
			TotalLessons:     len(lessons),
			TimeSpentMinutes: totalSeconds / 60,
		})
	}

	// Calculate average daily minutes
	avgDailyMinutes := 0.0
	if len(weekDailyStats) > 0 {
		totalMinutes := 0
		for _, stat := range weekDailyStats {
			totalMinutes += stat.TotalSeconds / 60
		}
		avgDailyMinutes = float64(totalMinutes) / 7.0
	}

	return &WeeklyReport{
		WeekStart:        weekStart.Format("2006-01-02"),
		WeekEnd:          weekEnd.Format("2006-01-02"),
		TotalHours:       float64(totalSeconds) / 3600.0,
		LessonsCompleted: totalLessonsCompleted,
		CoursesProgress:  coursesProgress,
		DailyStats:       weekDailyStats,
		AvgDailyMinutes:  avgDailyMinutes,
	}, nil
}

// GenerateMonthlyReport generates a monthly learning report
func (s *ProgressTrackingService) GenerateMonthlyReport(ctx context.Context, userID uuid.UUID, monthOffset int) (*MonthlyReport, error) {
	now := time.Now()
	// Calculate month start
	year := now.Year()
	month := int(now.Month()) - monthOffset
	if month <= 0 {
		month += 12
		year--
	}
	monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Second)

	// Get daily stats for the month
	dailyStats, err := s.GetDailyLearningStats(ctx, userID, 31+monthOffset*30)
	if err != nil {
		return nil, err
	}

	// Filter to only this month
	monthDailyStats := make([]DailyLearningStats, 0)
	totalSeconds := 0
	bestDay := ""
	bestDayMinutes := 0.0

	for _, stat := range dailyStats {
		statDate, _ := time.Parse("2006-01-02", stat.Date)
		if !statDate.Before(monthStart) && !statDate.After(monthEnd) {
			monthDailyStats = append(monthDailyStats, stat)
			totalSeconds += stat.TotalSeconds

			dayMinutes := float64(stat.TotalSeconds) / 60.0
			if dayMinutes > bestDayMinutes {
				bestDay = stat.Date
				bestDayMinutes = dayMinutes
			}
		}
	}

	// Get enrollments
	enrollments, err := s.enrollmentRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	coursesProgress := make([]CourseProgressSummary, 0)
	totalLessonsCompleted := 0
	coursesCompleted := 0

	for _, enrollment := range enrollments {
		course, err := s.courseRepo.GetByID(ctx, enrollment.CourseID)
		if err != nil {
			continue
		}

		lessons, err := s.lessonRepo.GetByCourseID(ctx, enrollment.CourseID)
		if err != nil {
			continue
		}

		completedLessons := 0
		for _, lesson := range lessons {
			progress, err := s.progressRepo.GetByUserAndLesson(ctx, userID, lesson.ID)
			if err == nil && progress != nil && progress.IsCompleted {
				if progress.CompletedAt != nil &&
					!progress.CompletedAt.Before(monthStart) &&
					!progress.CompletedAt.After(monthEnd) {
					completedLessons++
				}
			}
		}

		totalLessonsCompleted += completedLessons

		if enrollment.Status == "completed" && enrollment.CompletedAt != nil {
			if !enrollment.CompletedAt.Before(monthStart) && !enrollment.CompletedAt.After(monthEnd) {
				coursesCompleted++
			}
		}

		coursesProgress = append(coursesProgress, CourseProgressSummary{
			CourseID:         enrollment.CourseID,
			CourseTitle:      course.Title,
			ProgressPercent:  enrollment.ProgressPercentage,
			LessonsCompleted: completedLessons,
			TotalLessons:     len(lessons),
			TimeSpentMinutes: totalSeconds / 60,
		})
	}

	// Calculate average daily minutes
	avgDailyMinutes := 0.0
	daysInMonth := monthEnd.Day()
	if len(monthDailyStats) > 0 {
		totalMinutes := 0
		for _, stat := range monthDailyStats {
			totalMinutes += stat.TotalSeconds / 60
		}
		avgDailyMinutes = float64(totalMinutes) / float64(daysInMonth)
	}

	return &MonthlyReport{
		Month:            monthStart.Format("January"),
		Year:             year,
		TotalHours:       float64(totalSeconds) / 3600.0,
		LessonsCompleted: totalLessonsCompleted,
		CoursesCompleted: coursesCompleted,
		CoursesProgress:  coursesProgress,
		DailyStats:       monthDailyStats,
		AvgDailyMinutes:  avgDailyMinutes,
		BestDay:          bestDay,
		BestDayMinutes:   bestDayMinutes,
	}, nil
}

// ExportReportToCSV exports learning report to CSV format
func (s *ProgressTrackingService) ExportReportToCSV(report *MonthlyReport) (string, error) {
	// Create CSV in memory
	data := [][]string{
		{"AI Learning Platform - Monthly Learning Report"},
		{fmt.Sprintf("Month: %s %d", report.Month, report.Year)},
		{""},
		{"Summary"},
		{fmt.Sprintf("Total Hours: %.2f", report.TotalHours)},
		{fmt.Sprintf("Lessons Completed: %d", report.LessonsCompleted)},
		{fmt.Sprintf("Courses Completed: %d", report.CoursesCompleted)},
		{fmt.Sprintf("Average Daily Minutes: %.2f", report.AvgDailyMinutes)},
		{fmt.Sprintf("Best Day: %s (%.2f minutes)", report.BestDay, report.BestDayMinutes)},
		{""},
		{"Daily Statistics"},
		{"Date", "Total Seconds", "Lessons Completed", "Courses Accessed"},
	}

	for _, stat := range report.DailyStats {
		data = append(data, []string{
			stat.Date,
			fmt.Sprintf("%d", stat.TotalSeconds),
			fmt.Sprintf("%d", stat.LessonsCompleted),
			fmt.Sprintf("%d", stat.CoursesAccessed),
		})
	}

	data = append(data, []string{""}, []string{"Course Progress"})
	data = append(data, []string{"Course Title", "Progress %", "Lessons Completed", "Total Lessons", "Time Spent (min)"})

	for _, course := range report.CoursesProgress {
		data = append(data, []string{
			course.CourseTitle,
			fmt.Sprintf("%.2f", course.ProgressPercent),
			fmt.Sprintf("%d", course.LessonsCompleted),
			fmt.Sprintf("%d", course.TotalLessons),
			fmt.Sprintf("%d", course.TimeSpentMinutes),
		})
	}

	// Write CSV
	return s.writeCSV(data)
}

func (s *ProgressTrackingService) writeCSV(data [][]string) (string, error) {
	// In a real implementation, this would write to a file or return bytes
	// For now, we'll return a simple string representation
	result := ""
	for _, row := range data {
		for i, col := range row {
			if i > 0 {
				result += ","
			}
			result += col
		}
		result += "\n"
	}
	return result, nil
}

// GetLearningTimeStats retrieves total learning time statistics
func (s *ProgressTrackingService) GetLearningTimeStats(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	progresses, err := s.progressRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalSeconds := 0
	totalLessons := 0
	completedLessons := 0

	for _, progress := range progresses {
		if progress.VideoPosition > 0 {
			totalSeconds += progress.VideoPosition
		}
		totalLessons++
		if progress.IsCompleted {
			completedLessons++
		}
	}

	return map[string]interface{}{
		"total_learning_seconds": totalSeconds,
		"total_learning_hours":   float64(totalSeconds) / 3600.0,
		"total_lessons":          totalLessons,
		"completed_lessons":      completedLessons,
		"completion_rate":        float64(completedLessons) / float64(totalLessons) * 100,
	}, nil
}
