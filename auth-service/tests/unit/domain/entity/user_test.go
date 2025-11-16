package entity_test

import (
	"testing"
	"time"

	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
)

func TestNewUser(t *testing.T) {
	email, _ := valueobject.NewEmail("user@example.com")
	now := time.Now()

	user, err := entity.NewUser(email, "hashed-password", now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Email != email {
		t.Errorf("expected email %v, got %v", email, user.Email)
	}

	if user.Password != "hashed-password" {
		t.Errorf("expected password hash to be stored")
	}

	if user.Status != valueobject.UserStatusPending {
		t.Errorf("expected status PENDING, got %v", user.Status)
	}

	if !user.CreatedAt.Equal(now) {
		t.Errorf("expected CreatedAt to be set")
	}
	if !user.UpdatedAt.Equal(now) {
		t.Errorf("expected UpdatedAt to be set")
	}
}

func TestUser_IsLocked(t *testing.T) {
	email, _ := valueobject.NewEmail("user@example.com")
	now := time.Now()
	user, _ := entity.NewUser(email, "hash", now)

	// case: user not locked
	if user.IsLocked(now) {
		t.Error("expected user to NOT be locked")
	}

	// lock user
	future := now.Add(10 * time.Minute)
	user.LockedUntil = &future

	if !user.IsLocked(now) {
		t.Error("expected user to be locked")
	}
}

func TestUser_Activate(t *testing.T) {
	email, _ := valueobject.NewEmail("user@example.com")
	now := time.Now()
	user, _ := entity.NewUser(email, "hash", now)

	err := user.Activate(now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.Status != valueobject.UserStatusActive {
		t.Errorf("expected ACTIVE, got %v", user.Status)
	}

	if user.EmailVerifiedAt == nil {
		t.Errorf("expected EmailVerifiedAt to be set")
	}
}

func TestUser_ActivateShouldFail_WhenNotPending(t *testing.T) {
	email, _ := valueobject.NewEmail("user@example.com")
	now := time.Now()
	user, _ := entity.NewUser(email, "hash", now)

	user.Status = valueobject.UserStatusDisabled

	err := user.Activate(now)
	if err != valueobject.ErrUserNotInPending {
		t.Errorf("expected ErrUserNotInPending, got %v", err)
	}
}

func TestUser_RecordFailedLogin(t *testing.T) {
	email, _ := valueobject.NewEmail("user@example.com")
	now := time.Now()
	user, _ := entity.NewUser(email, "hash", now)

	locked := false
	for i := 0; i < entity.MaxFailedLoginAttempts; i++ {
		locked = user.RecordFailedLogin(now)
	}

	if !locked {
		t.Errorf("expected user to be locked after max failed attempts")
	}

	if user.LockedUntil == nil {
		t.Errorf("expected LockedUntil to be set")
	}
}

func TestUser_ResetFailedAttempts(t *testing.T) {
	email, _ := valueobject.NewEmail("john@example.com")
	now := time.Now()
	user, _ := entity.NewUser(email, "hash", now)

	user.FailedLoginAttempts = 3
	future := now.Add(5 * time.Minute)
	user.LockedUntil = &future

	user.ResetFailedAttempts(now)

	if user.FailedLoginAttempts != 0 {
		t.Errorf("expected attempts to be reset")
	}
	if user.LockedUntil != nil {
		t.Errorf("expected LockedUntil to be nil")
	}
}

func TestUser_EnsureCanLogin(t *testing.T) {
	email, _ := valueobject.NewEmail("john@example.com")
	now := time.Now()
	user, _ := entity.NewUser(email, "hash", now)

	// Case: pending
	user.Status = valueobject.UserStatusPending
	if err := user.EnsureCanLogin(now); err != valueobject.ErrEmailNotVerified {
		t.Errorf("expected ErrEmailNotVerified")
	}

	// Case: disabled
	user.Status = valueobject.UserStatusDisabled
	if err := user.EnsureCanLogin(now); err != valueobject.ErrAccountDisabled {
		t.Errorf("expected ErrAccountDisabled")
	}

	// Case: locked
	user.Status = valueobject.UserStatusActive
	future := now.Add(10 * time.Minute)
	user.LockedUntil = &future

	if err := user.EnsureCanLogin(now); err != valueobject.ErrAccountLocked {
		t.Errorf("expected ErrAccountLocked")
	}

	// Case: OK
	past := now.Add(-10 * time.Minute)
	user.LockedUntil = &past
	if err := user.EnsureCanLogin(now); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
