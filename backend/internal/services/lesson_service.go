package services

import (
	"context"
	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"

	"github.com/google/uuid"
)

// LessonService handles lesson-related business logic
type LessonService struct {
	lessonRepo *repository.LessonRepository
}

// NewLessonService creates a new LessonService
func NewLessonService(lessonRepo *repository.LessonRepository) *LessonService {
	return &LessonService{lessonRepo: lessonRepo}
}

// GetLesson retrieves a lesson by ID
func (s *LessonService) GetLesson(ctx context.Context, id uuid.UUID) (*models.Lesson, error) {
	return s.lessonRepo.GetByID(ctx, id)
}

// CreateLesson creates a new lesson
func (s *LessonService) CreateLesson(ctx context.Context, lesson *models.Lesson) error {
	return s.lessonRepo.Create(ctx, lesson)
}

// UpdateLesson updates an existing lesson
func (s *LessonService) UpdateLesson(ctx context.Context, lesson *models.Lesson) error {
	return s.lessonRepo.Update(ctx, lesson)
}

// DeleteLesson deletes a lesson
func (s *LessonService) DeleteLesson(ctx context.Context, id uuid.UUID) error {
	return s.lessonRepo.Delete(ctx, id)
}

// GetLessonExercises retrieves all exercises for a lesson
func (s *LessonService) GetLessonExercises(ctx context.Context, lessonID uuid.UUID) ([]*models.Exercise, error) {
	// TODO: Add exercise repo to lesson service or create separate method
	return []*models.Exercise{}, nil
}
