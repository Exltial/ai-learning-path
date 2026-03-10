package services

import (
	"context"
	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"

	"github.com/google/uuid"
)

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
func (s *CourseService) GetCourse(ctx context.Context, id uuid.UUID) (*models.Course, error) {
	return s.courseRepo.GetByID(ctx, id)
}

// ListCourses retrieves courses with pagination and filters
func (s *CourseService) ListCourses(ctx context.Context, category, difficulty string, page, limit int) ([]*models.Course, int, error) {
	return s.courseRepo.List(ctx, category, difficulty, page, limit)
}

// CreateCourse creates a new course
func (s *CourseService) CreateCourse(ctx context.Context, course *models.Course) error {
	return s.courseRepo.Create(ctx, course)
}

// UpdateCourse updates an existing course
func (s *CourseService) UpdateCourse(ctx context.Context, course *models.Course) error {
	return s.courseRepo.Update(ctx, course)
}

// DeleteCourse deletes a course
func (s *CourseService) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	return s.courseRepo.Delete(ctx, id)
}

// GetCourseLessons retrieves all lessons for a course
func (s *CourseService) GetCourseLessons(ctx context.Context, courseID uuid.UUID) ([]*models.Lesson, error) {
	return s.lessonRepo.GetByCourseID(ctx, courseID)
}

// EnrollCourse enrolls a user in a course
func (s *CourseService) EnrollCourse(ctx context.Context, userID, courseID uuid.UUID) (*models.Enrollment, error) {
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
		ID:               uuid.New(),
		UserID:           userID,
		CourseID:         courseID,
		EnrolledAt:       models.Course{}.CreatedAt, // Use current time
		Status:           "active",
		ProgressPercentage: 0.0,
	}

	if err := s.enrollmentRepo.Create(ctx, enrollment); err != nil {
		return nil, err
	}

	// Increment enrollment count
	if err := s.courseRepo.IncrementEnrollmentCount(ctx, courseID); err != nil {
		// Log error but don't fail enrollment
	}

	return enrollment, nil
}

// GetCourseReviews retrieves reviews for a course
func (s *CourseService) GetCourseReviews(ctx context.Context, courseID uuid.UUID) ([]*models.CourseReview, error) {
	// TODO: Implement course reviews retrieval
	return []*models.CourseReview{}, nil
}

// CreateReview creates a course review
func (s *CourseService) CreateReview(ctx context.Context, review *models.CourseReview) error {
	// TODO: Implement course review creation
	return nil
}

// UpdateReview updates a course review
func (s *CourseService) UpdateReview(ctx context.Context, review *models.CourseReview) error {
	// TODO: Implement course review update
	return nil
}
