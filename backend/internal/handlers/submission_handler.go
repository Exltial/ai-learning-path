package handlers

import (
	"net/http"

	"ai-learning-platform/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SubmissionHandler handles submission-related requests
type SubmissionHandler struct {
	submissionService *services.SubmissionService
}

// NewSubmissionHandler creates a new SubmissionHandler
func NewSubmissionHandler(submissionService *services.SubmissionService) *SubmissionHandler {
	return &SubmissionHandler{submissionService: submissionService}
}

// SubmitExerciseRequest represents a submission request
type SubmitExerciseRequest struct {
	// Answer for text-based exercises (fill_blank, true_false, essay)
	Answer string `json:"answer,omitempty"`
	
	// Code for coding exercises
	Code string `json:"code,omitempty"`
	
	// Selected options for multiple choice exercises
	SelectedOptions []string `json:"selected_options,omitempty"`
}

// SubmitExercise submits an answer for an exercise
// @Summary Submit exercise answer
// @Description Submit an answer for an exercise (auto-graded or manual review)
// @Tags Submissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param exercise_id path string true "Exercise ID" format(uuid)
// @Param request body SubmitExerciseRequest true "Submission data"
// @Success 201 {object} map[string]interface{} "Submission successful"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Exercise not found"
// @Failure 429 {object} map[string]interface{} "Max attempts reached"
// @Router /api/v1/exercises/{exercise_id}/submit [post]
func (h *SubmissionHandler) SubmitExercise(c *gin.Context) {
	exerciseIDStr := c.Param("exercise_id")
	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid exercise ID format",
			},
		})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "UNAUTHORIZED",
				"message": "User ID not found",
			},
		})
		return
	}
	
	userID, _ := uuid.Parse(userIDStr.(string))

	var req SubmitExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid input: " + err.Error(),
			},
		})
		return
	}

	// Submit the exercise
	submission, err := h.submissionService.SubmitExercise(
		c.Request.Context(),
		exerciseID,
		userID,
		req.Answer,
		req.Code,
		req.SelectedOptions,
	)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		message := err.Error()

		switch err {
		case services.ErrExerciseNotFound:
			status = http.StatusNotFound
			code = "NOT_FOUND"
			message = "Exercise not found"
		case services.ErrMaxAttemptsReached:
			status = http.StatusTooManyRequests
			code = "MAX_ATTEMPTS_REACHED"
			message = "Maximum attempts reached for this exercise"
		case services.ErrInvalidSubmission:
			status = http.StatusBadRequest
			code = "INVALID_SUBMISSION"
		}

		c.JSON(status, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    code,
				"message": message,
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Submission successful",
		"data": map[string]interface{}{
			"submission": submission,
			"is_correct": submission.IsCorrect,
			"score":      submission.Score,
			"feedback":   submission.Feedback,
			"attempt":    submission.AttemptNumber,
		},
	})
}

// GetSubmission retrieves a submission by ID
// @Summary Get submission details
// @Description Get detailed information about a specific submission
// @Tags Submissions
// @Produce json
// @Security BearerAuth
// @Param submission_id path string true "Submission ID" format(uuid)
// @Success 200 {object} map[string]interface{} "Submission details"
// @Failure 400 {object} map[string]interface{} "Invalid submission ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Submission not found"
// @Router /api/v1/submissions/{submission_id} [get]
func (h *SubmissionHandler) GetSubmission(c *gin.Context) {
	submissionIDStr := c.Param("submission_id")
	submissionID, err := uuid.Parse(submissionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid submission ID format",
			},
		})
		return
	}

	submission, err := h.submissionService.GetSubmission(c.Request.Context(), submissionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "NOT_FOUND",
				"message": "Submission not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": submission,
	})
}

// GetSubmissions retrieves submissions for an exercise
// @Summary Get exercise submissions
// @Description Get all submissions for an exercise by the current user
// @Tags Submissions
// @Produce json
// @Security BearerAuth
// @Param exercise_id path string true "Exercise ID" format(uuid)
// @Success 200 {object} map[string]interface{} "List of submissions"
// @Failure 400 {object} map[string]interface{} "Invalid exercise ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/exercises/{exercise_id}/submissions [get]
func (h *SubmissionHandler) GetSubmissions(c *gin.Context) {
	exerciseIDStr := c.Param("exercise_id")
	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid exercise ID format",
			},
		})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "UNAUTHORIZED",
				"message": "User ID not found",
			},
		})
		return
	}
	
	userID, _ := uuid.Parse(userIDStr.(string))

	submissions, err := h.submissionService.GetSubmissions(c.Request.Context(), exerciseID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"submissions": submissions,
			"count":       len(submissions),
		},
	})
}

// GradeSubmissionRequest represents a grading request
type GradeSubmissionRequest struct {
	// Score to assign (0-100 or based on exercise points)
	Score float64 `json:"score" binding:"required,min=0" example:"85.5"`
	
	// Feedback for the student
	Feedback string `json:"feedback" example:"Good work! But consider optimizing the algorithm."`
	
	// Whether the submission is correct
	IsCorrect bool `json:"is_correct" example:"true"`
}

// GradeSubmission grades a submission (for instructors)
// @Summary Grade a submission
// @Description Grade a submission manually (for essay or code exercises)
// @Tags Submissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param submission_id path string true "Submission ID" format(uuid)
// @Param request body GradeSubmissionRequest true "Grading data"
// @Success 200 {object} map[string]interface{} "Submission graded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (instructor only)"
// @Failure 404 {object} map[string]interface{} "Submission not found"
// @Router /api/v1/submissions/{submission_id}/grade [post]
func (h *SubmissionHandler) GradeSubmission(c *gin.Context) {
	submissionIDStr := c.Param("submission_id")
	submissionID, err := uuid.Parse(submissionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid submission ID format",
			},
		})
		return
	}

	// Get instructor ID from context
	gradedByStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "UNAUTHORIZED",
				"message": "User ID not found",
			},
		})
		return
	}
	
	gradedBy, _ := uuid.Parse(gradedByStr.(string))

	var req GradeSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid input: " + err.Error(),
			},
		})
		return
	}

	if err := h.submissionService.GradeSubmission(
		c.Request.Context(),
		submissionID,
		req.Score,
		req.Feedback,
		req.IsCorrect,
		gradedBy,
	); err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		message := err.Error()

		switch err {
		case services.ErrSubmissionNotFound:
			status = http.StatusNotFound
			code = "NOT_FOUND"
			message = "Submission not found"
		case services.ErrAlreadyGraded:
			status = http.StatusBadRequest
			code = "ALREADY_GRADED"
			message = "Submission already graded"
		}

		c.JSON(status, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    code,
				"message": message,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Submission graded successfully",
		"data": map[string]interface{}{
			"submission_id": submissionID,
			"score":         req.Score,
			"feedback":      req.Feedback,
			"is_correct":    req.IsCorrect,
			"graded_by":     gradedBy,
		},
	})
}

// GetMySubmissions retrieves all submissions for the current user
// @Summary Get my submissions
// @Description Get all submissions for the current user with optional filters
// @Tags Submissions
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status (pending, graded)" 
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{} "List of submissions"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/submissions [get]
func (h *SubmissionHandler) GetMySubmissions(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "UNAUTHORIZED",
				"message": "User ID not found",
			},
		})
		return
	}
	
	_ = userIDStr // Will be used when implementing pagination/filtering

	// TODO: Implement pagination and filtering
	// For now, return placeholder
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"submissions": []interface{}{},
			"pagination": map[string]interface{}{
				"page":  1,
				"limit": 20,
				"total": 0,
			},
		},
	})
}
