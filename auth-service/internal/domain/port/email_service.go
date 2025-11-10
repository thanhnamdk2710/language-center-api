package port

type EmailService interface {
	SendVerificationEmail(email, otp string) error
	SendPasswordResetEmail(email, token string) error
}
