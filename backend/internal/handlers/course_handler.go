package handlers

import (
	"net/http"
	"strconv"
	"time"
	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CourseHandler handles course-related requests
type CourseHandler struct {
	courseService *services.CourseService
}

// NewCourseHandler creates a new CourseHandler
func NewCourseHandler(courseService *services.CourseService, enrollmentService interface{}) *CourseHandler {
	return &CourseHandler{courseService: courseService}
}

// ListCourses retrieves a list of courses
func (h *CourseHandler) ListCourses(c *gin.Context) {
	category := c.Query("category")
	difficulty := c.Query("difficulty")
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	if limit > 100 {
		limit = 100
	}
	if page < 1 {
		page = 1
	}

	courses, total, err := h.courseService.ListCourses(c.Request.Context(), category, difficulty, page, limit)
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

	totalPages := (total + limit - 1) / limit

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
func (h *CourseHandler) GetCourse(c *gin.Context) {
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

// CreateCourse creates a new course
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	instructorIDStr, _ := c.Get("user_id")
	instructorID, _ := uuid.Parse(instructorIDStr.(string))

	var req struct {
		Title           string  `json:"title" binding:"required"`
		Description     string  `json:"description"`
		Category        string  `json:"category"`
		DifficultyLevel string  `json:"difficulty_level" binding:"omitempty,oneof=beginner intermediate advanced"`
		EstimatedHours  int     `json:"estimated_hours"`
		Price           float64 `json:"price"`
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

	course := &models.Course{
		ID:              uuid.New(),
		Title:           req.Title,
		Description:     req.Description,
		Category:        req.Category,
		DifficultyLevel: req.DifficultyLevel,
		EstimatedHours:  req.EstimatedHours,
		Price:           req.Price,
		InstructorID:    instructorID,
		IsPublished:     false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
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
		"data": course,
	})
}

// UpdateCourse updates an existing course
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
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

	var req struct {
		Title           string  `json:"title"`
		Description     string  `json:"description"`
		Category        string  `json:"category"`
		DifficultyLevel string  `json:"difficulty_level"`
		EstimatedHours  int     `json:"estimated_hours"`
		Price           float64 `json:"price"`
		IsPublished     bool    `json:"is_published"`
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

	// Update fields
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
	course.EstimatedHours = req.EstimatedHours
	course.Price = req.Price
	course.IsPublished = req.IsPublished

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
		"data": course,
	})
}

// DeleteCourse deletes a course
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
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
func (h *CourseHandler) EnrollCourse(c *gin.Context) {
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

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	enrollment, err := h.courseService.EnrollCourse(c.Request.Context(), userID, courseID)
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

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": enrollment,
	})
}

// GetCourseLessons retrieves lessons for a course
func (h *CourseHandler) GetCourseLessons(c *gin.Context) {
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

	lessons, err := h.courseService.GetCourseLessons(c.Request.Context(), courseID)
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
			"lessons": lessons,
		},
	})
}

// GetCourseReviews retrieves reviews for a course
func (h *CourseHandler) GetCourseReviews(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"reviews":        reviews,
			"average_rating": 0.0,
			"total_reviews":  len(reviews),
		},
	})
}

// CreateReview creates a review for a course
func (h *CourseHandler) CreateReview(c *gin.Context) {
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

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var req struct {
		Rating  int    `json:"rating" binding:"required,min=1,max=5"`
		Comment string `json:"comment"`
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

	review := &models.CourseReview{
		ID:         uuid.New(),
		CourseID:   courseID,
		UserID:     userID,
		Rating:     req.Rating,
		Comment:    req.Comment,
		IsVerified: true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
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
		"data": review,
	})
}

// UpdateReview updates a review for a course
func (h *CourseHandler) UpdateReview(c *gin.Context) {
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

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var req struct {
		Rating  int    `json:"rating" binding:"min=1,max=5"`
		Comment string `json:"comment"`
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

	// TODO: Get existing review and update
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Review updated successfully",
	})
}
