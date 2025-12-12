package auditlogs

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service handles business logic for audit logs
type Service struct {
	repo *Repository
}

// NewService creates a new audit log service
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// LogRequest logs an API request
func (s *Service) LogRequest(ctx context.Context, req *CreateAuditLogRequest) error {
	return s.repo.Create(ctx, req)
}

// GetByUserID retrieves audit logs for a user
func (s *Service) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*AuditLog, error) {
	return s.repo.GetByUserID(ctx, userID, limit, offset)
}

// GetByOperation retrieves audit logs for a specific operation
func (s *Service) GetByOperation(ctx context.Context, operation string, limit, offset int) ([]*AuditLog, error) {
	return s.repo.GetByOperation(ctx, operation, limit, offset)
}

// GetByDateRange retrieves audit logs within a date range
func (s *Service) GetByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*AuditLog, error) {
	return s.repo.GetByDateRange(ctx, startDate, endDate, limit, offset)
}

// GetAll retrieves all audit logs
func (s *Service) GetAll(ctx context.Context, limit, offset int) ([]*AuditLog, error) {
	return s.repo.GetAll(ctx, limit, offset)
}
