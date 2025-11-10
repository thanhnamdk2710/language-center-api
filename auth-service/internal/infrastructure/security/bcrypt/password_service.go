package bcrypt

import (
	"unicode"

	"github.com/thanhnamdk2710/auth-service/internal/domain/port"
	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
	"golang.org/x/crypto/bcrypt"
)

type passwordService struct {
	cost int
}

func NewPasswordService() port.PasswordService {
	return &passwordService{
		cost: bcrypt.DefaultCost,
	}
}

func (s *passwordService) Validate(pwd string) error {
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range pwd {
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
		return valueobject.ErrInvalidPassword
	}

	return nil
}

func (s *passwordService) Hash(pwd string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pwd), s.cost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func (s *passwordService) Verify(pwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}
