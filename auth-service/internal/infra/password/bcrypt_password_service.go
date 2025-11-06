package password

import (
	"unicode"

	domainPassword "github.com/thanhnamdk2710/auth-service/internal/domain/service/password"
	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
	"golang.org/x/crypto/bcrypt"
)

type bcryptPasswordService struct {
	cost int
}

func NewBcryptPasswordService() domainPassword.Service {
	return &bcryptPasswordService{
		cost: bcrypt.DefaultCost,
	}
}

func (s *bcryptPasswordService) Validate(password string) error {
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
		return valueobject.ErrInvalidPassword
	}

	return nil
}

func (s *bcryptPasswordService) Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func (s *bcryptPasswordService) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
