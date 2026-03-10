package repository

import (
	"context"
	"database/sql"
	"ai-learning-platform/internal/models"
	"github.com/google/uuid"
)

// StreakRepository handles user streak data access
type StreakRepository struct {
	db *sql.DB
}

// NewStreakRepository creates a new StreakRepository
func NewStreakRepository(db *sql.DB) *StreakRepository {
	return &StreakRepository{db: db}
}

// Create creates a new user streak
func (r *StreakRepository) Create(ctx context.Context, streak *models.UserStreak) error {
	query := `
		INSERT INTO user_streaks (id, user_id, current_streak, longest_streak, last_activity_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		streak.ID,
		streak.UserID,
		streak.CurrentStreak,
		streak.LongestStreak,
		streak.LastActivityAt,
		streak.UpdatedAt,
	)
	return err
}

// GetByUserID retrieves user streak by user ID
func (r *StreakRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserStreak, error) {
	query := `
		SELECT id, user_id, current_streak, longest_streak, last_activity_at, updated_at
		FROM user_streaks
		WHERE user_id = $1
	`
	streak := &models.UserStreak{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&streak.ID,
		&streak.UserID,
		&streak.CurrentStreak,
		&streak.LongestStreak,
		&streak.LastActivityAt,
		&streak.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return streak, nil
}

// Update updates user streak
func (r *StreakRepository) Update(ctx context.Context, streak *models.UserStreak) error {
	query := `
		UPDATE user_streaks
		SET current_streak = $2, longest_streak = $3, last_activity_at = $4, updated_at = $5
		WHERE user_id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		streak.UserID,
		streak.CurrentStreak,
		streak.LongestStreak,
		streak.LastActivityAt,
		streak.UpdatedAt,
	)
	return err
}

// GetTopStreaks retrieves users with longest streaks
func (r *StreakRepository) GetTopStreaks(ctx context.Context, limit int) ([]*models.UserStreak, error) {
	query := `
		SELECT id, user_id, current_streak, longest_streak, last_activity_at, updated_at
		FROM user_streaks
		ORDER BY current_streak DESC, longest_streak DESC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	streaks := []*models.UserStreak{}
	for rows.Next() {
		streak := &models.UserStreak{}
		err := rows.Scan(
			&streak.ID,
			&streak.UserID,
			&streak.CurrentStreak,
			&streak.LongestStreak,
			&streak.LastActivityAt,
			&streak.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		streaks = append(streaks, streak)
	}
	return streaks, nil
}
