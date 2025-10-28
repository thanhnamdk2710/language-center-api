package postgres

import (
	"context"
	"database/sql"

	"github.com/thanhnamdk2710/auth-service/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, status FROM users WHERE email=$1 LIMIT 1`
	var u domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Status)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (email, password, status, created_at, updated_at)
			  VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`
	return r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, domain.UserStatusPending).Scan(&user.ID)
}
