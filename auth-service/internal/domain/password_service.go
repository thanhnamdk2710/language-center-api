package domain

type PasswordService interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}
