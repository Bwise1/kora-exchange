package middleware

import (
	"net/http"
	"strings"

	"github.com/Bwise1/interstellar/internal/utils"
)

// AuthMiddleware validates JWT tokens and adds user info to context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			if err == utils.ErrExpiredToken {
				respondWithError(w, http.StatusUnauthorized, "Token has expired")
				return
			}
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Add user info to context using utils constants
		ctx := r.Context()
		ctx = utils.SetUserIDInContext(ctx, claims.UserID)
		ctx = utils.SetEmailInContext(ctx, claims.Email)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper function for error responses
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(`{"success":false,"error":"` + message + `"}`))
}
