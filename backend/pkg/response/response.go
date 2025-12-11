package response

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse is the standard success response format
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse is the standard error response format
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// Success sends a success response
func Success(w http.ResponseWriter, code int, message string, data interface{}) {
	response := SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// Error sends an error response
func Error(w http.ResponseWriter, code int, message string) {
	response := ErrorResponse{
		Success: false,
		Error:   message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
