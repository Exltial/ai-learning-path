package handlers

import (
	"ai-learning-platform/internal/services"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProgressTrackingHandler handles progress tracking requests
type ProgressTrackingHandler struct {
	progressTrackingService *services.ProgressTrackingService
}

// NewProgressTrackingHandler creates a new ProgressTrackingHandler
func NewProgressTrackingHandler(progressTrackingService *services.ProgressTrackingService) *ProgressTrackingHandler {
	return &ProgressTrackingHandler{progressTrackingService: progressTrackingService}
}

// GetCourseProgressRequest represents a request to get course progress
type GetCourseProgressRequest struct {
	CourseID string `json:"course_id" binding:"required"`
}

// UpdateVideoProgressRequest represents a request to update video progress
type UpdateVideoProgressRequest struct {
	LessonID string `json:"lesson_id" binding:"required"`
	Position int    `json:"position" binding:"required"`
	Duration int    `json:"duration"`
}

// GetHeatmapDataRequest represents a request to get heatmap data
type GetHeatmapDataRequest struct {
	Months int `json:"months" binding:"omitempty,min=1,max=24"`
}

// GetDailyStatsRequest represents a request to get daily statistics
type GetDailyStatsRequest struct {
	Days int `json:"days" binding:"omitempty,min=1,max=365"`
}

// GetReportRequest represents a request to get a report
type GetReportRequest struct {
	ReportType string `json:"report_type" binding:"required,oneof=weekly monthly"`
	Offset     int    `json:"offset" binding:"omitempty,min=-12,max=52"` // weeks or months offset
}

// ExportReportRequest represents a request to export a report
type ExportReportRequest struct {
	ReportType string `json:"report_type" binding:"required,oneof=weekly monthly"`
	Offset     int    `json:"offset" binding:"omitempty,min=-12,max=52"`
	Format     string `json:"format" binding:"omitempty,oneof=csv json"`
}

// GetCourseProgress retrieves detailed progress for a course
// @Summary Get course progress
// @Description Get detailed progress for a specific course including lesson completion and video progress
// @Tags Progress Tracking
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/progress/courses/{course_id} [get]
func (h *ProgressTrackingHandler) GetCourseProgress(c *gin.Context) {
	courseIDStr := c.Param("course_id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid course ID",
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

	progress, err := h.progressTrackingService.GetCourseProgress(c.Request.Context(), userID, courseID)
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
		"data":    progress,
	})
}

// UpdateVideoProgress updates video playback position in real-time
// @Summary Update video progress
// @Description Update video playback position (called periodically during video playback)
// @Tags Progress Tracking
// @Accept json
// @Produce json
// @Param request body UpdateVideoProgressRequest true "Video progress update"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/progress/video [put]
func (h *ProgressTrackingHandler) UpdateVideoProgress(c *gin.Context) {
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

	var req UpdateVideoProgressRequest
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

	lessonID, err := uuid.Parse(req.LessonID)
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

	err = h.progressTrackingService.UpdateVideoProgress(c.Request.Context(), userID, lessonID, req.Position, req.Duration)
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
		"message": "Video progress updated successfully",
	})
}

// MarkLessonCompleted marks a lesson as completed
// @Summary Mark lesson as completed
// @Description Mark a lesson as completed
// @Tags Progress Tracking
// @Accept json
// @Produce json
// @Param lesson_id path string true "Lesson ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/progress/lessons/{lesson_id}/complete [post]
func (h *ProgressTrackingHandler) MarkLessonCompleted(c *gin.Context) {
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

	err = h.progressTrackingService.MarkLessonCompleted(c.Request.Context(), userID, lessonID)
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
		"message": "Lesson marked as completed",
	})
}

// GetLearningHeatmap retrieves learning heatmap data
// @Summary Get learning heatmap
// @Description Get learning activity data for heatmap visualization
// @Tags Progress Tracking
// @Accept json
// @Produce json
// @Param months query int false "Number of months (default: 6)"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/progress/heatmap [get]
func (h *ProgressTrackingHandler) GetLearningHeatmap(c *gin.Context) {
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

	monthsStr := c.DefaultQuery("months", "6")
	months, err := strconv.Atoi(monthsStr)
	if err != nil || months < 1 || months > 24 {
		months = 6
	}

	heatmapData, err := h.progressTrackingService.GetLearningHeatmapData(c.Request.Context(), userID, months)
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
			"heatmap": heatmapData,
			"months":  months,
		},
	})
}

// GetDailyStats retrieves daily learning statistics
// @Summary Get daily statistics
// @Description Get daily learning statistics for the specified number of days
// @Tags Progress Tracking
// @Accept json
// @Produce json
// @Param days query int false "Number of days (default: 30)"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/progress/daily-stats [get]
func (h *ProgressTrackingHandler) GetDailyStats(c *gin.Context) {
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

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 365 {
		days = 30
	}

	dailyStats, err := h.progressTrackingService.GetDailyLearningStats(c.Request.Context(), userID, days)
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
			"daily_stats": dailyStats,
			"days":        days,
		},
	})
}

// GetWeeklyReport retrieves a weekly learning report
// @Summary Get weekly report
// @Description Get a weekly learning report
// @Tags Progress Tracking
// @Accept json
// @Produce json
// @Param offset query int false "Weeks offset (0=current week, -1=last week)"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/progress/reports/weekly [get]
func (h *ProgressTrackingHandler) GetWeeklyReport(c *gin.Context) {
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

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < -52 || offset > 52 {
		offset = 0
	}

	report, err := h.progressTrackingService.GenerateWeeklyReport(c.Request.Context(), userID, offset)
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
		"data":    report,
	})
}

// GetMonthlyReport retrieves a monthly learning report
// @Summary Get monthly report
// @Description Get a monthly learning report
// @Tags Progress Tracking
// @Accept json
// @Produce json
// @Param offset query int false "Months offset (0=current month, -1=last month)"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/progress/reports/monthly [get]
func (h *ProgressTrackingHandler) GetMonthlyReport(c *gin.Context) {
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

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < -12 || offset > 12 {
		offset = 0
	}

	report, err := h.progressTrackingService.GenerateMonthlyReport(c.Request.Context(), userID, offset)
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
		"data":    report,
	})
}

// ExportReport exports a learning report
// @Summary Export report
// @Description Export a learning report in CSV or JSON format
// @Tags Progress Tracking
// @Accept json
// @Produce json
// @Param request body ExportReportRequest true "Export request"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/progress/reports/export [post]
func (h *ProgressTrackingHandler) ExportReport(c *gin.Context) {
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

	var req ExportReportRequest
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

	if req.Format == "" {
		req.Format = "csv"
	}

	var csvData string
	var err error
	var reportDate string

	if req.ReportType == "weekly" {
		report, err := h.progressTrackingService.GenerateWeeklyReport(c.Request.Context(), userID, req.Offset)
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
		csvData, err = h.progressTrackingService.ExportReportToCSV(&services.MonthlyReport{
			Month:            report.WeekStart,
			TotalHours:       report.TotalHours,
			LessonsCompleted: report.LessonsCompleted,
			DailyStats:       report.DailyStats,
			CoursesProgress:  report.CoursesProgress,
		})
		reportDate = report.WeekStart
	} else {
		report, err := h.progressTrackingService.GenerateMonthlyReport(c.Request.Context(), userID, req.Offset)
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
		csvData, err = h.progressTrackingService.ExportReportToCSV(report)
		reportDate = time.Now().Format("2006-01")
	}

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

	// Set headers for file download
	filename := fmt.Sprintf("learning_report_%s_%s.%s", req.ReportType, reportDate, req.Format)
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	c.Data(http.StatusOK, "text/csv", []byte(csvData))
}

// GetLearningTimeStats retrieves total learning time statistics
// @Summary Get learning time stats
// @Description Get total learning time statistics
// @Tags Progress Tracking
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/progress/stats [get]
func (h *ProgressTrackingHandler) GetLearningTimeStats(c *gin.Context) {
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

	stats, err := h.progressTrackingService.GetLearningTimeStats(c.Request.Context(), userID)
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
		"data":    stats,
	})
}
