package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for users
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new user repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var user User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	var user User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// // Update updates a user's information
// func (r *Repository) Update(ctx context.Context, user *User) error {
// 	query := `
// 		UPDATE users
// 		SET name = $1, email = $2, updated_at = $3
// 		WHERE id = $4 AND deleted_at IS NULL
// 	`

// 	// TODO: Implement database update
// 	return nil
// }

// // Delete soft deletes a user
// func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
// 	query := `
// 		UPDATE users
// 		SET deleted_at = NOW()
// 		WHERE id = $1 AND deleted_at IS NULL
// 	`

// 	// TODO: Implement soft delete
// 	return nil
// }

// EmailExists checks if an email already exists
func (r *Repository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
