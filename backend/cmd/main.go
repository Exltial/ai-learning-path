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

	"ai-learning-platform/configs"
	"ai-learning-platform/internal/handlers"
	"ai-learning-platform/internal/middleware"
	"ai-learning-platform/internal/repository"
	"ai-learning-platform/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "ai-learning-platform/docs" // Import for Swagger docs
)

// @title AI Learning Platform API
// @version 2.0
// @description AI Learning Platform Backend API - Phase 2 Core Features
// @description 
// @description ## Features:
// @description - User Registration & Login with JWT Authentication
// @description - Course Management (List, Get, Create, Update, Delete)
// @description - Assignment Submission & Grading
// @description - Progress Tracking
// @description - User Achievements & Notifications
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@ai-learning-platform.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your JWT token in the format: Bearer {token}

// @tag.name Authentication
// @tag.description User authentication and authorization

// @tag.name Users
// @tag.description User management operations

// @tag.name Courses
// @tag.description Course management operations

// @tag.name Lessons
// @tag.description Lesson management operations

// @tag.name Exercises
// @tag.description Exercise management operations

// @tag.name Submissions
// @tag.description Assignment submission and grading

// @tag.name Progress
// @tag.description Learning progress tracking

func main() {
	// Load configuration from environment
	cfg := configs.Load()
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database connection with config
	dbConfig, err := pgxpool.ParseConfig(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}
	
	// Set connection pool settings
	dbConfig.MaxConns = int32(cfg.Database.MaxConnections)
	dbConfig.MinConns = int32(cfg.Database.MinConnections)
	dbConfig.MaxConnLifetime = time.Duration(cfg.Database.MaxLifetime) * time.Minute

	dbPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
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
	log.Println("✓ Successfully connected to database")

	// Initialize Redis client with config
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Unable to connect to Redis: %v", err)
	}
	log.Println("✓ Successfully connected to Redis")

	// Initialize repositories
	userRepo := repository.NewUserRepository(dbPool)
	courseRepo := repository.NewCourseRepository(dbPool)
	lessonRepo := repository.NewLessonRepository(dbPool)
	exerciseRepo := repository.NewExerciseRepository(dbPool)
	submissionRepo := repository.NewSubmissionRepository(dbPool)
	progressRepo := repository.NewProgressRepository(dbPool)
	enrollmentRepo := repository.NewEnrollmentRepository(dbPool)

	// Initialize services with JWT config
	authService := services.NewAuthService(userRepo, rdb, cfg.JWT.Secret, cfg.JWT.Expiration)
	userService := services.NewUserService(userRepo, progressRepo, enrollmentRepo)
	courseService := services.NewCourseService(courseRepo, lessonRepo, enrollmentRepo)
	lessonService := services.NewLessonService(lessonRepo)
	exerciseService := services.NewExerciseService(exerciseRepo)
	submissionService := services.NewSubmissionService(submissionRepo, exerciseRepo)
	progressService := services.NewProgressService(progressRepo, enrollmentRepo, lessonRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService, progressService)
	courseHandler := handlers.NewCourseHandler(courseService, enrollmentRepo)
	lessonHandler := handlers.NewLessonHandler(lessonService)
	exerciseHandler := handlers.NewExerciseHandler(exerciseService, submissionService)
	submissionHandler := handlers.NewSubmissionHandler(submissionService)
	progressHandler := handlers.NewProgressHandler(progressService)

	// Set up Gin router
	router := gin.Default()

	// Apply CORS middleware
	router.Use(middleware.CORS())

	// Health check endpoint
	// @Summary Health check
	// @Tags System
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /health [get]
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"time":      time.Now().Format(time.RFC3339),
			"version":   "2.0.0",
			"database":  "connected",
			"redis":     "connected",
		})
	})

	// API documentation (Swagger)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			// @Summary Register a new user
			// @Tags Authentication
			// @Accept json
			// @Produce json
			// @Param request body handlers.RegisterRequest true "Registration data"
			// @Success 201 {object} map[string]interface{} "User created successfully"
			// @Failure 400 {object} map[string]interface{} "Invalid input"
			// @Failure 409 {object} map[string]interface{} "Username or email already exists"
			// @Router /api/v1/auth/register [post]
			auth.POST("/register", authHandler.Register)
			
			// @Summary Login user
			// @Tags Authentication
			// @Accept json
			// @Produce json
			// @Param request body handlers.LoginRequest true "Login credentials"
			// @Success 200 {object} map[string]interface{} "Login successful"
			// @Failure 400 {object} map[string]interface{} "Invalid input"
			// @Failure 401 {object} map[string]interface{} "Invalid credentials"
			// @Router /api/v1/auth/login [post]
			auth.POST("/login", authHandler.Login)
			
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes (authentication required)
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
				courses.POST("", middleware.RequireRole("instructor", "admin"), courseHandler.CreateCourse)
				courses.PUT("/:course_id", middleware.RequireRole("instructor", "admin"), courseHandler.UpdateCourse)
				courses.DELETE("/:course_id", middleware.RequireRole("admin"), courseHandler.DeleteCourse)
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
				lessons.PUT("/:lesson_id", middleware.RequireRole("instructor", "admin"), lessonHandler.UpdateLesson)
				lessons.DELETE("/:lesson_id", middleware.RequireRole("instructor", "admin"), lessonHandler.DeleteLesson)
				lessons.PUT("/:lesson_id/progress", progressHandler.UpdateLessonProgress)
			}

			// Exercise routes
			exercises := protected.Group("/exercises")
			{
				exercises.GET("/:exercise_id", exerciseHandler.GetExercise)
				exercises.PUT("/:exercise_id", middleware.RequireRole("instructor", "admin"), exerciseHandler.UpdateExercise)
				exercises.DELETE("/:exercise_id", middleware.RequireRole("instructor", "admin"), exerciseHandler.DeleteExercise)
				exercises.POST("/:exercise_id/submit", exerciseHandler.SubmitExercise)
				exercises.GET("/:exercise_id/submissions", exerciseHandler.GetExerciseSubmissions)
			}

			// Submission routes
			submissions := protected.Group("/submissions")
			{
				submissions.GET("/:submission_id", submissionHandler.GetSubmission)
				submissions.POST("/:submission_id/grade", middleware.RequireRole("instructor", "admin"), submissionHandler.GradeSubmission)
			}

			// Discussion routes (placeholder)
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

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Starting server on port %s", cfg.Server.Port)
		log.Printf("📚 Swagger docs: http://localhost:%s/swagger/index.html", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✓ Server exited gracefully")
}
