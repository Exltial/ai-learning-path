package services

import (
	"context"
	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"

	"github.com/google/uuid"
)

// ExerciseService handles exercise-related business logic
type ExerciseService struct {
	exerciseRepo *repository.ExerciseRepository
}

// NewExerciseService creates a new ExerciseService
func NewExerciseService(exerciseRepo *repository.ExerciseRepository) *ExerciseService {
	return &ExerciseService{exerciseRepo: exerciseRepo}
}

// GetExercise retrieves an exercise by ID
func (s *ExerciseService) GetExercise(ctx context.Context, id uuid.UUID) (*models.Exercise, error) {
	return s.exerciseRepo.GetByID(ctx, id)
}

// CreateExercise creates a new exercise
func (s *ExerciseService) CreateExercise(ctx context.Context, exercise *models.Exercise) error {
	return s.exerciseRepo.Create(ctx, exercise)
}

// UpdateExercise updates an existing exercise
func (s *ExerciseService) UpdateExercise(ctx context.Context, exercise *models.Exercise) error {
	return s.exerciseRepo.Update(ctx, exercise)
}

// DeleteExercise deletes an exercise
func (s *ExerciseService) DeleteExercise(ctx context.Context, id uuid.UUID) error {
	return s.exerciseRepo.Delete(ctx, id)
}

// GetExerciseSubmissions retrieves submissions for an exercise
func (s *ExerciseService) GetExerciseSubmissions(ctx context.Context, exerciseID, userID uuid.UUID) ([]*models.Submission, error) {
	// This should be handled by submission service
	return []*models.Submission{}, nil
}
