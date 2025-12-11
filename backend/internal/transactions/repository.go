package transactions

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create
func (r *Repository) Create(ctx context.Context, tx *Transaction) error {
	query := `
		INSERT INTO transactions (
			id, transaction_type, status, wallet_id, user_id,
			recipient_wallet_id, from_currency, from_amount,
			to_currency, to_amount, exchange_rate, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		tx.ID,
		tx.TransactionType,
		tx.Status,
		tx.WalletID,
		tx.UserID,
		tx.RecipientWalletID,
		tx.FromCurrency,
		tx.FromAmount,
		tx.ToCurrency,
		tx.ToAmount,
		tx.ExchangeRate,
		tx.CreatedAt,
		tx.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// GetByID
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Transaction, error) {
	query := `
		SELECT id, transaction_type, status, wallet_id, user_id,
			   recipient_wallet_id, from_currency, from_amount,
			   to_currency, to_amount, exchange_rate, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`

	var tx Transaction
	err := r.db.QueryRow(ctx, query, id).Scan(
		&tx.ID,
		&tx.TransactionType,
		&tx.Status,
		&tx.WalletID,
		&tx.UserID,
		&tx.RecipientWalletID,
		&tx.FromCurrency,
		&tx.FromAmount,
		&tx.ToCurrency,
		&tx.ToAmount,
		&tx.ExchangeRate,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &tx, nil
}

// GetByWalletID
func (r *Repository) GetByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*Transaction, error) {
	query := `
		SELECT id, transaction_type, status, wallet_id, user_id,
			   recipient_wallet_id, from_currency, from_amount,
			   to_currency, to_amount, exchange_rate, created_at, updated_at
		FROM transactions
		WHERE wallet_id = $1 OR recipient_wallet_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, walletID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by wallet: %w", err)
	}
	defer rows.Close()

	var transactions []*Transaction
	for rows.Next() {
		var tx Transaction
		err := rows.Scan(
			&tx.ID,
			&tx.TransactionType,
			&tx.Status,
			&tx.WalletID,
			&tx.UserID,
			&tx.RecipientWalletID,
			&tx.FromCurrency,
			&tx.FromAmount,
			&tx.ToCurrency,
			&tx.ToAmount,
			&tx.ExchangeRate,
			&tx.CreatedAt,
			&tx.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, &tx)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

// GetByUserID
func (r *Repository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Transaction, error) {
	query := `
		SELECT id, transaction_type, status, wallet_id, user_id,
			   recipient_wallet_id, from_currency, from_amount,
			   to_currency, to_amount, exchange_rate, created_at, updated_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by user: %w", err)
	}
	defer rows.Close()

	var transactions []*Transaction
	for rows.Next() {
		var tx Transaction
		err := rows.Scan(
			&tx.ID,
			&tx.TransactionType,
			&tx.Status,
			&tx.WalletID,
			&tx.UserID,
			&tx.RecipientWalletID,
			&tx.FromCurrency,
			&tx.FromAmount,
			&tx.ToCurrency,
			&tx.ToAmount,
			&tx.ExchangeRate,
			&tx.CreatedAt,
			&tx.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, &tx)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

// UpdateStatus
func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status TransactionStatus) error {
	query := `
		UPDATE transactions
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}
