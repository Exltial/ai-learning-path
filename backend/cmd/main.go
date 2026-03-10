package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ai-learning-platform/internal/handlers"
	"ai-learning-platform/internal/middleware"
	"ai-learning-platform/internal/repository"
	"ai-learning-platform/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/ai_learning?sslmode=disable"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	// Initialize database connection
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Test database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	log.Println("Successfully connected to database")

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Unable to connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis")

	// Initialize repositories
	userRepo := repository.NewUserRepository(dbPool)
	courseRepo := repository.NewCourseRepository(dbPool)
	lessonRepo := repository.NewLessonRepository(dbPool)
	exerciseRepo := repository.NewExerciseRepository(dbPool)
	submissionRepo := repository.NewSubmissionRepository(dbPool)
	progressRepo := repository.NewProgressRepository(dbPool)
	enrollmentRepo := repository.NewEnrollmentRepository(dbPool)

	// Initialize services
	authService := services.NewAuthService(userRepo, rdb)
	userService := services.NewUserService(userRepo)
	courseService := services.NewCourseService(courseRepo, lessonRepo, enrollmentRepo)
	lessonService := services.NewLessonService(lessonRepo)
	exerciseService := services.NewExerciseService(exerciseRepo)
	submissionService := services.NewSubmissionService(submissionRepo, exerciseRepo)
	progressService := services.NewProgressService(progressRepo, enrollmentRepo, lessonRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService, progressService)
	courseHandler := handlers.NewCourseHandler(courseService, enrollmentService)
	lessonHandler := handlers.NewLessonHandler(lessonService)
	exerciseHandler := handlers.NewExerciseHandler(exerciseService, submissionService)
	submissionHandler := handlers.NewSubmissionHandler(submissionService)
	progressHandler := handlers.NewProgressHandler(progressService)

	// Set up Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.JWTMiddleware(authService))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetCurrentUser)
				users.PUT("/me", userHandler.UpdateUser)
				users.PUT("/me/password", userHandler.ChangePassword)
				users.GET("/me/stats", userHandler.GetUserStats)
				users.GET("/me/progress", progressHandler.GetUserProgress)
				users.GET("/me/achievements", userHandler.GetUserAchievements)
				users.GET("/me/notifications", userHandler.GetNotifications)
			}

			// Course routes
			courses := protected.Group("/courses")
			{
				courses.GET("", courseHandler.ListCourses)
				courses.GET("/:course_id", courseHandler.GetCourse)
				courses.POST("", courseHandler.CreateCourse) // Instructor/Admin
				courses.PUT("/:course_id", courseHandler.UpdateCourse)
				courses.DELETE("/:course_id", courseHandler.DeleteCourse)
				courses.POST("/:course_id/enroll", courseHandler.EnrollCourse)
				courses.GET("/:course_id/lessons", courseHandler.GetCourseLessons)
				courses.GET("/:course_id/progress", progressHandler.GetCourseProgress)
				courses.GET("/:course_id/reviews", courseHandler.GetCourseReviews)
				courses.POST("/:course_id/reviews", courseHandler.CreateReview)
				courses.PUT("/:course_id/reviews", courseHandler.UpdateReview)
			}

			// Lesson routes
			lessons := protected.Group("/lessons")
			{
				lessons.GET("/:lesson_id", lessonHandler.GetLesson)
				lessons.PUT("/:lesson_id", lessonHandler.UpdateLesson)
				lessons.DELETE("/:lesson_id", lessonHandler.DeleteLesson)
				lessons.PUT("/:lesson_id/progress", progressHandler.UpdateLessonProgress)
			}

			// Exercise routes
			exercises := protected.Group("/exercises")
			{
				exercises.GET("/:exercise_id", exerciseHandler.GetExercise)
				exercises.PUT("/:exercise_id", exerciseHandler.UpdateExercise)
				exercises.DELETE("/:exercise_id", exerciseHandler.DeleteExercise)
				exercises.POST("/:exercise_id/submit", exerciseHandler.SubmitExercise)
				exercises.GET("/:exercise_id/submissions", submissionHandler.GetSubmissions)
			}

			// Submission routes
			submissions := protected.Group("/submissions")
			{
				submissions.GET("/:submission_id", submissionHandler.GetSubmission)
				submissions.POST("/:submission_id/grade", submissionHandler.GradeSubmission) // Instructor/Admin
			}

			// Discussion routes
			discussions := protected.Group("/discussions")
			{
				// TODO: Implement discussion handlers
			}

			// Notification routes
			notifications := protected.Group("/notifications")
			{
				notifications.PATCH("/:notification_id/read", userHandler.MarkNotificationRead)
			}
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
