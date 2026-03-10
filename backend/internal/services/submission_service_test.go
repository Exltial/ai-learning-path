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

// MockSubmissionRepository is a mock implementation of SubmissionRepository
type MockSubmissionRepository struct {
	mock.Mock
}

func (m *MockSubmissionRepository) Create(ctx context.Context, submission *models.Submission) error {
	args := m.Called(ctx, submission)
	return args.Error(0)
}

func (m *MockSubmissionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Submission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Submission), args.Error(1)
}

func (m *MockSubmissionRepository) GetByExerciseIDAndUserID(ctx context.Context, exerciseID, userID uuid.UUID) ([]*models.Submission, error) {
	args := m.Called(ctx, exerciseID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Submission), args.Error(1)
}

func (m *MockSubmissionRepository) Update(ctx context.Context, submission *models.Submission) error {
	args := m.Called(ctx, submission)
	return args.Error(0)
}

func (m *MockSubmissionRepository) GetLatestSubmission(ctx context.Context, exerciseID, userID uuid.UUID) (*models.Submission, error) {
	args := m.Called(ctx, exerciseID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Submission), args.Error(1)
}

// MockExerciseRepository is a mock implementation of ExerciseRepository
type MockExerciseRepository struct {
	mock.Mock
}

func (m *MockExerciseRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Exercise, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Exercise), args.Error(1)
}

// TestSubmitExercise_MultipleChoice_Success tests successful multiple choice submission
func TestSubmitExercise_MultipleChoice_Success(t *testing.T) {
	// Setup
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	
	submissionService := services.NewSubmissionService(mockSubmissionRepo, mockExerciseRepo)
	
	exerciseID := uuid.New()
	userID := uuid.New()
	
	testExercise := &models.Exercise{
		ID:           exerciseID,
		Title:        "Test Question",
		ExerciseType: "multiple_choice",
		MaxAttempts:  3,
		Options: []models.ExerciseOption{
			{Text: "Option A", IsCorrect: true},
			{Text: "Option B", IsCorrect: false},
			{Text: "Option C", IsCorrect: false},
		},
	}
	
	// Mock expectations
	mockExerciseRepo.On("GetByID", mock.Anything, exerciseID).Return(testExercise, nil)
	mockSubmissionRepo.On("GetByExerciseIDAndUserID", mock.Anything, exerciseID, userID).Return([]*models.Submission{}, nil)
	mockSubmissionRepo.On("Create", mock.Anything, mock.MatchedBy(func(s *models.Submission) bool {
		return s.ExerciseID == exerciseID && s.UserID == userID && s.AttemptNumber == 1
	})).Return(nil)
	
	// Execute
	submission, err := submissionService.SubmitExercise(
		context.Background(),
		exerciseID,
		userID,
		"",
		"",
		[]string{"Option A"},
	)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, submission)
	assert.Equal(t, 1, submission.AttemptNumber)
	assert.True(t, *submission.IsCorrect)
	assert.Equal(t, 100.0, *submission.Score)
	
	mockExerciseRepo.AssertExpectations(t)
	mockSubmissionRepo.AssertExpectations(t)
}

// TestSubmitExercise_MultipleChoice_Wrong tests wrong multiple choice answer
func TestSubmitExercise_MultipleChoice_Wrong(t *testing.T) {
	// Setup
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	
	submissionService := services.NewSubmissionService(mockSubmissionRepo, mockExerciseRepo)
	
	exerciseID := uuid.New()
	userID := uuid.New()
	
	testExercise := &models.Exercise{
		ID:           exerciseID,
		ExerciseType: "multiple_choice",
		MaxAttempts:  3,
		Options: []models.ExerciseOption{
			{Text: "Option A", IsCorrect: true},
			{Text: "Option B", IsCorrect: false},
		},
	}
	
	// Mock expectations
	mockExerciseRepo.On("GetByID", mock.Anything, exerciseID).Return(testExercise, nil)
	mockSubmissionRepo.On("GetByExerciseIDAndUserID", mock.Anything, exerciseID, userID).Return([]*models.Submission{}, nil)
	mockSubmissionRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	
	// Execute - wrong answer
	submission, err := submissionService.SubmitExercise(
		context.Background(),
		exerciseID,
		userID,
		"",
		"",
		[]string{"Option B"},
	)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, submission)
	assert.False(t, *submission.IsCorrect)
	assert.Less(t, *submission.Score, 100.0)
	
	mockExerciseRepo.AssertExpectations(t)
}

// TestSubmitExercise_MaxAttemptsReached tests submission when max attempts reached
func TestSubmitExercise_MaxAttemptsReached(t *testing.T) {
	// Setup
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	
	submissionService := services.NewSubmissionService(mockSubmissionRepo, mockExerciseRepo)
	
	exerciseID := uuid.New()
	userID := uuid.New()
	
	testExercise := &models.Exercise{
		ID:           exerciseID,
		ExerciseType: "multiple_choice",
		MaxAttempts:  2,
	}
	
	previousSubmissions := []*models.Submission{
		{AttemptNumber: 2},
	}
	
	// Mock expectations
	mockExerciseRepo.On("GetByID", mock.Anything, exerciseID).Return(testExercise, nil)
	mockSubmissionRepo.On("GetByExerciseIDAndUserID", mock.Anything, exerciseID, userID).Return(previousSubmissions, nil)
	
	// Execute
	submission, err := submissionService.SubmitExercise(
		context.Background(),
		exerciseID,
		userID,
		"",
		"",
		[]string{"Option A"},
	)
	
	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrMaxAttemptsReached, err)
	assert.Nil(t, submission)
	
	mockExerciseRepo.AssertExpectations(t)
	mockSubmissionRepo.AssertExpectations(t)
}

// TestSubmitExercise_TrueFalse_Success tests successful true/false submission
func TestSubmitExercise_TrueFalse_Success(t *testing.T) {
	// Setup
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	
	submissionService := services.NewSubmissionService(mockSubmissionRepo, mockExerciseRepo)
	
	exerciseID := uuid.New()
	userID := uuid.New()
	
	testExercise := &models.Exercise{
		ID:           exerciseID,
		ExerciseType: "true_false",
		ExpectedAnswer: map[string]interface{}{
			"answer": "true",
		},
	}
	
	// Mock expectations
	mockExerciseRepo.On("GetByID", mock.Anything, exerciseID).Return(testExercise, nil)
	mockSubmissionRepo.On("GetByExerciseIDAndUserID", mock.Anything, exerciseID, userID).Return([]*models.Submission{}, nil)
	mockSubmissionRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	
	// Execute - correct answer
	submission, err := submissionService.SubmitExercise(
		context.Background(),
		exerciseID,
		userID,
		"true",
		"",
		nil,
	)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, submission)
	assert.True(t, *submission.IsCorrect)
	
	mockExerciseRepo.AssertExpectations(t)
}

// TestSubmitExercise_FillBlank_Success tests successful fill-in-the-blank submission
func TestSubmitExercise_FillBlank_Success(t *testing.T) {
	// Setup
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	
	submissionService := services.NewSubmissionService(mockSubmissionRepo, mockExerciseRepo)
	
	exerciseID := uuid.New()
	userID := uuid.New()
	
	testExercise := &models.Exercise{
		ID:           exerciseID,
		ExerciseType: "fill_blank",
		ExpectedAnswer: map[string]interface{}{
			"answers": []interface{}{"Go", "Golang", "GO"},
		},
	}
	
	// Mock expectations
	mockExerciseRepo.On("GetByID", mock.Anything, exerciseID).Return(testExercise, nil)
	mockSubmissionRepo.On("GetByExerciseIDAndUserID", mock.Anything, exerciseID, userID).Return([]*models.Submission{}, nil)
	mockSubmissionRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	
	// Execute - correct answer
	submission, err := submissionService.SubmitExercise(
		context.Background(),
		exerciseID,
		userID,
		"Go",
		"",
		nil,
	)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, submission)
	assert.True(t, *submission.IsCorrect)
	
	mockExerciseRepo.AssertExpectations(t)
}

// TestSubmitExercise_Essay tests essay submission (requires manual grading)
func TestSubmitExercise_Essay(t *testing.T) {
	// Setup
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	
	submissionService := services.NewSubmissionService(mockSubmissionRepo, mockExerciseRepo)
	
	exerciseID := uuid.New()
	userID := uuid.New()
	
	testExercise := &models.Exercise{
		ID:           exerciseID,
		ExerciseType: "essay",
	}
	
	// Mock expectations
	mockExerciseRepo.On("GetByID", mock.Anything, exerciseID).Return(testExercise, nil)
	mockSubmissionRepo.On("GetByExerciseIDAndUserID", mock.Anything, exerciseID, userID).Return([]*models.Submission{}, nil)
	mockSubmissionRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	
	// Execute
	submission, err := submissionService.SubmitExercise(
		context.Background(),
		exerciseID,
		userID,
		"This is my essay answer about the topic...",
		"",
		nil,
	)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, submission)
	assert.False(t, *submission.IsCorrect) // Not auto-graded
	assert.Equal(t, 0.0, *submission.Score)
	assert.Contains(t, submission.Feedback, "awaiting instructor review")
	
	mockExerciseRepo.AssertExpectations(t)
}

// TestSubmitExercise_InvalidSubmission tests submission with invalid data
func TestSubmitExercise_InvalidSubmission(t *testing.T) {
	// Setup
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	
	submissionService := services.NewSubmissionService(mockSubmissionRepo, mockExerciseRepo)
	
	exerciseID := uuid.New()
	userID := uuid.New()
	
	testExercise := &models.Exercise{
		ID:           exerciseID,
		ExerciseType: "multiple_choice",
	}
	
	// Mock expectations
	mockExerciseRepo.On("GetByID", mock.Anything, exerciseID).Return(testExercise, nil)
	
	// Execute - no options selected
	submission, err := submissionService.SubmitExercise(
		context.Background(),
		exerciseID,
		userID,
		"",
		"",
		[]string{}, // Empty options
	)
	
	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one option must be selected")
	assert.Nil(t, submission)
	
	mockExerciseRepo.AssertExpectations(t)
}

// TestGradeSubmission_Success tests successful submission grading
func TestGradeSubmission_Success(t *testing.T) {
	// Setup
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	
	submissionService := services.NewSubmissionService(mockSubmissionRepo, mockExerciseRepo)
	
	submissionID := uuid.New()
	graderID := uuid.New()
	
	testSubmission := &models.Submission{
		ID:             submissionID,
		ExerciseID:     uuid.New(),
		UserID:         uuid.New(),
		SubmissionType: "essay",
		GradedAt:       nil, // Not yet graded
	}
	
	// Mock expectations
	mockSubmissionRepo.On("GetByID", mock.Anything, submissionID).Return(testSubmission, nil)
	mockSubmissionRepo.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Submission) bool {
		return s.Score != nil && *s.Score == 85.5 && s.GradedBy != nil
	})).Return(nil)
	
	// Execute
	err := submissionService.GradeSubmission(
		context.Background(),
		submissionID,
		85.5,
		"Good work! Well structured answer.",
		true,
		graderID,
	)
	
	// Assert
	assert.NoError(t, err)
	
	mockSubmissionRepo.AssertExpectations(t)
}

// TestGradeSubmission_AlreadyGraded tests grading an already graded submission
func TestGradeSubmission_AlreadyGraded(t *testing.T) {
	// Setup
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	
	submissionService := services.NewSubmissionService(mockSubmissionRepo, mockExerciseRepo)
	
	submissionID := uuid.New()
	graderID := uuid.New()
	gradedAt := testingTime()
	
	testSubmission := &models.Submission{
		ID:         submissionID,
		GradedAt:   &gradedAt, // Already graded
	}
	
	// Mock expectations
	mockSubmissionRepo.On("GetByID", mock.Anything, submissionID).Return(testSubmission, nil)
	
	// Execute
	err := submissionService.GradeSubmission(
		context.Background(),
		submissionID,
		90.0,
		"Great!",
		true,
		graderID,
	)
	
	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrAlreadyGraded, err)
	
	mockSubmissionRepo.AssertExpectations(t)
}

// TestGradeSubmission_NotFound tests grading a non-existent submission
func TestGradeSubmission_NotFound(t *testing.T) {
	// Setup
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	
	submissionService := services.NewSubmissionService(mockSubmissionRepo, mockExerciseRepo)
	
	submissionID := uuid.New()
	graderID := uuid.New()
	
	// Mock expectations
	mockSubmissionRepo.On("GetByID", mock.Anything, submissionID).Return((*models.Submission)(nil), assert.AnError)
	
	// Execute
	err := submissionService.GradeSubmission(
		context.Background(),
		submissionID,
		90.0,
		"Great!",
		true,
		graderID,
	)
	
	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrSubmissionNotFound, err)
	
	mockSubmissionRepo.AssertExpectations(t)
}

// TestGetSubmission_Success tests successful submission retrieval
func TestGetSubmission_Success(t *testing.T) {
	// Setup
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	
	submissionService := services.NewSubmissionService(mockSubmissionRepo, mockExerciseRepo)
	
	submissionID := uuid.New()
	
	testSubmission := &models.Submission{
		ID:             submissionID,
		ExerciseID:     uuid.New(),
		UserID:         uuid.New(),
		SubmissionType: "multiple_choice",
		IsCorrect:      boolPtr(true),
		Score:          float64Ptr(100.0),
	}
	
	// Mock expectations
	mockSubmissionRepo.On("GetByID", mock.Anything, submissionID).Return(testSubmission, nil)
	
	// Execute
	submission, err := submissionService.GetSubmission(context.Background(), submissionID)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, submission)
	assert.Equal(t, submissionID, submission.ID)
	
	mockSubmissionRepo.AssertExpectations(t)
}

// TestGetSubmission_NotFound tests retrieval of non-existent submission
func TestGetSubmission_NotFound(t *testing.T) {
	// Setup
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	
	submissionService := services.NewSubmissionService(mockSubmissionRepo, mockExerciseRepo)
	
	submissionID := uuid.New()
	
	// Mock expectations
	mockSubmissionRepo.On("GetByID", mock.Anything, submissionID).Return((*models.Submission)(nil), assert.AnError)
	
	// Execute
	submission, err := submissionService.GetSubmission(context.Background(), submissionID)
	
	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrSubmissionNotFound, err)
	assert.Nil(t, submission)
	
	mockSubmissionRepo.AssertExpectations(t)
}

// Helper functions
func testingTime() *time.Time {
	t := time.Now()
	return &t
}

func boolPtr(b bool) *bool {
	return &b
}

func float64Ptr(f float64) *float64 {
	return &f
}
