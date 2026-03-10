-- Migration: Drop grading_history table
-- Created: 2026-03-10
-- Description: Rollback migration for grading_history table

DROP INDEX IF EXISTS idx_grading_history_user_graded_at;
DROP INDEX IF EXISTS idx_grading_history_exercise_graded_at;
DROP INDEX IF EXISTS idx_grading_history_grading_type;
DROP INDEX IF EXISTS idx_grading_history_graded_at;
DROP INDEX IF EXISTS idx_grading_history_user_id;
DROP INDEX IF EXISTS idx_grading_history_exercise_id;
DROP INDEX IF EXISTS idx_grading_history_submission_id;

DROP TABLE IF EXISTS grading_history;
