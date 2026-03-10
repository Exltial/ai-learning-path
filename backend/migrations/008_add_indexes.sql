-- Migration: 008_add_indexes.sql
-- Purpose: Add database indexes for performance optimization
-- Target: API response time < 200ms

-- +goose Up
-- +goose StatementBegin

-- ============================================
-- User Table Indexes
-- ============================================

-- Index on email for fast user lookup during authentication
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Index on created_at for sorting users by registration date
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- Index on status for filtering active/inactive users
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

-- Composite index for common user queries
CREATE INDEX IF NOT EXISTS idx_users_status_created ON users(status, created_at);

-- ============================================
-- Course Table Indexes
-- ============================================

-- Index on category for filtering courses by category
CREATE INDEX IF NOT EXISTS idx_courses_category ON courses(category);

-- Index on difficulty_level for filtering by difficulty
CREATE INDEX IF NOT EXISTS idx_courses_difficulty ON courses(difficulty_level);

-- Index on is_published for filtering published courses
CREATE INDEX IF NOT EXISTS idx_courses_is_published ON courses(is_published);

-- Index on instructor_id for fetching instructor's courses
CREATE INDEX IF NOT EXISTS idx_courses_instructor_id ON courses(instructor_id);

-- Index on created_at for sorting courses by date
CREATE INDEX IF NOT EXISTS idx_courses_created_at ON courses(created_at);

-- Index on enrollment_count for sorting popular courses
CREATE INDEX IF NOT EXISTS idx_courses_enrollment_count ON courses(enrollment_count);

-- Composite index for common course listing queries
CREATE INDEX IF NOT EXISTS idx_courses_published_category_difficulty 
ON courses(is_published, category, difficulty_level);

-- Composite index for instructor's published courses
CREATE INDEX IF NOT EXISTS idx_courses_instructor_published 
ON courses(instructor_id, is_published);

-- Full-text search index on title and description (PostgreSQL)
-- For MySQL, use FULLTEXT index instead
CREATE INDEX IF NOT EXISTS idx_courses_title_search ON courses USING gin(to_tsvector('simple', title));

-- ============================================
-- Lesson Table Indexes
-- ============================================

-- Index on course_id for fetching lessons by course
CREATE INDEX IF NOT EXISTS idx_lessons_course_id ON lessons(course_id);

-- Index on order_index for ordering lessons within a course
CREATE INDEX IF NOT EXISTS idx_lessons_order_index ON lessons(order_index);

-- Composite index for fetching ordered lessons
CREATE INDEX IF NOT EXISTS idx_lessons_course_order ON lessons(course_id, order_index);

-- Index on lesson_type for filtering by type
CREATE INDEX IF NOT EXISTS idx_lessons_type ON lessons(lesson_type);

-- ============================================
-- Enrollment Table Indexes
-- ============================================

-- Index on user_id for fetching user's enrollments
CREATE INDEX IF NOT EXISTS idx_enrollments_user_id ON enrollments(user_id);

-- Index on course_id for fetching enrolled users
CREATE INDEX IF NOT EXISTS idx_enrollments_course_id ON enrollments(course_id);

-- Index on status for filtering active enrollments
CREATE INDEX IF NOT EXISTS idx_enrollments_status ON enrollments(status);

-- Index on enrolled_at for sorting by enrollment date
CREATE INDEX IF NOT EXISTS idx_enrollments_enrolled_at ON enrollments(enrolled_at);

-- Unique composite index to prevent duplicate enrollments
CREATE UNIQUE INDEX IF NOT EXISTS idx_enrollments_user_course_unique 
ON enrollments(user_id, course_id);

-- Composite index for user's active enrollments
CREATE INDEX IF NOT EXISTS idx_enrollments_user_status 
ON enrollments(user_id, status);

-- ============================================
-- Progress Table Indexes
-- ============================================

-- Index on user_id for fetching user's progress
CREATE INDEX IF NOT EXISTS idx_progress_user_id ON progress(user_id);

-- Index on course_id for fetching course progress
CREATE INDEX IF NOT EXISTS idx_progress_course_id ON progress(course_id);

-- Index on lesson_id for fetching lesson progress
CREATE INDEX IF NOT EXISTS idx_progress_lesson_id ON progress(lesson_id);

-- Composite index for user's course progress
CREATE INDEX IF NOT EXISTS idx_progress_user_course ON progress(user_id, course_id);

-- Index on completed_at for filtering completed lessons
CREATE INDEX IF NOT EXISTS idx_progress_completed_at ON progress(completed_at);

-- Index on status for filtering progress status
CREATE INDEX IF NOT EXISTS idx_progress_status ON progress(status);

-- ============================================
-- Exercise Table Indexes
-- ============================================

-- Index on lesson_id for fetching exercises by lesson
CREATE INDEX IF NOT EXISTS idx_exercises_lesson_id ON exercises(lesson_id);

-- Index on exercise_type for filtering by type
CREATE INDEX IF NOT EXISTS idx_exercises_type ON exercises(exercise_type);

-- Index on difficulty for filtering by difficulty
CREATE INDEX IF NOT EXISTS idx_exercises_difficulty ON exercises(difficulty);

-- Composite index for lesson's exercises in order
CREATE INDEX IF NOT EXISTS idx_exercises_lesson_order ON exercises(lesson_id, order_index);

-- ============================================
-- Submission Table Indexes
-- ============================================

-- Index on user_id for fetching user's submissions
CREATE INDEX IF NOT EXISTS idx_submissions_user_id ON submissions(user_id);

-- Index on exercise_id for fetching exercise submissions
CREATE INDEX IF NOT EXISTS idx_submissions_exercise_id ON submissions(exercise_id);

-- Index on status for filtering submission status
CREATE INDEX IF NOT EXISTS idx_submissions_status ON submissions(status);

-- Index on submitted_at for sorting by submission date
CREATE INDEX IF NOT EXISTS idx_submissions_submitted_at ON submissions(submitted_at);

-- Composite index for user's exercise submissions
CREATE INDEX IF NOT EXISTS idx_submissions_user_exercise ON submissions(user_id, exercise_id);

-- Composite index for recent submissions
CREATE INDEX IF NOT EXISTS idx_submissions_status_submitted ON submissions(status, submitted_at);

-- ============================================
-- Achievement Table Indexes
-- ============================================

-- Index on achievement_type for filtering by type
CREATE INDEX IF NOT EXISTS idx_achievements_type ON achievements(achievement_type);

-- Index on category for filtering by category
CREATE INDEX IF NOT EXISTS idx_achievements_category ON achievements(category);

-- Index on points for sorting by point value
CREATE INDEX IF NOT EXISTS idx_achievements_points ON achievements(points);

-- ============================================
-- User Achievement Table Indexes
-- ============================================

-- Index on user_id for fetching user's achievements
CREATE INDEX IF NOT EXISTS idx_user_achievements_user_id ON user_achievements(user_id);

-- Index on achievement_id for fetching users with achievement
CREATE INDEX IF NOT EXISTS idx_user_achievements_achievement_id ON user_achievements(achievement_id);

-- Index on earned_at for sorting by earn date
CREATE INDEX IF NOT EXISTS idx_user_achievements_earned_at ON user_achievements(earned_at);

-- Unique composite index to prevent duplicate achievements
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_achievements_user_achievement_unique 
ON user_achievements(user_id, achievement_id);

-- ============================================
-- Points Transaction Table Indexes
-- ============================================

-- Index on user_id for fetching user's transactions
CREATE INDEX IF NOT EXISTS idx_points_transactions_user_id ON points_transactions(user_id);

-- Index on transaction_type for filtering by type
CREATE INDEX IF NOT EXISTS idx_points_transactions_type ON points_transactions(transaction_type);

-- Index on created_at for sorting by date
CREATE INDEX IF NOT EXISTS idx_points_transactions_created_at ON points_transactions(created_at);

-- Composite index for user's recent transactions
CREATE INDEX IF NOT EXISTS idx_points_transactions_user_created ON points_transactions(user_id, created_at);

-- ============================================
-- User Level Table Indexes
-- ============================================

-- Index on user_id for fetching user's level
CREATE INDEX IF NOT EXISTS idx_user_levels_user_id ON user_levels(user_id);

-- Index on level for filtering by level
CREATE INDEX IF NOT EXISTS idx_user_levels_level ON user_levels(level);

-- Index on total_points for sorting by points
CREATE INDEX IF NOT EXISTS idx_user_levels_total_points ON user_levels(total_points);

-- ============================================
-- Streak Table Indexes
-- ============================================

-- Index on user_id for fetching user's streak
CREATE INDEX IF NOT EXISTS idx_streaks_user_id ON streaks(user_id);

-- Index on current_streak for filtering by streak length
CREATE INDEX IF NOT EXISTS idx_streaks_current_streak ON streaks(current_streak);

-- Index on last_activity for filtering inactive users
CREATE INDEX IF NOT EXISTS idx_streaks_last_activity ON streaks(last_activity);

-- ============================================
-- Grading History Table Indexes
-- ============================================

-- Index on submission_id for fetching grading by submission
CREATE INDEX IF NOT EXISTS idx_grading_history_submission_id ON grading_history(submission_id);

-- Index on graded_by for fetching grader's history
CREATE INDEX IF NOT EXISTS idx_grading_history_graded_by ON grading_history(graded_by);

-- Index on graded_at for sorting by grading date
CREATE INDEX IF NOT EXISTS idx_grading_history_graded_at ON grading_history(graded_at);

-- Index on status for filtering grading status
CREATE INDEX IF NOT EXISTS idx_grading_history_status ON grading_history(status);

-- ============================================
-- Discussion Tables Indexes (if exists)
-- ============================================

-- Index on course_id for course discussions
CREATE INDEX IF NOT EXISTS idx_discussions_course_id ON discussions(course_id);

-- Index on user_id for user's discussions
CREATE INDEX IF NOT EXISTS idx_discussions_user_id ON discussions(user_id);

-- Index on parent_id for threaded discussions
CREATE INDEX IF NOT EXISTS idx_discussions_parent_id ON discussions(parent_id);

-- Index on created_at for sorting discussions
CREATE INDEX IF NOT EXISTS idx_discussions_created_at ON discussions(created_at);

-- Index on lesson_id for lesson-specific discussions
CREATE INDEX IF NOT EXISTS idx_discussions_lesson_id ON discussions(lesson_id);

-- ============================================
-- Discussion Comments Indexes
-- ============================================

-- Index on discussion_id for comments in a discussion
CREATE INDEX IF NOT EXISTS idx_discussion_comments_discussion_id ON discussion_comments(discussion_id);

-- Index on user_id for user's comments
CREATE INDEX IF NOT EXISTS idx_discussion_comments_user_id ON discussion_comments(user_id);

-- Index on parent_id for threaded comments
CREATE INDEX IF NOT EXISTS idx_discussion_comments_parent_id ON discussion_comments(parent_id);

-- Index on created_at for sorting comments
CREATE INDEX IF NOT EXISTS idx_discussion_comments_created_at ON discussion_comments(created_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop all indexes created in this migration
-- Note: In production, carefully evaluate before dropping indexes

-- User indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_status_created;

-- Course indexes
DROP INDEX IF EXISTS idx_courses_category;
DROP INDEX IF EXISTS idx_courses_difficulty;
DROP INDEX IF EXISTS idx_courses_is_published;
DROP INDEX IF EXISTS idx_courses_instructor_id;
DROP INDEX IF EXISTS idx_courses_created_at;
DROP INDEX IF EXISTS idx_courses_enrollment_count;
DROP INDEX IF EXISTS idx_courses_published_category_difficulty;
DROP INDEX IF EXISTS idx_courses_instructor_published;
DROP INDEX IF EXISTS idx_courses_title_search;

-- Lesson indexes
DROP INDEX IF EXISTS idx_lessons_course_id;
DROP INDEX IF EXISTS idx_lessons_order_index;
DROP INDEX IF EXISTS idx_lessons_course_order;
DROP INDEX IF EXISTS idx_lessons_type;

-- Enrollment indexes
DROP INDEX IF EXISTS idx_enrollments_user_id;
DROP INDEX IF EXISTS idx_enrollments_course_id;
DROP INDEX IF EXISTS idx_enrollments_status;
DROP INDEX IF EXISTS idx_enrollments_enrolled_at;
DROP INDEX IF EXISTS idx_enrollments_user_course_unique;
DROP INDEX IF EXISTS idx_enrollments_user_status;

-- Progress indexes
DROP INDEX IF EXISTS idx_progress_user_id;
DROP INDEX IF EXISTS idx_progress_course_id;
DROP INDEX IF EXISTS idx_progress_lesson_id;
DROP INDEX IF EXISTS idx_progress_user_course;
DROP INDEX IF EXISTS idx_progress_completed_at;
DROP INDEX IF EXISTS idx_progress_status;

-- Exercise indexes
DROP INDEX IF EXISTS idx_exercises_lesson_id;
DROP INDEX IF EXISTS idx_exercises_type;
DROP INDEX IF EXISTS idx_exercises_difficulty;
DROP INDEX IF EXISTS idx_exercises_lesson_order;

-- Submission indexes
DROP INDEX IF EXISTS idx_submissions_user_id;
DROP INDEX IF EXISTS idx_submissions_exercise_id;
DROP INDEX IF EXISTS idx_submissions_status;
DROP INDEX IF EXISTS idx_submissions_submitted_at;
DROP INDEX IF EXISTS idx_submissions_user_exercise;
DROP INDEX IF EXISTS idx_submissions_status_submitted;

-- Achievement indexes
DROP INDEX IF EXISTS idx_achievements_type;
DROP INDEX IF EXISTS idx_achievements_category;
DROP INDEX IF EXISTS idx_achievements_points;

-- User Achievement indexes
DROP INDEX IF EXISTS idx_user_achievements_user_id;
DROP INDEX IF EXISTS idx_user_achievements_achievement_id;
DROP INDEX IF EXISTS idx_user_achievements_earned_at;
DROP INDEX IF EXISTS idx_user_achievements_user_achievement_unique;

-- Points Transaction indexes
DROP INDEX IF EXISTS idx_points_transactions_user_id;
DROP INDEX IF EXISTS idx_points_transactions_type;
DROP INDEX IF EXISTS idx_points_transactions_created_at;
DROP INDEX IF EXISTS idx_points_transactions_user_created;

-- User Level indexes
DROP INDEX IF EXISTS idx_user_levels_user_id;
DROP INDEX IF EXISTS idx_user_levels_level;
DROP INDEX IF EXISTS idx_user_levels_total_points;

-- Streak indexes
DROP INDEX IF EXISTS idx_streaks_user_id;
DROP INDEX IF EXISTS idx_streaks_current_streak;
DROP INDEX IF EXISTS idx_streaks_last_activity;

-- Grading History indexes
DROP INDEX IF EXISTS idx_grading_history_submission_id;
DROP INDEX IF EXISTS idx_grading_history_graded_by;
DROP INDEX IF EXISTS idx_grading_history_graded_at;
DROP INDEX IF EXISTS idx_grading_history_status;

-- Discussion indexes
DROP INDEX IF EXISTS idx_discussions_course_id;
DROP INDEX IF EXISTS idx_discussions_user_id;
DROP INDEX IF EXISTS idx_discussions_parent_id;
DROP INDEX IF EXISTS idx_discussions_created_at;
DROP INDEX IF EXISTS idx_discussions_lesson_id;

-- Discussion Comments indexes
DROP INDEX IF EXISTS idx_discussion_comments_discussion_id;
DROP INDEX IF EXISTS idx_discussion_comments_user_id;
DROP INDEX IF EXISTS idx_discussion_comments_parent_id;
DROP INDEX IF EXISTS idx_discussion_comments_created_at;

-- +goose StatementEnd
