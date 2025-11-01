package domain

import (
	"context"
	"time"
)

type Session struct {
	UserID       string
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

type SessionService interface {
	// Store saves a refresh token session
	Store(ctx context.Context, session *Session) error

	// Get retrieves a session by refresh token
	Get(ctx context.Context, refreshToken string) (*Session, error)

	// Delete removes a session (logout)
	Delete(ctx context.Context, refreshToken string) error

	// DeleteAllByUserID removes all sessions for a user (logout all devices)
	DeleteAllByUserID(ctx context.Context, userID string) error
}
