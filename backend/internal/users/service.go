package users

import (
	"context"
	"errors"
	"time"

	"github.com/Bwise1/interstellar/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
)

// Service handles business logic for users
type Service struct {
	repo *Repository
}

// NewService creates a new user service
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
	// Check if email exists
	exists, err := s.repo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &User{
		ID:        uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns user with JWT token
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*User, string, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}
	if user == nil {
		return nil, "", ErrInvalidCredentials
	}

	// Compare password
	if err := s.comparePassword(user.Password, req.Password); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// generateToken creates a JWT token for a user
func (s *Service) generateToken(user *User) (string, error) {
	return utils.GenerateToken(user.ID, user.Email)
}

// // GetByID retrieves a user by ID
// func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
// 	user, err := s.repo.GetByID(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if user == nil {
// 		return nil, ErrUserNotFound
// 	}
// 	return user, nil
// }

// // Update updates user information
// func (s *Service) Update(ctx context.Context, user *User) error {
// 	user.UpdatedAt = time.Now()
// 	return s.repo.Update(ctx, user)
// }

// // Delete soft deletes a user
// func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
// 	return s.repo.Delete(ctx, id)
// }

// hashPassword hashes a plain text password
func (s *Service) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// comparePassword compares hashed password with plain text
func (s *Service) comparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
