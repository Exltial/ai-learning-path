package services

import (
	"context"
	"time"
	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"

	"github.com/google/uuid"
)

// SubmissionService handles submission-related business logic
type SubmissionService struct {
	submissionRepo *repository.SubmissionRepository
	exerciseRepo   *repository.ExerciseRepository
}

// NewSubmissionService creates a new SubmissionService
func NewSubmissionService(submissionRepo *repository.SubmissionRepository, exerciseRepo *repository.ExerciseRepository) *SubmissionService {
	return &SubmissionService{
		submissionRepo: submissionRepo,
		exerciseRepo:   exerciseRepo,
	}
}

// SubmitExercise submits an answer for an exercise
func (s *SubmissionService) SubmitExercise(ctx context.Context, exerciseID, userID uuid.UUID, answer, code string, selectedOptions []string) (*models.Submission, error) {
	// Get exercise to determine type and get correct answer
	exercise, err := s.exerciseRepo.GetByID(ctx, exerciseID)
	if err != nil {
		return nil, err
	}

	// Get previous submissions to determine attempt number
	previousSubmissions, err := s.submissionRepo.GetByExerciseIDAndUserID(ctx, exerciseID, userID)
	if err != nil {
		return nil, err
	}

	attemptNumber := 1
	if len(previousSubmissions) > 0 {
		attemptNumber = previousSubmissions[0].AttemptNumber + 1
	}

	// Check if max attempts reached
	if exercise.MaxAttempts > 0 && attemptNumber > exercise.MaxAttempts {
		return nil, nil // TODO: Return error for max attempts reached
	}

	// Determine submission type and correctness
	submissionType := exercise.ExerciseType
	var isCorrect bool
	var score float64
	var feedback string

	// Auto-grade based on exercise type
	switch exercise.ExerciseType {
	case "multiple_choice", "true_false":
		// TODO: Implement multiple choice grading
		isCorrect = false
		score = 0
		feedback = "Answer submitted"
	case "coding":
		// TODO: Implement code grading with test cases
		isCorrect = false
		score = 0
		feedback = "Code submitted, awaiting evaluation"
	case "fill_blank":
		// TODO: Implement fill-in-the-blank grading
		isCorrect = false
		score = 0
		feedback = "Answer submitted"
	case "essay":
		// Essay requires manual grading
		isCorrect = false
		score = 0
		feedback = "Essay submitted, awaiting instructor review"
	}

	// Create submission
	submission := &models.Submission{
		ID:             uuid.New(),
		ExerciseID:     exerciseID,
		UserID:         userID,
		SubmissionType: submissionType,
		Answer:         answer,
		Code:           code,
		IsCorrect:      &isCorrect,
		Score:          &score,
		Feedback:       feedback,
		AttemptNumber:  attemptNumber,
		SubmittedAt:    time.Now(),
	}

	if err := s.submissionRepo.Create(ctx, submission); err != nil {
		return nil, err
	}

	return submission, nil
}

// GetSubmission retrieves a submission by ID
func (s *SubmissionService) GetSubmission(ctx context.Context, id uuid.UUID) (*models.Submission, error) {
	return s.submissionRepo.GetByID(ctx, id)
}

// GetSubmissions retrieves submissions for an exercise by a user
func (s *SubmissionService) GetSubmissions(ctx context.Context, exerciseID, userID uuid.UUID) ([]*models.Submission, error) {
	return s.submissionRepo.GetByExerciseIDAndUserID(ctx, exerciseID, userID)
}

// GradeSubmission grades a submission (for essay type or manual review)
func (s *SubmissionService) GradeSubmission(ctx context.Context, submissionID uuid.UUID, score float64, feedback string, isCorrect bool, gradedBy uuid.UUID) error {
	submission, err := s.submissionRepo.GetByID(ctx, submissionID)
	if err != nil {
		return err
	}

	now := time.Now()
	submission.IsCorrect = &isCorrect
	submission.Score = &score
	submission.Feedback = feedback
	submission.GradedAt = &now
	submission.GradedBy = &gradedBy

	return s.submissionRepo.Update(ctx, submission)
}
