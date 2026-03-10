package handlers

import (
	"net/http"
	"strconv"

	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DiscussionHandler handles discussion-related requests
type DiscussionHandler struct {
	discussionService *services.DiscussionService
}

// NewDiscussionHandler creates a new DiscussionHandler
func NewDiscussionHandler(discussionService *services.DiscussionService) *DiscussionHandler {
	return &DiscussionHandler{discussionService: discussionService}
}

// CreateDiscussionRequest represents a request to create a discussion
type CreateDiscussionRequest struct {
	CourseID    string     `json:"course_id" binding:"required"`
	LessonID    *string    `json:"lesson_id,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Content     string     `json:"content" binding:"required,min=1,max=50000"`
	ParentID    *string    `json:"parent_id,omitempty"`
	IsAnonymous bool       `json:"is_anonymous,omitempty"`
	TagIDs      []string   `json:"tag_ids,omitempty"`
}

// CreateDiscussion creates a new discussion or reply
// @Summary Create a discussion
// @Description Create a new discussion thread or reply
// @Tags Discussions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateDiscussionRequest true "Discussion data"
// @Success 201 {object} map[string]interface{} "Discussion created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/discussions [post]
func (h *DiscussionHandler) CreateDiscussion(c *gin.Context) {
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

	var req CreateDiscussionRequest
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

	// Parse course ID
	courseID, err := uuid.Parse(req.CourseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid course ID format",
			},
		})
		return
	}

	// Parse lesson ID if provided
	var lessonID *uuid.UUID
	if req.LessonID != nil {
		id, err := uuid.Parse(*req.LessonID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": map[string]interface{}{
					"code":    "BAD_REQUEST",
					"message": "Invalid lesson ID format",
				},
			})
			return
		}
		lessonID = &id
	}

	// Parse parent ID if provided
	var parentID *uuid.UUID
	if req.ParentID != nil {
		id, err := uuid.Parse(*req.ParentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": map[string]interface{}{
					"code":    "BAD_REQUEST",
					"message": "Invalid parent ID format",
				},
			})
			return
		}
		parentID = &id
	}

	// Parse tag IDs if provided
	var tagIDs []uuid.UUID
	for _, tagIDStr := range req.TagIDs {
		tagID, err := uuid.Parse(tagIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": map[string]interface{}{
					"code":    "BAD_REQUEST",
					"message": "Invalid tag ID format: " + tagIDStr,
				},
			})
			return
		}
		tagIDs = append(tagIDs, tagID)
	}

	createReq := &models.CreateDiscussionRequest{
		CourseID:    courseID,
		LessonID:    lessonID,
		Title:       req.Title,
		Content:     req.Content,
		ParentID:    parentID,
		IsAnonymous: req.IsAnonymous,
		TagIDs:      tagIDs,
	}

	discussion, err := h.discussionService.CreateDiscussion(c.Request.Context(), createReq, userID)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"

		switch err {
		case services.ErrCourseNotFound:
			status = http.StatusNotFound
			code = "NOT_FOUND"
		case services.ErrParentNotFound:
			status = http.StatusNotFound
			code = "PARENT_NOT_FOUND"
		case services.ErrDiscussionLocked:
			status = http.StatusForbidden
			code = "DISCUSSION_LOCKED"
		case services.ErrInvalidDepth:
			status = http.StatusBadRequest
			code = "MAX_DEPTH_EXCEEDED"
		case services.ErrTitleRequired:
			status = http.StatusBadRequest
			code = "TITLE_REQUIRED"
		case services.ErrContentTooLong:
			status = http.StatusBadRequest
			code = "CONTENT_TOO_LONG"
		}

		c.JSON(status, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Discussion created successfully",
		"data":    discussion,
	})
}

// GetDiscussion retrieves a discussion by ID
// @Summary Get discussion details
// @Description Get detailed information about a specific discussion
// @Tags Discussions
// @Produce json
// @Param discussion_id path string true "Discussion ID" format(uuid)
// @Success 200 {object} map[string]interface{} "Discussion details"
// @Failure 400 {object} map[string]interface{} "Invalid discussion ID"
// @Failure 404 {object} map[string]interface{} "Discussion not found"
// @Router /api/v1/discussions/{discussion_id} [get]
func (h *DiscussionHandler) GetDiscussion(c *gin.Context) {
	discussionIDStr := c.Param("discussion_id")
	discussionID, err := uuid.Parse(discussionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid discussion ID format",
			},
		})
		return
	}

	// Get current user ID if available
	var currentUserID uuid.UUID
	userIDStr, exists := c.Get("user_id")
	if exists {
		currentUserID, _ = uuid.Parse(userIDStr.(string))
	}

	// Check if replies should be included
	withReplies := c.DefaultQuery("with_replies", "false") == "true"
	maxDepth, _ := strconv.Atoi(c.DefaultQuery("max_depth", "3"))
	if maxDepth < 1 {
		maxDepth = 1
	}
	if maxDepth > 10 {
		maxDepth = 10
	}

	var discussion interface{}
	if withReplies {
		discussion, err = h.discussionService.GetDiscussionWithReplies(c.Request.Context(), discussionID, maxDepth, currentUserID)
	} else {
		discussion, err = h.discussionService.GetDiscussion(c.Request.Context(), discussionID, currentUserID)
	}

	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"

		switch err {
		case services.ErrDiscussionNotFound:
			status = http.StatusNotFound
			code = "NOT_FOUND"
		case services.ErrDiscussionDeleted:
			status = http.StatusNotFound
			code = "DELETED"
		}

		c.JSON(status, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    discussion,
	})
}

// ListDiscussions retrieves discussions with pagination and filters
// @Summary List discussions
// @Description Get paginated list of discussions with optional filters
// @Tags Discussions
// @Produce json
// @Param course_id query string true "Course ID" format(uuid)
// @Param lesson_id query string false "Lesson ID" format(uuid)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param sort_by query string false "Sort by" enum(created_at,updated_at,upvotes,reply_count,hot)
// @Param search query string false "Search keyword"
// @Success 200 {object} map[string]interface{} "List of discussions"
// @Failure 400 {object} map[string]interface{} "Invalid parameters"
// @Router /api/v1/discussions [get]
func (h *DiscussionHandler) ListDiscussions(c *gin.Context) {
	courseIDStr := c.Query("course_id")
	if courseIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "course_id is required",
			},
		})
		return
	}

	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid course ID format",
			},
		})
		return
	}

	// Parse optional parameters
	var lessonID *uuid.UUID
	if lessonIDStr := c.Query("lesson_id"); lessonIDStr != "" {
		id, err := uuid.Parse(lessonIDStr)
		if err == nil {
			lessonID = &id
		}
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	search := c.Query("search")

	// Get current user ID if available
	var currentUserID uuid.UUID
	userIDStr, exists := c.Get("user_id")
	if exists {
		currentUserID, _ = uuid.Parse(userIDStr.(string))
	}

	opts := &models.DiscussionListOptions{
		CourseID:  courseID,
		LessonID:  lessonID,
		Page:      page,
		Limit:     limit,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Search:    search,
	}

	discussions, total, err := h.discussionService.ListDiscussions(c.Request.Context(), opts, currentUserID)
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

	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    discussions,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetHotDiscussions retrieves hot/trending discussions for a course
// @Summary Get hot discussions
// @Description Get hot/trending discussions for a course
// @Tags Discussions
// @Produce json
// @Param course_id query string true "Course ID" format(uuid)
// @Param limit query int false "Limit" default(10)
// @Success 200 {object} map[string]interface{} "List of hot discussions"
// @Failure 400 {object} map[string]interface{} "Invalid parameters"
// @Router /api/v1/discussions/hot [get]
func (h *DiscussionHandler) GetHotDiscussions(c *gin.Context) {
	courseIDStr := c.Query("course_id")
	if courseIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "course_id is required",
			},
		})
		return
	}

	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid course ID format",
			},
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	// Get current user ID if available
	var currentUserID uuid.UUID
	userIDStr, exists := c.Get("user_id")
	if exists {
		currentUserID, _ = uuid.Parse(userIDStr.(string))
	}

	hotDiscussions, err := h.discussionService.GetHotDiscussions(c.Request.Context(), courseID, limit, currentUserID)
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
		"data":    hotDiscussions,
	})
}

// UpdateDiscussionRequest represents a request to update a discussion
type UpdateDiscussionRequest struct {
	Title      *string `json:"title,omitempty"`
	Content    *string `json:"content,omitempty" binding:"omitempty,min=1,max=50000"`
	IsResolved *bool   `json:"is_resolved,omitempty"`
	IsLocked   *bool   `json:"is_locked,omitempty"`
	IsPinned   *bool   `json:"is_pinned,omitempty"`
	TagIDs     []string `json:"tag_ids,omitempty"`
}

// UpdateDiscussion updates an existing discussion
// @Summary Update a discussion
// @Description Update discussion information (author only)
// @Tags Discussions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param discussion_id path string true "Discussion ID" format(uuid)
// @Param request body UpdateDiscussionRequest true "Updated discussion data"
// @Success 200 {object} map[string]interface{} "Discussion updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Discussion not found"
// @Router /api/v1/discussions/{discussion_id} [put]
func (h *DiscussionHandler) UpdateDiscussion(c *gin.Context) {
	discussionIDStr := c.Param("discussion_id")
	discussionID, err := uuid.Parse(discussionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid discussion ID format",
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

	var req UpdateDiscussionRequest
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

	// Parse tag IDs if provided
	var tagIDs []uuid.UUID
	for _, tagIDStr := range req.TagIDs {
		tagID, err := uuid.Parse(tagIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": map[string]interface{}{
					"code":    "BAD_REQUEST",
					"message": "Invalid tag ID format: " + tagIDStr,
				},
			})
			return
		}
		tagIDs = append(tagIDs, tagID)
	}

	updateReq := &models.UpdateDiscussionRequest{
		Title:      req.Title,
		Content:    req.Content,
		IsResolved: req.IsResolved,
		IsLocked:   req.IsLocked,
		IsPinned:   req.IsPinned,
		TagIDs:     tagIDs,
	}

	discussion, err := h.discussionService.UpdateDiscussion(c.Request.Context(), discussionID, updateReq, userID)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"

		switch err {
		case services.ErrDiscussionNotFound:
			status = http.StatusNotFound
			code = "NOT_FOUND"
		case services.ErrDiscussionDeleted:
			status = http.StatusNotFound
			code = "DELETED"
		case services.ErrUnauthorized:
			status = http.StatusForbidden
			code = "FORBIDDEN"
		case services.ErrContentTooLong:
			status = http.StatusBadRequest
			code = "CONTENT_TOO_LONG"
		}

		c.JSON(status, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Discussion updated successfully",
		"data":    discussion,
	})
}

// DeleteDiscussion deletes a discussion
// @Summary Delete a discussion
// @Description Delete a discussion (author or admin only)
// @Tags Discussions
// @Produce json
// @Security BearerAuth
// @Param discussion_id path string true "Discussion ID" format(uuid)
// @Success 204 {object} nil "Discussion deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid discussion ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Discussion not found"
// @Router /api/v1/discussions/{discussion_id} [delete]
func (h *DiscussionHandler) DeleteDiscussion(c *gin.Context) {
	discussionIDStr := c.Param("discussion_id")
	discussionID, err := uuid.Parse(discussionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid discussion ID format",
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

	// Check if user is admin (you may need to adjust this based on your auth system)
	isAdmin := false
	if role, ok := c.Get("role"); ok {
		isAdmin = role == "admin" || role == "instructor"
	}

	if err := h.discussionService.DeleteDiscussion(c.Request.Context(), discussionID, userID, isAdmin); err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"

		switch err {
		case services.ErrDiscussionNotFound:
			status = http.StatusNotFound
			code = "NOT_FOUND"
		case services.ErrUnauthorized:
			status = http.StatusForbidden
			code = "FORBIDDEN"
		}

		c.JSON(status, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ToggleLike toggles a like on a discussion
// @Summary Toggle like on discussion
// @Description Toggle upvote/downvote on a discussion
// @Tags Discussions
// @Produce json
// @Security BearerAuth
// @Param discussion_id path string true "Discussion ID" format(uuid)
// @Param like_type query string false "Like type" enum(upvote,downvote) default(upvote)
// @Success 200 {object} map[string]interface{} "Like toggled successfully"
// @Failure 400 {object} map[string]interface{} "Invalid parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Discussion not found"
// @Router /api/v1/discussions/{discussion_id}/like [post]
func (h *DiscussionHandler) ToggleLike(c *gin.Context) {
	discussionIDStr := c.Param("discussion_id")
	discussionID, err := uuid.Parse(discussionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid discussion ID format",
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

	likeType := c.DefaultQuery("like_type", "upvote")
	if likeType != "upvote" && likeType != "downvote" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid like type. Must be 'upvote' or 'downvote'",
			},
		})
		return
	}

	if err := h.discussionService.ToggleLike(c.Request.Context(), discussionID, userID, likeType); err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"

		switch err {
		case services.ErrDiscussionNotFound:
			status = http.StatusNotFound
			code = "NOT_FOUND"
		case services.ErrDiscussionDeleted:
			status = http.StatusNotFound
			code = "DELETED"
		}

		c.JSON(status, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Like toggled successfully",
	})
}

// ToggleFavorite toggles a favorite on a discussion
// @Summary Toggle favorite on discussion
// @Description Toggle favorite/bookmark on a discussion
// @Tags Discussions
// @Produce json
// @Security BearerAuth
// @Param discussion_id path string true "Discussion ID" format(uuid)
// @Success 200 {object} map[string]interface{} "Favorite toggled successfully"
// @Failure 400 {object} map[string]interface{} "Invalid discussion ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Discussion not found"
// @Router /api/v1/discussions/{discussion_id}/favorite [post]
func (h *DiscussionHandler) ToggleFavorite(c *gin.Context) {
	discussionIDStr := c.Param("discussion_id")
	discussionID, err := uuid.Parse(discussionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid discussion ID format",
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

	isFavorited, err := h.discussionService.ToggleFavorite(c.Request.Context(), discussionID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"

		switch err {
		case services.ErrDiscussionNotFound:
			status = http.StatusNotFound
			code = "NOT_FOUND"
		case services.ErrDiscussionDeleted:
			status = http.StatusNotFound
			code = "DELETED"
		}

		c.JSON(status, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Favorite toggled successfully",
		"data": map[string]interface{}{
			"is_favorited": isFavorited,
		},
	})
}

// GetUserFavorites retrieves a user's favorited discussions
// @Summary Get user's favorites
// @Description Get a user's favorited discussions
// @Tags Discussions
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{} "List of favorited discussions"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/discussions/favorites [get]
func (h *DiscussionHandler) GetUserFavorites(c *gin.Context) {
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

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	discussions, total, err := h.discussionService.GetUserFavorites(c.Request.Context(), userID, page, limit)
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

	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    discussions,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetTags retrieves all available tags
// @Summary Get all tags
// @Description Get all available discussion tags
// @Tags Discussions
// @Produce json
// @Success 200 {object} map[string]interface{} "List of tags"
// @Router /api/v1/discussions/tags [get]
func (h *DiscussionHandler) GetTags(c *gin.Context) {
	tags, err := h.discussionService.GetAllTags(c.Request.Context())
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
		"data":    tags,
	})
}

// ResolveDiscussion marks a discussion as resolved
// @Summary Resolve discussion
// @Description Mark a discussion as resolved (author only)
// @Tags Discussions
// @Produce json
// @Security BearerAuth
// @Param discussion_id path string true "Discussion ID" format(uuid)
// @Success 200 {object} map[string]interface{} "Discussion resolved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid discussion ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Discussion not found"
// @Router /api/v1/discussions/{discussion_id}/resolve [post]
func (h *DiscussionHandler) ResolveDiscussion(c *gin.Context) {
	discussionIDStr := c.Param("discussion_id")
	discussionID, err := uuid.Parse(discussionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid discussion ID format",
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

	discussion, err := h.discussionService.ResolveDiscussion(c.Request.Context(), discussionID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"

		switch err {
		case services.ErrDiscussionNotFound:
			status = http.StatusNotFound
			code = "NOT_FOUND"
		case services.ErrUnauthorized:
			status = http.StatusForbidden
			code = "FORBIDDEN"
		}

		c.JSON(status, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Discussion resolved successfully",
		"data":    discussion,
	})
}
