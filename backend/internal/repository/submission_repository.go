package repository

import (
	"context"

	"ai-learning-platform/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SubmissionRepository handles database operations for submissions
type SubmissionRepository struct {
	db *pgxpool.Pool
}

// NewSubmissionRepository creates a new SubmissionRepository
func NewSubmissionRepository(db *pgxpool.Pool) *SubmissionRepository {
	return &SubmissionRepository{db: db}
}

// Create creates a new submission
func (r *SubmissionRepository) Create(ctx context.Context, submission *models.Submission) error {
	query := `
		INSERT INTO submissions (id, exercise_id, user_id, submission_type, answer, code, is_correct, 
			score, feedback, attempt_number, submitted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(ctx, query,
		submission.ID,
		submission.ExerciseID,
		submission.UserID,
		submission.SubmissionType,
		submission.Answer,
		submission.Code,
		submission.IsCorrect,
		submission.Score,
		submission.Feedback,
		submission.AttemptNumber,
		submission.SubmittedAt,
	)
	return err
}

// GetByID retrieves a submission by ID
func (r *SubmissionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Submission, error) {
	query := `
		SELECT id, exercise_id, user_id, submission_type, answer, code, is_correct, score, 
			feedback, attempt_number, submitted_at, graded_at, graded_by
		FROM submissions WHERE id = $1
	`
	submission := &models.Submission{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&submission.ID,
		&submission.ExerciseID,
		&submission.UserID,
		&submission.SubmissionType,
		&submission.Answer,
		&submission.Code,
		&submission.IsCorrect,
		&submission.Score,
		&submission.Feedback,
		&submission.AttemptNumber,
		&submission.SubmittedAt,
		&submission.GradedAt,
		&submission.GradedBy,
	)
	if err != nil {
		return nil, err
	}
	return submission, nil
}

// GetByExerciseIDAndUserID retrieves submissions by exercise and user
func (r *SubmissionRepository) GetByExerciseIDAndUserID(ctx context.Context, exerciseID, userID uuid.UUID) ([]*models.Submission, error) {
	query := `
		SELECT id, exercise_id, user_id, submission_type, answer, code, is_correct, score, 
			feedback, attempt_number, submitted_at, graded_at, graded_by
		FROM submissions WHERE exercise_id = $1 AND user_id = $2 ORDER BY attempt_number DESC
	`
	rows, err := r.db.Query(ctx, query, exerciseID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	submissions := make([]*models.Submission, 0)
	for rows.Next() {
		submission := &models.Submission{}
		err := rows.Scan(
			&submission.ID,
			&submission.ExerciseID,
			&submission.UserID,
			&submission.SubmissionType,
			&submission.Answer,
			&submission.Code,
			&submission.IsCorrect,
			&submission.Score,
			&submission.Feedback,
			&submission.AttemptNumber,
			&submission.SubmittedAt,
			&submission.GradedAt,
			&submission.GradedBy,
		)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, submission)
	}

	return submissions, rows.Err()
}

// Update updates an existing submission (for grading)
func (r *SubmissionRepository) Update(ctx context.Context, submission *models.Submission) error {
	query := `
		UPDATE submissions SET is_correct = $2, score = $3, feedback = $4, graded_at = $5, graded_by = $6
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		submission.ID,
		submission.IsCorrect,
		submission.Score,
		submission.Feedback,
		submission.GradedAt,
		submission.GradedBy,
	)
	return err
}

// GetLatestSubmission retrieves the latest submission for an exercise by a user
func (r *SubmissionRepository) GetLatestSubmission(ctx context.Context, exerciseID, userID uuid.UUID) (*models.Submission, error) {
	query := `
		SELECT id, exercise_id, user_id, submission_type, answer, code, is_correct, score, 
			feedback, attempt_number, submitted_at, graded_at, graded_by
		FROM submissions WHERE exercise_id = $1 AND user_id = $2 ORDER BY attempt_number DESC LIMIT 1
	`
	submission := &models.Submission{}
	err := r.db.QueryRow(ctx, query, exerciseID, userID).Scan(
		&submission.ID,
		&submission.ExerciseID,
		&submission.UserID,
		&submission.SubmissionType,
		&submission.Answer,
		&submission.Code,
		&submission.IsCorrect,
		&submission.Score,
		&submission.Feedback,
		&submission.AttemptNumber,
		&submission.SubmittedAt,
		&submission.GradedAt,
		&submission.GradedBy,
	)
	if err != nil {
		return nil, err
	}
	return submission, nil
}

// GetByUserID retrieves all submissions for a user
func (r *SubmissionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Submission, error) {
	query := `
		SELECT id, exercise_id, user_id, submission_type, answer, code, is_correct, score, 
			feedback, attempt_number, submitted_at, graded_at, graded_by
		FROM submissions WHERE user_id = $1 ORDER BY submitted_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	submissions := make([]*models.Submission, 0)
	for rows.Next() {
		submission := &models.Submission{}
		err := rows.Scan(
			&submission.ID,
			&submission.ExerciseID,
			&submission.UserID,
			&submission.SubmissionType,
			&submission.Answer,
			&submission.Code,
			&submission.IsCorrect,
			&submission.Score,
			&submission.Feedback,
			&submission.AttemptNumber,
			&submission.SubmittedAt,
			&submission.GradedAt,
			&submission.GradedBy,
		)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, submission)
	}

	return submissions, rows.Err()
}
