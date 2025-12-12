package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/Bwise1/interstellar/internal/auditlogs"
	"github.com/google/uuid"
)

// AuditLoggerService defines the interface for audit logging
type AuditLoggerService interface {
	LogRequest(ctx context.Context, req *auditlogs.CreateAuditLogRequest) error
}

// AuditMiddleware logs all API requests for security and compliance
func AuditMiddleware(auditService AuditLoggerService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip audit logging for health check and public endpoints that don't need auditing
			if shouldSkipAudit(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Extract user ID from context if user is authenticated
			var userID *uuid.UUID
			id, ok := r.Context().Value(UserIDKey).(uuid.UUID)
			if ok && id != uuid.Nil {
				userID = &id
			}

			// Get client IP
			clientIP := getClientIP(r)

			// Get user agent
			userAgent := r.UserAgent()

			// Read request body for logging (if applicable)
			var requestBody string
			if shouldLogBody(r.Method) {
				bodyBytes, err := io.ReadAll(r.Body)
				if err == nil {
					// Restore the body for the next handler
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

					// Sanitize sensitive data before logging
					requestBody = sanitizeBody(string(bodyBytes))

					// Limit body size for logging
					if len(requestBody) > 1000 {
						requestBody = requestBody[:1000] + "... (truncated)"
					}
				}
			}

			operation := determineOperation(r.Method, r.URL.Path)

			go func() {
				_ = auditService.LogRequest(context.Background(), &auditlogs.CreateAuditLogRequest{
					UserID:        userID,
					Operation:     operation,
					ClientIP:      clientIP,
					UserAgent:     userAgent,
					RequestMethod: r.Method,
					RequestPath:   r.URL.Path,
					RequestBody:   requestBody,
				})
			}()

			// Continue with the request
			next.ServeHTTP(w, r)
		})
	}
}

// determine if a path should skip audit logging
func shouldSkipAudit(path string) bool {
	skipPaths := []string{
		"/health",
		"/api/fx-rates", // Public endpoint
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	return false
}

// determine if request body should be logged
func shouldLogBody(method string) bool {
	return method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch
}

// remove sensitive information from request body
func sanitizeBody(body string) string {

	sanitized := body

	if strings.Contains(sanitized, "password") {

		sanitized = strings.ReplaceAll(sanitized, "\"password\":", "\"password\":\"***REDACTED***\"")
	}

	return sanitized
}

// determineOperation maps HTTP method and path to operation name
func determineOperation(method, path string) string {
	if strings.Contains(path, "/auth/login") {
		return "LOGIN"
	}
	if strings.Contains(path, "/auth/register") {
		return "REGISTER"
	}
	if strings.Contains(path, "/transactions/deposit") {
		return "DEPOSIT"
	}
	if strings.Contains(path, "/transactions/swap") {
		return "SWAP"
	}
	if strings.Contains(path, "/transactions/transfer") {
		return "TRANSFER"
	}
	if strings.Contains(path, "/transactions") && method == http.MethodGet {
		return "VIEW_TRANSACTIONS"
	}
	if strings.Contains(path, "/wallets") && method == http.MethodGet {
		return "VIEW_WALLET"
	}

	return method + " " + path
}

// getClientIP extracts the real client IP from the request
// Supports Cloudflare, standard proxies, and IPv6
func getClientIP(r *http.Request) string {
	// Priority 1: Cloudflare specific header (most reliable when behind Cloudflare)
	cfConnectingIP := r.Header.Get("CF-Connecting-IP")
	if cfConnectingIP != "" {
		return sanitizeIP(cfConnectingIP)
	}

	// Priority 2: True-Client-IP (used by some CDNs including Cloudflare Enterprise)
	trueClientIP := r.Header.Get("True-Client-IP")
	if trueClientIP != "" {
		return sanitizeIP(trueClientIP)
	}

	// Priority 3: X-Forwarded-For header (standard proxy header)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Get the first IP in the comma-separated list (original client)
		parts := strings.Split(forwarded, ",")
		if len(parts) > 0 {
			return sanitizeIP(parts[0])
		}
	}

	// Priority 4: X-Real-IP header (used by nginx and others)
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return sanitizeIP(realIP)
	}

	// Priority 5: Fall back to RemoteAddr
	ip := r.RemoteAddr
	return sanitizeIP(ip)
}

// sanitizeIP cleans up the IP address
// Handles IPv4, IPv6, and removes port numbers
func sanitizeIP(ip string) string {
	// Trim whitespace
	ip = strings.TrimSpace(ip)

	// Handle IPv6 addresses with port (e.g., [::1]:8080)
	if strings.HasPrefix(ip, "[") {
		// Extract IP from [ip]:port format
		endBracket := strings.Index(ip, "]")
		if endBracket != -1 {
			ip = ip[1:endBracket]
		}
	} else {
		// Handle IPv4 with port (e.g., 192.168.1.1:8080)
		// For IPv6 without brackets, don't split on colons
		if strings.Count(ip, ":") == 1 {
			// IPv4 with port
			parts := strings.Split(ip, ":")
			ip = parts[0]
		}
	}

	// Convert localhost IPv6 to readable format
	if ip == "::1" {
		return "localhost"
	}

	// Convert localhost IPv4 to readable format
	if ip == "127.0.0.1" {
		return "localhost"
	}

	return ip
}
