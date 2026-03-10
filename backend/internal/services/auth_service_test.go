package services_test

import (
	"context"
	"testing"

	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"
	"ai-learning-platform/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	args := m.Called(ctx, id, passwordHash)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

// TestRegister_Success tests successful user registration
func TestRegister_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockRedis := &MockRedisClient{}
	
	authService := services.NewAuthService(mockUserRepo, mockRedis, "test-secret", 24)
	
	// Mock expectations
	mockUserRepo.On("ExistsByUsername", mock.Anything, "testuser").Return(false, nil)
	mockUserRepo.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false, nil)
	mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
		return u.Username == "testuser" && u.Email == "test@example.com"
	})).Return(nil)
	
	// Execute
	user, token, err := authService.Register(context.Background(), "testuser", "test@example.com", "Password123")
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NotEmpty(t, token)
	
	mockUserRepo.AssertExpectations(t)
}

// TestRegister_UsernameExists tests registration with existing username
func TestRegister_UsernameExists(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockRedis := &MockRedisClient{}
	
	authService := services.NewAuthService(mockUserRepo, mockRedis, "test-secret", 24)
	
	// Mock expectations
	mockUserRepo.On("ExistsByUsername", mock.Anything, "existinguser").Return(true, nil)
	
	// Execute
	user, token, err := authService.Register(context.Background(), "existinguser", "new@example.com", "Password123")
	
	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrUsernameExists, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	
	mockUserRepo.AssertExpectations(t)
}

// TestRegister_EmailExists tests registration with existing email
func TestRegister_EmailExists(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockRedis := &MockRedisClient{}
	
	authService := services.NewAuthService(mockUserRepo, mockRedis, "test-secret", 24)
	
	// Mock expectations
	mockUserRepo.On("ExistsByUsername", mock.Anything, "newuser").Return(false, nil)
	mockUserRepo.On("ExistsByEmail", mock.Anything, "existing@example.com").Return(true, nil)
	
	// Execute
	user, token, err := authService.Register(context.Background(), "newuser", "existing@example.com", "Password123")
	
	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrEmailExists, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	
	mockUserRepo.AssertExpectations(t)
}

// TestRegister_WeakPassword tests registration with weak password
func TestRegister_WeakPassword(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockRedis := &MockRedisClient{}
	
	authService := services.NewAuthService(mockUserRepo, mockRedis, "test-secret", 24)
	
	// Execute - password too short
	user, token, err := authService.Register(context.Background(), "testuser", "test@example.com", "short")
	
	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "password must be at least 8 characters")
	assert.Nil(t, user)
	assert.Empty(t, token)
}

// TestLogin_Success tests successful login
func TestLogin_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockRedis := &MockRedisClient{}
	
	authService := services.NewAuthService(mockUserRepo, mockRedis, "test-secret", 24)
	
	// Create a test user with hashed password
	testUser := &models.User{
		ID:           uuid.New(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // "Password123"
		Role:         "student",
		IsActive:     true,
	}
	
	// Mock expectations
	mockUserRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, testUser.ID).Return(nil)
	
	// Execute
	user, token, err := authService.Login(context.Background(), "test@example.com", "Password123")
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.NotEmpty(t, token)
	
	mockUserRepo.AssertExpectations(t)
}

// TestLogin_InvalidCredentials tests login with wrong password
func TestLogin_InvalidCredentials(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockRedis := &MockRedisClient{}
	
	authService := services.NewAuthService(mockUserRepo, mockRedis, "test-secret", 24)
	
	// Create a test user
	testUser := &models.User{
		ID:           uuid.New(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // "Password123"
		IsActive:     true,
	}
	
	// Mock expectations
	mockUserRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
	
	// Execute - wrong password
	user, token, err := authService.Login(context.Background(), "test@example.com", "WrongPassword")
	
	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidCredentials, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	
	mockUserRepo.AssertExpectations(t)
}

// TestLogin_UserNotFound tests login with non-existent user
func TestLogin_UserNotFound(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockRedis := &MockRedisClient{}
	
	authService := services.NewAuthService(mockUserRepo, mockRedis, "test-secret", 24)
	
	// Mock expectations
	mockUserRepo.On("GetByEmail", mock.Anything, "nonexistent@example.com").Return((*models.User)(nil), assert.AnError)
	
	// Execute
	user, token, err := authService.Login(context.Background(), "nonexistent@example.com", "Password123")
	
	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidCredentials, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	
	mockUserRepo.AssertExpectations(t)
}

// TestGenerateToken tests JWT token generation
func TestGenerateToken(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockRedis := &MockRedisClient{}
	
	authService := services.NewAuthService(mockUserRepo, mockRedis, "test-secret", 24)
	
	testUser := &models.User{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "student",
	}
	
	// Execute
	token, err := authService.GenerateToken(testUser)
	
	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

// TestValidateToken_Valid tests validation of a valid token
func TestValidateToken_Valid(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockRedis := &MockRedisClient{}
	
	authService := services.NewAuthService(mockUserRepo, mockRedis, "test-secret", 24)
	
	testUser := &models.User{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "student",
	}
	
	// Generate a token
	token, _ := authService.GenerateToken(testUser)
	
	// Execute
	claims, err := authService.ValidateToken(token)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, testUser.Email, claims["email"])
	assert.Equal(t, testUser.Role, claims["role"])
}

// TestValidateToken_Invalid tests validation of an invalid token
func TestValidateToken_Invalid(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockRedis := &MockRedisClient{}
	
	authService := services.NewAuthService(mockUserRepo, mockRedis, "test-secret", 24)
	
	// Execute - invalid token
	claims, err := authService.ValidateToken("invalid.token.here")
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// MockRedisClient is a mock implementation of Redis client
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration interface{}) interface{} {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0)
}

func (m *MockRedisClient) Get(ctx context.Context, key string) interface{} {
	args := m.Called(ctx, key)
	return args.Get(0)
}

func (m *MockRedisClient) HSet(ctx context.Context, key string, values ...interface{}) interface{} {
	args := m.Called(ctx, key, values)
	return args.Get(0)
}

func (m *MockRedisClient) HGetAll(ctx context.Context, key string) interface{} {
	args := m.Called(ctx, key)
	return args.Get(0)
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration interface{}) interface{} {
	args := m.Called(ctx, key, expiration)
	return args.Get(0)
}

func (m *MockRedisClient) Ping(ctx context.Context) interface{} {
	args := m.Called(ctx)
	return args.Get(0)
}
