package repository

import (
	"context"
	"time"

	"ai-learning-platform/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ExerciseRepository handles database operations for exercises
type ExerciseRepository struct {
	db *pgxpool.Pool
}

// NewExerciseRepository creates a new ExerciseRepository
func NewExerciseRepository(db *pgxpool.Pool) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

// Create creates a new exercise
func (r *ExerciseRepository) Create(ctx context.Context, exercise *models.Exercise) error {
	query := `
		INSERT INTO exercises (id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts, 
			time_limit, starter_code, test_cases, expected_answer, options, order_index, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	_, err := r.db.Exec(ctx, query,
		exercise.ID,
		exercise.LessonID,
		exercise.Title,
		exercise.Description,
		exercise.ExerciseType,
		exercise.Difficulty,
		exercise.Points,
		exercise.MaxAttempts,
		exercise.TimeLimit,
		exercise.StarterCode,
		exercise.TestCases,
		exercise.ExpectedAnswer,
		exercise.Options,
		exercise.OrderIndex,
		exercise.CreatedAt,
		exercise.UpdatedAt,
	)
	return err
}

// GetByID retrieves an exercise by ID
func (r *ExerciseRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Exercise, error) {
	query := `
		SELECT id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts,
			time_limit, starter_code, test_cases, expected_answer, options, order_index, created_at, updated_at
		FROM exercises WHERE id = $1
	`
	exercise := &models.Exercise{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&exercise.ID,
		&exercise.LessonID,
		&exercise.Title,
		&exercise.Description,
		&exercise.ExerciseType,
		&exercise.Difficulty,
		&exercise.Points,
		&exercise.MaxAttempts,
		&exercise.TimeLimit,
		&exercise.StarterCode,
		&exercise.TestCases,
		&exercise.ExpectedAnswer,
		&exercise.Options,
		&exercise.OrderIndex,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return exercise, nil
}

// GetByLessonID retrieves all exercises for a lesson
func (r *ExerciseRepository) GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*models.Exercise, error) {
	query := `
		SELECT id, lesson_id, title, description, exercise_type, difficulty, points, max_attempts,
			time_limit, starter_code, test_cases, expected_answer, options, order_index, created_at, updated_at
		FROM exercises WHERE lesson_id = $1 ORDER BY order_index ASC
	`
	rows, err := r.db.Query(ctx, query, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exercises := make([]*models.Exercise, 0)
	for rows.Next() {
		exercise := &models.Exercise{}
		err := rows.Scan(
			&exercise.ID,
			&exercise.LessonID,
			&exercise.Title,
			&exercise.Description,
			&exercise.ExerciseType,
			&exercise.Difficulty,
			&exercise.Points,
			&exercise.MaxAttempts,
			&exercise.TimeLimit,
			&exercise.StarterCode,
			&exercise.TestCases,
			&exercise.ExpectedAnswer,
			&exercise.Options,
			&exercise.OrderIndex,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}

	return exercises, rows.Err()
}

// Update updates an existing exercise
func (r *ExerciseRepository) Update(ctx context.Context, exercise *models.Exercise) error {
	query := `
		UPDATE exercises SET title = $2, description = $3, exercise_type = $4, difficulty = $5,
			points = $6, max_attempts = $7, time_limit = $8, starter_code = $9, test_cases = $10,
			expected_answer = $11, options = $12, order_index = $13, updated_at = $14
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		exercise.ID,
		exercise.Title,
		exercise.Description,
		exercise.ExerciseType,
		exercise.Difficulty,
		exercise.Points,
		exercise.MaxAttempts,
		exercise.TimeLimit,
		exercise.StarterCode,
		exercise.TestCases,
		exercise.ExpectedAnswer,
		exercise.Options,
		exercise.OrderIndex,
		time.Now(),
	)
	return err
}

// Delete deletes an exercise
func (r *ExerciseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM exercises WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
