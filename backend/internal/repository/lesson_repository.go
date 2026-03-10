package repository

import (
	"context"
	"time"

	"ai-learning-platform/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// LessonRepository handles database operations for lessons
type LessonRepository struct {
	db *pgxpool.Pool
}

// NewLessonRepository creates a new LessonRepository
func NewLessonRepository(db *pgxpool.Pool) *LessonRepository {
	return &LessonRepository{db: db}
}

// Create creates a new lesson
func (r *LessonRepository) Create(ctx context.Context, lesson *models.Lesson) error {
	query := `
		INSERT INTO lessons (id, course_id, title, description, content, video_url, video_duration, order_index, is_free_preview, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(ctx, query,
		lesson.ID,
		lesson.CourseID,
		lesson.Title,
		lesson.Description,
		lesson.Content,
		lesson.VideoURL,
		lesson.VideoDuration,
		lesson.OrderIndex,
		lesson.IsFreePreview,
		lesson.CreatedAt,
		lesson.UpdatedAt,
	)
	return err
}

// GetByID retrieves a lesson by ID
func (r *LessonRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Lesson, error) {
	query := `
		SELECT id, course_id, title, description, content, video_url, video_duration, order_index, is_free_preview, created_at, updated_at
		FROM lessons WHERE id = $1
	`
	lesson := &models.Lesson{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&lesson.ID,
		&lesson.CourseID,
		&lesson.Title,
		&lesson.Description,
		&lesson.Content,
		&lesson.VideoURL,
		&lesson.VideoDuration,
		&lesson.OrderIndex,
		&lesson.IsFreePreview,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return lesson, nil
}

// GetByCourseID retrieves all lessons for a course
func (r *LessonRepository) GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]*models.Lesson, error) {
	query := `
		SELECT id, course_id, title, description, content, video_url, video_duration, order_index, is_free_preview, created_at, updated_at
		FROM lessons WHERE course_id = $1 ORDER BY order_index ASC
	`
	rows, err := r.db.Query(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lessons := make([]*models.Lesson, 0)
	for rows.Next() {
		lesson := &models.Lesson{}
		err := rows.Scan(
			&lesson.ID,
			&lesson.CourseID,
			&lesson.Title,
			&lesson.Description,
			&lesson.Content,
			&lesson.VideoURL,
			&lesson.VideoDuration,
			&lesson.OrderIndex,
			&lesson.IsFreePreview,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}

	return lessons, rows.Err()
}

// Update updates an existing lesson
func (r *LessonRepository) Update(ctx context.Context, lesson *models.Lesson) error {
	query := `
		UPDATE lessons SET title = $2, description = $3, content = $4, video_url = $5,
			video_duration = $6, order_index = $7, is_free_preview = $8, updated_at = $9
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		lesson.ID,
		lesson.Title,
		lesson.Description,
		lesson.Content,
		lesson.VideoURL,
		lesson.VideoDuration,
		lesson.OrderIndex,
		lesson.IsFreePreview,
		time.Now(),
	)
	return err
}

// Delete deletes a lesson
func (r *LessonRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM lessons WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
