package wallets

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Wallet
type Wallet struct {
	ID            uuid.UUID          `json:"id"`
	UserID        uuid.UUID          `json:"user_id"`
	WalletAddress string             `json:"wallet_address"`
	Balances      map[string]float64 `json:"balances"` // {"cNGN": 5000, "USDx": 100}
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

// WalletResponse
type WalletResponse struct {
	ID            uuid.UUID          `json:"id"`
	UserID        uuid.UUID          `json:"user_id"`
	WalletAddress string             `json:"wallet_address"`
	Balances      map[string]float64 `json:"balances"`
	TotalUSD      float64            `json:"total_usd,omitempty"`
	CreatedAt     time.Time          `json:"created_at"`
}

// GetBalanceRequest
type GetBalanceRequest struct {
	Currency string `json:"currency"`
}

// BalanceResponse
type BalanceResponse struct {
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}

// convert Wallet to WalletResponse
func (w *Wallet) ToResponse() *WalletResponse {
	return &WalletResponse{
		ID:            w.ID,
		UserID:        w.UserID,
		WalletAddress: w.WalletAddress,
		Balances:      w.Balances,
		CreatedAt:     w.CreatedAt,
	}
}

func (w *Wallet) MarshalBalances() ([]byte, error) {
	return json.Marshal(w.Balances)
}

func (w *Wallet) UnmarshalBalances(data []byte) error {
	return json.Unmarshal(data, &w.Balances)
}

// GetBalance
func (w *Wallet) GetBalance(currency string) float64 {
	if balance, exists := w.Balances[currency]; exists {
		return balance
	}
	return 0
}

// SetBalance
func (w *Wallet) SetBalance(currency string, amount float64) {
	if w.Balances == nil {
		w.Balances = make(map[string]float64)
	}
	w.Balances[currency] = amount
}

// UpdateBalance
func (w *Wallet) UpdateBalance(currency string, amount float64) {
	if w.Balances == nil {
		w.Balances = make(map[string]float64)
	}
	currentBalance := w.GetBalance(currency)
	w.Balances[currency] = currentBalance + amount
}
