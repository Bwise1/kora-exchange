package users

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Bwise1/interstellar/internal/utils"
	"github.com/Bwise1/interstellar/pkg/response"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for users
type Handler struct {
	service       *Service
	walletService WalletService
	auditPassword string
}

// WalletService interface for wallet operations
type WalletService interface {
	CreateWallet(ctx context.Context, userID uuid.UUID) error
}

// NewHandler creates a new user handler
func NewHandler(service *Service, walletService WalletService, auditPassword string) *Handler {
	return &Handler{
		service:       service,
		walletService: walletService,
		auditPassword: auditPassword,
	}
}

// Register handles user registration
// POST /api/auth/register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// TODO: Validate request

	user, err := h.service.Register(r.Context(), &req)
	if err != nil {
		if err == ErrEmailExists {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	// Auto-create wallet for new user
	if err := h.walletService.CreateWallet(r.Context(), user.ID); err != nil {
		// Log error but don't fail registration
		log.Printf("Failed to create wallet for user %s: %v", user.ID, err)
	}

	response.Success(w, http.StatusCreated, "User registered successfully", user.ToUserResponse())
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, token, err := h.service.Login(r.Context(), &req)
	if err != nil {
		if err == ErrInvalidCredentials {
			response.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to login")
		return
	}

	loginResponse := LoginResponse{
		Token: token,
		User:  user.ToUserResponse(),
	}

	response.Success(w, http.StatusOK, "Login successful", loginResponse)
}

// VerifyPassword verifies the audit logs access password
// POST /api/users/verify-password
func (h *Handler) VerifyPassword(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT context to ensure user is authenticated
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok || userID == uuid.Nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Password == "" {
		response.Error(w, http.StatusBadRequest, "Password is required")
		return
	}

	// Verify against generic audit password from env
	if req.Password != h.auditPassword {
		response.Error(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	response.Success(w, http.StatusOK, "Password verified successfully", map[string]bool{"verified": true})
}

// GetProfile gets the current user's profile
// GET /api/users/profile
// func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
// 	// TODO: Get user ID from JWT token in context
// 	userID := uuid.New() // Placeholder

// 	user, err := h.service.GetByID(r.Context(), userID)
// 	if err != nil {
// 		if err == ErrUserNotFound {
// 			respondWithError(w, http.StatusNotFound, err.Error())
// 			return
// 		}
// 		respondWithError(w, http.StatusInternalServerError, "Failed to get user")
// 		return
// 	}

// 	respondWithJSON(w, http.StatusOK, user.ToUserResponse())
// }

// UpdateProfile updates user profile
// PUT /api/users/profile
// func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
// 	// TODO: Get user ID from JWT token in context
// 	userID := uuid.New() // Placeholder

// 	var req struct {
// 		Name  string `json:"name"`
// 		Email string `json:"email"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
// 		return
// 	}

// 	user, err := h.service.GetByID(r.Context(), userID)
// 	if err != nil {
// 		respondWithError(w, http.StatusNotFound, "User not found")
// 		return
// 	}

// 	// Update fields
// 	user.Name = req.Name
// 	user.Email = req.Email

// 	if err := h.service.Update(r.Context(), user); err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
// 		return
// 	}

// 	respondWithJSON(w, http.StatusOK, user.ToUserResponse())
// }

// DeleteAccount deletes user account
// DELETE /api/users/profile
// func (h *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
// 	// TODO: Get user ID from JWT token in context
// 	userID := uuid.New() // Placeholder

// 	if err := h.service.Delete(r.Context(), userID); err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Failed to delete user")
// 		return
// 	}

// 	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Account deleted successfully"})
// }

