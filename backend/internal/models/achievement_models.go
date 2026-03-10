package models

import (
	"time"

	"github.com/google/uuid"
)

// AchievementType represents the type of achievement
type AchievementType string

const (
	AchievementTypeGeneral      AchievementType = "general"       // General achievements
	AchievementTypeCourse       AchievementType = "course"        // Course-related achievements
	AchievementTypeStreak       AchievementType = "streak"        // Streak-based achievements
	AchievementTypeExercise     AchievementType = "exercise"      // Exercise-related achievements
	AchievementTypeSocial       AchievementType = "social"        // Social achievements
	AchievementTypeMilestone    AchievementType = "milestone"     // Milestone achievements
)

// AchievementTier represents the tier/rarity of an achievement
type AchievementTier string

const (
	AchievementTierBronze   AchievementTier = "bronze"   // Common
	AchievementTierSilver   AchievementTier = "silver"   // Uncommon
	AchievementTierGold     AchievementTier = "gold"     // Rare
	AchievementTierPlatinum AchievementTier = "platinum" // Epic
	AchievementTierDiamond  AchievementTier = "diamond"  // Legendary
)

// Achievement represents an achievement badge in the system
type Achievement struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	Name            string          `json:"name" db:"name"`
	Description     string          `json:"description" db:"description"`
	IconURL         string          `json:"icon_url,omitempty" db:"icon_url"`
	Points          int             `json:"points" db:"points"`
	AchievementType AchievementType `json:"achievement_type" db:"achievement_type"`
	Tier            AchievementTier `json:"tier" db:"tier"`
	Criteria        AchievementCriteria `json:"criteria" db:"criteria"`
	IsEnabled       bool            `json:"is_enabled" db:"is_enabled"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

// AchievementCriteria defines the conditions to unlock an achievement
type AchievementCriteria struct {
	Type           string      `json:"type" db:"type"`                     // e.g., "complete_lesson", "complete_course", "streak_days", "complete_exercises"
	TargetID       *uuid.UUID  `json:"target_id,omitempty" db:"target_id"` // Optional: specific course/lesson ID
	Threshold      int         `json:"threshold" db:"threshold"`           // Required count/days
	Metadata       interface{} `json:"metadata,omitempty" db:"metadata"`   // Additional criteria data
}

// UserAchievement represents an achievement earned by a user
type UserAchievement struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	AchievementID uuid.UUID `json:"achievement_id" db:"achievement_id"`
	EarnedAt      time.Time `json:"earned_at" db:"earned_at"`
	IsNotified    bool      `json:"is_notified" db:"is_notified"` // Whether user has been notified
}

// UserLevel represents a user's level in the gamification system
type UserLevel struct {
	ID              uuid.UUID `json:"id" db:"id"`
	UserID          uuid.UUID `json:"user_id" db:"user_id"`
	Level           int       `json:"level" db:"level"`
	CurrentPoints   int       `json:"current_points" db:"current_points"`
	TotalPoints     int       `json:"total_points" db:"total_points"`
	Experience      int       `json:"experience" db:"experience"`
	NextLevelExp    int       `json:"next_level_exp" db:"next_level_exp"`
	Title           string    `json:"title" db:"title"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// LevelTitle represents a title available at a specific level
type LevelTitle struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Level       int       `json:"level" db:"level"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	IconURL     string    `json:"icon_url,omitempty" db:"icon_url"`
}

// UserPointsTransaction represents a points transaction for a user
type UserPointsTransaction struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	Amount        int       `json:"amount" db:"amount"` // Positive for earned, negative for spent
	BalanceAfter  int       `json:"balance_after" db:"balance_after"`
	SourceType    string    `json:"source_type" db:"source_type"` // achievement, course_complete, exercise, daily_bonus, etc.
	SourceID      *uuid.UUID `json:"source_id,omitempty" db:"source_id"`
	Description   string    `json:"description" db:"description"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// LeaderboardEntry represents a user's position on the leaderboard
type LeaderboardEntry struct {
	Rank          int       `json:"rank" db:"rank"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	Username      string    `json:"username" db:"username"`
	AvatarURL     string    `json:"avatar_url,omitempty" db:"avatar_url"`
	TotalPoints   int       `json:"total_points" db:"total_points"`
	Level         int       `json:"level" db:"level"`
	Title         string    `json:"title" db:"title"`
	AchievementsCount int   `json:"achievements_count" db:"achievements_count"`
}

// LeaderboardType represents different leaderboard categories
type LeaderboardType string

const (
	LeaderboardTypeWeekly   LeaderboardType = "weekly"
	LeaderboardTypeMonthly  LeaderboardType = "monthly"
	LeaderboardTypeAllTime  LeaderboardType = "all_time"
	LeaderboardTypeFriends  LeaderboardType = "friends"
)

// UserStreak represents a user's learning streak
type UserStreak struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	CurrentStreak   int        `json:"current_streak" db:"current_streak"`
	LongestStreak   int        `json:"longest_streak" db:"longest_streak"`
	LastActivityAt  time.Time  `json:"last_activity_at" db:"last_activity_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// DailyChallenge represents a daily challenge for users
type DailyChallenge struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Title         string    `json:"title" db:"title"`
	Description   string    `json:"description" db:"description"`
	ChallengeType string    `json:"challenge_type" db:"challenge_type"`
	TargetCount   int       `json:"target_count" db:"target_count"`
	RewardPoints  int       `json:"reward_points" db:"reward_points"`
	Date          time.Time `json:"date" db:"date"`
	IsActive      bool      `json:"is_active" db:"is_active"`
}

// UserDailyChallengeProgress tracks user progress on daily challenges
type UserDailyChallengeProgress struct {
	ID              uuid.UUID `json:"id" db:"id"`
	UserID          uuid.UUID `json:"user_id" db:"user_id"`
	ChallengeID     uuid.UUID `json:"challenge_id" db:"challenge_id"`
	CurrentCount    int       `json:"current_count" db:"current_count"`
	IsCompleted     bool      `json:"is_completed" db:"is_completed"`
	CompletedAt     *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// AchievementWithProgress wraps an achievement with user's progress towards it
type AchievementWithProgress struct {
	Achievement   Achievement `json:"achievement"`
	CurrentCount  int         `json:"current_count"`
	ProgressPercent float64   `json:"progress_percent"`
	IsUnlocked    bool        `json:"is_unlocked"`
	EarnedAt      *time.Time  `json:"earned_at,omitempty"`
}

// UserAchievementSummary provides a summary of user's achievement stats
type UserAchievementSummary struct {
	UserID              uuid.UUID `json:"user_id"`
	TotalAchievements   int       `json:"total_achievements"`
	UnlockedAchievements int      `json:"unlocked_achievements"`
	TotalPoints         int       `json:"total_points"`
	CurrentLevel        int       `json:"current_level"`
	CurrentTitle        string    `json:"current_title"`
	CurrentStreak       int       `json:"current_streak"`
	LongestStreak       int       `json:"longest_streak"`
	Rank                int       `json:"rank"`
}
