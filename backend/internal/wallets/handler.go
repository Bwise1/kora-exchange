package wallets

import (
	"net/http"

	"github.com/Bwise1/interstellar/internal/middleware"
	"github.com/Bwise1/interstellar/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for wallets
type Handler struct {
	service *Service
}

// NewHandler creates a new wallet handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GET /api/wallets
func (h *Handler) GetWallet(w http.ResponseWriter, r *http.Request) {

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	wallet, err := h.service.GetWalletByUserID(r.Context(), userID)
	if err != nil {
		if err == ErrWalletNotFound {
			response.Error(w, http.StatusNotFound, "Wallet not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve wallet")
		return
	}

	response.Success(w, http.StatusOK, "Wallet retrieved successfully", wallet.ToResponse())
}

// GET /api/wallets/:id
func (h *Handler) GetWalletByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	walletID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	wallet, err := h.service.GetWalletByID(r.Context(), walletID)
	if err != nil {
		if err == ErrWalletNotFound {
			response.Error(w, http.StatusNotFound, "Wallet not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve wallet")
		return
	}

	response.Success(w, http.StatusOK, "Wallet retrieved successfully", wallet.ToResponse())
}

// GET /api/wallets/balance/:currency
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	currency := chi.URLParam(r, "currency")
	if currency == "" {
		response.Error(w, http.StatusBadRequest, "Currency is required")
		return
	}

	balance, err := h.service.GetBalance(r.Context(), userID, currency)
	if err != nil {
		if err == ErrWalletNotFound {
			response.Error(w, http.StatusNotFound, "Wallet not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve balance")
		return
	}

	balanceResponse := BalanceResponse{
		Currency: currency,
		Balance:  balance,
	}

	response.Success(w, http.StatusOK, "Balance retrieved successfully", balanceResponse)
}

// GET /api/wallets/balances
func (h *Handler) GetAllBalances(w http.ResponseWriter, r *http.Request) {

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	balances, err := h.service.GetAllBalances(r.Context(), userID)
	if err != nil {
		if err == ErrWalletNotFound {
			response.Error(w, http.StatusNotFound, "Wallet not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve balances")
		return
	}

	response.Success(w, http.StatusOK, "Balances retrieved successfully", balances)
}
