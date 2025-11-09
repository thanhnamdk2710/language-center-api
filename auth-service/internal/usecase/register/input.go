package register

import (
	"errors"
	"fmt"

	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
)

type Input struct {
	Email    string
	Password string
}

func (i Input) Validate() error {
	if i.Email == "" {
		return errors.New("email is required")
	}

	if _, err := valueobject.NewEmail(i.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	if i.Password == "" {
		return errors.New("password is required")
	}

	if len(i.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return nil
}
