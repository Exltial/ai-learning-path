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

// GetSubmission retrieves a submission by ID
func (h *SubmissionHandler) GetSubmission(c *gin.Context) {
	submissionIDStr := c.Param("submission_id")
	submissionID, err := uuid.Parse(submissionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid submission ID",
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
func (h *SubmissionHandler) GetSubmissions(c *gin.Context) {
	exerciseIDStr := c.Param("exercise_id")
	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid exercise ID",
			},
		})
		return
	}

	userIDStr, _ := c.Get("user_id")
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
			"pagination": map[string]interface{}{
				"page":  1,
				"limit": len(submissions),
				"total": len(submissions),
			},
		},
	})
}

// GradeSubmission grades a submission (for instructors)
func (h *SubmissionHandler) GradeSubmission(c *gin.Context) {
	submissionIDStr := c.Param("submission_id")
	submissionID, err := uuid.Parse(submissionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid submission ID",
			},
		})
		return
	}

	gradedByStr, _ := c.Get("user_id")
	gradedBy, _ := uuid.Parse(gradedByStr.(string))

	var req struct {
		Score    float64 `json:"score" binding:"required,min=0"`
		Feedback string  `json:"feedback"`
		IsCorrect bool   `json:"is_correct"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	if err := h.submissionService.GradeSubmission(c.Request.Context(), submissionID, req.Score, req.Feedback, req.IsCorrect, gradedBy); err != nil {
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
		"message": "Submission graded successfully",
	})
}
