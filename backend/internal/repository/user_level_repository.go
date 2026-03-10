package repository

import (
	"context"
	"database/sql"
	"ai-learning-platform/internal/models"
	"github.com/google/uuid"
)

// UserLevelRepository handles user level data access
type UserLevelRepository struct {
	db *sql.DB
}

// NewUserLevelRepository creates a new UserLevelRepository
func NewUserLevelRepository(db *sql.DB) *UserLevelRepository {
	return &UserLevelRepository{db: db}
}

// Create creates a new user level entry
func (r *UserLevelRepository) Create(ctx context.Context, level *models.UserLevel) error {
	query := `
		INSERT INTO user_levels (id, user_id, level, current_points, total_points, experience, next_level_exp, title, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		level.ID,
		level.UserID,
		level.Level,
		level.CurrentPoints,
		level.TotalPoints,
		level.Experience,
		level.NextLevelExp,
		level.Title,
		level.UpdatedAt,
	)
	return err
}

// GetByUserID retrieves user level by user ID
func (r *UserLevelRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserLevel, error) {
	query := `
		SELECT id, user_id, level, current_points, total_points, experience, next_level_exp, title, updated_at
		FROM user_levels
		WHERE user_id = $1
	`
	level := &models.UserLevel{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&level.ID,
		&level.UserID,
		&level.Level,
		&level.CurrentPoints,
		&level.TotalPoints,
		&level.Experience,
		&level.NextLevelExp,
		&level.Title,
		&level.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return level, nil
}

// Update updates user level
func (r *UserLevelRepository) Update(ctx context.Context, level *models.UserLevel) error {
	query := `
		UPDATE user_levels
		SET level = $2, current_points = $3, total_points = $4, experience = $5, next_level_exp = $6, title = $7, updated_at = $8
		WHERE user_id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		level.UserID,
		level.Level,
		level.CurrentPoints,
		level.TotalPoints,
		level.Experience,
		level.NextLevelExp,
		level.Title,
		level.UpdatedAt,
	)
	return err
}

// GetTopUsers retrieves top users by total points
func (r *UserLevelRepository) GetTopUsers(ctx context.Context, limit int) ([]*models.UserLevel, error) {
	query := `
		SELECT id, user_id, level, current_points, total_points, experience, next_level_exp, title, updated_at
		FROM user_levels
		ORDER BY total_points DESC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	levels := []*models.UserLevel{}
	for rows.Next() {
		level := &models.UserLevel{}
		err := rows.Scan(
			&level.ID,
			&level.UserID,
			&level.Level,
			&level.CurrentPoints,
			&level.TotalPoints,
			&level.Experience,
			&level.NextLevelExp,
			&level.Title,
			&level.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		levels = append(levels, level)
	}
	return levels, nil
}
