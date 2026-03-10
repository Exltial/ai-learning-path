package services

import (
	"context"
	"errors"
	"time"

	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Error definitions for user service
var (
	ErrInvalidPassword     = errors.New("invalid password")
	ErrPasswordTooWeak     = errors.New("password too weak")
	ErrUsernameUnavailable = errors.New("username unavailable")
	ErrEmailUnavailable    = errors.New("email unavailable")
)

// UserService handles user-related business logic
type UserService struct {
	userRepo       *repository.UserRepository
	progressRepo   *repository.ProgressRepository
	enrollmentRepo *repository.EnrollmentRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo *repository.UserRepository, progressRepo *repository.ProgressRepository, enrollmentRepo *repository.EnrollmentRepository) *UserService {
	return &UserService{
		userRepo:       userRepo,
		progressRepo:   progressRepo,
		enrollmentRepo: enrollmentRepo,
	}
}

// GetUserByID retrieves a user by ID (without password hash)
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	// Clear password hash for security
	user.PasswordHash = ""
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	user.PasswordHash = ""
	return user, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(ctx context.Context, userID uuid.UUID, username, avatarURL string) (*models.User, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Validate and update username if provided
	if username != "" {
		if username != user.Username {
			// Check if new username is available
			exists, err := s.userRepo.ExistsByUsername(ctx, username)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, ErrUsernameUnavailable
			}
			user.Username = username
		}
	}

	// Update avatar URL if provided
	if avatarURL != "" {
		user.AvatarURL = avatarURL
	}

	user.UpdatedAt = time.Now()

	// Save changes
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	user.PasswordHash = ""
	return user, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return ErrInvalidPassword
	}

	// Validate new password
	if err := validatePasswordStrength(newPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	return s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword))
}

// GetUserStats retrieves comprehensive user statistics
func (s *UserService) GetUserStats(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	// Get user's enrollments
	enrollments, err := s.enrollmentRepo.GetByUserID(ctx, userID)
	if err != nil {
		enrollments = []*models.Enrollment{}
	}

	// Count completed courses
	completedCourses := 0
	for _, enrollment := range enrollments {
		if enrollment.Status == "completed" {
			completedCourses++
		}
	}

	// Get user's progress
	progresses, err := s.progressRepo.GetByUserID(ctx, userID)
	if err != nil {
		progresses = []*models.Progress{}
	}

	// Count completed lessons
	completedLessons := 0
	for _, progress := range progresses {
		if progress.IsCompleted {
			completedLessons++
		}
	}

	// Calculate learning streak (simplified)
	learningStreak := s.calculateLearningStreak(ctx, userID)

	stats := map[string]interface{}{
		"total_courses":        len(enrollments),
		"completed_courses":    completedCourses,
		"in_progress_courses":  len(enrollments) - completedCourses,
		"total_lessons":        len(progresses),
		"completed_lessons":    completedLessons,
		"learning_streak":      learningStreak,
		"total_points":         0, // TODO: Implement points system
		"achievements_count":   0, // TODO: Implement achievements
		"last_activity":        time.Now().Format(time.RFC3339),
	}

	return stats, nil
}

// GetUserAchievements retrieves user achievements
func (s *UserService) GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]*models.UserAchievement, error) {
	// TODO: Implement achievements retrieval from database
	achievements := []*models.UserAchievement{}
	return achievements, nil
}

// GetNotifications retrieves user notifications
func (s *UserService) GetNotifications(ctx context.Context, userID uuid.UUID) ([]*models.Notification, error) {
	// TODO: Implement notifications retrieval from database
	notifications := []*models.Notification{}
	return notifications, nil
}

// MarkNotificationRead marks a notification as read
func (s *UserService) MarkNotificationRead(ctx context.Context, notificationID uuid.UUID) error {
	// TODO: Implement notification read marking in database
	return nil
}

// GetUnreadNotificationCount returns the count of unread notifications
func (s *UserService) GetUnreadNotificationCount(ctx context.Context, userID uuid.UUID) (int, error) {
	notifications, err := s.GetNotifications(ctx, userID)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, n := range notifications {
		if !n.IsRead {
			count++
		}
	}
	return count, nil
}

// DeleteUser soft-deletes a user account
func (s *UserService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Soft delete by deactivating
	user.IsActive = false
	user.UpdatedAt = time.Now()
	
	return s.userRepo.Update(ctx, user)
}

// calculateLearningStreak calculates the user's learning streak
func (s *UserService) calculateLearningStreak(ctx context.Context, userID uuid.UUID) int {
	// TODO: Implement proper streak calculation based on activity
	// For now, return a placeholder
	return 0
}

// validatePasswordStrength validates password strength
func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooWeak
	}
	
	// Check for at least one uppercase letter
	hasUpper := false
	hasLower := false
	hasDigit := false
	
	for _, c := range password {
		if c >= 'A' && c <= 'Z' {
			hasUpper = true
		}
		if c >= 'a' && c <= 'z' {
			hasLower = true
		}
		if c >= '0' && c <= '9' {
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return errors.New("password must contain uppercase, lowercase, and digit")
	}

	return nil
}

// GetUserEnrollments retrieves all enrollments for a user
func (s *UserService) GetUserEnrollments(ctx context.Context, userID uuid.UUID) ([]*models.Enrollment, error) {
	return s.enrollmentRepo.GetByUserID(ctx, userID)
}

// GetUserProgress retrieves learning progress for a user
func (s *UserService) GetUserProgress(ctx context.Context, userID uuid.UUID) ([]*models.Progress, error) {
	return s.progressRepo.GetByUserID(ctx, userID)
}
