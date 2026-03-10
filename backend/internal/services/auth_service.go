package services

import (
	"context"
	"errors"
	"time"

	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// Error definitions for authentication service
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUsernameExists     = errors.New("username already exists")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenGeneration    = errors.New("failed to generate token")
)

// AuthService handles authentication and authorization logic
type AuthService struct {
	userRepo    *repository.UserRepository
	redis       *redis.Client
	jwtSecret   []byte
	jwtExpiration int // hours
}

// NewAuthService creates a new AuthService with JWT configuration
func NewAuthService(userRepo *repository.UserRepository, redis *redis.Client, jwtSecret string, jwtExpiration int) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		redis:         redis,
		jwtSecret:     []byte(jwtSecret),
		jwtExpiration: jwtExpiration,
	}
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	AvatarURL string `json:"avatar_url"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Register creates a new user and returns JWT token
// It performs validation, checks for existing users, hashes password, and generates token
func (s *AuthService) Register(ctx context.Context, username, email, password string) (*models.User, string, error) {
	// Validate input
	if err := validateUsername(username); err != nil {
		return nil, "", err
	}
	if err := validateEmail(email); err != nil {
		return nil, "", err
	}
	if err := validatePassword(password); err != nil {
		return nil, "", err
	}

	// Check if username exists (case-insensitive)
	exists, err := s.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, "", err
	}
	if exists {
		return nil, "", ErrUsernameExists
	}

	// Check if email exists (case-insensitive)
	exists, err = s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if exists {
		return nil, "", ErrEmailExists
	}

	// Hash password with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// Create user object
	user := &models.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "student", // Default role
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save user to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	// Generate JWT token
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	// Cache user in Redis for faster lookups (optional optimization)
	s.cacheUser(ctx, user)

	return user, token, nil
}

// Login authenticates a user and returns JWT token
func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	// Validate input
	if err := validateEmail(email); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, "", errors.New("user account is deactivated")
	}

	// Verify password using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Update last login timestamp (non-blocking)
	go func() {
		updateCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := s.userRepo.UpdateLastLogin(updateCtx, user.ID); err != nil {
			// Log error but don't fail login
		}
	}()

	// Generate JWT token
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	// Update cache
	s.cacheUser(ctx, user)

	return user, token, nil
}

// GenerateToken generates a JWT token for a user
func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	// Create claims with user information
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Duration(s.jwtExpiration) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign token with secret
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", ErrTokenGeneration
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns claims
func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	// Parse token with validation
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	// Extract and validate claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, ErrInvalidToken
			}
		}
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	// Try to get from cache first
	if cachedUser, err := s.getCachedUser(ctx, id); err == nil {
		return cachedUser, nil
	}

	// Get from database
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	
	return user, nil
}

// RefreshToken refreshes an existing JWT token
func (s *AuthService) RefreshToken(ctx context.Context, oldToken string) (string, error) {
	// Validate old token
	claims, err := s.ValidateToken(oldToken)
	if err != nil {
		return "", ErrInvalidToken
	}

	// Get user ID from claims
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return "", ErrInvalidToken
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", ErrInvalidToken
	}

	// Get user
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return "", ErrUserNotFound
	}

	// Generate new token
	return s.GenerateToken(user)
}

// Logout invalidates a token (adds to Redis blacklist)
func (s *AuthService) Logout(ctx context.Context, token string) error {
	// Parse token to get expiration
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	
	if err != nil {
		return err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return ErrInvalidToken
	}

	// Get expiration time
	exp, ok := claims["exp"].(float64)
	if !ok {
		return ErrInvalidToken
	}

	// Add token to blacklist in Redis
	ttl := time.Duration(exp) * time.Second
	return s.redis.Set(ctx, "blacklist:"+token, "1", ttl-ttl).Err()
}

// IsTokenBlacklisted checks if a token is blacklisted
func (s *AuthService) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	result, err := s.redis.Get(ctx, "blacklist:"+token).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return result == "1", nil
}

// cacheUser caches user data in Redis for faster lookups
func (s *AuthService) cacheUser(ctx context.Context, user *models.User) {
	// Cache for 1 hour
	key := "user:" + user.ID.String()
	s.redis.HSet(ctx, key, 
		"id", user.ID.String(),
		"username", user.Username,
		"email", user.Email,
		"role", user.Role,
	)
	s.redis.Expire(ctx, key, time.Hour)
}

// getCachedUser retrieves user from Redis cache
func (s *AuthService) getCachedUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	key := "user:" + id.String()
	result, err := s.redis.HGetAll(ctx, key).Result()
	if err != nil || len(result) == 0 {
		return nil, err
	}

	userID, _ := uuid.Parse(result["id"])
	return &models.User{
		ID:       userID,
		Username: result["username"],
		Email:    result["email"],
		Role:     result["role"],
	}, nil
}

// validateUsername validates username format
func validateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if len(username) > 50 {
		return errors.New("username must be less than 50 characters")
	}
	return nil
}

// validateEmail validates email format
func validateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	// Basic email validation (more thorough validation can be added)
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			return nil
		}
	}
	return errors.New("invalid email format")
}

// validatePassword validates password strength
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	// Additional password strength checks can be added here
	return nil
}
