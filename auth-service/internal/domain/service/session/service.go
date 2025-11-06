package session

import (
	"context"

	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
)

type Service interface {
	// Store saves a refresh token session
	Store(ctx context.Context, session *entity.Session) error

	// Get retrieves a session by refresh token
	Get(ctx context.Context, refreshToken string) (*entity.Session, error)

	// Delete removes a session (logout)
	Delete(ctx context.Context, refreshToken string) error

	// DeleteAllByUserID removes all sessions for a user (logout all devices)
	DeleteAllByUserID(ctx context.Context, userID string) error
}
