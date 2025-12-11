package wallets

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrWalletNotFound    = errors.New("wallet not found")
	ErrWalletExists      = errors.New("user already has a wallet")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidCurrency   = errors.New("invalid currency")
	ErrInvalidAmount     = errors.New("invalid amount")
)

// Service handles business logic for wallets
type Service struct {
	repo *Repository
}

// NewService creates a new wallet service
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateWallet creates a new wallet for a user (satisfies users.WalletService interface)
func (s *Service) CreateWallet(ctx context.Context, userID uuid.UUID) error {
	_, err := s.createWalletInternal(ctx, userID)
	return err
}

// create a wallet
func (s *Service) createWalletInternal(ctx context.Context, userID uuid.UUID) (*Wallet, error) {

	exists, err := s.repo.UserHasWallet(ctx, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrWalletExists
	}

	// Generate wallet address
	walletAddress := generateWalletAddress(userID)

	wallet := &Wallet{
		ID:            uuid.New(),
		UserID:        userID,
		WalletAddress: walletAddress,
		Balances:      make(map[string]float64),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.Create(ctx, wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

// GetWalletByID
func (s *Service) GetWalletByID(ctx context.Context, id uuid.UUID) (*Wallet, error) {
	wallet, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, ErrWalletNotFound
	}
	return wallet, nil
}

// GetWalletByUserID
func (s *Service) GetWalletByUserID(ctx context.Context, userID uuid.UUID) (*Wallet, error) {
	wallet, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, ErrWalletNotFound
	}
	return wallet, nil
}

// GetWalletByAddress
func (s *Service) GetWalletByAddress(ctx context.Context, address string) (*Wallet, error) {
	wallet, err := s.repo.GetByAddress(ctx, address)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, ErrWalletNotFound
	}
	return wallet, nil
}

// GetBalance
func (s *Service) GetBalance(ctx context.Context, userID uuid.UUID, currency string) (float64, error) {
	wallet, err := s.GetWalletByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return wallet.GetBalance(currency), nil
}

// GetAllBalances
func (s *Service) GetAllBalances(ctx context.Context, userID uuid.UUID) (map[string]float64, error) {
	wallet, err := s.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return wallet.Balances, nil
}

// UpdateBalance
func (s *Service) UpdateBalance(ctx context.Context, walletID uuid.UUID, currency string, amount float64) error {
	wallet, err := s.GetWalletByID(ctx, walletID)
	if err != nil {
		return err
	}

	wallet.UpdateBalance(currency, amount)

	// Check for negative balance
	if wallet.GetBalance(currency) < 0 {
		return ErrInsufficientFunds
	}

	wallet.UpdatedAt = time.Now()

	return s.repo.UpdateBalances(ctx, wallet)
}

// SetBalance
func (s *Service) SetBalance(ctx context.Context, walletID uuid.UUID, currency string, amount float64) error {
	if amount < 0 {
		return ErrInvalidAmount
	}

	wallet, err := s.GetWalletByID(ctx, walletID)
	if err != nil {
		return err
	}

	wallet.SetBalance(currency, amount)
	wallet.UpdatedAt = time.Now()

	return s.repo.UpdateBalances(ctx, wallet)
}

// AddBalance
func (s *Service) AddBalance(ctx context.Context, walletID uuid.UUID, currency string, amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	return s.UpdateBalance(ctx, walletID, currency, amount)
}

// DeductBalance
func (s *Service) DeductBalance(ctx context.Context, walletID uuid.UUID, currency string, amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	return s.UpdateBalance(ctx, walletID, currency, -amount)
}

// generate wallet address
func generateWalletAddress(userID uuid.UUID) string {
	// extract first 8 chars
	idStr := strings.ReplaceAll(userID.String(), "-", "")
	return fmt.Sprintf("WLT-%s", strings.ToUpper(idStr[:8]))
}
