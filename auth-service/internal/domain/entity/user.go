package entity

import (
	"errors"
	"time"

	"github.com/thanhnamdk2710/auth-service/internal/domain/service/password"
	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
)

type User struct {
	ID                  valueobject.ID
	Email               valueobject.Email
	Password            string
	Status              valueobject.UserStatus
	EmailVerifiedAt     *time.Time
	FailedLoginAttempts int
	LockedUntil         *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func NewUser(email, passwordHash string) (*User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	if passwordHash == "" {
		return nil, errors.New("password hash is required")
	}

	emailValid, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, errors.New("email is required")
	}

	now := time.Now()
	return &User{
		ID:                  valueobject.NewID(),
		Email:               emailValid,
		Password:            passwordHash,
		Status:              valueobject.UserStatusPending,
		FailedLoginAttempts: 0,
		CreatedAt:           now,
		UpdatedAt:           now,
	}, nil
}

func (u *User) IsLocked(now time.Time) bool {
	return u.LockedUntil != nil && now.Before(*u.LockedUntil)
}

func (u *User) Activate() error {
	if u.Status != valueobject.UserStatusPending {
		return errors.New("user is not in pending status")
	}
	u.Status = valueobject.UserStatusActive
	now := time.Now()
	u.EmailVerifiedAt = &now
	return nil
}

func (u *User) RecordFailedLogin() error {
	u.FailedLoginAttempts++
	if u.FailedLoginAttempts >= 5 {
		return u.Lock(30 * time.Minute)
	}
	return nil
}

func (u *User) Lock(duration time.Duration) error {
	lockedUntil := time.Now().Add(duration)
	u.LockedUntil = &lockedUntil
	return nil
}

func (u *User) ResetFailedAttempts() {
	u.FailedLoginAttempts = 0
	u.LockedUntil = nil
}

func (u *User) CanLogin() error {
	if u.Status == valueobject.UserStatusPending {
		return valueobject.ErrEmailNotVerified
	}
	if u.Status == valueobject.UserStatusDisabled {
		return valueobject.ErrAccountDisabled
	}
	if u.IsLocked(time.Now()) {
		return valueobject.ErrAccountLocked
	}
	return nil
}

func (u *User) VerifyPassword(plainPassword string, passwordSvc password.Service) bool {
	return passwordSvc.Verify(plainPassword, u.Password)
}
