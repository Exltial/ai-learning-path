package services

import (
	"context"
	"time"
	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"

	"github.com/google/uuid"
)

// ProgressService handles progress-related business logic
type ProgressService struct {
	progressRepo   *repository.ProgressRepository
	enrollmentRepo *repository.EnrollmentRepository
	lessonRepo     *repository.LessonRepository
}

// NewProgressService creates a new ProgressService
func NewProgressService(progressRepo *repository.ProgressRepository, enrollmentRepo *repository.EnrollmentRepository, lessonRepo *repository.LessonRepository) *ProgressService {
	return &ProgressService{
		progressRepo:   progressRepo,
		enrollmentRepo: enrollmentRepo,
		lessonRepo:     lessonRepo,
	}
}

// GetCourseProgress retrieves progress for a user in a course
func (s *ProgressService) GetCourseProgress(ctx context.Context, userID, courseID uuid.UUID) (map[string]interface{}, error) {
	// Get enrollment
	enrollment, err := s.enrollmentRepo.GetByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}

	// Get all lessons for the course
	lessons, err := s.lessonRepo.GetByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	// Get progress for each lesson
	progressList := make([]map[string]interface{}, 0)
	completedLessons := 0

	for _, lesson := range lessons {
		progress, err := s.progressRepo.GetByUserAndLesson(ctx, userID, lesson.ID)
		if err != nil {
			// Progress doesn't exist yet
			progressList = append(progressList, map[string]interface{}{
				"lesson_id":     lesson.ID,
				"title":         lesson.Title,
				"is_completed":  false,
				"completed_at":  nil,
			})
		} else {
			if progress.IsCompleted {
				completedLessons++
			}
			progressList = append(progressList, map[string]interface{}{
				"lesson_id":     progress.LessonID,
				"title":         lesson.Title,
				"is_completed":  progress.IsCompleted,
				"completed_at":  progress.CompletedAt,
			})
		}
	}

	// Calculate progress percentage
	totalLessons := len(lessons)
	progressPercentage := 0.0
	if totalLessons > 0 {
		progressPercentage = float64(completedLessons) / float64(totalLessons) * 100
	}

	return map[string]interface{}{
		"course_id":          courseID,
		"enrollment_id":      enrollment.ID,
		"progress_percentage": progressPercentage,
		"completed_lessons":  completedLessons,
		"total_lessons":      totalLessons,
		"lessons":            progressList,
	}, nil
}

// GetUserProgress retrieves all progress for a user
func (s *ProgressService) GetUserProgress(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	// Get all enrollments for the user
	enrollments, err := s.enrollmentRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	activeCourses := make([]map[string]interface{}, 0)
	completedCourses := make([]map[string]interface{}, 0)

	for _, enrollment := range enrollments {
		// Get course details (TODO: implement course retrieval)
		courseInfo := map[string]interface{}{
			"course_id":          enrollment.CourseID,
			"progress_percentage": enrollment.ProgressPercentage,
			"enrolled_at":        enrollment.EnrolledAt,
		}

		if enrollment.Status == "completed" {
			courseInfo["completed_at"] = enrollment.CompletedAt
			completedCourses = append(completedCourses, courseInfo)
		} else {
			activeCourses = append(activeCourses, courseInfo)
		}
	}

	return map[string]interface{}{
		"active_courses":   activeCourses,
		"completed_courses": completedCourses,
	}, nil
}

// UpdateLessonProgress updates progress for a lesson
func (s *ProgressService) UpdateLessonProgress(ctx context.Context, userID, lessonID uuid.UUID, isCompleted bool, videoPosition int) (*models.Progress, error) {
	// Try to get existing progress
	progress, err := s.progressRepo.GetByUserAndLesson(ctx, userID, lessonID)
	
	now := time.Now()
	
	if err != nil {
		// Create new progress record
		// TODO: Get enrollment ID
		progress = &models.Progress{
			ID:             uuid.New(),
			UserID:         userID,
			LessonID:       lessonID,
			EnrollmentID:   uuid.Nil, // TODO: Get actual enrollment ID
			IsCompleted:    isCompleted,
			IsWatching:     !isCompleted && videoPosition > 0,
			VideoPosition:  videoPosition,
			LastAccessedAt: now,
		}
		
		if isCompleted {
			progress.CompletedAt = &now
		}
		
		if err := s.progressRepo.Create(ctx, progress); err != nil {
			return nil, err
		}
	} else {
		// Update existing progress
		progress.IsCompleted = isCompleted
		progress.VideoPosition = videoPosition
		progress.LastAccessedAt = now
		
		if isCompleted && !progress.IsCompleted {
			progress.CompletedAt = &now
		}
		
		if err := s.progressRepo.Update(ctx, progress); err != nil {
			return nil, err
		}
	}

	return progress, nil
}

// MarkLessonCompleted marks a lesson as completed
func (s *ProgressService) MarkLessonCompleted(ctx context.Context, userID, lessonID uuid.UUID) error {
	return s.progressRepo.MarkCompleted(ctx, userID, lessonID)
}

// UpdateVideoPosition updates video playback position
func (s *ProgressService) UpdateVideoPosition(ctx context.Context, userID, lessonID uuid.UUID, position int) error {
	return s.progressRepo.UpdateVideoPosition(ctx, userID, lessonID, position)
}
