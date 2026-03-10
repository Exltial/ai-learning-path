package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"

	"github.com/google/uuid"
)

// GradingService handles auto-grading logic for all exercise types
type GradingService struct {
	gradingRepo    *repository.GradingRepository
	submissionRepo *repository.SubmissionRepository
	exerciseRepo   *repository.ExerciseRepository
}

// GradingResult represents the result of grading a submission
type GradingResult struct {
	IsCorrect        bool              `json:"is_correct"`
	Score            float64           `json:"score"`
	MaxScore         float64           `json:"max_score"`
	Percentage       float64           `json:"percentage"`
	Feedback         string            `json:"feedback"`
	DetailedFeedback *DetailedFeedback `json:"detailed_feedback,omitempty"`
	PassRate         float64           `json:"pass_rate,omitempty"` // For coding exercises
	TestResults      []*TestResult     `json:"test_results,omitempty"`
}

// DetailedFeedback provides granular feedback for each part of a submission
type DetailedFeedback struct {
	CorrectParts   []string `json:"correct_parts,omitempty"`
	IncorrectParts []string `json:"incorrect_parts,omitempty"`
	MissingParts   []string `json:"missing_parts,omitempty"`
	Suggestions    []string `json:"suggestions,omitempty"`
}

// TestResult represents the result of running a single test case
type TestResult struct {
	TestName      string `json:"test_name"`
	Passed        bool   `json:"passed"`
	Expected      string `json:"expected,omitempty"`
	Actual        string `json:"actual,omitempty"`
	ErrorMessage  string `json:"error_message,omitempty"`
	ExecutionTime int64  `json:"execution_time_ms,omitempty"`
}

// NewGradingService creates a new GradingService
func NewGradingService(
	gradingRepo *repository.GradingRepository,
	submissionRepo *repository.SubmissionRepository,
	exerciseRepo *repository.ExerciseRepository,
) *GradingService {
	return &GradingService{
		gradingRepo:    gradingRepo,
		submissionRepo: submissionRepo,
		exerciseRepo:   exerciseRepo,
	}
}

// GradeSubmission grades a submission based on exercise type
func (s *GradingService) GradeSubmission(ctx context.Context, submission *models.Submission, exercise *models.Exercise) (*GradingResult, error) {
	var result *GradingResult
	var err error

	switch exercise.ExerciseType {
	case "multiple_choice":
		result, err = s.gradeMultipleChoice(submission, exercise)
	case "true_false":
		result, err = s.gradeTrueFalse(submission, exercise)
	case "fill_blank":
		result, err = s.gradeFillBlank(submission, exercise)
	case "coding":
		result, err = s.gradeCoding(submission, exercise)
	case "essay":
		result, err = s.gradeEssay(submission, exercise)
	default:
		return nil, fmt.Errorf("unsupported exercise type: %s", exercise.ExerciseType)
	}

	if err != nil {
		return nil, err
	}

	// Record grading history
	if err := s.recordGradingHistory(ctx, submission, exercise, result, "auto"); err != nil {
		// Log error but don't fail the grading
		fmt.Printf("Warning: Failed to record grading history: %v\n", err)
	}

	return result, nil
}

// gradeMultipleChoice grades multiple choice exercises
// Supports single and multiple correct answers
func (s *GradingService) gradeMultipleChoice(submission *models.Submission, exercise *models.Exercise) (*GradingResult, error) {
	if exercise.Options == nil || len(exercise.Options) == 0 {
		return nil, fmt.Errorf("exercise has no options configured")
	}

	// Parse selected options from submission.Answer (JSON array or comma-separated)
	selectedOptions := s.parseSelectedOptions(submission.Answer)
	if len(selectedOptions) == 0 {
		return &GradingResult{
			IsCorrect:  false,
			Score:      0,
			MaxScore:   float64(exercise.Points),
			Percentage: 0,
			Feedback:   "No options selected",
		}, nil
	}

	// Find correct options
	var correctOptions []string
	for _, option := range exercise.Options {
		if option.IsCorrect {
			correctOptions = append(correctOptions, option.Text)
		}
	}

	// Grade the submission
	allCorrect := true
	correctCount := 0
	incorrectCount := 0
	missingCount := 0

	detailedFeedback := &DetailedFeedback{}

	// Check selected options
	for _, selected := range selectedOptions {
		found := false
		for _, option := range exercise.Options {
			if strings.EqualFold(option.Text, selected) {
				found = true
				if option.IsCorrect {
					correctCount++
					detailedFeedback.CorrectParts = append(detailedFeedback.CorrectParts, selected)
				} else {
					incorrectCount++
					allCorrect = false
					detailedFeedback.IncorrectParts = append(detailedFeedback.IncorrectParts, selected)
				}
				break
			}
		}
		if !found {
			incorrectCount++
			allCorrect = false
			detailedFeedback.IncorrectParts = append(detailedFeedback.IncorrectParts, selected)
		}
	}

	// Check for missing correct options
	for _, correct := range correctOptions {
		found := false
		for _, selected := range selectedOptions {
			if strings.EqualFold(selected, correct) {
				found = true
				break
			}
		}
		if !found {
			missingCount++
			allCorrect = false
			detailedFeedback.MissingParts = append(detailedFeedback.MissingParts, correct)
		}
	}

	// Calculate score with partial credit support
	maxScore := float64(exercise.Points)
	var score float64

	if allCorrect && incorrectCount == 0 && missingCount == 0 {
		// Perfect answer
		score = maxScore
	} else {
		// Partial credit: (correct selections / total correct options) * maxScore
		// Penalty for incorrect selections
		correctRatio := float64(correctCount) / float64(len(correctOptions))
		penalty := float64(incorrectCount) * 0.25 // 25% penalty per incorrect selection
		score = maxScore * (correctRatio - penalty)
		if score < 0 {
			score = 0
		}
	}

	percentage := (score / maxScore) * 100

	// Generate feedback
	feedback := s.generateMultipleChoiceFeedback(allCorrect, correctCount, incorrectCount, missingCount, len(correctOptions))

	return &GradingResult{
		IsCorrect:        allCorrect,
		Score:            score,
		MaxScore:         maxScore,
		Percentage:       percentage,
		Feedback:         feedback,
		DetailedFeedback: detailedFeedback,
	}, nil
}

// parseSelectedOptions parses selected options from submission answer
func (s *GradingService) parseSelectedOptions(answer string) []string {
	answer = strings.TrimSpace(answer)
	if answer == "" {
		return []string{}
	}

	// Try JSON array format first
	if strings.HasPrefix(answer, "[") && strings.HasSuffix(answer, "]") {
		// Simple JSON parsing without external library
		answer = strings.TrimPrefix(answer, "[")
		answer = strings.TrimSuffix(answer, "]")
	}

	// Split by comma
	parts := strings.Split(answer, ",")
	options := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		part = strings.Trim(part, "\"'") // Remove quotes
		if part != "" {
			options = append(options, part)
		}
	}

	return options
}

// generateMultipleChoiceFeedback generates detailed feedback for multiple choice
func (s *GradingService) generateMultipleChoiceFeedback(allCorrect bool, correctCount, incorrectCount, missingCount, totalCorrect int) string {
	if allCorrect {
		return "Correct! All right answers selected."
	}

	feedback := "Incorrect. "
	if correctCount > 0 {
		feedback += fmt.Sprintf("You selected %d correct option(s). ", correctCount)
	}
	if incorrectCount > 0 {
		feedback += fmt.Sprintf("You selected %d incorrect option(s). ", incorrectCount)
	}
	if missingCount > 0 {
		feedback += fmt.Sprintf("You missed %d correct option(s). ", missingCount)
	}
	feedback += "Review and try again."

	return feedback
}

// gradeTrueFalse grades true/false exercises
func (s *GradingService) gradeTrueFalse(submission *models.Submission, exercise *models.Exercise) (*GradingResult, error) {
	userAnswer := strings.TrimSpace(submission.Answer)
	if userAnswer == "" {
		return &GradingResult{
			IsCorrect:  false,
			Score:      0,
			MaxScore:   float64(exercise.Points),
			Percentage: 0,
			Feedback:   "No answer provided",
		}, nil
	}

	// Normalize user answer
	userAnswer = strings.ToLower(userAnswer)
	userIsTrue := userAnswer == "true" || userAnswer == "t" || userAnswer == "1" || userAnswer == "yes" || userAnswer == "y"

	// Get correct answer
	correctAnswer := false
	if exercise.ExpectedAnswer != nil {
		if expectedStr, ok := exercise.ExpectedAnswer["answer"].(string); ok {
			correctAnswer = strings.ToLower(expectedStr) == "true" ||
				strings.ToLower(expectedStr) == "t" ||
				expectedStr == "1" ||
				strings.ToLower(expectedStr) == "yes"
		}
	}

	maxScore := float64(exercise.Points)
	var score float64
	var feedback string

	if userIsTrue == correctAnswer {
		score = maxScore
		feedback = "Correct!"
	} else {
		score = 0
		feedback = fmt.Sprintf("Incorrect. The correct answer is: %s",
			map[bool]string{true: "True", false: "False"}[correctAnswer])
	}

	return &GradingResult{
		IsCorrect:  userIsTrue == correctAnswer,
		Score:      score,
		MaxScore:   maxScore,
		Percentage: (score / maxScore) * 100,
		Feedback:   feedback,
	}, nil
}

// gradeFillBlank grades fill-in-the-blank exercises
// Supports multiple correct answers and fuzzy matching
func (s *GradingService) gradeFillBlank(submission *models.Submission, exercise *models.Exercise) (*GradingResult, error) {
	userAnswer := strings.TrimSpace(submission.Answer)
	if userAnswer == "" {
		return &GradingResult{
			IsCorrect:  false,
			Score:      0,
			MaxScore:   float64(exercise.Points),
			Percentage: 0,
			Feedback:   "No answer provided",
		}, nil
	}

	if exercise.ExpectedAnswer == nil {
		return nil, fmt.Errorf("exercise has no expected answer configured")
	}

	// Get expected answers (can have multiple valid answers)
	var expectedAnswers []string
	if answers, ok := exercise.ExpectedAnswer["answers"].([]interface{}); ok {
		for _, a := range answers {
			if str, ok := a.(string); ok {
				expectedAnswers = append(expectedAnswers, str)
			}
		}
	} else if answer, ok := exercise.ExpectedAnswer["answer"].(string); ok {
		expectedAnswers = []string{answer}
	} else if answers, ok := exercise.ExpectedAnswer["answers"].([]string); ok {
		expectedAnswers = answers
	}

	if len(expectedAnswers) == 0 {
		return nil, fmt.Errorf("no valid expected answers found")
	}

	maxScore := float64(exercise.Points)
	userAnswerNormalized := strings.ToLower(strings.TrimSpace(userAnswer))

	// Check for exact match
	for _, expected := range expectedAnswers {
		expectedNormalized := strings.ToLower(strings.TrimSpace(expected))
		if userAnswerNormalized == expectedNormalized {
			return &GradingResult{
				IsCorrect:  true,
				Score:      maxScore,
				MaxScore:   maxScore,
				Percentage: 100,
				Feedback:   "Correct!",
			}, nil
		}
	}

	// Check for fuzzy match (case-insensitive, trim whitespace)
	for _, expected := range expectedAnswers {
		if s.fuzzyMatch(userAnswer, expected) {
			return &GradingResult{
				IsCorrect:  true,
				Score:      maxScore * 0.9, // Slight penalty for fuzzy match
				MaxScore:   maxScore,
				Percentage: 90,
				Feedback:   "Correct! (Accepted with minor variation)",
			}, nil
		}
	}

	// Check for partial credit (contains key terms)
	partialScore := s.calculatePartialCredit(userAnswer, expectedAnswers)
	if partialScore > 0 {
		return &GradingResult{
			IsCorrect:  false,
			Score:      partialScore,
			MaxScore:   maxScore,
			Percentage: (partialScore / maxScore) * 100,
			Feedback:   "Partially correct. Review the key concepts and try again.",
		}, nil
	}

	return &GradingResult{
		IsCorrect:  false,
		Score:      0,
		MaxScore:   maxScore,
		Percentage: 0,
		Feedback: fmt.Sprintf("Incorrect. Accepted answers include: %s",
			strings.Join(expectedAnswers[:min(len(expectedAnswers), 3)], ", ")),
	}, nil
}

// fuzzyMatch performs fuzzy matching between two strings
func (s *GradingService) fuzzyMatch(userAnswer, expected string) bool {
	// Normalize both strings
	userNorm := strings.ToLower(strings.TrimSpace(userAnswer))
	expectedNorm := strings.ToLower(strings.TrimSpace(expected))

	// Remove common punctuation
	userNorm = regexp.MustCompile(`[.,!?;:'"]`).ReplaceAllString(userNorm, "")
	expectedNorm = regexp.MustCompile(`[.,!?;:'"]`).ReplaceAllString(expectedNorm, "")

	// Check if one contains the other (for longer answers)
	if len(userNorm) > 5 && len(expectedNorm) > 5 {
		if strings.Contains(userNorm, expectedNorm) || strings.Contains(expectedNorm, userNorm) {
			return true
		}
	}

	// Check Levenshtein distance for short answers (allow 1-2 character difference)
	if len(userNorm) <= 10 && len(expectedNorm) <= 10 {
		distance := s.levenshteinDistance(userNorm, expectedNorm)
		if distance <= 2 {
			return true
		}
	}

	return false
}

// levenshteinDistance calculates the edit distance between two strings
func (s *GradingService) levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Create matrix
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			matrix[i][j] = min3(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

// calculatePartialCredit calculates partial credit based on keyword matching
func (s *GradingService) calculatePartialCredit(userAnswer string, expectedAnswers []string) float64 {
	userWords := strings.Fields(strings.ToLower(userAnswer))

	maxMatchRatio := 0.0
	for _, expected := range expectedAnswers {
		expectedWords := strings.Fields(strings.ToLower(expected))
		if len(expectedWords) == 0 {
			continue
		}

		matchCount := 0
		for _, expectedWord := range expectedWords {
			for _, userWord := range userWords {
				if userWord == expectedWord {
					matchCount++
					break
				}
			}
		}

		matchRatio := float64(matchCount) / float64(len(expectedWords))
		if matchRatio > maxMatchRatio {
			maxMatchRatio = matchRatio
		}
	}

	// Give partial credit if at least 50% of keywords match
	if maxMatchRatio >= 0.5 {
		return float64(100) * maxMatchRatio * 0.5 // Max 50% partial credit
	}

	return 0
}

// gradeCoding grades coding exercises using test cases
func (s *GradingService) gradeCoding(submission *models.Submission, exercise *models.Exercise) (*GradingResult, error) {
	if strings.TrimSpace(submission.Code) == "" {
		return &GradingResult{
			IsCorrect:  false,
			Score:      0,
			MaxScore:   float64(exercise.Points),
			Percentage: 0,
			Feedback:   "No code submitted",
		}, nil
	}

	maxScore := float64(exercise.Points)

	// If no test cases configured, return pending
	if exercise.TestCases == nil {
		return &GradingResult{
			IsCorrect:  false,
			Score:      0,
			MaxScore:   maxScore,
			Percentage: 0,
			Feedback:   "Code submitted, awaiting evaluation (no test cases configured)",
		}, nil
	}

	// Execute test cases
	testResults := s.executeTestCases(submission.Code, exercise.TestCases)

	// Calculate score based on passed tests
	passedCount := 0
	for _, result := range testResults {
		if result.Passed {
			passedCount++
		}
	}

	passRate := float64(passedCount) / float64(len(testResults))
	score := maxScore * passRate

	// Generate feedback
	feedback := s.generateCodingFeedback(passedCount, len(testResults), testResults)

	return &GradingResult{
		IsCorrect:   passRate >= 0.8, // Consider correct if 80%+ tests pass
		Score:       score,
		MaxScore:    maxScore,
		Percentage:  passRate * 100,
		Feedback:    feedback,
		PassRate:    passRate * 100,
		TestResults: testResults,
	}, nil
}

// executeTestCases executes test cases against submitted code
// This is a placeholder - in production, this would run code in a sandbox
func (s *GradingService) executeTestCases(code string, testCases map[string]interface{}) []*TestResult {
	results := make([]*TestResult, 0)

	// Parse test cases from the map
	// Expected format: {"tests": [{"name": "...", "input": "...", "expected": "..."}, ...]}
	if tests, ok := testCases["tests"].([]interface{}); ok {
		for i, test := range tests {
			if testMap, ok := test.(map[string]interface{}); ok {
				testName := fmt.Sprintf("Test %d", i+1)
				if name, ok := testMap["name"].(string); ok {
					testName = name
				}

				expected := ""
				if exp, ok := testMap["expected"].(string); ok {
					expected = exp
				}

				// In production, execute the code with the input and compare output
				// For now, simulate test execution
				actual := s.simulateCodeExecution(code, testMap)

				passed := strings.TrimSpace(actual) == strings.TrimSpace(expected)

				results = append(results, &TestResult{
					TestName:      testName,
					Passed:        passed,
					Expected:      expected,
					Actual:        actual,
					ExecutionTime: 10, // Simulated execution time
				})
			}
		}
	}

	// If no tests found in expected format, create a default test
	if len(results) == 0 {
		results = append(results, &TestResult{
			TestName:     "Default Test",
			Passed:       false,
			ErrorMessage: "Test cases not properly configured",
		})
	}

	return results
}

// simulateCodeExecution simulates code execution (placeholder for real execution)
func (s *GradingService) simulateCodeExecution(code string, testCase map[string]interface{}) string {
	// In production, this would:
	// 1. Write code to a temporary file
	// 2. Execute in a sandboxed environment (Docker container)
	// 3. Capture stdout/stderr
	// 4. Return the output

	// Placeholder: return a simulated output
	if input, ok := testCase["input"].(string); ok {
		// Simple simulation: echo the input
		return input
	}
	return ""
}

// generateCodingFeedback generates detailed feedback for coding exercises
func (s *GradingService) generateCodingFeedback(passedCount, totalCount int, testResults []*TestResult) string {
	if passedCount == totalCount {
		return "Excellent! All test cases passed."
	}

	feedback := fmt.Sprintf("Passed %d/%d test cases. ", passedCount, totalCount)

	// Add details about failed tests
	failedTests := make([]string, 0)
	for _, result := range testResults {
		if !result.Passed {
			failedTests = append(failedTests, result.TestName)
		}
	}

	if len(failedTests) > 0 {
		feedback += fmt.Sprintf("Failed tests: %s. ", strings.Join(failedTests, ", "))
		feedback += "Review your code and try again."
	}

	return feedback
}

// gradeEssay provides a framework for essay grading (requires manual review)
func (s *GradingService) gradeEssay(submission *models.Submission, exercise *models.Exercise) (*GradingResult, error) {
	userAnswer := strings.TrimSpace(submission.Answer)
	if userAnswer == "" {
		return &GradingResult{
			IsCorrect:  false,
			Score:      0,
			MaxScore:   float64(exercise.Points),
			Percentage: 0,
			Feedback:   "No answer provided",
		}, nil
	}

	// Essay requires manual grading, but we can do basic validation
	maxScore := float64(exercise.Points)

	// Check minimum length
	minLength := 50
	if exercise.ExpectedAnswer != nil {
		if minLen, ok := exercise.ExpectedAnswer["min_length"].(float64); ok {
			minLength = int(minLen)
		}
	}

	if len(userAnswer) < minLength {
		return &GradingResult{
			IsCorrect:  false,
			Score:      0,
			MaxScore:   maxScore,
			Percentage: 0,
			Feedback: fmt.Sprintf("Essay is too short. Minimum length: %d characters. Current length: %d characters.",
				minLength, len(userAnswer)),
		}, nil
	}

	// Check for basic structure (paragraphs)
	paragraphs := strings.Split(userAnswer, "\n\n")
	if len(paragraphs) < 2 {
		return &GradingResult{
			IsCorrect:  false,
			Score:      maxScore * 0.3, // Partial credit for attempting
			MaxScore:   maxScore,
			Percentage: 30,
			Feedback:   "Essay submitted but lacks proper structure. Consider organizing into multiple paragraphs. Awaiting instructor review.",
		}, nil
	}

	// Essay submitted successfully, awaiting manual review
	return &GradingResult{
		IsCorrect:  false, // Will be determined by manual review
		Score:      0,     // Will be assigned by instructor
		MaxScore:   maxScore,
		Percentage: 0,
		Feedback:   "Essay submitted successfully. Awaiting instructor review for final grading.",
	}, nil
}

// ManualGradeEssay allows instructors to manually grade an essay submission
func (s *GradingService) ManualGradeEssay(
	ctx context.Context,
	submissionID uuid.UUID,
	score float64,
	feedback string,
	gradedBy uuid.UUID,
) error {
	submission, err := s.submissionRepo.GetByID(ctx, submissionID)
	if err != nil {
		return fmt.Errorf("submission not found: %w", err)
	}

	// Validate score
	if score < 0 {
		return fmt.Errorf("score cannot be negative")
	}

	// Update submission
	now := time.Now()
	submission.Score = &score
	submission.Feedback = feedback
	submission.IsCorrect = func() *bool { b := score >= 60; return &b }()
	submission.GradedAt = &now
	submission.GradedBy = &gradedBy

	if err := s.submissionRepo.Update(ctx, submission); err != nil {
		return fmt.Errorf("failed to update submission: %w", err)
	}

	// Record grading history
	exercise, err := s.exerciseRepo.GetByID(ctx, submission.ExerciseID)
	if err != nil {
		return fmt.Errorf("exercise not found: %w", err)
	}

	gradingResult := &GradingResult{
		IsCorrect:  *submission.IsCorrect,
		Score:      score,
		MaxScore:   float64(exercise.Points),
		Percentage: (score / float64(exercise.Points)) * 100,
		Feedback:   feedback,
	}

	return s.recordGradingHistory(ctx, submission, exercise, gradingResult, "manual")
}

// recordGradingHistory records a grading event in the history
func (s *GradingService) recordGradingHistory(
	ctx context.Context,
	submission *models.Submission,
	exercise *models.Exercise,
	result *GradingResult,
	gradingType string,
) error {
	var previousScore *float64
	if submission.Score != nil {
		previousScore = submission.Score
	}

	scoreChange := result.Score
	if previousScore != nil {
		scoreChange = result.Score - *previousScore
	}

	history := &repository.GradingHistory{
		ID:            uuid.New(),
		SubmissionID:  submission.ID,
		ExerciseID:    exercise.ID,
		UserID:        submission.UserID,
		GradingType:   gradingType,
		PreviousScore: previousScore,
		NewScore:      result.Score,
		ScoreChange:   scoreChange,
		GradedBy:      submission.GradedBy,
		GradedAt:      time.Now(),
		Metadata:      "", // Could include detailed feedback as JSON
	}

	return s.gradingRepo.Create(ctx, history)
}

// GetGradingHistory retrieves grading history for a submission
func (s *GradingService) GetGradingHistory(ctx context.Context, submissionID uuid.UUID) ([]*repository.GradingHistory, error) {
	return s.gradingRepo.GetBySubmissionID(ctx, submissionID)
}

// GetGradingHistoryByUser retrieves grading history for a user
func (s *GradingService) GetGradingHistoryByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*repository.GradingHistory, error) {
	return s.gradingRepo.GetByUserID(ctx, userID, limit)
}

// GetGradingStats retrieves grading statistics for an exercise
func (s *GradingService) GetGradingStats(ctx context.Context, exerciseID uuid.UUID) (map[string]interface{}, error) {
	return s.gradingRepo.GetStatsByExercise(ctx, exerciseID)
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper function for min of three values
func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
