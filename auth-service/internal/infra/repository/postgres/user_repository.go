package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/thanhnamdk2710/auth-service/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, status, password, failed_login_attempts, locked_until
			  FROM users WHERE email = $1 LIMIT 1`
	var u domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.Status,
		&u.Password,
		&u.FailedLoginAttempts,
		&u.LockedUntil,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) VerifyEmail(ctx context.Context, userID string) error {
	query := `UPDATE users SET status = $1, email_verified_at = NOW(), updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, domain.UserStatusActive, userID)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}
	return nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (email, password, status, created_at, updated_at)
			  VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`
	return r.db.QueryRowContext(ctx, query, user.Email, user.Password, domain.UserStatusPending).Scan(&user.ID)
}

func (r *userRepository) IncrementFailedAttempts(ctx context.Context, userID string) error {
	query := `UPDATE users SET failed_login_attemps = failed_login_attemps + 1, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to increment failed attempts: %w", err)
	}
	return nil
}

func (r *userRepository) ResetFailedAttempts(ctx context.Context, userID string) error {
	query := `UPDATE users SET failed_login_attemps = 0, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to reset failed attempts: %w", err)
	}
	return nil
}

func (r *userRepository) LockAccount(ctx context.Context, userID string, duration time.Duration) error {
	lockedUntil := time.Now().Add(duration)
	query := `UPDATE users SET locked_until = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, lockedUntil, userID)
	if err != nil {
		return fmt.Errorf("failed to lock account: %w", err)
	}
	return nil
}
