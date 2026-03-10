package repository

import (
	"context"
	"database/sql"

	"ai-learning-platform/internal/models"

	"github.com/google/uuid"
)

// AchievementRepository handles achievement data access
type AchievementRepository struct {
	db *sql.DB
}

// NewAchievementRepository creates a new AchievementRepository
func NewAchievementRepository(db *sql.DB) *AchievementRepository {
	return &AchievementRepository{db: db}
}

// Create creates a new achievement
func (r *AchievementRepository) Create(ctx context.Context, achievement *models.Achievement) error {
	query := `
		INSERT INTO achievements (id, name, description, icon_url, points, achievement_type, tier, criteria, is_enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		achievement.ID,
		achievement.Name,
		achievement.Description,
		achievement.IconURL,
		achievement.Points,
		achievement.AchievementType,
		achievement.Tier,
		achievement.Criteria,
		achievement.IsEnabled,
		achievement.CreatedAt,
		achievement.UpdatedAt,
	)
	return err
}

// GetByID retrieves an achievement by ID
func (r *AchievementRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Achievement, error) {
	query := `
		SELECT id, name, description, icon_url, points, achievement_type, tier, criteria, is_enabled, created_at, updated_at
		FROM achievements
		WHERE id = $1
	`
	achievement := &models.Achievement{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&achievement.ID,
		&achievement.Name,
		&achievement.Description,
		&achievement.IconURL,
		&achievement.Points,
		&achievement.AchievementType,
		&achievement.Tier,
		&achievement.Criteria,
		&achievement.IsEnabled,
		&achievement.CreatedAt,
		&achievement.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return achievement, nil
}

// GetByName retrieves an achievement by name
func (r *AchievementRepository) GetByName(ctx context.Context, name string) (*models.Achievement, error) {
	query := `
		SELECT id, name, description, icon_url, points, achievement_type, tier, criteria, is_enabled, created_at, updated_at
		FROM achievements
		WHERE name = $1
	`
	achievement := &models.Achievement{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&achievement.ID,
		&achievement.Name,
		&achievement.Description,
		&achievement.IconURL,
		&achievement.Points,
		&achievement.AchievementType,
		&achievement.Tier,
		&achievement.Criteria,
		&achievement.IsEnabled,
		&achievement.CreatedAt,
		&achievement.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return achievement, nil
}

// GetAllEnabled retrieves all enabled achievements
func (r *AchievementRepository) GetAllEnabled(ctx context.Context) ([]*models.Achievement, error) {
	query := `
		SELECT id, name, description, icon_url, points, achievement_type, tier, criteria, is_enabled, created_at, updated_at
		FROM achievements
		WHERE is_enabled = true
		ORDER BY points DESC, created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	achievements := []*models.Achievement{}
	for rows.Next() {
		achievement := &models.Achievement{}
		err := rows.Scan(
			&achievement.ID,
			&achievement.Name,
			&achievement.Description,
			&achievement.IconURL,
			&achievement.Points,
			&achievement.AchievementType,
			&achievement.Tier,
			&achievement.Criteria,
			&achievement.IsEnabled,
			&achievement.CreatedAt,
			&achievement.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, achievement)
	}
	return achievements, nil
}

// GetCount returns total count of achievements
func (r *AchievementRepository) GetCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM achievements WHERE is_enabled = true`
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

// Update updates an achievement
func (r *AchievementRepository) Update(ctx context.Context, achievement *models.Achievement) error {
	query := `
		UPDATE achievements
		SET name = $2, description = $3, icon_url = $4, points = $5, achievement_type = $6, tier = $7, criteria = $8, is_enabled = $9, updated_at = $10
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		achievement.ID,
		achievement.Name,
		achievement.Description,
		achievement.IconURL,
		achievement.Points,
		achievement.AchievementType,
		achievement.Tier,
		achievement.Criteria,
		achievement.IsEnabled,
		achievement.UpdatedAt,
	)
	return err
}

// Delete deletes an achievement
func (r *AchievementRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM achievements WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
