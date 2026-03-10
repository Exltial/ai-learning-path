package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Username     string     `json:"username" db:"username"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	AvatarURL    string     `json:"avatar_url,omitempty" db:"avatar_url"`
	Role         string     `json:"role" db:"role"` // student, instructor, admin
	IsActive     bool       `json:"is_active" db:"is_active"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}

// Course represents a course
type Course struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	Title            string     `json:"title" db:"title"`
	Description      string     `json:"description" db:"description"`
	ThumbnailURL     string     `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	InstructorID     uuid.UUID  `json:"instructor_id" db:"instructor_id"`
	Instructor       *User      `json:"instructor,omitempty" db:"-"`
	Category         string     `json:"category,omitempty" db:"category"`
	DifficultyLevel  string     `json:"difficulty_level" db:"difficulty_level"` // beginner, intermediate, advanced
	EstimatedHours   int        `json:"estimated_hours" db:"estimated_hours"`
	Price            float64    `json:"price" db:"price"`
	IsPublished      bool       `json:"is_published" db:"is_published"`
	EnrollmentCount  int        `json:"enrollment_count" db:"enrollment_count"`
	Rating           float64    `json:"rating" db:"rating"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// Lesson represents a lesson within a course
type Lesson struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	CourseID      uuid.UUID  `json:"course_id" db:"course_id"`
	Title         string     `json:"title" db:"title"`
	Description   string     `json:"description,omitempty" db:"description"`
	Content       string     `json:"content,omitempty" db:"content"`
	VideoURL      string     `json:"video_url,omitempty" db:"video_url"`
	VideoDuration int        `json:"video_duration,omitempty" db:"video_duration"` // in seconds
	OrderIndex    int        `json:"order_index" db:"order_index"`
	IsFreePreview bool       `json:"is_free_preview" db:"is_free_preview"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// Exercise represents an exercise within a lesson
type Exercise struct {
	ID           uuid.UUID          `json:"id" db:"id"`
	LessonID     uuid.UUID          `json:"lesson_id" db:"lesson_id"`
	Title        string             `json:"title" db:"title"`
	Description  string             `json:"description,omitempty" db:"description"`
	ExerciseType string             `json:"exercise_type" db:"exercise_type"` // multiple_choice, coding, fill_blank, true_false, essay
	Difficulty   string             `json:"difficulty" db:"difficulty"`       // easy, medium, hard
	Points       int                `json:"points" db:"points"`
	MaxAttempts  int                `json:"max_attempts" db:"max_attempts"`
	TimeLimit    *int               `json:"time_limit,omitempty" db:"time_limit"` // in seconds
	StarterCode  string             `json:"starter_code,omitempty" db:"starter_code"`
	TestCases    map[string]interface{} `json:"test_cases,omitempty" db:"test_cases"`
	ExpectedAnswer map[string]interface{} `json:"expected_answer,omitempty" db:"expected_answer"`
	Options      []ExerciseOption   `json:"options,omitempty" db:"options"`
	OrderIndex   int                `json:"order_index" db:"order_index"`
	CreatedAt    time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" db:"updated_at"`
}

// ExerciseOption represents an option for multiple choice exercises
type ExerciseOption struct {
	Text       string `json:"text" db:"text"`
	IsCorrect  bool   `json:"is_correct" db:"is_correct"`
}

// Enrollment represents a user's enrollment in a course
type Enrollment struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	UserID           uuid.UUID  `json:"user_id" db:"user_id"`
	CourseID         uuid.UUID  `json:"course_id" db:"course_id"`
	EnrolledAt       time.Time  `json:"enrolled_at" db:"enrolled_at"`
	CompletedAt      *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	Status           string     `json:"status" db:"status"` // active, completed, dropped
	ProgressPercentage float64  `json:"progress_percentage" db:"progress_percentage"`
}

// Submission represents a user's submission for an exercise
type Submission struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	ExerciseID     uuid.UUID  `json:"exercise_id" db:"exercise_id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	SubmissionType string     `json:"submission_type" db:"submission_type"`
	Answer         string     `json:"answer,omitempty" db:"answer"`
	Code           string     `json:"code,omitempty" db:"code"`
	IsCorrect      *bool      `json:"is_correct,omitempty" db:"is_correct"`
	Score          *float64   `json:"score,omitempty" db:"score"`
	Feedback       string     `json:"feedback,omitempty" db:"feedback"`
	AttemptNumber  int        `json:"attempt_number" db:"attempt_number"`
	SubmittedAt    time.Time  `json:"submitted_at" db:"submitted_at"`
	GradedAt       *time.Time `json:"graded_at,omitempty" db:"graded_at"`
	GradedBy       *uuid.UUID `json:"graded_by,omitempty" db:"graded_by"`
}

// Progress represents a user's progress in a lesson
type Progress struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	LessonID        uuid.UUID  `json:"lesson_id" db:"lesson_id"`
	EnrollmentID    uuid.UUID  `json:"enrollment_id" db:"enrollment_id"`
	IsCompleted     bool       `json:"is_completed" db:"is_completed"`
	IsWatching      bool       `json:"is_watching" db:"is_watching"`
	VideoPosition   int        `json:"video_position" db:"video_position"` // in seconds
	LastAccessedAt  time.Time  `json:"last_accessed_at" db:"last_accessed_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}

// Achievement represents an achievement badge
type Achievement struct {
	ID              uuid.UUID              `json:"id" db:"id"`
	Name            string                 `json:"name" db:"name"`
	Description     string                 `json:"description" db:"description"`
	IconURL         string                 `json:"icon_url,omitempty" db:"icon_url"`
	Points          int                    `json:"points" db:"points"`
	AchievementType string                 `json:"achievement_type" db:"achievement_type"` // general, course, streak, social
	Criteria        map[string]interface{} `json:"criteria" db:"criteria"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
}

// UserAchievement represents an achievement earned by a user
type UserAchievement struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	AchievementID uuid.UUID `json:"achievement_id" db:"achievement_id"`
	EarnedAt      time.Time `json:"earned_at" db:"earned_at"`
}

// Discussion represents a discussion thread
type Discussion struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	CourseID   uuid.UUID  `json:"course_id" db:"course_id"`
	LessonID   *uuid.UUID `json:"lesson_id,omitempty" db:"lesson_id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	ParentID   *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
	Title      string     `json:"title,omitempty" db:"title"`
	Content    string     `json:"content" db:"content"`
	IsResolved bool       `json:"is_resolved" db:"is_resolved"`
	Upvotes    int        `json:"upvotes" db:"upvotes"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

// Notification represents a user notification
type Notification struct {
	ID               uuid.UUID `json:"id" db:"id"`
	UserID           uuid.UUID `json:"user_id" db:"user_id"`
	Title            string    `json:"title" db:"title"`
	Message          string    `json:"message" db:"message"`
	NotificationType string    `json:"notification_type" db:"notification_type"` // info, success, warning, error
	IsRead           bool      `json:"is_read" db:"is_read"`
	ActionURL        string    `json:"action_url,omitempty" db:"action_url"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// CourseReview represents a review for a course
type CourseReview struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CourseID   uuid.UUID `json:"course_id" db:"course_id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	User       *User     `json:"user,omitempty" db:"-"`
	Rating     int       `json:"rating" db:"rating"` // 1-5
	Comment    string    `json:"comment,omitempty" db:"comment"`
	IsVerified bool      `json:"is_verified" db:"is_verified"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// QuizResult represents a quiz result
type QuizResult struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	LessonID       *uuid.UUID `json:"lesson_id,omitempty" db:"lesson_id"`
	CourseID       *uuid.UUID `json:"course_id,omitempty" db:"course_id"`
	TotalQuestions int        `json:"total_questions" db:"total_questions"`
	CorrectAnswers int        `json:"correct_answers" db:"correct_answers"`
	ScorePercentage float64   `json:"score_percentage" db:"score_percentage"`
	TimeTaken      *int       `json:"time_taken,omitempty" db:"time_taken"` // in seconds
	StartedAt      time.Time  `json:"started_at" db:"started_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}
