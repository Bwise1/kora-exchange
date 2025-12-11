package transactions

import (
	"context"
	"fmt"
	"time"

	"github.com/Bwise1/interstellar/internal/wallets"
	"github.com/google/uuid"
)

// WalletRepository defines the interface for wallet operations
type WalletRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*wallets.Wallet, error)
	GetByAddress(ctx context.Context, address string) (*wallets.Wallet, error)
	UpdateBalances(ctx context.Context, wallet *wallets.Wallet) error
}

// Service handles business logic for transactions
type Service struct {
	repo       *Repository
	walletRepo WalletRepository
}

// NewService creates a new transaction service
func NewService(repo *Repository, walletRepo WalletRepository) *Service {
	return &Service{
		repo:       repo,
		walletRepo: walletRepo,
	}
}

// ProcessDeposit
func (s *Service) ProcessDeposit(ctx context.Context, userID uuid.UUID, req *DepositRequest) (*Transaction, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	if wallet == nil {
		return nil, fmt.Errorf("wallet not found for user")
	}

	currentBalance := wallet.GetBalance(req.Currency)
	newBalance := currentBalance + req.Amount
	wallet.SetBalance(req.Currency, newBalance)
	wallet.SetUpdatedAt(time.Now())

	if err := s.walletRepo.UpdateBalances(ctx, wallet); err != nil {
		return nil, fmt.Errorf("failed to update wallet balance: %w", err)
	}

	tx := &Transaction{
		ID:              uuid.New(),
		TransactionType: TransactionTypeDeposit,
		Status:          TransactionStatusCompleted,
		WalletID:        wallet.ID,
		UserID:          userID,
		FromCurrency:    req.Currency,
		FromAmount:      req.Amount,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	return tx, nil
}

// GetTransactionsByWallet
func (s *Service) GetTransactionsByWallet(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*Transaction, error) {
	return s.repo.GetByWalletID(ctx, walletID, limit, offset)
}

// GetTransactionsByUser
func (s *Service) GetTransactionsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Transaction, error) {
	return s.repo.GetByUserID(ctx, userID, limit, offset)
}

// GetTransaction
func (s *Service) GetTransaction(ctx context.Context, id uuid.UUID) (*Transaction, error) {
	return s.repo.GetByID(ctx, id)
}
