package handlers

import (
	"net/http"
	"strconv"

	"ai-learning-platform/internal/repository"
	"ai-learning-platform/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GradingHandler handles grading-related HTTP requests
type GradingHandler struct {
	gradingService    *services.GradingService
	submissionService *services.SubmissionService
	exerciseRepo      *repository.ExerciseRepository
}

// NewGradingHandler creates a new GradingHandler
func NewGradingHandler(
	gradingService *services.GradingService,
	submissionService *services.SubmissionService,
	exerciseRepo *repository.ExerciseRepository,
) *GradingHandler {
	return &GradingHandler{
		gradingService:    gradingService,
		submissionService: submissionService,
		exerciseRepo:      exerciseRepo,
	}
}

// AutoGradeRequest represents a request to auto-grade a submission
type AutoGradeRequest struct {
	// Force re-grading even if already graded
	Force bool `json:"force,omitempty"`
}

// AutoGradeSubmission auto-grades a submission
// @Summary Auto-grade a submission
// @Description Trigger automatic grading for a submission based on exercise type
// @Tags Grading
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param submission_id path string true "Submission ID" format(uuid)
// @Param request body AutoGradeRequest false "Grading options"
// @Success 200 {object} map[string]interface{} "Grading result"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Submission not found"
// @Failure 500 {object} map[string]interface{} "Internal error"
// @Router /api/v1/submissions/{submission_id}/auto-grade [post]
func (h *GradingHandler) AutoGradeSubmission(c *gin.Context) {
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

	var req AutoGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid request body",
			},
		})
		return
	}

	// Get submission
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

	// Check if already graded (unless force)
	if !req.Force && submission.GradedAt != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "ALREADY_GRADED",
				"message": "Submission already graded. Use force=true to re-grade.",
			},
		})
		return
	}

	// Get exercise
	exercise, err := h.exerciseRepo.GetByID(c.Request.Context(), submission.ExerciseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "NOT_FOUND",
				"message": "Exercise not found",
			},
		})
		return
	}

	// Skip auto-grading for essay type
	if exercise.ExerciseType == "essay" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "MANUAL_GRADING_REQUIRED",
				"message": "Essay submissions require manual grading by an instructor",
			},
		})
		return
	}

	// Perform auto-grading
	result, err := h.gradingService.GradeSubmission(c.Request.Context(), submission, exercise)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "GRADING_FAILED",
				"message": "Failed to grade submission: " + err.Error(),
			},
		})
		return
	}

	// Update submission with grading result
	now := submission.GradedAt
	if now == nil {
		t := submission.SubmittedAt
		now = &t
	}
	submission.IsCorrect = &result.IsCorrect
	submission.Score = &result.Score
	submission.Feedback = result.Feedback
	submission.GradedAt = now

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Submission graded successfully",
		"data": map[string]interface{}{
			"submission_id":     submissionID,
			"is_correct":        result.IsCorrect,
			"score":             result.Score,
			"max_score":         result.MaxScore,
			"percentage":        result.Percentage,
			"feedback":          result.Feedback,
			"detailed_feedback": result.DetailedFeedback,
			"pass_rate":         result.PassRate,
			"test_results":      result.TestResults,
		},
	})
}

// ManualGradeRequest represents a manual grading request
type ManualGradeRequest struct {
	// Score to assign (0-100 or based on exercise points)
	Score float64 `json:"score" binding:"required,min=0" example:"85.5"`

	// Feedback for the student
	Feedback string `json:"feedback" example:"Good work! Consider improving the structure."`

	// Whether the submission is correct (optional, auto-calculated from score)
	IsCorrect *bool `json:"is_correct,omitempty" example:"true"`

	// Reason for manual grading (optional)
	Reason string `json:"reason,omitempty" example:"Re-grade after student appeal"`
}

// ManualGradeSubmission manually grades a submission (for instructors)
// @Summary Manually grade a submission
// @Description Manually grade a submission (for essay or code review)
// @Tags Grading
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param submission_id path string true "Submission ID" format(uuid)
// @Param request body ManualGradeRequest true "Grading data"
// @Success 200 {object} map[string]interface{} "Grading successful"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (instructor only)"
// @Failure 404 {object} map[string]interface{} "Submission not found"
// @Router /api/v1/submissions/{submission_id}/manual-grade [post]
func (h *GradingHandler) ManualGradeSubmission(c *gin.Context) {
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

	var req ManualGradeRequest
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

	// Get submission to determine exercise type
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

	// Get exercise to validate score range
	exercise, err := h.exerciseRepo.GetByID(c.Request.Context(), submission.ExerciseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "NOT_FOUND",
				"message": "Exercise not found",
			},
		})
		return
	}

	// Validate score range
	maxScore := float64(exercise.Points)
	if req.Score > maxScore {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "INVALID_SCORE",
				"message": "Score cannot exceed maximum points for this exercise",
			},
		})
		return
	}

	// For essay type, use the manual grading service method
	if exercise.ExerciseType == "essay" {
		if err := h.gradingService.ManualGradeEssay(
			c.Request.Context(),
			submissionID,
			req.Score,
			req.Feedback,
			gradedBy,
		); err != nil {
			status := http.StatusInternalServerError
			code := "INTERNAL_ERROR"
			message := err.Error()

			if err.Error() == "submission not found" {
				status = http.StatusNotFound
				code = "NOT_FOUND"
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
	} else {
		// For other types, update submission directly
		isCorrect := req.IsCorrect != nil && *req.IsCorrect
		if req.IsCorrect == nil {
			// Auto-determine based on score (60% passing threshold)
			isCorrect = req.Score >= maxScore*0.6
		}

		if err := h.submissionService.GradeSubmission(
			c.Request.Context(),
			submissionID,
			req.Score,
			req.Feedback,
			isCorrect,
			gradedBy,
		); err != nil {
			status := http.StatusInternalServerError
			code := "INTERNAL_ERROR"
			message := err.Error()

			if err.Error() == "submission already graded" {
				status = http.StatusBadRequest
				code = "ALREADY_GRADED"
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
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Submission graded successfully",
		"data": map[string]interface{}{
			"submission_id": submissionID,
			"score":         req.Score,
			"max_score":     maxScore,
			"percentage":    (req.Score / maxScore) * 100,
			"feedback":      req.Feedback,
			"is_correct":    req.IsCorrect,
			"graded_by":     gradedBy,
		},
	})
}

// GetGradingHistory retrieves grading history for a submission
// @Summary Get grading history
// @Description Get the grading history for a specific submission
// @Tags Grading
// @Produce json
// @Security BearerAuth
// @Param submission_id path string true "Submission ID" format(uuid)
// @Success 200 {object} map[string]interface{} "Grading history"
// @Failure 400 {object} map[string]interface{} "Invalid submission ID"
// @Failure 404 {object} map[string]interface{} "Submission not found"
// @Router /api/v1/submissions/{submission_id}/grading-history [get]
func (h *GradingHandler) GetGradingHistory(c *gin.Context) {
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

	histories, err := h.gradingService.GetGradingHistory(c.Request.Context(), submissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve grading history",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"submission_id": submissionID,
			"history":       histories,
			"count":         len(histories),
		},
	})
}

// GetMyGradingHistory retrieves grading history for the current user
// @Summary Get my grading history
// @Description Get recent grading history for the current user
// @Tags Grading
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of records" default(20)
// @Success 200 {object} map[string]interface{} "Grading history"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/grading/my-history [get]
func (h *GradingHandler) GetMyGradingHistory(c *gin.Context) {
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

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	histories, err := h.gradingService.GetGradingHistoryByUser(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve grading history",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"user_id": userID,
			"history": histories,
			"count":   len(histories),
		},
	})
}

// GetExerciseGradingStats retrieves grading statistics for an exercise
// @Summary Get exercise grading statistics
// @Description Get grading statistics for a specific exercise
// @Tags Grading
// @Produce json
// @Security BearerAuth
// @Param exercise_id path string true "Exercise ID" format(uuid)
// @Success 200 {object} map[string]interface{} "Grading statistics"
// @Failure 400 {object} map[string]interface{} "Invalid exercise ID"
// @Failure 404 {object} map[string]interface{} "Exercise not found"
// @Router /api/v1/exercises/{exercise_id}/grading-stats [get]
func (h *GradingHandler) GetExerciseGradingStats(c *gin.Context) {
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

	stats, err := h.gradingService.GetGradingStats(c.Request.Context(), exerciseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve grading statistics",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// BatchGradeRequest represents a batch grading request
type BatchGradeRequest struct {
	// Submission IDs to grade
	SubmissionIDs []uuid.UUID `json:"submission_ids" binding:"required"`

	// Score to assign to all
	Score float64 `json:"score" binding:"required,min=0"`

	// Feedback to add to all
	Feedback string `json:"feedback"`
}

// BatchGradeSubmissions grades multiple submissions at once (for instructors)
// @Summary Batch grade submissions
// @Description Grade multiple submissions at once (for instructors)
// @Tags Grading
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body BatchGradeRequest true "Batch grading data"
// @Success 200 {object} map[string]interface{} "Batch grading result"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (instructor only)"
// @Router /api/v1/grading/batch-grade [post]
func (h *GradingHandler) BatchGradeSubmissions(c *gin.Context) {
	var req BatchGradeRequest
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

	if len(req.SubmissionIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "No submission IDs provided",
			},
		})
		return
	}

	// Get instructor ID
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

	// Process each submission
	results := make([]map[string]interface{}, 0)
	successCount := 0
	failCount := 0

	for _, submissionID := range req.SubmissionIDs {
		submission, err := h.submissionService.GetSubmission(c.Request.Context(), submissionID)
		if err != nil {
			failCount++
			results = append(results, map[string]interface{}{
				"submission_id": submissionID,
				"success":       false,
				"error":         "Submission not found",
			})
			continue
		}

		// Get exercise for max score validation
		exercise, err := h.exerciseRepo.GetByID(c.Request.Context(), submission.ExerciseID)
		if err != nil {
			failCount++
			results = append(results, map[string]interface{}{
				"submission_id": submissionID,
				"success":       false,
				"error":         "Exercise not found",
			})
			continue
		}

		// Validate score
		maxScore := float64(exercise.Points)
		if req.Score > maxScore {
			failCount++
			results = append(results, map[string]interface{}{
				"submission_id": submissionID,
				"success":       false,
				"error":         "Score exceeds maximum points",
			})
			continue
		}

		// Grade the submission
		isCorrect := req.Score >= maxScore*0.6
		if err := h.submissionService.GradeSubmission(
			c.Request.Context(),
			submissionID,
			req.Score,
			req.Feedback,
			isCorrect,
			gradedBy,
		); err != nil {
			failCount++
			results = append(results, map[string]interface{}{
				"submission_id": submissionID,
				"success":       false,
				"error":         err.Error(),
			})
		} else {
			successCount++
			results = append(results, map[string]interface{}{
				"submission_id": submissionID,
				"success":       true,
				"score":         req.Score,
				"percentage":    (req.Score / maxScore) * 100,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Batch grading completed",
		"data": map[string]interface{}{
			"total":      len(req.SubmissionIDs),
			"successful": successCount,
			"failed":     failCount,
			"results":    results,
		},
	})
}

// ReGradeSubmission requests re-grading of a submission
// @Summary Re-grade a submission
// @Description Request re-grading of a submission (for students to appeal or instructors to re-evaluate)
// @Tags Grading
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param submission_id path string true "Submission ID" format(uuid)
// @Param request body AutoGradeRequest false "Re-grading options"
// @Success 200 {object} map[string]interface{} "Re-grading result"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Submission not found"
// @Router /api/v1/submissions/{submission_id}/regrade [post]
func (h *GradingHandler) ReGradeSubmission(c *gin.Context) {
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

	// Get submission
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

	// Get exercise
	exercise, err := h.exerciseRepo.GetByID(c.Request.Context(), submission.ExerciseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "NOT_FOUND",
				"message": "Exercise not found",
			},
		})
		return
	}

	// Skip for essay type
	if exercise.ExerciseType == "essay" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "MANUAL_GRADING_REQUIRED",
				"message": "Essay submissions require manual re-grading by an instructor",
			},
		})
		return
	}

	// Perform re-grading with force=true
	result, err := h.gradingService.GradeSubmission(c.Request.Context(), submission, exercise)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "GRADING_FAILED",
				"message": "Failed to re-grade submission: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Submission re-graded successfully",
		"data": map[string]interface{}{
			"submission_id":  submissionID,
			"previous_score": submission.Score,
			"new_score":      result.Score,
			"score_change":   result.Score - *submission.Score,
			"is_correct":     result.IsCorrect,
			"feedback":       result.Feedback,
		},
	})
}

// GetGradingAnalytics retrieves comprehensive grading analytics
// @Summary Get grading analytics
// @Description Get comprehensive grading analytics for a course or exercise
// @Tags Grading
// @Produce json
// @Security BearerAuth
// @Param exercise_id query string false "Filter by exercise ID" format(uuid)
// @Param days query int false "Number of days to analyze" default(30)
// @Success 200 {object} map[string]interface{} "Grading analytics"
// @Failure 400 {object} map[string]interface{} "Invalid parameters"
// @Router /api/v1/grading/analytics [get]
func (h *GradingHandler) GetGradingAnalytics(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 365 {
		days = 30
	}

	// TODO: Implement comprehensive analytics
	// This would include:
	// - Grade distribution
	// - Average scores over time
	// - Question difficulty analysis
	// - Student performance trends

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"period_days": days,
			"analytics":   "Coming soon",
		},
	})
}
