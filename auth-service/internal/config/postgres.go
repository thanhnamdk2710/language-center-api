package config

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func loadPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     getString("POSTGRES_HOST", "localhost"),
		Port:     getString("POSTGRES_PORT", "5432"),
		User:     getString("POSTGRES_USER", "postgres"),
		Password: getString("POSTGRES_PASSWORD", "secret"),
		DBName:   getString("POSTGRES_DB", "auth_db"),
		SSLMode:  getString("POSTGRES_SSLMODE", "disable"),
	}
}
