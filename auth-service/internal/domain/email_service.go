package domain

type EmailService interface {
	SendVerificationEmail(email, otp string) error
	SendPasswordResetEmail(email, token string) error
}
