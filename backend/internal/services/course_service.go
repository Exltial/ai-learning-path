package services

import (
	"context"
	"errors"
	"time"

	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"

	"github.com/google/uuid"
)

// Error definitions for course service
var (
	ErrCourseNotFound     = errors.New("course not found")
	ErrCourseNotPublished = errors.New("course is not published")
	ErrAlreadyEnrolled    = errors.New("already enrolled in this course")
	ErrInvalidDifficulty  = errors.New("invalid difficulty level")
	ErrInvalidCategory    = errors.New("invalid category")
)

// Valid difficulty levels
var validDifficulties = map[string]bool{
	"beginner":     true,
	"intermediate": true,
	"advanced":     true,
}

// CourseService handles course-related business logic
type CourseService struct {
	courseRepo     *repository.CourseRepository
	lessonRepo     *repository.LessonRepository
	enrollmentRepo *repository.EnrollmentRepository
}

// NewCourseService creates a new CourseService
func NewCourseService(courseRepo *repository.CourseRepository, lessonRepo *repository.LessonRepository, enrollmentRepo *repository.EnrollmentRepository) *CourseService {
	return &CourseService{
		courseRepo:     courseRepo,
		lessonRepo:     lessonRepo,
		enrollmentRepo: enrollmentRepo,
	}
}

// GetCourse retrieves a course by ID
// If includeLessons is true, also fetches lessons for the course
func (s *CourseService) GetCourse(ctx context.Context, id uuid.UUID) (*models.Course, error) {
	course, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrCourseNotFound
	}
	return course, nil
}

// ListCourses retrieves courses with pagination and filters
// Supports filtering by category and difficulty level
func (s *CourseService) ListCourses(ctx context.Context, category, difficulty string, page, limit int) ([]*models.Course, int, error) {
	// Validate difficulty if provided
	if difficulty != "" && !validDifficulties[difficulty] {
		return nil, 0, ErrInvalidDifficulty
	}

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100 // Max limit to prevent abuse
	}

	courses, total, err := s.courseRepo.List(ctx, category, difficulty, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

// CreateCourse creates a new course
func (s *CourseService) CreateCourse(ctx context.Context, course *models.Course) error {
	// Validate course data
	if err := s.validateCourse(course); err != nil {
		return err
	}

	// Set initial values
	course.ID = uuid.New()
	course.CreatedAt = time.Now()
	course.UpdatedAt = time.Now()
	course.IsPublished = false // Default to unpublished
	course.EnrollmentCount = 0
	course.Rating = 0.0

	return s.courseRepo.Create(ctx, course)
}

// UpdateCourse updates an existing course
func (s *CourseService) UpdateCourse(ctx context.Context, course *models.Course) error {
	// Check if course exists
	existing, err := s.courseRepo.GetByID(ctx, course.ID)
	if err != nil {
		return ErrCourseNotFound
	}

	// Validate updated data
	if err := s.validateCourse(course); err != nil {
		return err
	}

	// Preserve certain fields
	course.CreatedAt = existing.CreatedAt
	course.EnrollmentCount = existing.EnrollmentCount
	course.Rating = existing.Rating
	course.UpdatedAt = time.Now()

	return s.courseRepo.Update(ctx, course)
}

// DeleteCourse deletes a course
// Note: Should check for existing enrollments before deletion
func (s *CourseService) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	// Check if course exists
	_, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return ErrCourseNotFound
	}

	// TODO: Check for existing enrollments and handle appropriately
	// For now, just delete
	return s.courseRepo.Delete(ctx, id)
}

// GetCourseLessons retrieves all lessons for a course, ordered by order_index
func (s *CourseService) GetCourseLessons(ctx context.Context, courseID uuid.UUID) ([]*models.Lesson, error) {
	// Verify course exists
	_, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, ErrCourseNotFound
	}

	lessons, err := s.lessonRepo.GetByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	return lessons, nil
}

// EnrollCourse enrolls a user in a course
func (s *CourseService) EnrollCourse(ctx context.Context, userID, courseID uuid.UUID) (*models.Enrollment, error) {
	// Verify course exists and is published
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, ErrCourseNotFound
	}
	if !course.IsPublished {
		return nil, ErrCourseNotPublished
	}

	// Check if already enrolled
	exists, err := s.enrollmentRepo.Exists(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}
	if exists {
		// Return existing enrollment
		return s.enrollmentRepo.GetByUserAndCourse(ctx, userID, courseID)
	}

	// Create new enrollment
	enrollment := &models.Enrollment{
		ID:                 uuid.New(),
		UserID:             userID,
		CourseID:           courseID,
		EnrolledAt:         time.Now(),
		Status:             "active",
		ProgressPercentage: 0.0,
	}

	if err := s.enrollmentRepo.Create(ctx, enrollment); err != nil {
		return nil, err
	}

	// Increment enrollment count (non-blocking)
	go func() {
		updateCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := s.courseRepo.IncrementEnrollmentCount(updateCtx, courseID); err != nil {
			// Log error but don't fail enrollment
		}
	}()

	return enrollment, nil
}

// GetCourseReviews retrieves reviews for a course
func (s *CourseService) GetCourseReviews(ctx context.Context, courseID uuid.UUID) ([]*models.CourseReview, error) {
	// Verify course exists
	_, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, ErrCourseNotFound
	}

	// TODO: Implement review retrieval from repository
	// For now, return empty slice
	return []*models.CourseReview{}, nil
}

// CreateReview creates a course review
func (s *CourseService) CreateReview(ctx context.Context, review *models.CourseReview) error {
	// Validate review
	if review.Rating < 1 || review.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	// Verify course exists
	_, err := s.courseRepo.GetByID(ctx, review.CourseID)
	if err != nil {
		return ErrCourseNotFound
	}

	// Check if user already reviewed
	// TODO: Implement duplicate review check

	review.ID = uuid.New()
	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()
	review.IsVerified = true // TODO: Verify enrollment

	// TODO: Save review to database
	// For now, return nil (placeholder)
	return nil
}

// UpdateReview updates a course review
func (s *CourseService) UpdateReview(ctx context.Context, review *models.CourseReview) error {
	// Validate review
	if review.Rating < 1 || review.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	review.UpdatedAt = time.Now()

	// TODO: Update review in database
	return nil
}

// GetCoursesByInstructor retrieves all courses by an instructor
func (s *CourseService) GetCoursesByInstructor(ctx context.Context, instructorID uuid.UUID) ([]*models.Course, error) {
	return s.courseRepo.GetByInstructorID(ctx, instructorID)
}

// PublishCourse publishes a course (makes it visible to students)
func (s *CourseService) PublishCourse(ctx context.Context, courseID uuid.UUID) error {
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return ErrCourseNotFound
	}

	course.IsPublished = true
	course.UpdatedAt = time.Now()

	return s.courseRepo.Update(ctx, course)
}

// UnpublishCourse unpublishes a course (hides it from students)
func (s *CourseService) UnpublishCourse(ctx context.Context, courseID uuid.UUID) error {
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return ErrCourseNotFound
	}

	course.IsPublished = false
	course.UpdatedAt = time.Now()

	return s.courseRepo.Update(ctx, course)
}

// validateCourse validates course data
func (s *CourseService) validateCourse(course *models.Course) error {
	if course.Title == "" {
		return errors.New("title is required")
	}
	if len(course.Title) > 200 {
		return errors.New("title must be less than 200 characters")
	}
	if course.DifficultyLevel != "" && !validDifficulties[course.DifficultyLevel] {
		return ErrInvalidDifficulty
	}
	if course.EstimatedHours < 0 {
		return errors.New("estimated hours cannot be negative")
	}
	if course.Price < 0 {
		return errors.New("price cannot be negative")
	}
	return nil
}

// GetCourseStats retrieves statistics for a course
func (s *CourseService) GetCourseStats(ctx context.Context, courseID uuid.UUID) (map[string]interface{}, error) {
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, ErrCourseNotFound
	}

	// Get enrollment count
	enrollmentCount, err := s.enrollmentRepo.CountByCourse(ctx, courseID)
	if err != nil {
		enrollmentCount = int64(course.EnrollmentCount)
	}

	// Get lesson count
	lessons, err := s.lessonRepo.GetByCourseID(ctx, courseID)
	if err != nil {
		lessons = []*models.Lesson{}
	}

	stats := map[string]interface{}{
		"enrollment_count": enrollmentCount,
		"lesson_count":     len(lessons),
		"rating":           course.Rating,
		"is_published":     course.IsPublished,
	}

	return stats, nil
}
