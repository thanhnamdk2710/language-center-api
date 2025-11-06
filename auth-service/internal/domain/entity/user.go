package entity

import (
	"time"

	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
)

type User struct {
	ID                  string
	Email               string
	Password            string
	Status              valueobject.UserStatus
	EmailVerifiedAt     *time.Time
	FailedLoginAttempts int
	LockedUntil         *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (u *User) IsLocked(now time.Time) bool {
	return u.LockedUntil != nil && now.Before(*u.LockedUntil)
}
