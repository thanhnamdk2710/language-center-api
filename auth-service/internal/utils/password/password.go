package password

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct{}

func (s *PasswordService) HashPassword(password string) (string, error) {
	// Validate password
	if err := s.Validate(password); err != nil {
		return "", err
	}

	// Generate hash
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (s *PasswordService) Validate(password string) error {
	// Check for required character types
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("Password must contain uppercase, lowercase, number, and special character")
	}

	return nil
}
