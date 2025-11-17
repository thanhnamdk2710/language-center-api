package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

type Email string

var (
	ErrInvalidEmailFormat = errors.New("invalid email format")
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9_%+\-]*[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9_%+\-]*[a-zA-Z0-9])?)*@` +
	`([a-zA-Z0-9]([a-zA-Z0-9\-]*[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

func NewEmail(raw string) (Email, error) {
	normalized := strings.TrimSpace(strings.ToLower(raw))

	if !emailRegex.MatchString(normalized) {
		return "", ErrInvalidEmailFormat
	}
	return Email(normalized), nil
}

func (e Email) String() string {
	return string(e)
}
