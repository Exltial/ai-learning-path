package repository

import (
	"context"
	"database/sql"
	"ai-learning-platform/internal/models"
	"github.com/google/uuid"
)

// PointsTransactionRepository handles points transaction data access
type PointsTransactionRepository struct {
	db *sql.DB
}

// NewPointsTransactionRepository creates a new PointsTransactionRepository
func NewPointsTransactionRepository(db *sql.DB) *PointsTransactionRepository {
	return &PointsTransactionRepository{db: db}
}

// Create creates a new points transaction
func (r *PointsTransactionRepository) Create(ctx context.Context, transaction *models.UserPointsTransaction) error {
	query := `
		INSERT INTO points_transactions (id, user_id, amount, balance_after, source_type, source_id, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		transaction.ID,
		transaction.UserID,
		transaction.Amount,
		transaction.BalanceAfter,
		transaction.SourceType,
		transaction.SourceID,
		transaction.Description,
		transaction.CreatedAt,
	)
	return err
}

// GetByUserID retrieves transactions for a user
func (r *PointsTransactionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.UserPointsTransaction, error) {
	query := `
		SELECT id, user_id, amount, balance_after, source_type, source_id, description, created_at
		FROM points_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []*models.UserPointsTransaction{}
	for rows.Next() {
		t := &models.UserPointsTransaction{}
		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Amount,
			&t.BalanceAfter,
			&t.SourceType,
			&t.SourceID,
			&t.Description,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

// GetTotalByUserID retrieves total points for a user
func (r *PointsTransactionRepository) GetTotalByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM points_transactions WHERE user_id = $1`
	var total int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&total)
	return total, err
}
