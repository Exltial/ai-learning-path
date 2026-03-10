package repository

import (
	"context"
	"database/sql"
	"ai-learning-platform/internal/models"
	"github.com/google/uuid"
)

// LeaderboardRepository handles leaderboard data access
type LeaderboardRepository struct {
	db *sql.DB
}

// NewLeaderboardRepository creates a new LeaderboardRepository
func NewLeaderboardRepository(db *sql.DB) *LeaderboardRepository {
	return &LeaderboardRepository{db: db}
}

// GetAllTime retrieves all-time leaderboard
func (r *LeaderboardRepository) GetAllTime(ctx context.Context, limit int) ([]*models.LeaderboardEntry, error) {
	query := `
		SELECT 
			ROW_NUMBER() OVER (ORDER BY ul.total_points DESC) as rank,
			ul.user_id,
			u.username,
			u.avatar_url,
			ul.total_points,
			ul.level,
			ul.title,
			COUNT(ua.id) as achievements_count
		FROM user_levels ul
		INNER JOIN users u ON ul.user_id = u.id
		LEFT JOIN user_achievements ua ON ul.user_id = ua.user_id
		WHERE u.is_active = true
		GROUP BY ul.user_id, ul.total_points, ul.level, ul.title, u.username, u.avatar_url
		ORDER BY ul.total_points DESC
		LIMIT $1
	`
	return r.queryLeaderboard(ctx, query, limit)
}

// GetWeekly retrieves weekly leaderboard
func (r *LeaderboardRepository) GetWeekly(ctx context.Context, limit int) ([]*models.LeaderboardEntry, error) {
	query := `
		SELECT 
			ROW_NUMBER() OVER (ORDER BY COALESCE(weekly_points.points, 0) DESC) as rank,
			u.id as user_id,
			u.username,
			u.avatar_url,
			COALESCE(weekly_points.points, 0) as total_points,
			ul.level,
			ul.title,
			COUNT(ua.id) as achievements_count
		FROM users u
		LEFT JOIN user_levels ul ON u.id = ul.user_id
		LEFT JOIN (
			SELECT user_id, SUM(amount) as points
			FROM points_transactions
			WHERE created_at >= NOW() - INTERVAL '7 days'
			GROUP BY user_id
		) weekly_points ON u.id = weekly_points.user_id
		LEFT JOIN user_achievements ua ON u.id = ua.user_id
		WHERE u.is_active = true
		GROUP BY u.id, u.username, u.avatar_url, ul.level, ul.title, weekly_points.points
		ORDER BY weekly_points.points DESC NULLS LAST
		LIMIT $1
	`
	return r.queryLeaderboard(ctx, query, limit)
}

// GetMonthly retrieves monthly leaderboard
func (r *LeaderboardRepository) GetMonthly(ctx context.Context, limit int) ([]*models.LeaderboardEntry, error) {
	query := `
		SELECT 
			ROW_NUMBER() OVER (ORDER BY COALESCE(monthly_points.points, 0) DESC) as rank,
			u.id as user_id,
			u.username,
			u.avatar_url,
			COALESCE(monthly_points.points, 0) as total_points,
			ul.level,
			ul.title,
			COUNT(ua.id) as achievements_count
		FROM users u
		LEFT JOIN user_levels ul ON u.id = ul.user_id
		LEFT JOIN (
			SELECT user_id, SUM(amount) as points
			FROM points_transactions
			WHERE created_at >= NOW() - INTERVAL '30 days'
			GROUP BY user_id
		) monthly_points ON u.id = monthly_points.user_id
		LEFT JOIN user_achievements ua ON u.id = ua.user_id
		WHERE u.is_active = true
		GROUP BY u.id, u.username, u.avatar_url, ul.level, ul.title, monthly_points.points
		ORDER BY monthly_points.points DESC NULLS LAST
		LIMIT $1
	`
	return r.queryLeaderboard(ctx, query, limit)
}

// GetFriends retrieves friends leaderboard
func (r *LeaderboardRepository) GetFriends(ctx context.Context, userID uuid.UUID, limit int) ([]*models.LeaderboardEntry, error) {
	// This would require a friends/following table in a real implementation
	// For now, return all-time leaderboard
	return r.GetAllTime(ctx, limit)
}

// GetUserRank retrieves a user's rank
func (r *LeaderboardRepository) GetUserRank(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT rank FROM (
			SELECT 
				user_id,
				ROW_NUMBER() OVER (ORDER BY total_points DESC) as rank
			FROM user_levels
		) ranked
		WHERE user_id = $1
	`
	var rank int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&rank)
	if err != nil {
		return 0, err
	}
	return rank, nil
}

func (r *LeaderboardRepository) queryLeaderboard(ctx context.Context, query string, limit int) ([]*models.LeaderboardEntry, error) {
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []*models.LeaderboardEntry{}
	for rows.Next() {
		entry := &models.LeaderboardEntry{}
		err := rows.Scan(
			&entry.Rank,
			&entry.UserID,
			&entry.Username,
			&entry.AvatarURL,
			&entry.TotalPoints,
			&entry.Level,
			&entry.Title,
			&entry.AchievementsCount,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
