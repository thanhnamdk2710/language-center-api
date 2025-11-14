package di

import (
	"database/sql"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/thanhnamdk2710/auth-service/internal/config"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/auth/jwt"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/cache/redis"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/database/postgres"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/handler"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/mailer/smtp"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/security/bcrypt"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/login"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/register"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/verify_otp"
)

type Container struct {
	HealthHandler *handler.HealthHandler
	AuthHandler   *handler.AuthHandler
}

// NewContainer creates and wires all dependencies
func NewContainer(db *sql.DB, redisClient *goredis.Client, cfg *config.Config) *Container {
	// Interface Adapters - Concrete implementations
	userRepo := postgres.NewUserRepository(db)
	passwordSvc := bcrypt.NewPasswordService()
	emailSvc := smtp.NewEmailService(cfg.SMTP)
	otpSvc := redis.NewOTPService(redisClient, 10*time.Minute)
	sessionSvc := redis.NewSessionService(redisClient)
	tokenSvc := jwt.NewTokenService(
		cfg.JWT.SecretKey,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
	)

	// Use Cases - Application business rules
	registerUC := register.NewRegisterUsecase(userRepo, passwordSvc, emailSvc, otpSvc)
	verifyOTPUC := verify_otp.NewVerifyOTPUsecase(userRepo, otpSvc, tokenSvc)
	loginUC := login.NewLoginUsecase(userRepo, passwordSvc, tokenSvc, sessionSvc)

	// HTTP Handlers - Delivery mechanism
	healthHandler := &handler.HealthHandler{}
	authHandler := handler.NewAuthHandler(registerUC, verifyOTPUC, loginUC)

	return &Container{
		HealthHandler: healthHandler,
		AuthHandler:   authHandler,
	}
}
