package transactions

import (
	"time"

	"github.com/google/uuid"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "DEPOSIT"
	TransactionTypeSwap     TransactionType = "SWAP"
	TransactionTypeTransfer TransactionType = "TRANSFER"
	TransactionTypeWithdraw TransactionType = "WITHDRAW"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID                uuid.UUID         `json:"id"`
	TransactionType   TransactionType   `json:"transaction_type"`
	Status            TransactionStatus `json:"status"`
	WalletID          uuid.UUID         `json:"wallet_id"`
	UserID            uuid.UUID         `json:"user_id"`
	RecipientWalletID *uuid.UUID        `json:"recipient_wallet_id,omitempty"`
	FromCurrency      string            `json:"from_currency"`
	FromAmount        float64           `json:"from_amount"`
	ToCurrency        *string           `json:"to_currency,omitempty"`
	ToAmount          *float64          `json:"to_amount,omitempty"`
	ExchangeRate      *float64          `json:"exchange_rate,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

// DepositRequest represents a request to deposit funds
type DepositRequest struct {
	Currency string  `json:"currency" validate:"required"`
	Amount   float64 `json:"amount" validate:"required,gt=0"`
}

// SwapRequest represents a request to swap currencies
type SwapRequest struct {
	FromCurrency string  `json:"from_currency" validate:"required"`
	ToCurrency   string  `json:"to_currency" validate:"required"`
	Amount       float64 `json:"amount" validate:"required,gt=0"`
}

// TransferRequest represents a request to transfer funds
type TransferRequest struct {
	RecipientWalletAddress string  `json:"recipient_wallet_address" validate:"required"`
	FromCurrency           string  `json:"from_currency" validate:"required"`
	Amount                 float64 `json:"amount" validate:"required,gt=0"`
	ToCurrency             *string `json:"to_currency,omitempty"`
}
