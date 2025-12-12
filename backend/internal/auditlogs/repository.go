package auditlogs

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for audit logs
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new audit log repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create inserts a new audit log entry
func (r *Repository) Create(ctx context.Context, req *CreateAuditLogRequest) error {
	query := `
		INSERT INTO audit_logs (
			user_id, operation, client_ip, user_agent,
			request_method, request_path, request_body, timestamp
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		req.UserID,
		req.Operation,
		req.ClientIP,
		req.UserAgent,
		req.RequestMethod,
		req.RequestPath,
		req.RequestBody,
		time.Now(),
	)

	return err
}

// GetByUserID retrieves audit logs for a specific user
func (r *Repository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*AuditLog, error) {
	query := `
		SELECT id, user_id, operation, client_ip, user_agent,
		       request_method, request_path, request_body, timestamp
		FROM audit_logs
		WHERE user_id = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*AuditLog
	for rows.Next() {
		var log AuditLog
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Operation,
			&log.ClientIP,
			&log.UserAgent,
			&log.RequestMethod,
			&log.RequestPath,
			&log.RequestBody,
			&log.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	return logs, rows.Err()
}

// GetByOperation retrieves audit logs for a specific operation
func (r *Repository) GetByOperation(ctx context.Context, operation string, limit, offset int) ([]*AuditLog, error) {
	query := `
		SELECT id, user_id, operation, client_ip, user_agent,
		       request_method, request_path, request_body, timestamp
		FROM audit_logs
		WHERE operation = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, operation, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*AuditLog
	for rows.Next() {
		var log AuditLog
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Operation,
			&log.ClientIP,
			&log.UserAgent,
			&log.RequestMethod,
			&log.RequestPath,
			&log.RequestBody,
			&log.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	return logs, rows.Err()
}

// GetByDateRange retrieves audit logs within a date range
func (r *Repository) GetByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*AuditLog, error) {
	query := `
		SELECT id, user_id, operation, client_ip, user_agent,
		       request_method, request_path, request_body, timestamp
		FROM audit_logs
		WHERE timestamp BETWEEN $1 AND $2
		ORDER BY timestamp DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(ctx, query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*AuditLog
	for rows.Next() {
		var log AuditLog
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Operation,
			&log.ClientIP,
			&log.UserAgent,
			&log.RequestMethod,
			&log.RequestPath,
			&log.RequestBody,
			&log.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	return logs, rows.Err()
}

// GetAll retrieves all audit logs with pagination
func (r *Repository) GetAll(ctx context.Context, limit, offset int) ([]*AuditLog, error) {
	query := `
		SELECT id, user_id, operation, client_ip, user_agent,
		       request_method, request_path, request_body, timestamp
		FROM audit_logs
		ORDER BY timestamp DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*AuditLog
	for rows.Next() {
		var log AuditLog
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Operation,
			&log.ClientIP,
			&log.UserAgent,
			&log.RequestMethod,
			&log.RequestPath,
			&log.RequestBody,
			&log.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	return logs, rows.Err()
}
