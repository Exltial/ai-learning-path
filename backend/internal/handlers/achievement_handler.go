package handlers

import (
	"net/http"
	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AchievementHandler handles achievement-related requests
type AchievementHandler struct {
	achievementService *services.AchievementService
}

// NewAchievementHandler creates a new AchievementHandler
func NewAchievementHandler(achievementService *services.AchievementService) *AchievementHandler {
	return &AchievementHandler{achievementService: achievementService}
}

// GetUserAchievements retrieves all achievements for the current user
func (h *AchievementHandler) GetUserAchievements(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid user ID",
			},
		})
		return
	}

	achievements, err := h.achievementService.GetUserAchievementsWithProgress(c.Request.Context(), userID)
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
		"data":    achievements,
	})
}

// GetUserAchievementSummary retrieves achievement summary for the current user
func (h *AchievementHandler) GetUserAchievementSummary(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid user ID",
			},
		})
		return
	}

	summary, err := h.achievementService.GetUserAchievementSummary(c.Request.Context(), userID)
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
		"data":    summary,
	})
}

// GetLeaderboard retrieves the leaderboard
func (h *AchievementHandler) GetLeaderboard(c *gin.Context) {
	leaderboardType := c.DefaultQuery("type", "all_time")
	limit := 100

	var lt models.LeaderboardType
	switch leaderboardType {
	case "weekly":
		lt = models.LeaderboardTypeWeekly
	case "monthly":
		lt = models.LeaderboardTypeMonthly
	case "friends":
		lt = models.LeaderboardTypeFriends
	default:
		lt = models.LeaderboardTypeAllTime
	}

	var userID *uuid.UUID
	if lt == models.LeaderboardTypeFriends {
		userIDStr, exists := c.Get("user_id")
		if exists {
			id, err := uuid.Parse(userIDStr.(string))
			if err == nil {
				userID = &id
			}
		}
	}

	entries, err := h.achievementService.GetLeaderboard(c.Request.Context(), lt, limit, userID)
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
		"data":    entries,
	})
}

// GetUserLevel retrieves the current user's level information
func (h *AchievementHandler) GetUserLevel(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid user ID",
			},
		})
		return
	}

	level, err := h.achievementService.GetUserLevel(c.Request.Context(), userID)
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
		"data":    level,
	})
}

// GetUserStreak retrieves the current user's learning streak
func (h *AchievementHandler) GetUserStreak(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid user ID",
			},
		})
		return
	}

	streak, err := h.achievementService.GetUserStreak(c.Request.Context(), userID)
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
		"data":    streak,
	})
}

// GetPointsHistory retrieves the user's points transaction history
func (h *AchievementHandler) GetPointsHistory(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid user ID",
			},
		})
		return
	}

	limit := 50
	offset := 0

	history, err := h.achievementService.GetPointsHistory(c.Request.Context(), userID, limit, offset)
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
		"data":    history,
	})
}

// TriggerAchievementCheck manually triggers an achievement check (for testing/debugging)
func (h *AchievementHandler) TriggerAchievementCheck(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid user ID",
			},
		})
		return
	}

	var req struct {
		EventType string                 `json:"event_type"`
		EventData map[string]interface{} `json:"event_data"`
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

	unlocked, err := h.achievementService.CheckAndUnlockAchievements(c.Request.Context(), userID, req.EventType, req.EventData)
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
		"success":  true,
		"unlocked": len(unlocked),
		"data":     unlocked,
	})
}

// UpdateUserActivity updates user activity (for streak tracking)
func (h *AchievementHandler) UpdateUserActivity(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid user ID",
			},
		})
		return
	}

	err = h.achievementService.UpdateStreak(c.Request.Context(), userID)
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

	// Check for streak-based achievements
	eventData := map[string]interface{}{
		"activity_type": "daily_login",
	}
	
	unlocked, err := h.achievementService.CheckAndUnlockAchievements(c.Request.Context(), userID, "daily_activity", eventData)
	if err != nil {
		// Don't fail the request if achievement check fails
		unlocked = []*models.Achievement{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"unlocked": len(unlocked),
		"data":     unlocked,
	})
}

// AwardPoints awards points to the user (admin only or for specific events)
func (h *AchievementHandler) AwardPoints(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid user ID",
			},
		})
		return
	}

	var req struct {
		Amount      int     `json:"amount"`
		SourceType  string  `json:"source_type"`
		SourceID    *string `json:"source_id"`
		Description string  `json:"description"`
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

	var sourceID *uuid.UUID
	if req.SourceID != nil {
		id, err := uuid.Parse(*req.SourceID)
		if err == nil {
			sourceID = &id
		}
	}

	err = h.achievementService.AwardPoints(c.Request.Context(), userID, req.Amount, req.SourceType, sourceID, req.Description)
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
		"message": "Points awarded successfully",
	})
}
