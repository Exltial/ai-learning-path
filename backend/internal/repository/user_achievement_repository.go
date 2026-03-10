package repository

import (
	"context"
	"database/sql"
	"ai-learning-platform/internal/models"
	"github.com/google/uuid"
)

// UserAchievementRepository handles user achievement data access
type UserAchievementRepository struct {
	db *sql.DB
}

// NewUserAchievementRepository creates a new UserAchievementRepository
func NewUserAchievementRepository(db *sql.DB) *UserAchievementRepository {
	return &UserAchievementRepository{db: db}
}

// Create creates a new user achievement
func (r *UserAchievementRepository) Create(ctx context.Context, userAchievement *models.UserAchievement) error {
	query := `
		INSERT INTO user_achievements (id, user_id, achievement_id, earned_at, is_notified)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		userAchievement.ID,
		userAchievement.UserID,
		userAchievement.AchievementID,
		userAchievement.EarnedAt,
		userAchievement.IsNotified,
	)
	return err
}

// GetByUserID retrieves all achievements for a user
func (r *UserAchievementRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.UserAchievement, error) {
	query := `
		SELECT id, user_id, achievement_id, earned_at, is_notified
		FROM user_achievements
		WHERE user_id = $1
		ORDER BY earned_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	achievements := []*models.UserAchievement{}
	for rows.Next() {
		ua := &models.UserAchievement{}
		err := rows.Scan(
			&ua.ID,
			&ua.UserID,
			&ua.AchievementID,
			&ua.EarnedAt,
			&ua.IsNotified,
		)
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, ua)
	}
	return achievements, nil
}

// GetByUserAndAchievement retrieves a specific user achievement
func (r *UserAchievementRepository) GetByUserAndAchievement(ctx context.Context, userID, achievementID uuid.UUID) (*models.UserAchievement, error) {
	query := `
		SELECT id, user_id, achievement_id, earned_at, is_notified
		FROM user_achievements
		WHERE user_id = $1 AND achievement_id = $2
	`
	ua := &models.UserAchievement{}
	err := r.db.QueryRowContext(ctx, query, userID, achievementID).Scan(
		&ua.ID,
		&ua.UserID,
		&ua.AchievementID,
		&ua.EarnedAt,
		&ua.IsNotified,
	)
	if err != nil {
		return nil, err
	}
	return ua, nil
}

// GetUserAchievementIDs retrieves all achievement IDs for a user
func (r *UserAchievementRepository) GetUserAchievementIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT achievement_id FROM user_achievements WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []uuid.UUID{}
	for rows.Next() {
		var id uuid.UUID
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// GetCountByUserID returns count of achievements for a user
func (r *UserAchievementRepository) GetCountByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM user_achievements WHERE user_id = $1`
	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

// UpdateNotificationStatus updates the notification status
func (r *UserAchievementRepository) UpdateNotificationStatus(ctx context.Context, id uuid.UUID, isNotified bool) error {
	query := `UPDATE user_achievements SET is_notified = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, isNotified)
	return err
}

// Delete deletes a user achievement
func (r *UserAchievementRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM user_achievements WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
