package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"

	"github.com/google/uuid"
)

// Error definitions for submission service
var (
	ErrSubmissionNotFound   = errors.New("submission not found")
	ErrExerciseNotFound     = errors.New("exercise not found")
	ErrMaxAttemptsReached   = errors.New("maximum attempts reached")
	ErrTimeLimitExceeded    = errors.New("time limit exceeded")
	ErrNotEnrolled          = errors.New("not enrolled in this course")
	ErrAlreadyGraded        = errors.New("submission already graded")
	ErrInvalidSubmission    = errors.New("invalid submission data")
)

// SubmissionService handles submission-related business logic
type SubmissionService struct {
	submissionRepo *repository.SubmissionRepository
	exerciseRepo   *repository.ExerciseRepository
	enrollmentRepo *repository.EnrollmentRepository
}

// NewSubmissionService creates a new SubmissionService
func NewSubmissionService(submissionRepo *repository.SubmissionRepository, exerciseRepo *repository.ExerciseRepository) *SubmissionService {
	return &SubmissionService{
		submissionRepo: submissionRepo,
		exerciseRepo:   exerciseRepo,
	}
}

// SubmitExercise submits an answer for an exercise
// Handles different exercise types: multiple_choice, coding, fill_blank, true_false, essay
func (s *SubmissionService) SubmitExercise(ctx context.Context, exerciseID, userID uuid.UUID, answer, code string, selectedOptions []string) (*models.Submission, error) {
	// Get exercise details
	exercise, err := s.exerciseRepo.GetByID(ctx, exerciseID)
	if err != nil {
		return nil, ErrExerciseNotFound
	}

	// Validate submission based on exercise type
	if err := s.validateSubmission(exercise, answer, code, selectedOptions); err != nil {
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
		return nil, ErrMaxAttemptsReached
	}

	// Check time limit if applicable
	if exercise.TimeLimit != nil && *exercise.TimeLimit > 0 {
		// TODO: Implement time limit check based on first attempt time
	}

	// Determine submission type and auto-grade if applicable
	submissionType := exercise.ExerciseType
	var isCorrect bool
	var score float64
	var feedback string

	// Auto-grade based on exercise type
	switch exercise.ExerciseType {
	case "multiple_choice":
		isCorrect, score, feedback = s.gradeMultipleChoice(exercise, selectedOptions)
	case "true_false":
		isCorrect, score, feedback = s.gradeTrueFalse(exercise, answer)
	case "fill_blank":
		isCorrect, score, feedback = s.gradeFillBlank(exercise, answer)
	case "coding":
		// Code submissions require test case execution
		isCorrect, score, feedback = s.gradeCoding(exercise, code)
	case "essay":
		// Essay requires manual grading
		isCorrect = false
		score = 0
		feedback = "Essay submitted, awaiting instructor review"
	}

	// Create submission record
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
	submission, err := s.submissionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrSubmissionNotFound
	}
	return submission, nil
}

// GetSubmissions retrieves submissions for an exercise by a user
func (s *SubmissionService) GetSubmissions(ctx context.Context, exerciseID, userID uuid.UUID) ([]*models.Submission, error) {
	submissions, err := s.submissionRepo.GetByExerciseIDAndUserID(ctx, exerciseID, userID)
	if err != nil {
		return nil, err
	}
	return submissions, nil
}

// GradeSubmission grades a submission (for essay type or manual review by instructors)
func (s *SubmissionService) GradeSubmission(ctx context.Context, submissionID uuid.UUID, score float64, feedback string, isCorrect bool, gradedBy uuid.UUID) error {
	submission, err := s.submissionRepo.GetByID(ctx, submissionID)
	if err != nil {
		return ErrSubmissionNotFound
	}

	// Check if already graded
	if submission.GradedAt != nil {
		return ErrAlreadyGraded
	}

	// Validate score
	if score < 0 {
		return errors.New("score cannot be negative")
	}
	if submission.Score != nil && *submission.Score > 0 {
		// If exercise has a max score, validate against it
		// TODO: Get max score from exercise
	}

	// Update submission
	now := time.Now()
	submission.IsCorrect = &isCorrect
	submission.Score = &score
	submission.Feedback = feedback
	submission.GradedAt = &now
	submission.GradedBy = &gradedBy

	return s.submissionRepo.Update(ctx, submission)
}

// GetLatestSubmission retrieves the latest submission for an exercise by a user
func (s *SubmissionService) GetLatestSubmission(ctx context.Context, exerciseID, userID uuid.UUID) (*models.Submission, error) {
	submission, err := s.submissionRepo.GetLatestSubmission(ctx, exerciseID, userID)
	if err != nil {
		return nil, ErrSubmissionNotFound
	}
	return submission, nil
}

// GetSubmissionStats retrieves statistics for a user's submissions
func (s *SubmissionService) GetSubmissionStats(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	// TODO: Implement submission statistics
	return map[string]interface{}{
		"total_submissions":   0,
		"correct_submissions": 0,
		"average_score":       0.0,
	}, nil
}

// validateSubmission validates submission data based on exercise type
func (s *SubmissionService) validateSubmission(exercise *models.Exercise, answer, code string, selectedOptions []string) error {
	switch exercise.ExerciseType {
	case "multiple_choice":
		if len(selectedOptions) == 0 {
			return errors.New("at least one option must be selected")
		}
	case "true_false":
		if answer == "" {
			return errors.New("answer is required for true/false questions")
		}
		lowerAnswer := strings.ToLower(answer)
		if lowerAnswer != "true" && lowerAnswer != "false" && lowerAnswer != "t" && lowerAnswer != "f" {
			return errors.New("answer must be true or false")
		}
	case "fill_blank":
		if strings.TrimSpace(answer) == "" {
			return errors.New("answer cannot be empty")
		}
	case "coding":
		if strings.TrimSpace(code) == "" {
			return errors.New("code cannot be empty")
		}
	case "essay":
		if strings.TrimSpace(answer) == "" {
			return errors.New("essay answer cannot be empty")
		}
		if len(answer) < 50 {
			return errors.New("essay must be at least 50 characters")
		}
	}
	return nil
}

// gradeMultipleChoice auto-grades multiple choice submissions
func (s *SubmissionService) gradeMultipleChoice(exercise *models.Exercise, selectedOptions []string) (bool, float64, string) {
	if exercise.Options == nil || len(exercise.Options) == 0 {
		return false, 0, "Exercise has no options configured"
	}

	// Count correct selections
	correctCount := 0
	incorrectCount := 0

	for _, selected := range selectedOptions {
		found := false
		for _, option := range exercise.Options {
			if option.Text == selected {
				found = true
				if option.IsCorrect {
					correctCount++
				} else {
					incorrectCount++
				}
				break
			}
		}
		if !found {
			incorrectCount++
		}
	}

	// Check if all correct options were selected and no incorrect ones
	allCorrectSelected := correctCount > 0 && incorrectCount == 0
	
	// Calculate score (percentage of correct options selected)
	var score float64
	if len(exercise.Options) > 0 {
		score = float64(correctCount) / float64(len(exercise.Options)) * 100
	}

	if allCorrectSelected {
		return true, score, "Correct! All right answers selected."
	}
	return false, score, "Incorrect. Review the options and try again."
}

// gradeTrueFalse auto-grades true/false submissions
func (s *SubmissionService) gradeTrueFalse(exercise *models.Exercise, answer string) (bool, float64, string) {
	lowerAnswer := strings.ToLower(strings.TrimSpace(answer))
	
	// Determine user's answer
	userAnswer := false
	if lowerAnswer == "true" || lowerAnswer == "t" {
		userAnswer = true
	}

	// Get correct answer from exercise
	correctAnswer := false
	if exercise.ExpectedAnswer != nil {
		if expectedStr, ok := exercise.ExpectedAnswer["answer"].(string); ok {
			correctAnswer = strings.ToLower(expectedStr) == "true"
		}
	}

	if userAnswer == correctAnswer {
		return true, 100, "Correct!"
	}
	return false, 0, "Incorrect. The correct answer is different."
}

// gradeFillBlank auto-grades fill-in-the-blank submissions
func (s *SubmissionService) gradeFillBlank(exercise *models.Exercise, answer string) (bool, float64, string) {
	if exercise.ExpectedAnswer == nil {
		return false, 0, "Exercise has no expected answer configured"
	}

	// Get expected answers (can have multiple valid answers)
	var expectedAnswers []string
	if answers, ok := exercise.ExpectedAnswer["answers"].([]interface{}); ok {
		for _, a := range answers {
			if str, ok := a.(string); ok {
				expectedAnswers = append(expectedAnswers, strings.ToLower(strings.TrimSpace(str)))
			}
		}
	} else if answer, ok := exercise.ExpectedAnswer["answer"].(string); ok {
		expectedAnswers = []string{strings.ToLower(strings.TrimSpace(answer))}
	}

	userAnswer := strings.ToLower(strings.TrimSpace(answer))

	for _, expected := range expectedAnswers {
		if userAnswer == expected {
			return true, 100, "Correct!"
		}
	}

	return false, 0, "Incorrect. Please try again."
}

// gradeCoding auto-grades coding submissions using test cases
func (s *SubmissionService) gradeCoding(exercise *models.Exercise, code string) (bool, float64, string) {
	// TODO: Implement code execution and test case validation
	// This would typically involve:
	// 1. Running the code in a sandbox environment
	// 2. Comparing output against test cases
	// 3. Checking for correctness and efficiency

	// Placeholder implementation
	if exercise.TestCases == nil {
		return false, 0, "Code submitted, awaiting evaluation (no test cases configured)"
	}

	// For now, return pending evaluation
	return false, 0, "Code submitted, awaiting automated evaluation"
}

// ResubmitExercise allows a user to resubmit an exercise
func (s *SubmissionService) ResubmitExercise(ctx context.Context, exerciseID, userID uuid.UUID, answer, code string, selectedOptions []string) (*models.Submission, error) {
	// Get exercise
	exercise, err := s.exerciseRepo.GetByID(ctx, exerciseID)
	if err != nil {
		return nil, ErrExerciseNotFound
	}

	// Check if resubmission is allowed
	if exercise.MaxAttempts > 0 {
		previousSubmissions, err := s.submissionRepo.GetByExerciseIDAndUserID(ctx, exerciseID, userID)
		if err != nil {
			return nil, err
		}
		if len(previousSubmissions) >= exercise.MaxAttempts {
			return nil, ErrMaxAttemptsReached
		}
	}

	// Submit new attempt
	return s.SubmitExercise(ctx, exerciseID, userID, answer, code, selectedOptions)
}
