package valueobject_test

import (
	"errors"
	"testing"

	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
)

func TestNewEmail_ValidEmails(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"john@example.com", "john@example.com"},
		{"John.Doe@Example.COM", "john.doe@example.com"}, // lowercase normalization
		{"   user123@gmail.com   ", "user123@gmail.com"}, // trimming
		{"a+b@test.co", "a+b@test.co"},
		{"hello.world@test-domain.com", "hello.world@test-domain.com"},
	}

	for _, tt := range tests {
		got, err := valueobject.NewEmail(tt.input)
		if err != nil {
			t.Errorf("unexpected error for input %q: %v", tt.input, err)
			continue
		}

		if got.String() != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, got.String())
		}
	}
}

func TestNewEmail_InvalidEmails(t *testing.T) {
	tests := []string{
		"",
		"  ",
		"plainaddress",
		"@missinglocal.com",
		"missingdomain@",
		"user@.com",
		"user@domain",
		"user@domain..com",
		"john@exa mple.com",
		"user@@domain.com",
		"user@domain.c",         // TLD too short
		"user@-domain.com",      // invalid domain start
		"john..doe@example.com", // double dots in local part
	}

	for _, input := range tests {
		_, err := valueobject.NewEmail(input)
		if err == nil {
			t.Errorf("expected error for invalid email %q, got nil", input)
		}

		if !errors.Is(err, valueobject.ErrInvalidEmailFormat) {
			t.Errorf("expected ErrInvalidEmailFormat, got %v", err)
		}
	}
}

func TestEmail_String(t *testing.T) {
	email, err := valueobject.NewEmail("John@Example.COM")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "john@example.com"
	if email.String() != expected {
		t.Errorf("expected %q, got %q", expected, email.String())
	}
}
