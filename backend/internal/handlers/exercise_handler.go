package handlers

import (
	"net/http"
	"ai-learning-platform/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ExerciseHandler handles exercise-related requests
type ExerciseHandler struct {
	exerciseService  *services.ExerciseService
	submissionService *services.SubmissionService
}

// NewExerciseHandler creates a new ExerciseHandler
func NewExerciseHandler(exerciseService *services.ExerciseService, submissionService *services.SubmissionService) *ExerciseHandler {
	return &ExerciseHandler{
		exerciseService:   exerciseService,
		submissionService: submissionService,
	}
}

// GetExercise retrieves an exercise by ID
func (h *ExerciseHandler) GetExercise(c *gin.Context) {
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

	exercise, err := h.exerciseService.GetExercise(c.Request.Context(), exerciseID)
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": exercise,
	})
}

// UpdateExercise updates an existing exercise
func (h *ExerciseHandler) UpdateExercise(c *gin.Context) {
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

	var req struct {
		Title        string      `json:"title"`
		Description  string      `json:"description"`
		ExerciseType string      `json:"exercise_type"`
		Difficulty   string      `json:"difficulty"`
		Points       int         `json:"points"`
		MaxAttempts  int         `json:"max_attempts"`
		TimeLimit    *int        `json:"time_limit"`
		Options      interface{} `json:"options"`
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

	exercise, err := h.exerciseService.GetExercise(c.Request.Context(), exerciseID)
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

	// Update fields
	if req.Title != "" {
		exercise.Title = req.Title
	}
	if req.Description != "" {
		exercise.Description = req.Description
	}
	if req.ExerciseType != "" {
		exercise.ExerciseType = req.ExerciseType
	}
	if req.Difficulty != "" {
		exercise.Difficulty = req.Difficulty
	}
	if req.Points > 0 {
		exercise.Points = req.Points
	}
	if req.MaxAttempts > 0 {
		exercise.MaxAttempts = req.MaxAttempts
	}
	exercise.TimeLimit = req.TimeLimit
	exercise.Options = req.Options

	if err := h.exerciseService.UpdateExercise(c.Request.Context(), exercise); err != nil {
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
		"data": exercise,
	})
}

// DeleteExercise deletes an exercise
func (h *ExerciseHandler) DeleteExercise(c *gin.Context) {
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

	if err := h.exerciseService.DeleteExercise(c.Request.Context(), exerciseID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// SubmitExercise submits an answer for an exercise
func (h *ExerciseHandler) SubmitExercise(c *gin.Context) {
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

	var req struct {
		Answer         string   `json:"answer"`
		Code           string   `json:"code"`
		SelectedOptions []string `json:"selected_options"`
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

	submission, err := h.submissionService.SubmitExercise(c.Request.Context(), exerciseID, userID, req.Answer, req.Code, req.SelectedOptions)
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

	remainingAttempts := 0 // TODO: Calculate remaining attempts

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"submission_id":      submission.ID,
			"is_correct":         submission.IsCorrect,
			"score":              submission.Score,
			"feedback":           submission.Feedback,
			"attempt_number":     submission.AttemptNumber,
			"remaining_attempts": remainingAttempts,
		},
	})
}
