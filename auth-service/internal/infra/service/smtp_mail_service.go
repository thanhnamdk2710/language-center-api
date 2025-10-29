package service

import (
	"fmt"

	"github.com/thanhnamdk2710/auth-service/internal/config"
	"github.com/thanhnamdk2710/auth-service/internal/domain"
)

type smtpEmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSMTPEmailService(cfg config.SMTPConfig) domain.EmailService {
	return &smtpEmailService{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
		from:     cfg.From,
	}
}

func (s *smtpEmailService) SendVerificationEmail(email, otp string) error {
	subject := "Email Verification"
	body := fmt.Sprintf("Your verification code is: %s\n\nThis code will expire in 10 minutes.", otp)
	return s.sendMail(email, subject, body)
}

func (s *smtpEmailService) SendPasswordResetEmail(email, otp string) error {
	subject := "Email Verification"
	body := fmt.Sprintf("Your verification code is: %s\n\nThis code will expire in 10 minutes.", otp)
	return s.sendMail(email, subject, body)
}

func (s *smtpEmailService) sendMail(to, subject, body string) error {
	// Setup authentication
	return nil
}
