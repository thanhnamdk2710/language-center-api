package service

import (
	"fmt"
	"net/smtp"

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
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	// Compose message with header
	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/plain; charset=UTF-8\r\n"+
			"\r\n"+
			"%s\r\n",
		s.from, to, subject, body,
	))

	// Send email
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	err := smtp.SendMail(addr, auth, s.from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}

	return nil
}
