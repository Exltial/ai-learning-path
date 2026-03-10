package handlers

import (
	"net/http"
	"ai-learning-platform/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LessonHandler handles lesson-related requests
type LessonHandler struct {
	lessonService *services.LessonService
}

// NewLessonHandler creates a new LessonHandler
func NewLessonHandler(lessonService *services.LessonService) *LessonHandler {
	return &LessonHandler{lessonService: lessonService}
}

// GetLesson retrieves a lesson by ID
func (h *LessonHandler) GetLesson(c *gin.Context) {
	lessonIDStr := c.Param("lesson_id")
	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid lesson ID",
			},
		})
		return
	}

	lesson, err := h.lessonService.GetLesson(c.Request.Context(), lessonID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "NOT_FOUND",
				"message": "Lesson not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": lesson,
	})
}

// UpdateLesson updates an existing lesson
func (h *LessonHandler) UpdateLesson(c *gin.Context) {
	lessonIDStr := c.Param("lesson_id")
	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid lesson ID",
			},
		})
		return
	}

	var req struct {
		Title         string `json:"title"`
		Description   string `json:"description"`
		Content       string `json:"content"`
		VideoURL      string `json:"video_url"`
		VideoDuration int    `json:"video_duration"`
		IsFreePreview bool   `json:"is_free_preview"`
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

	lesson, err := h.lessonService.GetLesson(c.Request.Context(), lessonID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "NOT_FOUND",
				"message": "Lesson not found",
			},
		})
		return
	}

	// Update fields
	if req.Title != "" {
		lesson.Title = req.Title
	}
	if req.Description != "" {
		lesson.Description = req.Description
	}
	if req.Content != "" {
		lesson.Content = req.Content
	}
	if req.VideoURL != "" {
		lesson.VideoURL = req.VideoURL
	}
	if req.VideoDuration > 0 {
		lesson.VideoDuration = req.VideoDuration
	}
	lesson.IsFreePreview = req.IsFreePreview

	if err := h.lessonService.UpdateLesson(c.Request.Context(), lesson); err != nil {
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
		"data": lesson,
	})
}

// DeleteLesson deletes a lesson
func (h *LessonHandler) DeleteLesson(c *gin.Context) {
	lessonIDStr := c.Param("lesson_id")
	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid lesson ID",
			},
		})
		return
	}

	if err := h.lessonService.DeleteLesson(c.Request.Context(), lessonID); err != nil {
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
