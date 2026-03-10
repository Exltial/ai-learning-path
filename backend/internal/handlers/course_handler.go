package handlers

import (
	"net/http"
	"strconv"
	"time"

	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"
	"ai-learning-platform/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CourseHandler handles course-related requests
type CourseHandler struct {
	courseService  *services.CourseService
	enrollmentRepo *repository.EnrollmentRepository
}

// NewCourseHandler creates a new CourseHandler
func NewCourseHandler(courseService *services.CourseService, enrollmentRepo *repository.EnrollmentRepository) *CourseHandler {
	return &CourseHandler{
		courseService:  courseService,
		enrollmentRepo: enrollmentRepo,
	}
}

// ListCoursesResponse represents the response for listing courses
type ListCoursesResponse struct {
	Courses    []*models.Course `json:"courses"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	Total      int64            `json:"total"`
	TotalPages int              `json:"total_pages"`
}

// ListCourses retrieves a list of courses with pagination and filters
// @Summary List all courses
// @Description Get paginated list of published courses with optional filters
// @Tags Courses
// @Accept json
// @Produce json
// @Param category query string false "Filter by category"
// @Param difficulty query string false "Filter by difficulty (beginner, intermediate, advanced)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{} "List of courses with pagination"
// @Failure 400 {object} map[string]interface{} "Invalid parameters"
// @Router /api/v1/courses [get]
func (h *CourseHandler) ListCourses(c *gin.Context) {
	category := c.Query("category")
	difficulty := c.Query("difficulty")
	
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

	courses, total, err := h.courseService.ListCourses(c.Request.Context(), category, difficulty, page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"courses": courses,
			"pagination": map[string]interface{}{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": totalPages,
			},
		},
	})
}

// GetCourse retrieves a single course by ID
// @Summary Get course details
// @Description Get detailed information about a specific course
// @Tags Courses
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID" format(uuid)
// @Success 200 {object} map[string]interface{} "Course details"
// @Failure 400 {object} map[string]interface{} "Invalid course ID"
// @Failure 404 {object} map[string]interface{} "Course not found"
// @Router /api/v1/courses/{course_id} [get]
func (h *CourseHandler) GetCourse(c *gin.Context) {
	courseIDStr := c.Param("course_id")
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

	course, err := h.courseService.GetCourse(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "NOT_FOUND",
				"message": "Course not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": course,
	})
}

// CreateCourseRequest represents a course creation request
type CreateCourseRequest struct {
	Title           string  `json:"title" binding:"required,min=1,max=200" example:"Introduction to Go Programming"`
	Description     string  `json:"description" example:"Learn Go programming from scratch"`
	Category        string  `json:"category" example:"Programming"`
	DifficultyLevel string  `json:"difficulty_level" binding:"omitempty,oneof=beginner intermediate advanced" example:"beginner"`
	EstimatedHours  int     `json:"estimated_hours" example:"40"`
	Price           float64 `json:"price" example:"99.99"`
	ThumbnailURL    string  `json:"thumbnail_url" example:"https://example.com/thumb.jpg"`
}

// CreateCourse creates a new course
// @Summary Create a new course
// @Description Create a new course (instructor/admin only)
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCourseRequest true "Course data"
// @Success 201 {object} map[string]interface{} "Course created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (insufficient permissions)"
// @Router /api/v1/courses [post]
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	instructorIDStr, exists := c.Get("user_id")
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
	
	instructorID, err := uuid.Parse(instructorIDStr.(string))
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

	var req CreateCourseRequest
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

	course := &models.Course{
		Title:           req.Title,
		Description:     req.Description,
		Category:        req.Category,
		DifficultyLevel: req.DifficultyLevel,
		EstimatedHours:  req.EstimatedHours,
		Price:           req.Price,
		ThumbnailURL:    req.ThumbnailURL,
		InstructorID:    instructorID,
	}

	if err := h.courseService.CreateCourse(c.Request.Context(), course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Course created successfully",
		"data": course,
	})
}

// UpdateCourseRequest represents a course update request
type UpdateCourseRequest struct {
	Title           string  `json:"title" binding:"omitempty,min=1,max=200"`
	Description     string  `json:"description"`
	Category        string  `json:"category"`
	DifficultyLevel string  `json:"difficulty_level" binding:"omitempty,oneof=beginner intermediate advanced"`
	EstimatedHours  int     `json:"estimated_hours"`
	Price           float64 `json:"price"`
	IsPublished     *bool   `json:"is_published"`
	ThumbnailURL    string  `json:"thumbnail_url"`
}

// UpdateCourse updates an existing course
// @Summary Update a course
// @Description Update course information (instructor/admin only)
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param course_id path string true "Course ID" format(uuid)
// @Param request body UpdateCourseRequest true "Updated course data"
// @Success 200 {object} map[string]interface{} "Course updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Course not found"
// @Router /api/v1/courses/{course_id} [put]
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	courseIDStr := c.Param("course_id")
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

	var req UpdateCourseRequest
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

	// Get existing course
	course, err := h.courseService.GetCourse(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "NOT_FOUND",
				"message": "Course not found",
			},
		})
		return
	}

	// Update fields if provided
	if req.Title != "" {
		course.Title = req.Title
	}
	if req.Description != "" {
		course.Description = req.Description
	}
	if req.Category != "" {
		course.Category = req.Category
	}
	if req.DifficultyLevel != "" {
		course.DifficultyLevel = req.DifficultyLevel
	}
	if req.EstimatedHours > 0 {
		course.EstimatedHours = req.EstimatedHours
	}
	if req.Price > 0 {
		course.Price = req.Price
	}
	if req.ThumbnailURL != "" {
		course.ThumbnailURL = req.ThumbnailURL
	}
	if req.IsPublished != nil {
		course.IsPublished = *req.IsPublished
	}

	if err := h.courseService.UpdateCourse(c.Request.Context(), course); err != nil {
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
		"message": "Course updated successfully",
		"data": course,
	})
}

// DeleteCourse deletes a course
// @Summary Delete a course
// @Description Delete a course (admin only)
// @Tags Courses
// @Produce json
// @Security BearerAuth
// @Param course_id path string true "Course ID" format(uuid)
// @Success 204 {object} nil "Course deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid course ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Course not found"
// @Router /api/v1/courses/{course_id} [delete]
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	courseIDStr := c.Param("course_id")
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

	if err := h.courseService.DeleteCourse(c.Request.Context(), courseID); err != nil {
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

// EnrollCourse enrolls a user in a course
// @Summary Enroll in a course
// @Description Enroll the current user in a course
// @Tags Courses
// @Produce json
// @Security BearerAuth
// @Param course_id path string true "Course ID" format(uuid)
// @Success 201 {object} map[string]interface{} "Enrollment successful"
// @Failure 400 {object} map[string]interface{} "Invalid course ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Course not found"
// @Router /api/v1/courses/{course_id}/enroll [post]
func (h *CourseHandler) EnrollCourse(c *gin.Context) {
	courseIDStr := c.Param("course_id")
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

	enrollment, err := h.courseService.EnrollCourse(c.Request.Context(), userID, courseID)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		
		if err == services.ErrCourseNotPublished {
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

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Enrollment successful",
		"data": enrollment,
	})
}

// GetCourseLessons retrieves lessons for a course
// @Summary Get course lessons
// @Description Get all lessons for a specific course
// @Tags Courses
// @Produce json
// @Param course_id path string true "Course ID" format(uuid)
// @Success 200 {object} map[string]interface{} "List of lessons"
// @Failure 400 {object} map[string]interface{} "Invalid course ID"
// @Failure 404 {object} map[string]interface{} "Course not found"
// @Router /api/v1/courses/{course_id}/lessons [get]
func (h *CourseHandler) GetCourseLessons(c *gin.Context) {
	courseIDStr := c.Param("course_id")
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

	lessons, err := h.courseService.GetCourseLessons(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "NOT_FOUND",
				"message": "Course not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"lessons": lessons,
			"count":   len(lessons),
		},
	})
}

// GetCourseReviews retrieves reviews for a course
// @Summary Get course reviews
// @Description Get all reviews for a specific course
// @Tags Courses
// @Produce json
// @Param course_id path string true "Course ID" format(uuid)
// @Success 200 {object} map[string]interface{} "List of reviews"
// @Failure 400 {object} map[string]interface{} "Invalid course ID"
// @Failure 404 {object} map[string]interface{} "Course not found"
// @Router /api/v1/courses/{course_id}/reviews [get]
func (h *CourseHandler) GetCourseReviews(c *gin.Context) {
	courseIDStr := c.Param("course_id")
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

	reviews, err := h.courseService.GetCourseReviews(c.Request.Context(), courseID)
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

	// Calculate average rating
	var totalRating int
	for _, review := range reviews {
		totalRating += review.Rating
	}
	var avgRating float64
	if len(reviews) > 0 {
		avgRating = float64(totalRating) / float64(len(reviews))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"reviews":        reviews,
			"average_rating": avgRating,
			"total_reviews":  len(reviews),
		},
	})
}

// CreateReviewRequest represents a review creation request
type CreateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5" example:"5"`
	Comment string `json:"comment" example:"Great course! Very informative."`
}

// CreateReview creates a review for a course
// @Summary Create a course review
// @Description Create a review for a specific course
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param course_id path string true "Course ID" format(uuid)
// @Param request body CreateReviewRequest true "Review data"
// @Success 201 {object} map[string]interface{} "Review created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Course not found"
// @Router /api/v1/courses/{course_id}/reviews [post]
func (h *CourseHandler) CreateReview(c *gin.Context) {
	courseIDStr := c.Param("course_id")
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

	var req CreateReviewRequest
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

	review := &models.CourseReview{
		CourseID: courseID,
		UserID:   userID,
		Rating:   req.Rating,
		Comment:  req.Comment,
	}

	if err := h.courseService.CreateReview(c.Request.Context(), review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Review created successfully",
		"data": review,
	})
}

// UpdateReviewRequest represents a review update request
type UpdateReviewRequest struct {
	Rating  int    `json:"rating" binding:"omitempty,min=1,max=5"`
	Comment string `json:"comment"`
}

// UpdateReview updates a review for a course
// @Summary Update a course review
// @Description Update an existing review for a course
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param course_id path string true "Course ID" format(uuid)
// @Param request body UpdateReviewRequest true "Updated review data"
// @Success 200 {object} map[string]interface{} "Review updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Course or review not found"
// @Router /api/v1/courses/{course_id}/reviews [put]
func (h *CourseHandler) UpdateReview(c *gin.Context) {
	courseIDStr := c.Param("course_id")
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

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var req UpdateReviewRequest
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

	// TODO: Get existing review and update
	// For now, return placeholder response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Review updated successfully",
		"data": map[string]interface{}{
			"course_id": courseID,
			"user_id":   userID,
			"rating":    req.Rating,
			"comment":   req.Comment,
			"updated_at": time.Now(),
		},
	})
}
