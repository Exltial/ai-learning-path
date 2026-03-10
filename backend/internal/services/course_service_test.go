package services_test

import (
	"context"
	"testing"
	"time"

	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"
	"ai-learning-platform/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCourseRepository is a mock implementation of CourseRepository
type MockCourseRepository struct {
	mock.Mock
}

func (m *MockCourseRepository) Create(ctx context.Context, course *models.Course) error {
	args := m.Called(ctx, course)
	return args.Error(0)
}

func (m *MockCourseRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Course, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Course), args.Error(1)
}

func (m *MockCourseRepository) List(ctx context.Context, category, difficulty string, page, limit int) ([]*models.Course, int, error) {
	args := m.Called(ctx, category, difficulty, page, limit)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*models.Course), args.Int(1), args.Error(2)
}

func (m *MockCourseRepository) Update(ctx context.Context, course *models.Course) error {
	args := m.Called(ctx, course)
	return args.Error(0)
}

func (m *MockCourseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCourseRepository) IncrementEnrollmentCount(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCourseRepository) UpdateRating(ctx context.Context, id uuid.UUID, rating float64) error {
	args := m.Called(ctx, id, rating)
	return args.Error(0)
}

func (m *MockCourseRepository) GetByInstructorID(ctx context.Context, instructorID uuid.UUID) ([]*models.Course, error) {
	args := m.Called(ctx, instructorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Course), args.Error(1)
}

// MockLessonRepository is a mock implementation of LessonRepository
type MockLessonRepository struct {
	mock.Mock
}

func (m *MockLessonRepository) GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]*models.Lesson, error) {
	args := m.Called(ctx, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Lesson), args.Error(1)
}

// MockEnrollmentRepository is a mock implementation of EnrollmentRepository
type MockEnrollmentRepository struct {
	mock.Mock
}

func (m *MockEnrollmentRepository) Create(ctx context.Context, enrollment *models.Enrollment) error {
	args := m.Called(ctx, enrollment)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) Exists(ctx context.Context, userID, courseID uuid.UUID) (bool, error) {
	args := m.Called(ctx, userID, courseID)
	return args.Bool(0), args.Error(1)
}

func (m *MockEnrollmentRepository) GetByUserAndCourse(ctx context.Context, userID, courseID uuid.UUID) (*models.Enrollment, error) {
	args := m.Called(ctx, userID, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Enrollment), args.Error(1)
}

func (m *MockEnrollmentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Enrollment, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Enrollment), args.Error(1)
}

func (m *MockEnrollmentRepository) CountByCourse(ctx context.Context, courseID uuid.UUID) (int, error) {
	args := m.Called(ctx, courseID)
	return args.Int(0), args.Error(1)
}

// TestCreateCourse_Success tests successful course creation
func TestCreateCourse_Success(t *testing.T) {
	// Setup
	mockCourseRepo := new(MockCourseRepository)
	mockLessonRepo := new(MockLessonRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	
	courseService := services.NewCourseService(mockCourseRepo, mockLessonRepo, mockEnrollmentRepo)
	
	// Mock expectations
	mockCourseRepo.On("Create", mock.Anything, mock.MatchedBy(func(c *models.Course) bool {
		return c.Title == "Test Course" && c.DifficultyLevel == "beginner"
	})).Return(nil)
	
	// Execute
	course := &models.Course{
		Title:           "Test Course",
		Description:     "A test course",
		DifficultyLevel: "beginner",
		EstimatedHours:  10,
		Price:           99.99,
	}
	
	err := courseService.CreateCourse(context.Background(), course)
	
	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, course.ID)
	assert.False(t, course.IsPublished) // Should default to unpublished
	
	mockCourseRepo.AssertExpectations(t)
}

// TestCreateCourse_InvalidTitle tests course creation with invalid title
func TestCreateCourse_InvalidTitle(t *testing.T) {
	// Setup
	mockCourseRepo := new(MockCourseRepository)
	mockLessonRepo := new(MockLessonRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	
	courseService := services.NewCourseService(mockCourseRepo, mockLessonRepo, mockEnrollmentRepo)
	
	// Execute - empty title
	course := &models.Course{
		Title: "",
	}
	
	err := courseService.CreateCourse(context.Background(), course)
	
	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title is required")
}

// TestCreateCourse_InvalidDifficulty tests course creation with invalid difficulty
func TestCreateCourse_InvalidDifficulty(t *testing.T) {
	// Setup
	mockCourseRepo := new(MockCourseRepository)
	mockLessonRepo := new(MockLessonRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	
	courseService := services.NewCourseService(mockCourseRepo, mockLessonRepo, mockEnrollmentRepo)
	
	// Execute - invalid difficulty
	course := &models.Course{
		Title:           "Test Course",
		DifficultyLevel: "invalid_level",
	}
	
	err := courseService.CreateCourse(context.Background(), course)
	
	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidDifficulty, err)
}

// TestGetCourse_Success tests successful course retrieval
func TestGetCourse_Success(t *testing.T) {
	// Setup
	mockCourseRepo := new(MockCourseRepository)
	mockLessonRepo := new(MockLessonRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	
	courseService := services.NewCourseService(mockCourseRepo, mockLessonRepo, mockEnrollmentRepo)
	
	courseID := uuid.New()
	testCourse := &models.Course{
		ID:          courseID,
		Title:       "Test Course",
		Description: "A test course",
	}
	
	// Mock expectations
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(testCourse, nil)
	
	// Execute
	course, err := courseService.GetCourse(context.Background(), courseID)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, course)
	assert.Equal(t, "Test Course", course.Title)
	
	mockCourseRepo.AssertExpectations(t)
}

// TestGetCourse_NotFound tests course retrieval for non-existent course
func TestGetCourse_NotFound(t *testing.T) {
	// Setup
	mockCourseRepo := new(MockCourseRepository)
	mockLessonRepo := new(MockLessonRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	
	courseService := services.NewCourseService(mockCourseRepo, mockLessonRepo, mockEnrollmentRepo)
	
	courseID := uuid.New()
	
	// Mock expectations
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return((*models.Course)(nil), assert.AnError)
	
	// Execute
	course, err := courseService.GetCourse(context.Background(), courseID)
	
	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrCourseNotFound, err)
	assert.Nil(t, course)
	
	mockCourseRepo.AssertExpectations(t)
}

// TestListCourses_Success tests successful course listing
func TestListCourses_Success(t *testing.T) {
	// Setup
	mockCourseRepo := new(MockCourseRepository)
	mockLessonRepo := new(MockLessonRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	
	courseService := services.NewCourseService(mockCourseRepo, mockLessonRepo, mockEnrollmentRepo)
	
	testCourses := []*models.Course{
		{ID: uuid.New(), Title: "Course 1"},
		{ID: uuid.New(), Title: "Course 2"},
	}
	
	// Mock expectations
	mockCourseRepo.On("List", mock.Anything, "", "", 1, 20).Return(testCourses, 2, nil)
	
	// Execute
	courses, total, err := courseService.ListCourses(context.Background(), "", "", 1, 20)
	
	// Assert
	assert.NoError(t, err)
	assert.Len(t, courses, 2)
	assert.Equal(t, 2, total)
	
	mockCourseRepo.AssertExpectations(t)
}

// TestListCourses_InvalidDifficulty tests listing with invalid difficulty filter
func TestListCourses_InvalidDifficulty(t *testing.T) {
	// Setup
	mockCourseRepo := new(MockCourseRepository)
	mockLessonRepo := new(MockLessonRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	
	courseService := services.NewCourseService(mockCourseRepo, mockLessonRepo, mockEnrollmentRepo)
	
	// Execute - invalid difficulty
	courses, total, err := courseService.ListCourses(context.Background(), "", "invalid", 1, 20)
	
	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidDifficulty, err)
	assert.Nil(t, courses)
	assert.Equal(t, 0, total)
}

// TestEnrollCourse_Success tests successful course enrollment
func TestEnrollCourse_Success(t *testing.T) {
	// Setup
	mockCourseRepo := new(MockCourseRepository)
	mockLessonRepo := new(MockLessonRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	
	courseService := services.NewCourseService(mockCourseRepo, mockLessonRepo, mockEnrollmentRepo)
	
	userID := uuid.New()
	courseID := uuid.New()
	
	testCourse := &models.Course{
		ID:           courseID,
		Title:        "Test Course",
		IsPublished:  true,
		CreatedAt:    time.Now(),
	}
	
	// Mock expectations
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(testCourse, nil)
	mockEnrollmentRepo.On("Exists", mock.Anything, userID, courseID).Return(false, nil)
	mockEnrollmentRepo.On("Create", mock.Anything, mock.MatchedBy(func(e *models.Enrollment) bool {
		return e.UserID == userID && e.CourseID == courseID
	})).Return(nil)
	
	// Execute
	enrollment, err := courseService.EnrollCourse(context.Background(), userID, courseID)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, enrollment)
	assert.Equal(t, userID, enrollment.UserID)
	assert.Equal(t, courseID, enrollment.CourseID)
	assert.Equal(t, "active", enrollment.Status)
	
	mockCourseRepo.AssertExpectations(t)
	mockEnrollmentRepo.AssertExpectations(t)
}

// TestEnrollCourse_NotPublished tests enrollment in unpublished course
func TestEnrollCourse_NotPublished(t *testing.T) {
	// Setup
	mockCourseRepo := new(MockCourseRepository)
	mockLessonRepo := new(MockLessonRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	
	courseService := services.NewCourseService(mockCourseRepo, mockLessonRepo, mockEnrollmentRepo)
	
	userID := uuid.New()
	courseID := uuid.New()
	
	testCourse := &models.Course{
		ID:          courseID,
		Title:       "Test Course",
		IsPublished: false, // Not published
	}
	
	// Mock expectations
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(testCourse, nil)
	
	// Execute
	enrollment, err := courseService.EnrollCourse(context.Background(), userID, courseID)
	
	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrCourseNotPublished, err)
	assert.Nil(t, enrollment)
	
	mockCourseRepo.AssertExpectations(t)
}

// TestEnrollCourse_AlreadyEnrolled tests enrollment when already enrolled
func TestEnrollCourse_AlreadyEnrolled(t *testing.T) {
	// Setup
	mockCourseRepo := new(MockCourseRepository)
	mockLessonRepo := new(MockLessonRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	
	courseService := services.NewCourseService(mockCourseRepo, mockLessonRepo, mockEnrollmentRepo)
	
	userID := uuid.New()
	courseID := uuid.New()
	
	testCourse := &models.Course{
		ID:          courseID,
		IsPublished: true,
	}
	
	existingEnrollment := &models.Enrollment{
		ID:      uuid.New(),
		UserID:  userID,
		CourseID: courseID,
		Status:  "active",
	}
	
	// Mock expectations
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(testCourse, nil)
	mockEnrollmentRepo.On("Exists", mock.Anything, userID, courseID).Return(true, nil)
	mockEnrollmentRepo.On("GetByUserAndCourse", mock.Anything, userID, courseID).Return(existingEnrollment, nil)
	
	// Execute
	enrollment, err := courseService.EnrollCourse(context.Background(), userID, courseID)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, enrollment)
	assert.Equal(t, existingEnrollment.ID, enrollment.ID)
	
	mockCourseRepo.AssertExpectations(t)
	mockEnrollmentRepo.AssertExpectations(t)
}

// TestUpdateCourse_Success tests successful course update
func TestUpdateCourse_Success(t *testing.T) {
	// Setup
	mockCourseRepo := new(MockCourseRepository)
	mockLessonRepo := new(MockLessonRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	
	courseService := services.NewCourseService(mockCourseRepo, mockLessonRepo, mockEnrollmentRepo)
	
	courseID := uuid.New()
	existingCourse := &models.Course{
		ID:               courseID,
		Title:            "Original Title",
		Description:      "Original description",
		EnrollmentCount:  10,
		Rating:           4.5,
		CreatedAt:        time.Now(),
	}
	
	// Mock expectations
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(existingCourse, nil)
	mockCourseRepo.On("Update", mock.Anything, mock.MatchedBy(func(c *models.Course) bool {
		return c.Title == "Updated Title"
	})).Return(nil)
	
	// Execute
	existingCourse.Title = "Updated Title"
	err := courseService.UpdateCourse(context.Background(), existingCourse)
	
	// Assert
	assert.NoError(t, err)
	
	mockCourseRepo.AssertExpectations(t)
}

// TestPublishCourse tests publishing a course
func TestPublishCourse(t *testing.T) {
	// Setup
	mockCourseRepo := new(MockCourseRepository)
	mockLessonRepo := new(MockLessonRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	
	courseService := services.NewCourseService(mockCourseRepo, mockLessonRepo, mockEnrollmentRepo)
	
	courseID := uuid.New()
	testCourse := &models.Course{
		ID:           courseID,
		Title:        "Test Course",
		IsPublished:  false,
	}
	
	// Mock expectations
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(testCourse, nil)
	mockCourseRepo.On("Update", mock.Anything, mock.MatchedBy(func(c *models.Course) bool {
		return c.IsPublished == true
	})).Return(nil)
	
	// Execute
	err := courseService.PublishCourse(context.Background(), courseID)
	
	// Assert
	assert.NoError(t, err)
	
	mockCourseRepo.AssertExpectations(t)
}
