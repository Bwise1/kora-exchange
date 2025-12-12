package auditlogs

import (
	"net/http"
	"strconv"

	"github.com/Bwise1/interstellar/internal/utils"
	"github.com/Bwise1/interstellar/pkg/response"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for audit logs
type Handler struct {
	service *Service
}

// NewHandler creates a new audit log handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GET /api/audit-logs - Get audit logs for authenticated user
func (h *Handler) GetUserAuditLogs(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT context
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok || userID == uuid.Nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse query parameters for pagination
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

	// Get audit logs for the user
	logs, err := h.service.GetByUserID(r.Context(), userID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch audit logs")
		return
	}

	response.Success(w, http.StatusOK, "Audit logs retrieved successfully", logs)
}
