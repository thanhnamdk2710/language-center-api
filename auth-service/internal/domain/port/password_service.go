package port

type PasswordService interface {
	Validate(password string) error
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}
