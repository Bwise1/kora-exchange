package fxrates

import (
	"encoding/json"
	"net/http"

	"github.com/Bwise1/interstellar/pkg/response"
	"github.com/go-chi/chi/v5"
)

// Handler handles HTTP requests for FX rates
type Handler struct {
	service *Service
}

// NewHandler creates a new FX rates handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GET /api/fx-rates/:currency
func (h *Handler) GetRates(w http.ResponseWriter, r *http.Request) {
	baseCurrency := chi.URLParam(r, "currency")
	if baseCurrency == "" {
		baseCurrency = "USD" // Default to USD
	}

	rates, err := h.service.GetRates(baseCurrency)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch exchange rates: "+err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Exchange rates retrieved successfully", rates)
}

// GET /api/fx-rates
func (h *Handler) GetAllRates(w http.ResponseWriter, r *http.Request) {
	// Default to USD as base
	baseCurrency := r.URL.Query().Get("base")
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	rates, err := h.service.GetRates(baseCurrency)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch exchange rates: "+err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Exchange rates retrieved successfully", rates)
}

// POST /api/fx-rates/convert
func (h *Handler) Convert(w http.ResponseWriter, r *http.Request) {
	var req ConversionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.From == "" || req.To == "" {
		response.Error(w, http.StatusBadRequest, "From and To currencies are required")
		return
	}

	if req.Amount <= 0 {
		response.Error(w, http.StatusBadRequest, "Amount must be greater than 0")
		return
	}

	result, err := h.service.Convert(req.From, req.To, req.Amount)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to convert: "+err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Conversion successful", result)
}

// POST /api/fx-rates/refresh
func (h *Handler) RefreshRates(w http.ResponseWriter, r *http.Request) {
	baseCurrency := r.URL.Query().Get("base")
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	if err := h.service.RefreshCache(baseCurrency); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to refresh rates: "+err.Error())
		return
	}

	rates := h.service.getCachedRates()
	response.Success(w, http.StatusOK, "Exchange rates refreshed successfully", rates)
}
