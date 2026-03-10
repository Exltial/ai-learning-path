package repository

import (
	"context"
	"fmt"
	"time"

	"ai-learning-platform/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CourseRepository handles database operations for courses
type CourseRepository struct {
	db *pgxpool.Pool
}

// NewCourseRepository creates a new CourseRepository
func NewCourseRepository(db *pgxpool.Pool) *CourseRepository {
	return &CourseRepository{db: db}
}

// Create creates a new course
func (r *CourseRepository) Create(ctx context.Context, course *models.Course) error {
	query := `
		INSERT INTO courses (id, title, description, thumbnail_url, instructor_id, category, difficulty_level, 
			estimated_hours, price, is_published, enrollment_count, rating, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := r.db.Exec(ctx, query,
		course.ID,
		course.Title,
		course.Description,
		course.ThumbnailURL,
		course.InstructorID,
		course.Category,
		course.DifficultyLevel,
		course.EstimatedHours,
		course.Price,
		course.IsPublished,
		course.EnrollmentCount,
		course.Rating,
		course.CreatedAt,
		course.UpdatedAt,
	)
	return err
}

// GetByID retrieves a course by ID
func (r *CourseRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Course, error) {
	query := `
		SELECT id, title, description, thumbnail_url, instructor_id, category, difficulty_level,
			estimated_hours, price, is_published, enrollment_count, rating, created_at, updated_at
		FROM courses WHERE id = $1
	`
	course := &models.Course{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.ThumbnailURL,
		&course.InstructorID,
		&course.Category,
		&course.DifficultyLevel,
		&course.EstimatedHours,
		&course.Price,
		&course.IsPublished,
		&course.EnrollmentCount,
		&course.Rating,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return course, nil
}

// List retrieves courses with pagination and filters
func (r *CourseRepository) List(ctx context.Context, category, difficulty string, page, limit int) ([]*models.Course, int, error) {
	// Count query
	countQuery := `SELECT COUNT(*) FROM courses WHERE is_published = true`
	var args []interface{}
	argIndex := 1

	if category != "" {
		countQuery += fmt.Sprintf(` AND category = $%d`, argIndex)
		args = append(args, category)
		argIndex++
	}
	if difficulty != "" {
		countQuery += fmt.Sprintf(` AND difficulty_level = $%d`, argIndex)
		args = append(args, difficulty)
		argIndex++
	}

	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Data query
	offset := (page - 1) * limit
	query := `
		SELECT id, title, description, thumbnail_url, instructor_id, category, difficulty_level,
			estimated_hours, price, is_published, enrollment_count, rating, created_at, updated_at
		FROM courses WHERE is_published = true
	`
	
	if category != "" {
		query += fmt.Sprintf(` AND category = $%d`, argIndex)
		args = append(args, category)
		argIndex++
	}
	if difficulty != "" {
		query += fmt.Sprintf(` AND difficulty_level = $%d`, argIndex)
		args = append(args, difficulty)
		argIndex++
	}

	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	courses := make([]*models.Course, 0)
	for rows.Next() {
		course := &models.Course{}
		err := rows.Scan(
			&course.ID,
			&course.Title,
			&course.Description,
			&course.ThumbnailURL,
			&course.InstructorID,
			&course.Category,
			&course.DifficultyLevel,
			&course.EstimatedHours,
			&course.Price,
			&course.IsPublished,
			&course.EnrollmentCount,
			&course.Rating,
			&course.CreatedAt,
			&course.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		courses = append(courses, course)
	}

	return courses, total, rows.Err()
}

// Update updates an existing course
func (r *CourseRepository) Update(ctx context.Context, course *models.Course) error {
	query := `
		UPDATE courses SET title = $2, description = $3, thumbnail_url = $4, category = $5,
			difficulty_level = $6, estimated_hours = $7, price = $8, is_published = $9, updated_at = $10
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		course.ID,
		course.Title,
		course.Description,
		course.ThumbnailURL,
		course.Category,
		course.DifficultyLevel,
		course.EstimatedHours,
		course.Price,
		course.IsPublished,
		time.Now(),
	)
	return err
}

// Delete deletes a course
func (r *CourseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM courses WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// IncrementEnrollmentCount increments the enrollment count
func (r *CourseRepository) IncrementEnrollmentCount(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE courses SET enrollment_count = enrollment_count + 1 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// UpdateRating updates the course rating
func (r *CourseRepository) UpdateRating(ctx context.Context, id uuid.UUID, rating float64) error {
	query := `UPDATE courses SET rating = $2 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, rating)
	return err
}

// GetByInstructorID retrieves courses by instructor ID
func (r *CourseRepository) GetByInstructorID(ctx context.Context, instructorID uuid.UUID) ([]*models.Course, error) {
	query := `
		SELECT id, title, description, thumbnail_url, instructor_id, category, difficulty_level,
			estimated_hours, price, is_published, enrollment_count, rating, created_at, updated_at
		FROM courses WHERE instructor_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, instructorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := make([]*models.Course, 0)
	for rows.Next() {
		course := &models.Course{}
		err := rows.Scan(
			&course.ID,
			&course.Title,
			&course.Description,
			&course.ThumbnailURL,
			&course.InstructorID,
			&course.Category,
			&course.DifficultyLevel,
			&course.EstimatedHours,
			&course.Price,
			&course.IsPublished,
			&course.EnrollmentCount,
			&course.Rating,
			&course.CreatedAt,
			&course.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, rows.Err()
}
