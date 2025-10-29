package service

import (
	"github.com/thanhnamdk2710/auth-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type bcryptPasswordService struct {
	cost int
}

func NewBcryptPasswordService() domain.PasswordService {
	return &bcryptPasswordService{
		cost: bcrypt.DefaultCost,
	}
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
