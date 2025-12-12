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

// FXRateService defines the interface for FX rate operations
type FXRateService interface {
	GetRate(from, to string) (float64, error)
}

// mapToRealCurrency maps stablecoin codes to their real currency equivalents
func mapToRealCurrency(currency string) string {
	mapping := map[string]string{
		"cNGN": "NGN",
		"cXAF": "XAF",
		"USDx": "USD",
		"EURx": "EUR",
		"cGHS": "GHS",
		"cKES": "KES",
	}
	if realCurrency, exists := mapping[currency]; exists {
		return realCurrency
	}
	return currency // Return as-is if not in mapping
}

// Service handles business logic for transactions
type Service struct {
	repo       *Repository
	walletRepo WalletRepository
	fxService  FXRateService
}

// NewService creates a new transaction service
func NewService(repo *Repository, walletRepo WalletRepository, fxService FXRateService) *Service {
	return &Service{
		repo:       repo,
		walletRepo: walletRepo,
		fxService:  fxService,
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

// ProcessSwap handles swapping between currencies
func (s *Service) ProcessSwap(ctx context.Context, userID uuid.UUID, req *SwapRequest) (*Transaction, error) {
	// Get user's wallet
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	if wallet == nil {
		return nil, fmt.Errorf("wallet not found for user")
	}

	// Check if user has sufficient balance
	currentBalance := wallet.GetBalance(req.FromCurrency)
	if currentBalance < req.Amount {
		return nil, fmt.Errorf("insufficient balance: have %f, need %f", currentBalance, req.Amount)
	}

	// Map stablecoin codes to real currency codes for FX service
	fromCurrency := mapToRealCurrency(req.FromCurrency)
	toCurrency := mapToRealCurrency(req.ToCurrency)

	// Get exchange rate between currencies
	rate, err := s.fxService.GetRate(fromCurrency, toCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	if rate <= 0 {
		return nil, fmt.Errorf("invalid exchange rate: %f", rate)
	}

	// Calculate converted amount
	convertedAmount := req.Amount * rate

	// Update balances
	wallet.SetBalance(req.FromCurrency, currentBalance-req.Amount)
	wallet.SetBalance(req.ToCurrency, wallet.GetBalance(req.ToCurrency)+convertedAmount)
	wallet.SetUpdatedAt(time.Now())

	// Save updated wallet
	if err := s.walletRepo.UpdateBalances(ctx, wallet); err != nil {
		return nil, fmt.Errorf("failed to update wallet balance: %w", err)
	}

	// Create transaction record
	toCurrency = req.ToCurrency
	toAmount := convertedAmount
	exchangeRate := rate

	tx := &Transaction{
		ID:              uuid.New(),
		TransactionType: TransactionTypeSwap,
		Status:          TransactionStatusCompleted,
		WalletID:        wallet.ID,
		UserID:          userID,
		FromCurrency:    req.FromCurrency,
		FromAmount:      req.Amount,
		ToCurrency:      &toCurrency,
		ToAmount:        &toAmount,
		ExchangeRate:    &exchangeRate,
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

// ProcessTransfer handles transferring funds between wallets
func (s *Service) ProcessTransfer(ctx context.Context, senderUserID uuid.UUID, req *TransferRequest) (*Transaction, error) {
	// Get sender's wallet
	senderWallet, err := s.walletRepo.GetByUserID(ctx, senderUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender wallet: %w", err)
	}
	if senderWallet == nil {
		return nil, fmt.Errorf("sender wallet not found")
	}

	// Get recipient's wallet by address
	recipientWallet, err := s.walletRepo.GetByAddress(ctx, req.RecipientWalletAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipient wallet: %w", err)
	}
	if recipientWallet == nil {
		return nil, fmt.Errorf("recipient wallet not found")
	}

	// Prevent sending to self
	if senderWallet.ID == recipientWallet.ID {
		return nil, fmt.Errorf("cannot transfer to your own wallet")
	}

	// Check if sender has sufficient balance
	senderBalance := senderWallet.GetBalance(req.FromCurrency)
	if senderBalance < req.Amount {
		return nil, fmt.Errorf("insufficient balance: have %f, need %f", senderBalance, req.Amount)
	}

	// Determine target currency (default to same as source if not specified)
	toCurrency := req.FromCurrency
	if req.ToCurrency != nil && *req.ToCurrency != "" {
		toCurrency = *req.ToCurrency
	}

	var receivedAmount float64
	var exchangeRate *float64

	// If currencies are different, perform conversion
	if req.FromCurrency != toCurrency {
		// Map stablecoin codes to real currency codes for FX service
		fromReal := mapToRealCurrency(req.FromCurrency)
		toReal := mapToRealCurrency(toCurrency)

		// Get exchange rate
		rate, err := s.fxService.GetRate(fromReal, toReal)
		if err != nil {
			return nil, fmt.Errorf("failed to get exchange rate: %w", err)
		}

		if rate <= 0 {
			return nil, fmt.Errorf("invalid exchange rate: %f", rate)
		}

		receivedAmount = req.Amount * rate
		exchangeRate = &rate
	} else {
		// Same currency, no conversion needed
		receivedAmount = req.Amount
	}

	// Update sender's balance
	senderWallet.SetBalance(req.FromCurrency, senderBalance-req.Amount)
	senderWallet.SetUpdatedAt(time.Now())

	// Update recipient's balance
	recipientBalance := recipientWallet.GetBalance(toCurrency)
	recipientWallet.SetBalance(toCurrency, recipientBalance+receivedAmount)
	recipientWallet.SetUpdatedAt(time.Now())

	// Save both wallets
	if err := s.walletRepo.UpdateBalances(ctx, senderWallet); err != nil {
		return nil, fmt.Errorf("failed to update sender wallet: %w", err)
	}

	if err := s.walletRepo.UpdateBalances(ctx, recipientWallet); err != nil {
		// TODO: Rollback sender wallet update
		return nil, fmt.Errorf("failed to update recipient wallet: %w", err)
	}

	// Create transaction record
	toCurrencyPtr := &toCurrency
	receivedAmountPtr := &receivedAmount

	tx := &Transaction{
		ID:                uuid.New(),
		TransactionType:   TransactionTypeTransfer,
		Status:            TransactionStatusCompleted,
		WalletID:          senderWallet.ID,
		UserID:            senderUserID,
		RecipientWalletID: &recipientWallet.ID,
		FromCurrency:      req.FromCurrency,
		FromAmount:        req.Amount,
		ToCurrency:        toCurrencyPtr,
		ToAmount:          receivedAmountPtr,
		ExchangeRate:      exchangeRate,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.repo.Create(ctx, tx); err != nil {
		// TODO: Rollback wallet updates
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	return tx, nil
}
