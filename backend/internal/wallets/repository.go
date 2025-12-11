package wallets

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for wallets
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new wallet repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

// Create wallet
func (r *Repository) Create(ctx context.Context, wallet *Wallet) error {
	balancesJSON, err := json.Marshal(wallet.Balances)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO wallets (id, user_id, wallet_address, balances, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err = r.db.QueryRow(
		ctx,
		query,
		wallet.ID,
		wallet.UserID,
		wallet.WalletAddress,
		balancesJSON,
		wallet.CreatedAt,
		wallet.UpdatedAt,
	).Scan(&wallet.ID, &wallet.CreatedAt, &wallet.UpdatedAt)

	return err
}

// GetByID
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Wallet, error) {
	query := `
		SELECT id, user_id, wallet_address, balances, created_at, updated_at
		FROM wallets
		WHERE id = $1
	`

	var wallet Wallet
	var balancesJSON []byte

	err := r.db.QueryRow(ctx, query, id).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.WalletAddress,
		&balancesJSON,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Unmarshal JSONB to map
	if err := json.Unmarshal(balancesJSON, &wallet.Balances); err != nil {
		return nil, err
	}

	return &wallet, nil
}

// GetByUserID
func (r *Repository) GetByUserID(ctx context.Context, userID uuid.UUID) (*Wallet, error) {
	query := `
		SELECT id, user_id, wallet_address, balances, created_at, updated_at
		FROM wallets
		WHERE user_id = $1
	`

	var wallet Wallet
	var balancesJSON []byte

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.WalletAddress,
		&balancesJSON,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(balancesJSON, &wallet.Balances); err != nil {
		return nil, err
	}

	return &wallet, nil
}

// GetByAddress
func (r *Repository) GetByAddress(ctx context.Context, address string) (*Wallet, error) {
	query := `
		SELECT id, user_id, wallet_address, balances, created_at, updated_at
		FROM wallets
		WHERE wallet_address = $1
	`
	var wallet Wallet
	var balancesJSON []byte

	err := r.db.QueryRow(ctx, query, address).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.WalletAddress,
		&balancesJSON,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(balancesJSON, &wallet.Balances); err != nil {
		return nil, err
	}

	return &wallet, nil
}

// UpdateBalances
func (r *Repository) UpdateBalances(ctx context.Context, wallet *Wallet) error {
	balancesJSON, err := json.Marshal(wallet.Balances)
	if err != nil {
		return err
	}

	query := `
		UPDATE wallets
		SET balances = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err = r.db.Exec(ctx, query, balancesJSON, wallet.ID)
	return err
}

// Delete soft deletes a wallet (if you add deleted_at column later)
// For now, it's a hard delete
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM wallets WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// UserHasWallet
func (r *Repository) UserHasWallet(ctx context.Context, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM wallets WHERE user_id = $1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID).Scan(&exists)
	return exists, err
}
