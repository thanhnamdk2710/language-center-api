package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
	"github.com/thanhnamdk2710/auth-service/internal/domain/repository"
	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, email, status, password, failed_login_attempts, locked_until
			  FROM users WHERE email = $1 LIMIT 1`
	var u entity.User
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
	_, err := r.db.ExecContext(ctx, query, valueobject.UserStatusActive, userID)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}
	return nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (email, password, status, created_at, updated_at)
			  VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`
	return r.db.QueryRowContext(ctx, query, user.Email, user.Password, valueobject.UserStatusPending).Scan(&user.ID)
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users 
		SET 
			email = $1,
			password = $2,
			status = $3,
			email_verified_at = $4,
			failed_login_attempts = $5,
			locked_until = $6,
			updated_at = NOW()
		WHERE id = $7
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Email.String(),
		user.Password,
		user.Status,
		user.EmailVerifiedAt,
		user.FailedLoginAttempts,
		user.LockedUntil,
		user.ID.String(),
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
