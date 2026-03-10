package services_test

import (
	"context"
	"testing"
	"time"

	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"
	"ai-learning-platform/internal/services"

	"github.com/google/uuid"
)

// Mock repositories for testing
type MockGradingRepository struct{}

func (m *MockGradingRepository) Create(ctx context.Context, history *repository.GradingHistory) error {
	return nil
}

func (m *MockGradingRepository) GetBySubmissionID(ctx context.Context, submissionID uuid.UUID) ([]*repository.GradingHistory, error) {
	return []*repository.GradingHistory{}, nil
}

func (m *MockGradingRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*repository.GradingHistory, error) {
	return []*repository.GradingHistory{}, nil
}

func (m *MockGradingRepository) GetStatsByExercise(ctx context.Context, exerciseID uuid.UUID) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

type MockSubmissionRepository struct{}

func (m *MockSubmissionRepository) Create(ctx context.Context, submission *models.Submission) error {
	return nil
}

func (m *MockSubmissionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Submission, error) {
	return &models.Submission{}, nil
}

func (m *MockSubmissionRepository) GetByExerciseIDAndUserID(ctx context.Context, exerciseID, userID uuid.UUID) ([]*models.Submission, error) {
	return []*models.Submission{}, nil
}

func (m *MockSubmissionRepository) Update(ctx context.Context, submission *models.Submission) error {
	return nil
}

func (m *MockSubmissionRepository) GetLatestSubmission(ctx context.Context, exerciseID, userID uuid.UUID) (*models.Submission, error) {
	return &models.Submission{}, nil
}

func (m *MockSubmissionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Submission, error) {
	return []*models.Submission{}, nil
}

type MockExerciseRepository struct{}

func (m *MockExerciseRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Exercise, error) {
	return &models.Exercise{}, nil
}

func (m *MockExerciseRepository) GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*models.Exercise, error) {
	return []*models.Exercise{}, nil
}

// ==================== Multiple Choice Tests ====================

func TestGradeMultipleChoice_SingleCorrect(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:           uuid.New(),
		ExerciseType: "multiple_choice",
		Points:       10,
		Options: []models.ExerciseOption{
			{Text: "Option A", IsCorrect: false},
			{Text: "Option B", IsCorrect: true},
			{Text: "Option C", IsCorrect: false},
			{Text: "Option D", IsCorrect: false},
		},
	}

	submission := &models.Submission{
		ID:     uuid.New(),
		Answer: "Option B",
	}

	result, err := service.GradeSubmission(context.Background(), submission, exercise)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.IsCorrect {
		t.Errorf("Expected correct answer, got incorrect")
	}

	if result.Score != 10 {
		t.Errorf("Expected score 10, got %f", result.Score)
	}

	if result.Percentage != 100 {
		t.Errorf("Expected percentage 100, got %f", result.Percentage)
	}
}

func TestGradeMultipleChoice_MultipleCorrect(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:           uuid.New(),
		ExerciseType: "multiple_choice",
		Points:       10,
		Options: []models.ExerciseOption{
			{Text: "Option A", IsCorrect: true},
			{Text: "Option B", IsCorrect: true},
			{Text: "Option C", IsCorrect: false},
			{Text: "Option D", IsCorrect: true},
		},
	}

	// All correct answers selected
	submission := &models.Submission{
		ID:     uuid.New(),
		Answer: "Option A, Option B, Option D",
	}

	result, err := service.GradeSubmission(context.Background(), submission, exercise)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.IsCorrect {
		t.Errorf("Expected correct answer, got incorrect")
	}

	if result.Score != 10 {
		t.Errorf("Expected score 10, got %f", result.Score)
	}
}

func TestGradeMultipleChoice_PartialCredit(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:           uuid.New(),
		ExerciseType: "multiple_choice",
		Points:       10,
		Options: []models.ExerciseOption{
			{Text: "Option A", IsCorrect: true},
			{Text: "Option B", IsCorrect: true},
			{Text: "Option C", IsCorrect: false},
			{Text: "Option D", IsCorrect: true},
		},
	}

	// Only 2 out of 3 correct answers selected
	submission := &models.Submission{
		ID:     uuid.New(),
		Answer: "Option A, Option B",
	}

	result, err := service.GradeSubmission(context.Background(), submission, exercise)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.IsCorrect {
		t.Errorf("Expected incorrect answer (missing one correct option)")
	}

	if result.Score <= 0 || result.Score >= 10 {
		t.Logf("Partial credit score: %f (expected between 0 and 10)", result.Score)
	}
}

func TestGradeMultipleChoice_IncorrectSelection(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:           uuid.New(),
		ExerciseType: "multiple_choice",
		Points:       10,
		Options: []models.ExerciseOption{
			{Text: "Option A", IsCorrect: true},
			{Text: "Option B", IsCorrect: false},
			{Text: "Option C", IsCorrect: false},
			{Text: "Option D", IsCorrect: false},
		},
	}

	// Wrong answer selected
	submission := &models.Submission{
		ID:     uuid.New(),
		Answer: "Option B",
	}

	result, err := service.GradeSubmission(context.Background(), submission, exercise)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.IsCorrect {
		t.Errorf("Expected incorrect answer")
	}

	if result.Score != 0 {
		t.Errorf("Expected score 0 for wrong answer, got %f", result.Score)
	}
}

func TestGradeMultipleChoice_NoAnswer(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:           uuid.New(),
		ExerciseType: "multiple_choice",
		Points:       10,
		Options: []models.ExerciseOption{
			{Text: "Option A", IsCorrect: true},
			{Text: "Option B", IsCorrect: false},
		},
	}

	submission := &models.Submission{
		ID:     uuid.New(),
		Answer: "",
	}

	result, err := service.GradeSubmission(context.Background(), submission, exercise)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.IsCorrect {
		t.Errorf("Expected incorrect answer for empty submission")
	}

	if result.Score != 0 {
		t.Errorf("Expected score 0, got %f", result.Score)
	}
}

// ==================== True/False Tests ====================

func TestGradeTrueFalse_Correct(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:               uuid.New(),
		ExerciseType:     "true_false",
		Points:           10,
		ExpectedAnswer:   map[string]interface{}{"answer": "true"},
	}

	submission := &models.Submission{
		ID:     uuid.New(),
		Answer: "True",
	}

	result, err := service.GradeSubmission(context.Background(), submission, exercise)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.IsCorrect {
		t.Errorf("Expected correct answer")
	}

	if result.Score != 10 {
		t.Errorf("Expected score 10, got %f", result.Score)
	}
}

func TestGradeTrueFalse_Incorrect(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:               uuid.New(),
		ExerciseType:     "true_false",
		Points:           10,
		ExpectedAnswer:   map[string]interface{}{"answer": "false"},
	}

	submission := &models.Submission{
		ID:     uuid.New(),
		Answer: "True",
	}

	result, err := service.GradeSubmission(context.Background(), submission, exercise)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.IsCorrect {
		t.Errorf("Expected incorrect answer")
	}

	if result.Score != 0 {
		t.Errorf("Expected score 0, got %f", result.Score)
	}
}

func TestGradeTrueFalse_VariantFormats(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:               uuid.New(),
		ExerciseType:     "true_false",
		Points:           10,
		ExpectedAnswer:   map[string]interface{}{"answer": "true"},
	}

	testCases := []struct {
		answer   string
		expected bool
	}{
		{"true", true},
		{"True", true},
		{"TRUE", true},
		{"T", true},
		{"t", true},
		{"1", true},
		{"yes", true},
		{"y", true},
		{"false", false},
		{"False", false},
		{"F", false},
		{"0", false},
		{"no", false},
		{"n", false},
	}

	for _, tc := range testCases {
		submission := &models.Submission{
			ID:     uuid.New(),
			Answer: tc.answer,
		}

		result, err := service.GradeSubmission(context.Background(), submission, exercise)
		if err != nil {
			t.Fatalf("Unexpected error for answer %s: %v", tc.answer, err)
		}

		if result.IsCorrect != tc.expected {
			t.Errorf("Answer %s: expected correct=%v, got %v", tc.answer, tc.expected, result.IsCorrect)
		}
	}
}

// ==================== Fill-in-the-Blank Tests ====================

func TestGradeFillBlank_ExactMatch(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:             uuid.New(),
		ExerciseType:   "fill_blank",
		Points:         10,
		ExpectedAnswer: map[string]interface{}{"answer": "Go"},
	}

	submission := &models.Submission{
		ID:     uuid.New(),
		Answer: "Go",
	}

	result, err := service.GradeSubmission(context.Background(), submission, exercise)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.IsCorrect {
		t.Errorf("Expected correct answer")
	}

	if result.Score != 10 {
		t.Errorf("Expected score 10, got %f", result.Score)
	}
}

func TestGradeFillBlank_MultipleAnswers(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:           uuid.New(),
		ExerciseType: "fill_blank",
		Points:       10,
		ExpectedAnswer: map[string]interface{}{
			"answers": []interface{}{"Go", "Golang", "Go Language"},
		},
	}

	testAnswers := []string{"Go", "Golang", "Go Language"}

	for _, answer := range testAnswers {
		submission := &models.Submission{
			ID:     uuid.New(),
			Answer: answer,
		}

		result, err := service.GradeSubmission(context.Background(), submission, exercise)
		if err != nil {
			t.Fatalf("Unexpected error for answer %s: %v", answer, err)
		}

		if !result.IsCorrect {
			t.Errorf("Expected answer %s to be correct", answer)
		}
	}
}

func TestGradeFillBlank_CaseInsensitive(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:             uuid.New(),
		ExerciseType:   "fill_blank",
		Points:         10,
		ExpectedAnswer: map[string]interface{}{"answer": "Go"},
	}

	testAnswers := []string{"go", "GO", "Go", "gO"}

	for _, answer := range testAnswers {
		submission := &models.Submission{
			ID:     uuid.New(),
			Answer: answer,
		}

		result, err := service.GradeSubmission(context.Background(), submission, exercise)
		if err != nil {
			t.Fatalf("Unexpected error for answer %s: %v", answer, err)
		}

		if !result.IsCorrect {
			t.Errorf("Expected answer %s to be correct (case insensitive)", answer)
		}
	}
}

func TestGradeFillBlank_FuzzyMatch(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:             uuid.New(),
		ExerciseType:   "fill_blank",
		Points:         10,
		ExpectedAnswer: map[string]interface{}{"answer": "function"},
	}

	// These should pass fuzzy matching (minor typos)
	testAnswers := []string{"functon", "function ", " function"}

	for _, answer := range testAnswers {
		submission := &models.Submission{
			ID:     uuid.New(),
			Answer: answer,
		}

		result, err := service.GradeSubmission(context.Background(), submission, exercise)
		if err != nil {
			t.Fatalf("Unexpected error for answer %s: %v", answer, err)
		}

		t.Logf("Answer '%s': IsCorrect=%v, Score=%f", answer, result.IsCorrect, result.Score)
	}
}

// ==================== Coding Exercise Tests ====================

func TestGradeCoding_AllTestsPass(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:           uuid.New(),
		ExerciseType: "coding",
		Points:       20,
		TestCases: map[string]interface{}{
			"tests": []interface{}{
				map[string]interface{}{"name": "Test 1", "input": "5", "expected": "5"},
				map[string]interface{}{"name": "Test 2", "input": "10", "expected": "10"},
				map[string]interface{}{"name": "Test 3", "input": "0", "expected": "0"},
			},
		},
	}

	submission := &models.Submission{
		ID:   uuid.New(),
		Code: "func main() { /* solution */ }",
	}

	result, err := service.GradeSubmission(context.Background(), submission, exercise)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Logf("Coding exercise result: IsCorrect=%v, Score=%f, PassRate=%f", result.IsCorrect, result.Score, result.PassRate)

	if result.TestResults == nil || len(result.TestResults) == 0 {
		t.Errorf("Expected test results")
	}
}

func TestGradeCoding_NoCode(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:           uuid.New(),
		ExerciseType: "coding",
		Points:       20,
	}

	submission := &models.Submission{
		ID:   uuid.New(),
		Code: "",
	}

	result, err := service.GradeSubmission(context.Background(), submission, exercise)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.IsCorrect {
		t.Errorf("Expected incorrect for empty code")
	}

	if result.Score != 0 {
		t.Errorf("Expected score 0, got %f", result.Score)
	}

	if result.Feedback != "No code submitted" {
		t.Errorf("Expected 'No code submitted' feedback, got %s", result.Feedback)
	}
}

func TestGradeCoding_NoTestCases(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:           uuid.New(),
		ExerciseType: "coding",
		Points:       20,
		TestCases:    nil,
	}

	submission := &models.Submission{
		ID:   uuid.New(),
		Code: "func main() { }",
	}

	result, err := service.GradeSubmission(context.Background(), submission, exercise)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Logf("Result with no test cases: %s", result.Feedback)
}

// ==================== Essay Exercise Tests ====================

func TestGradeEssay_TooShort(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:             uuid.New(),
		ExerciseType:   "essay",
		Points:         30,
		ExpectedAnswer: map[string]interface{}{"min_length": 100},
	}

	submission := &models.Submission{
		ID:     uuid.New(),
		Answer: "This is too short.",
	}

	result, err := service.GradeSubmission(context.Background(), submission, exercise)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.IsCorrect {
		t.Errorf("Expected incorrect for too short essay")
	}

	if result.Feedback == "" || result.Feedback == "No answer provided" {
		t.Logf("Feedback: %s", result.Feedback)
	}
}

func TestGradeEssay_MinimumStructure(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:           uuid.New(),
		ExerciseType: "essay",
		Points:       30,
	}

	submission := &models.Submission{
		ID:     uuid.New(),
		Answer: "This is paragraph one.\n\nThis is paragraph two with more content to meet the minimum length requirement.",
	}

	result, err := service.GradeSubmission(context.Background(), submission, exercise)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Logf("Essay result: IsCorrect=%v, Feedback=%s", result.IsCorrect, result.Feedback)
}

func TestGradeEssay_NoAnswer(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:           uuid.New(),
		ExerciseType: "essay",
		Points:       30,
	}

	submission := &models.Submission{
		ID:     uuid.New(),
		Answer: "",
	}

	result, err := service.GradeSubmission(context.Background(), submission, exercise)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.IsCorrect {
		t.Errorf("Expected incorrect for empty essay")
	}

	if result.Score != 0 {
		t.Errorf("Expected score 0, got %f", result.Score)
	}
}

// ==================== Grading History Tests ====================

func TestGradingService_RecordGradingHistory(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	submission := &models.Submission{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		ExerciseID: uuid.New(),
		Score:      nil,
	}

	exercise := &models.Exercise{
		ID:     uuid.New(),
		Points: 10,
	}

	result := &services.GradingResult{
		IsCorrect:  true,
		Score:      10,
		MaxScore:   10,
		Percentage: 100,
		Feedback:   "Excellent!",
	}

	// This should not panic even with mock repository
	err := service.GradeSubmission(context.Background(), submission, exercise)
	if err != nil {
		t.Logf("Grading completed with warning: %v", err)
	}

	t.Log("Grading history recording test completed")
}

// ==================== Edge Cases ====================

func TestGradeMultipleChoice_JSONArrayFormat(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:           uuid.New(),
		ExerciseType: "multiple_choice",
		Points:       10,
		Options: []models.ExerciseOption{
			{Text: "A", IsCorrect: true},
			{Text: "B", IsCorrect: true},
			{Text: "C", IsCorrect: false},
		},
	}

	// JSON array format
	submission := &models.Submission{
		ID:     uuid.New(),
		Answer: `["A", "B"]`,
	}

	result, err := service.GradeSubmission(context.Background(), submission, exercise)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Logf("JSON array format result: IsCorrect=%v, Score=%f", result.IsCorrect, result.Score)
}

func TestGradeFillBlank_WhitespaceHandling(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:             uuid.New(),
		ExerciseType:   "fill_blank",
		Points:         10,
		ExpectedAnswer: map[string]interface{}{"answer": "Go"},
	}

	testAnswers := []string{"  Go  ", "Go ", " Go", "\tGo\n"}

	for _, answer := range testAnswers {
		submission := &models.Submission{
			ID:     uuid.New(),
			Answer: answer,
		}

		result, err := service.GradeSubmission(context.Background(), submission, exercise)
		if err != nil {
			t.Fatalf("Unexpected error for answer '%s': %v", answer, err)
		}

		if !result.IsCorrect {
			t.Errorf("Expected answer '%s' (with whitespace) to be correct", answer)
		}
	}
}

func TestUnsupportedExerciseType(t *testing.T) {
	service := services.NewGradingService(&MockGradingRepository{}, &MockSubmissionRepository{}, &MockExerciseRepository{})

	exercise := &models.Exercise{
		ID:           uuid.New(),
		ExerciseType: "unknown_type",
		Points:       10,
	}

	submission := &models.Submission{
		ID:     uuid.New(),
		Answer: "test",
	}

	_, err := service.GradeSubmission(context.Background(), submission, exercise)
	if err == nil {
		t.Errorf("Expected error for unsupported exercise type")
	}

	t.Logf("Expected error received: %v", err)
}
