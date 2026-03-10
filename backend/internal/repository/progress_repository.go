package repository

import (
	"context"
	"time"

	"ai-learning-platform/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProgressRepository handles database operations for progress
type ProgressRepository struct {
	db *pgxpool.Pool
}

// NewProgressRepository creates a new ProgressRepository
func NewProgressRepository(db *pgxpool.Pool) *ProgressRepository {
	return &ProgressRepository{db: db}
}

// Create creates a new progress record
func (r *ProgressRepository) Create(ctx context.Context, progress *models.Progress) error {
	query := `
		INSERT INTO progress (id, user_id, lesson_id, enrollment_id, is_completed, is_watching, 
			video_position, last_accessed_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id, lesson_id) DO UPDATE SET
			is_completed = EXCLUDED.is_completed,
			is_watching = EXCLUDED.is_watching,
			video_position = EXCLUDED.video_position,
			last_accessed_at = EXCLUDED.last_accessed_at,
			completed_at = EXCLUDED.completed_at
	`
	_, err := r.db.Exec(ctx, query,
		progress.ID,
		progress.UserID,
		progress.LessonID,
		progress.EnrollmentID,
		progress.IsCompleted,
		progress.IsWatching,
		progress.VideoPosition,
		progress.LastAccessedAt,
		progress.CompletedAt,
	)
	return err
}

// GetByID retrieves a progress record by ID
func (r *ProgressRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Progress, error) {
	query := `
		SELECT id, user_id, lesson_id, enrollment_id, is_completed, is_watching, video_position, 
			last_accessed_at, completed_at
		FROM progress WHERE id = $1
	`
	progress := &models.Progress{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&progress.ID,
		&progress.UserID,
		&progress.LessonID,
		&progress.EnrollmentID,
		&progress.IsCompleted,
		&progress.IsWatching,
		&progress.VideoPosition,
		&progress.LastAccessedAt,
		&progress.CompletedAt,
	)
	if err != nil {
		return nil, err
	}
	return progress, nil
}

// GetByUserAndLesson retrieves progress for a user and lesson
func (r *ProgressRepository) GetByUserAndLesson(ctx context.Context, userID, lessonID uuid.UUID) (*models.Progress, error) {
	query := `
		SELECT id, user_id, lesson_id, enrollment_id, is_completed, is_watching, video_position, 
			last_accessed_at, completed_at
		FROM progress WHERE user_id = $1 AND lesson_id = $2
	`
	progress := &models.Progress{}
	err := r.db.QueryRow(ctx, query, userID, lessonID).Scan(
		&progress.ID,
		&progress.UserID,
		&progress.LessonID,
		&progress.EnrollmentID,
		&progress.IsCompleted,
		&progress.IsWatching,
		&progress.VideoPosition,
		&progress.LastAccessedAt,
		&progress.CompletedAt,
	)
	if err != nil {
		return nil, err
	}
	return progress, nil
}

// GetByUser retrieves all progress records for a user
func (r *ProgressRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]*models.Progress, error) {
	query := `
		SELECT id, user_id, lesson_id, enrollment_id, is_completed, is_watching, video_position, 
			last_accessed_at, completed_at
		FROM progress WHERE user_id = $1 ORDER BY last_accessed_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	progresses := make([]*models.Progress, 0)
	for rows.Next() {
		progress := &models.Progress{}
		err := rows.Scan(
			&progress.ID,
			&progress.UserID,
			&progress.LessonID,
			&progress.EnrollmentID,
			&progress.IsCompleted,
			&progress.IsWatching,
			&progress.VideoPosition,
			&progress.LastAccessedAt,
			&progress.CompletedAt,
		)
		if err != nil {
			return nil, err
		}
		progresses = append(progresses, progress)
	}

	return progresses, rows.Err()
}

// Update updates an existing progress record
func (r *ProgressRepository) Update(ctx context.Context, progress *models.Progress) error {
	query := `
		UPDATE progress SET is_completed = $2, is_watching = $3, video_position = $4, 
			last_accessed_at = $5, completed_at = $6
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		progress.ID,
		progress.IsCompleted,
		progress.IsWatching,
		progress.VideoPosition,
		progress.LastAccessedAt,
		progress.CompletedAt,
	)
	return err
}

// MarkCompleted marks a lesson as completed
func (r *ProgressRepository) MarkCompleted(ctx context.Context, userID, lessonID uuid.UUID) error {
	now := time.Now()
	query := `
		UPDATE progress SET is_completed = true, completed_at = $3, last_accessed_at = $3
		WHERE user_id = $1 AND lesson_id = $2
	`
	_, err := r.db.Exec(ctx, query, userID, lessonID, now)
	return err
}

// UpdateVideoPosition updates video playback position
func (r *ProgressRepository) UpdateVideoPosition(ctx context.Context, userID, lessonID uuid.UUID, position int) error {
	query := `
		UPDATE progress SET video_position = $3, is_watching = true, last_accessed_at = $4
		WHERE user_id = $1 AND lesson_id = $2
	`
	_, err := r.db.Exec(ctx, query, userID, lessonID, position, time.Now())
	return err
}

// GetByUserID retrieves all progress records for a user
func (r *ProgressRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Progress, error) {
	query := `
		SELECT id, user_id, lesson_id, is_completed, is_watching, video_position, completed_at, last_accessed_at
		FROM progress WHERE user_id = $1
		ORDER BY last_accessed_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	progresses := make([]*models.Progress, 0)
	for rows.Next() {
		progress := &models.Progress{}
		err := rows.Scan(
			&progress.ID,
			&progress.UserID,
			&progress.LessonID,
			&progress.IsCompleted,
			&progress.IsWatching,
			&progress.VideoPosition,
			&progress.CompletedAt,
			&progress.LastAccessedAt,
		)
		if err != nil {
			return nil, err
		}
		progresses = append(progresses, progress)
	}
	return progresses, rows.Err()
}
