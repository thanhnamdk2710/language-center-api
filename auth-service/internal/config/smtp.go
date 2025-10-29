package config

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	Enabled  bool
}

func loadSMTPConfig() SMTPConfig {
	return SMTPConfig{
		Host:     getString("SMTP_HOST", "smtp.gmail.com"),
		Port:     getString("SMTP_PORT", "587"),
		Username: getString("SMTP_USERNAME", ""),
		Password: getString("SMTP_PASSWORD", ""),
		From:     getString("SMTP_FROM", "noreply@languagecenter.com"),
		Enabled:  getBool("SMTP_ENABLED", false),
	}
}
