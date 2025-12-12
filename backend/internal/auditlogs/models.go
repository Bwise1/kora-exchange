package auditlogs

import (
	"time"

	"github.com/google/uuid"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID            uuid.UUID  `json:"id"`
	UserID        *uuid.UUID `json:"user_id,omitempty"`
	Operation     string     `json:"operation"`
	ClientIP      string     `json:"client_ip"`
	UserAgent     string     `json:"user_agent,omitempty"`
	RequestMethod string     `json:"request_method,omitempty"`
	RequestPath   string     `json:"request_path,omitempty"`
	RequestBody   string     `json:"request_body,omitempty"`
	Timestamp     time.Time  `json:"timestamp"`
}

// CreateAuditLogRequest represents a request to create an audit log entry
type CreateAuditLogRequest struct {
	UserID        *uuid.UUID
	Operation     string
	ClientIP      string
	UserAgent     string
	RequestMethod string
	RequestPath   string
	RequestBody   string
}
