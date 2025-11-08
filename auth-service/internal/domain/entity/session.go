package entity

import (
	"errors"
	"time"

	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
)

type Session struct {
	UserID       valueobject.ID
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

func NewSession(userID valueobject.ID, refreshToken string, expiresAt time.Time) (*Session, error) {
	if userID == "" || refreshToken == "" {
		return nil, errors.New("userID and refreshToken are required")
	}
	if expiresAt.Before(time.Now()) {
		return nil, errors.New("expiresAt must be in the future")
	}

	return &Session{
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		CreatedAt:    time.Now(),
	}, nil
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) Renew(duration time.Duration) {
	s.ExpiresAt = time.Now().Add(duration)
}

func (s *Session) IsValid() error {
	if s.IsExpired() {
		return valueobject.ErrSessionExpired
	}
	if s.RefreshToken == "" {
		return valueobject.ErrSessionInvalid
	}
	return nil
}
