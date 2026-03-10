package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GradingHistory represents a record of grading activity
type GradingHistory struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	SubmissionID  uuid.UUID  `json:"submission_id" db:"submission_id"`
	ExerciseID    uuid.UUID  `json:"exercise_id" db:"exercise_id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	GradingType   string     `json:"grading_type" db:"grading_type"` // auto, manual, semi_auto
	PreviousScore *float64   `json:"previous_score,omitempty" db:"previous_score"`
	NewScore      float64    `json:"new_score" db:"new_score"`
	ScoreChange   float64    `json:"score_change" db:"score_change"`
	Reason        string     `json:"reason,omitempty" db:"reason"`
	GradedBy      *uuid.UUID `json:"graded_by,omitempty" db:"graded_by"` // nil for auto-grading
	GradedAt      time.Time  `json:"graded_at" db:"graded_at"`
	Metadata      string     `json:"metadata,omitempty" db:"metadata"` // JSON metadata
}

// GradingRepository handles database operations for grading history
type GradingRepository struct {
	db *pgxpool.Pool
}

// NewGradingRepository creates a new GradingRepository
func NewGradingRepository(db *pgxpool.Pool) *GradingRepository {
	return &GradingRepository{db: db}
}

// Create creates a new grading history record
func (r *GradingRepository) Create(ctx context.Context, history *GradingHistory) error {
	query := `
		INSERT INTO grading_history (id, submission_id, exercise_id, user_id, grading_type,
			previous_score, new_score, score_change, reason, graded_by, graded_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	var metadataJSON interface{}
	if history.Metadata != "" {
		if err := json.Unmarshal([]byte(history.Metadata), &metadataJSON); err != nil {
			metadataJSON = history.Metadata
		}
	}

	_, err := r.db.Exec(ctx, query,
		history.ID,
		history.SubmissionID,
		history.ExerciseID,
		history.UserID,
		history.GradingType,
		history.PreviousScore,
		history.NewScore,
		history.ScoreChange,
		history.Reason,
		history.GradedBy,
		history.GradedAt,
		metadataJSON,
	)
	return err
}

// GetByID retrieves a grading history record by ID
func (r *GradingRepository) GetByID(ctx context.Context, id uuid.UUID) (*GradingHistory, error) {
	query := `
		SELECT id, submission_id, exercise_id, user_id, grading_type,
			previous_score, new_score, score_change, reason, graded_by, graded_at, metadata
		FROM grading_history WHERE id = $1
	`

	history := &GradingHistory{}
	var metadataJSON interface{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&history.ID,
		&history.SubmissionID,
		&history.ExerciseID,
		&history.UserID,
		&history.GradingType,
		&history.PreviousScore,
		&history.NewScore,
		&history.ScoreChange,
		&history.Reason,
		&history.GradedBy,
		&history.GradedAt,
		&metadataJSON,
	)

	if err != nil {
		return nil, err
	}

	// Convert metadata back to JSON string
	if metadataJSON != nil {
		if metadataBytes, err := json.Marshal(metadataJSON); err == nil {
			history.Metadata = string(metadataBytes)
		} else {
			history.Metadata = fmt.Sprintf("%v", metadataJSON)
		}
	}

	return history, nil
}

// GetBySubmissionID retrieves grading history for a specific submission
func (r *GradingRepository) GetBySubmissionID(ctx context.Context, submissionID uuid.UUID) ([]*GradingHistory, error) {
	query := `
		SELECT id, submission_id, exercise_id, user_id, grading_type,
			previous_score, new_score, score_change, reason, graded_by, graded_at, metadata
		FROM grading_history 
		WHERE submission_id = $1 
		ORDER BY graded_at DESC
	`

	rows, err := r.db.Query(ctx, query, submissionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	histories := make([]*GradingHistory, 0)
	for rows.Next() {
		history := &GradingHistory{}
		var metadataJSON interface{}

		err := rows.Scan(
			&history.ID,
			&history.SubmissionID,
			&history.ExerciseID,
			&history.UserID,
			&history.GradingType,
			&history.PreviousScore,
			&history.NewScore,
			&history.ScoreChange,
			&history.Reason,
			&history.GradedBy,
			&history.GradedAt,
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}

		// Convert metadata back to JSON string
		if metadataJSON != nil {
			if metadataBytes, err := json.Marshal(metadataJSON); err == nil {
				history.Metadata = string(metadataBytes)
			} else {
				history.Metadata = fmt.Sprintf("%v", metadataJSON)
			}
		}

		histories = append(histories, history)
	}

	return histories, rows.Err()
}

// GetByUserID retrieves grading history for a specific user
func (r *GradingRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*GradingHistory, error) {
	query := `
		SELECT id, submission_id, exercise_id, user_id, grading_type,
			previous_score, new_score, score_change, reason, graded_by, graded_at, metadata
		FROM grading_history 
		WHERE user_id = $1 
		ORDER BY graded_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	histories := make([]*GradingHistory, 0)
	for rows.Next() {
		history := &GradingHistory{}
		var metadataJSON interface{}

		err := rows.Scan(
			&history.ID,
			&history.SubmissionID,
			&history.ExerciseID,
			&history.UserID,
			&history.GradingType,
			&history.PreviousScore,
			&history.NewScore,
			&history.ScoreChange,
			&history.Reason,
			&history.GradedBy,
			&history.GradedAt,
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}

		// Convert metadata back to JSON string
		if metadataJSON != nil {
			if metadataBytes, err := json.Marshal(metadataJSON); err == nil {
				history.Metadata = string(metadataBytes)
			} else {
				history.Metadata = fmt.Sprintf("%v", metadataJSON)
			}
		}

		histories = append(histories, history)
	}

	return histories, rows.Err()
}

// GetByExerciseID retrieves grading history for a specific exercise
func (r *GradingRepository) GetByExerciseID(ctx context.Context, exerciseID uuid.UUID) ([]*GradingHistory, error) {
	query := `
		SELECT id, submission_id, exercise_id, user_id, grading_type,
			previous_score, new_score, score_change, reason, graded_by, graded_at, metadata
		FROM grading_history 
		WHERE exercise_id = $1 
		ORDER BY graded_at DESC
	`

	rows, err := r.db.Query(ctx, query, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	histories := make([]*GradingHistory, 0)
	for rows.Next() {
		history := &GradingHistory{}
		var metadataJSON interface{}

		err := rows.Scan(
			&history.ID,
			&history.SubmissionID,
			&history.ExerciseID,
			&history.UserID,
			&history.GradingType,
			&history.PreviousScore,
			&history.NewScore,
			&history.ScoreChange,
			&history.Reason,
			&history.GradedBy,
			&history.GradedAt,
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}

		// Convert metadata back to JSON string
		if metadataJSON != nil {
			if metadataBytes, err := json.Marshal(metadataJSON); err == nil {
				history.Metadata = string(metadataBytes)
			} else {
				history.Metadata = fmt.Sprintf("%v", metadataJSON)
			}
		}

		histories = append(histories, history)
	}

	return histories, rows.Err()
}

// GetStatsByExercise retrieves grading statistics for an exercise
func (r *GradingRepository) GetStatsByExercise(ctx context.Context, exerciseID uuid.UUID) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_submissions,
			COUNT(CASE WHEN grading_type = 'auto' THEN 1 END) as auto_graded,
			COUNT(CASE WHEN grading_type = 'manual' THEN 1 END) as manually_graded,
			AVG(new_score) as average_score,
			MIN(new_score) as min_score,
			MAX(new_score) as max_score,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY new_score) as median_score
		FROM grading_history 
		WHERE exercise_id = $1
	`

	var totalSubmissions, autoGraded, manuallyGraded int64
	var avgScore, minScore, maxScore, medianScore float64

	err := r.db.QueryRow(ctx, query, exerciseID).Scan(
		&totalSubmissions,
		&autoGraded,
		&manuallyGraded,
		&avgScore,
		&minScore,
		&maxScore,
		&medianScore,
	)

	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_submissions": totalSubmissions,
		"auto_graded":       autoGraded,
		"manually_graded":   manuallyGraded,
		"average_score":     avgScore,
		"min_score":         minScore,
		"max_score":         maxScore,
		"median_score":      medianScore,
		"pass_rate":         0.0,
	}

	// Calculate pass rate (assuming 60% is passing)
	passQuery := `
		SELECT COUNT(CASE WHEN new_score >= 60 THEN 1 END)::float / COUNT(*) * 100
		FROM grading_history 
		WHERE exercise_id = $1 AND new_score IS NOT NULL
	`
	var passRate float64
	if err := r.db.QueryRow(ctx, passQuery, exerciseID).Scan(&passRate); err == nil {
		stats["pass_rate"] = passRate
	}

	return stats, nil
}

// GetRecentGradingActivity retrieves recent grading activity
func (r *GradingRepository) GetRecentGradingActivity(ctx context.Context, limit int) ([]*GradingHistory, error) {
	query := `
		SELECT id, submission_id, exercise_id, user_id, grading_type,
			previous_score, new_score, score_change, reason, graded_by, graded_at, metadata
		FROM grading_history 
		ORDER BY graded_at DESC
		LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	histories := make([]*GradingHistory, 0)
	for rows.Next() {
		history := &GradingHistory{}
		var metadataJSON interface{}

		err := rows.Scan(
			&history.ID,
			&history.SubmissionID,
			&history.ExerciseID,
			&history.UserID,
			&history.GradingType,
			&history.PreviousScore,
			&history.NewScore,
			&history.ScoreChange,
			&history.Reason,
			&history.GradedBy,
			&history.GradedAt,
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}

		// Convert metadata back to JSON string
		if metadataJSON != nil {
			if metadataBytes, err := json.Marshal(metadataJSON); err == nil {
				history.Metadata = string(metadataBytes)
			} else {
				history.Metadata = fmt.Sprintf("%v", metadataJSON)
			}
		}

		histories = append(histories, history)
	}

	return histories, rows.Err()
}

// DeleteBySubmissionID deletes grading history for a submission (cascade on submission delete)
func (r *GradingRepository) DeleteBySubmissionID(ctx context.Context, submissionID uuid.UUID) error {
	query := `DELETE FROM grading_history WHERE submission_id = $1`
	_, err := r.db.Exec(ctx, query, submissionID)
	return err
}

// GetGradingTrends retrieves grading trends over time for an exercise
func (r *GradingRepository) GetGradingTrends(ctx context.Context, exerciseID uuid.UUID, days int) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			DATE(graded_at) as date,
			COUNT(*) as submissions,
			AVG(new_score) as avg_score,
			COUNT(CASE WHEN new_score >= 60 THEN 1 END) as passed
		FROM grading_history 
		WHERE exercise_id = $1 
			AND graded_at >= NOW() - INTERVAL '1 day' * $2
		GROUP BY DATE(graded_at)
		ORDER BY date ASC
	`

	rows, err := r.db.Query(ctx, query, exerciseID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trends := make([]map[string]interface{}, 0)
	for rows.Next() {
		var date time.Time
		var submissions int64
		var avgScore float64
		var passed int64

		err := rows.Scan(&date, &submissions, &avgScore, &passed)
		if err != nil {
			return nil, err
		}

		trends = append(trends, map[string]interface{}{
			"date":        date.Format("2006-01-02"),
			"submissions": submissions,
			"avg_score":   avgScore,
			"passed":      passed,
			"pass_rate":   float64(passed) / float64(submissions) * 100,
		})
	}

	return trends, rows.Err()
}
