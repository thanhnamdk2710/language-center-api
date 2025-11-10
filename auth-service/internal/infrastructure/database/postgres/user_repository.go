package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
	"github.com/thanhnamdk2710/auth-service/internal/domain/repository"
	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
)

const (
	// SQL queries
	queryGetByEmail = `
		SELECT 
			id, 
			email, 
			password, 
			status, 
			email_verified_at,
			failed_login_attempts, 
			locked_until,
			created_at,
			updated_at
		FROM users 
		WHERE email = $1 
		LIMIT 1
	`

	queryCreate = `
		INSERT INTO users (
			id, 
			email, 
			password, 
			status, 
			email_verified_at,
			failed_login_attempts, 
			locked_until, 
			created_at, 
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	queryUpdate = `
		UPDATE users 
		SET 
			email = $1,
			password = $2,
			status = $3,
			email_verified_at = $4,
			failed_login_attempts = $5,
			locked_until = $6,
			updated_at = $7
		WHERE id = $8
	`

	queryDelete = `DELETE FROM users WHERE id = $1`
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	var u entity.User
	err := r.db.QueryRowContext(ctx, queryGetByEmail, email).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Status,
		&u.EmailVerifiedAt,
		&u.FailedLoginAttempts,
		&u.LockedUntil,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &u, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	if user == nil {
		return errors.New("user is required")
	}

	if err := r.validateUser(user); err != nil {
		return fmt.Errorf("invalid user: %w", err)
	}

	result, err := r.db.ExecContext(
		ctx,
		queryCreate,
		user.ID.String(),
		user.Email.String(),
		user.Password,
		user.Status,
		user.EmailVerifiedAt,
		user.FailedLoginAttempts,
		user.LockedUntil,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	if user == nil {
		return errors.New("user is required")
	}

	if err := r.validateUser(user); err != nil {
		return fmt.Errorf("invalid user: %w", err)
	}

	result, err := r.db.ExecContext(
		ctx,
		queryUpdate,
		user.Email.String(),
		user.Password,
		user.Status,
		user.EmailVerifiedAt,
		user.FailedLoginAttempts,
		user.LockedUntil,
		user.UpdatedAt,
		user.ID.String(),
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return valueobject.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}

	result, err := r.db.ExecContext(ctx, queryDelete, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return valueobject.ErrUserNotFound
	}

	return nil
}

// validateUser validates the user entity before database operations
func (r *userRepository) validateUser(user *entity.User) error {
	if user.ID.String() == "" {
		return errors.New("user ID is required")
	}

	if user.Email.String() == "" {
		return errors.New("user email is required")
	}

	if user.Password == "" {
		return errors.New("user password is required")
	}

	return nil
}
