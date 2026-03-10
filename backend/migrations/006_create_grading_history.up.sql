-- Migration: Create grading_history table
-- Created: 2026-03-10
-- Description: Stores grading history for submissions to track score changes and grading events

CREATE TABLE IF NOT EXISTS grading_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    grading_type VARCHAR(20) NOT NULL CHECK (grading_type IN ('auto', 'manual', 'semi_auto')),
    previous_score DECIMAL(5,2),
    new_score DECIMAL(5,2) NOT NULL,
    score_change DECIMAL(5,2) NOT NULL,
    reason TEXT,
    graded_by UUID REFERENCES users(id),
    graded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX idx_grading_history_submission_id ON grading_history(submission_id);
CREATE INDEX idx_grading_history_exercise_id ON grading_history(exercise_id);
CREATE INDEX idx_grading_history_user_id ON grading_history(user_id);
CREATE INDEX idx_grading_history_graded_at ON grading_history(graded_at DESC);
CREATE INDEX idx_grading_history_grading_type ON grading_history(grading_type);

-- Composite indexes for analytics
CREATE INDEX idx_grading_history_exercise_graded_at ON grading_history(exercise_id, graded_at DESC);
CREATE INDEX idx_grading_history_user_graded_at ON grading_history(user_id, graded_at DESC);

-- Comment on table
COMMENT ON TABLE grading_history IS 'Stores historical grading records for submissions, tracking score changes and grading events';
COMMENT ON COLUMN grading_history.grading_type IS 'Type of grading: auto (automated), manual (instructor), semi_auto (assisted)';
COMMENT ON COLUMN grading_history.metadata IS 'Additional metadata about the grading event (e.g., detailed feedback, test results)';
