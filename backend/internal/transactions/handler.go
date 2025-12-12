package transactions

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bwise1/interstellar/internal/middleware"
	"github.com/Bwise1/interstellar/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for transactions
type Handler struct {
	service *Service
}

// NewHandler creates a new transaction handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// POST /api/transactions/deposit
func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT context
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Currency == "" {
		response.Error(w, http.StatusBadRequest, "Currency is required")
		return
	}
	if req.Amount <= 0 {
		response.Error(w, http.StatusBadRequest, "Amount must be greater than 0")
		return
	}

	tx, err := h.service.ProcessDeposit(r.Context(), userID, &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusCreated, "Deposit successful", tx)
}

// POST /api/transactions/swap
func (h *Handler) Swap(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT context
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req SwapRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.FromCurrency == "" || req.ToCurrency == "" {
		response.Error(w, http.StatusBadRequest, "FromCurrency and ToCurrency are required")
		return
	}

	if req.FromCurrency == req.ToCurrency {
		response.Error(w, http.StatusBadRequest, "Cannot swap same currency")
		return
	}

	if req.Amount <= 0 {
		response.Error(w, http.StatusBadRequest, "Amount must be greater than 0")
		return
	}

	tx, err := h.service.ProcessSwap(r.Context(), userID, &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusCreated, "Swap successful", tx)
}

// GET /api/transactions?limit=10&offset=0
func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT context
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // default limit
	offset := 0 // default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	transactions, err := h.service.GetTransactionsByUser(r.Context(), userID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get transactions")
		return
	}

	response.Success(w, http.StatusOK, "Transactions retrieved successfully", transactions)
}

// GET /api/transactions/{id}
func (h *Handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT context
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get transaction ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	// Get transaction
	tx, err := h.service.GetTransaction(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Transaction not found")
		return
	}

	// Verify the transaction belongs to the user
	if tx.UserID != userID {
		response.Error(w, http.StatusForbidden, "Access denied")
		return
	}

	response.Success(w, http.StatusOK, "Transaction retrieved successfully", tx)
}
