package email

type Service interface {
	SendVerificationEmail(email, otp string) error
	SendPasswordResetEmail(email, token string) error
}
