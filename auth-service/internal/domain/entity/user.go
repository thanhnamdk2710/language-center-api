package entity

import (
	"errors"
	"time"

	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
)

const (
	MaxFailedLoginAttempts = 5
	LockDuration           = 30 * time.Minute
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

func NewUser(email valueobject.Email, passwordHash string, now time.Time) (*User, error) {
	if passwordHash == "" {
		return nil, errors.New("password hash is required")
	}

	return &User{
		ID:                  valueobject.NewID(),
		Email:               email,
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

func (u *User) Activate(now time.Time) error {
	if u.Status != valueobject.UserStatusPending {
		return valueobject.ErrUserNotInPending
	}
	u.Status = valueobject.UserStatusActive
	u.EmailVerifiedAt = &now
	u.UpdatedAt = now
	return nil
}

func (u *User) RecordFailedLogin(now time.Time) bool {
	u.FailedLoginAttempts++
	if u.FailedLoginAttempts >= MaxFailedLoginAttempts {
		lockedUntil := now.Add(LockDuration)
		u.LockedUntil = &lockedUntil
		u.UpdatedAt = now
		return true
	}
	return false
}

func (u *User) ResetFailedAttempts() {
	u.FailedLoginAttempts = 0
	u.LockedUntil = nil
}

func (u *User) ValidateLoginEligibility(now time.Time) error {
	if u.Status == valueobject.UserStatusPending {
		return valueobject.ErrEmailNotVerified
	}
	if u.Status == valueobject.UserStatusDisabled {
		return valueobject.ErrAccountDisabled
	}
	if u.IsLocked(now) {
		return valueobject.ErrAccountLocked
	}
	return nil
}
