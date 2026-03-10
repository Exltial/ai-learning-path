package repository

import (
	"context"
	"time"

	"ai-learning-platform/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EnrollmentRepository handles database operations for enrollments
type EnrollmentRepository struct {
	db *pgxpool.Pool
}

// NewEnrollmentRepository creates a new EnrollmentRepository
func NewEnrollmentRepository(db *pgxpool.Pool) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

// Create creates a new enrollment
func (r *EnrollmentRepository) Create(ctx context.Context, enrollment *models.Enrollment) error {
	query := `
		INSERT INTO enrollments (id, user_id, course_id, enrolled_at, status, progress_percentage)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, course_id) DO UPDATE SET
			status = EXCLUDED.status,
			enrolled_at = EXCLUDED.enrolled_at
	`
	_, err := r.db.Exec(ctx, query,
		enrollment.ID,
		enrollment.UserID,
		enrollment.CourseID,
		enrollment.EnrolledAt,
		enrollment.Status,
		enrollment.ProgressPercentage,
	)
	return err
}

// GetByID retrieves an enrollment by ID
func (r *EnrollmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Enrollment, error) {
	query := `
		SELECT id, user_id, course_id, enrolled_at, completed_at, status, progress_percentage
		FROM enrollments WHERE id = $1
	`
	enrollment := &models.Enrollment{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&enrollment.ID,
		&enrollment.UserID,
		&enrollment.CourseID,
		&enrollment.EnrolledAt,
		&enrollment.CompletedAt,
		&enrollment.Status,
		&enrollment.ProgressPercentage,
	)
	if err != nil {
		return nil, err
	}
	return enrollment, nil
}

// GetByUserAndCourse retrieves an enrollment by user and course
func (r *EnrollmentRepository) GetByUserAndCourse(ctx context.Context, userID, courseID uuid.UUID) (*models.Enrollment, error) {
	query := `
		SELECT id, user_id, course_id, enrolled_at, completed_at, status, progress_percentage
		FROM enrollments WHERE user_id = $1 AND course_id = $2
	`
	enrollment := &models.Enrollment{}
	err := r.db.QueryRow(ctx, query, userID, courseID).Scan(
		&enrollment.ID,
		&enrollment.UserID,
		&enrollment.CourseID,
		&enrollment.EnrolledAt,
		&enrollment.CompletedAt,
		&enrollment.Status,
		&enrollment.ProgressPercentage,
	)
	if err != nil {
		return nil, err
	}
	return enrollment, nil
}

// GetByUser retrieves all enrollments for a user
func (r *EnrollmentRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]*models.Enrollment, error) {
	query := `
		SELECT id, user_id, course_id, enrolled_at, completed_at, status, progress_percentage
		FROM enrollments WHERE user_id = $1 ORDER BY enrolled_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	enrollments := make([]*models.Enrollment, 0)
	for rows.Next() {
		enrollment := &models.Enrollment{}
		err := rows.Scan(
			&enrollment.ID,
			&enrollment.UserID,
			&enrollment.CourseID,
			&enrollment.EnrolledAt,
			&enrollment.CompletedAt,
			&enrollment.Status,
			&enrollment.ProgressPercentage,
		)
		if err != nil {
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
	}

	return enrollments, rows.Err()
}

// GetByCourse retrieves all enrollments for a course
func (r *EnrollmentRepository) GetByCourse(ctx context.Context, courseID uuid.UUID) ([]*models.Enrollment, error) {
	query := `
		SELECT id, user_id, course_id, enrolled_at, completed_at, status, progress_percentage
		FROM enrollments WHERE course_id = $1 ORDER BY enrolled_at DESC
	`
	rows, err := r.db.Query(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	enrollments := make([]*models.Enrollment, 0)
	for rows.Next() {
		enrollment := &models.Enrollment{}
		err := rows.Scan(
			&enrollment.ID,
			&enrollment.UserID,
			&enrollment.CourseID,
			&enrollment.EnrolledAt,
			&enrollment.CompletedAt,
			&enrollment.Status,
			&enrollment.ProgressPercentage,
		)
		if err != nil {
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
	}

	return enrollments, rows.Err()
}

// Update updates an existing enrollment
func (r *EnrollmentRepository) Update(ctx context.Context, enrollment *models.Enrollment) error {
	query := `
		UPDATE enrollments SET status = $2, progress_percentage = $3, completed_at = $4
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		enrollment.ID,
		enrollment.Status,
		enrollment.ProgressPercentage,
		enrollment.CompletedAt,
	)
	return err
}

// UpdateProgress updates enrollment progress
func (r *EnrollmentRepository) UpdateProgress(ctx context.Context, id uuid.UUID, progressPercentage float64) error {
	query := `UPDATE enrollments SET progress_percentage = $2 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, progressPercentage)
	return err
}

// MarkCompleted marks an enrollment as completed
func (r *EnrollmentRepository) MarkCompleted(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	query := `UPDATE enrollments SET status = 'completed', completed_at = $2 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, now)
	return err
}

// Exists checks if an enrollment exists
func (r *EnrollmentRepository) Exists(ctx context.Context, userID, courseID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM enrollments WHERE user_id = $1 AND course_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, courseID).Scan(&exists)
	return exists, err
}

// GetByUserID retrieves all enrollments for a user
func (r *EnrollmentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Enrollment, error) {
	query := `
		SELECT id, user_id, course_id, enrolled_at, completed_at, status, progress_percentage
		FROM enrollments WHERE user_id = $1
		ORDER BY enrolled_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	enrollments := make([]*models.Enrollment, 0)
	for rows.Next() {
		enrollment := &models.Enrollment{}
		err := rows.Scan(
			&enrollment.ID,
			&enrollment.UserID,
			&enrollment.CourseID,
			&enrollment.EnrolledAt,
			&enrollment.CompletedAt,
			&enrollment.Status,
			&enrollment.ProgressPercentage,
		)
		if err != nil {
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
	}
	return enrollments, rows.Err()
}

// CountByCourse counts enrollments for a course
func (r *EnrollmentRepository) CountByCourse(ctx context.Context, courseID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM enrollments WHERE course_id = $1`
	var count int64
	err := r.db.QueryRow(ctx, query, courseID).Scan(&count)
	return count, err
}
