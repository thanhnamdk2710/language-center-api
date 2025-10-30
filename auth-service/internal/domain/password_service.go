package domain

import "errors"

var (
	ErrInvalidPassword = errors.New("password must contain uppercase, lowercase, number, and special character")
)

type PasswordService interface {
	Validate(password string) error
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}
