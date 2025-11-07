package container

import (
	"database/sql"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/thanhnamdk2710/auth-service/internal/adapter/auth/jwt"
	"github.com/thanhnamdk2710/auth-service/internal/adapter/cache/redis"
	"github.com/thanhnamdk2710/auth-service/internal/adapter/http/handler"
	"github.com/thanhnamdk2710/auth-service/internal/adapter/mailer/smtp"
	"github.com/thanhnamdk2710/auth-service/internal/adapter/persistence/postgres"
	"github.com/thanhnamdk2710/auth-service/internal/adapter/security/bcrypt"
	"github.com/thanhnamdk2710/auth-service/internal/config"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/forgot_password"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/login"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/refresh_token"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/register"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/verify_otp"
)

type Container struct {
	RegisterHandler       *handler.RegisterHandler
	VerifyOTPHandler      *handler.VerifyOTPHandler
	LoginHandler          *handler.LoginHandler
	ForgotPasswordHandler *handler.ForgotPasswordHandler
	RefreshTokenHandler   *handler.RefreshTokenHandler
}

func NewContainer(db *sql.DB, redisClient *goredis.Client, cfg *config.Config) *Container {
	// Infrastructure - Adapters
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

	// Usecases
	registerUC := register.NewRegisterUsecase(userRepo, passwordSvc, emailSvc, otpSvc)
	verifyOTPUC := verify_otp.NewVerifyOTPUsecase(userRepo, otpSvc, tokenSvc)
	loginUC := login.NewLoginUsecase(userRepo, passwordSvc, tokenSvc, sessionSvc)
	forgotPasswordUC := forgot_password.NewForgotPasswordUsecase()
	refreshTokenUC := refresh_token.NewRefreshTokenUsecase()

	// Handlers
	registerHandler := handler.NewRegisterHandler(registerUC)
	verifyOTPHandler := handler.NewVerifyOTPHandler(verifyOTPUC)
	loginHandler := handler.NewLoginHandler(loginUC)
	forgotPasswordHandler := handler.NewForgotPasswordHandler(forgotPasswordUC)
	refreshTokenHandler := handler.NewRefreshTokenHandler(refreshTokenUC)

	return &Container{
		RegisterHandler:       registerHandler,
		VerifyOTPHandler:      verifyOTPHandler,
		LoginHandler:          loginHandler,
		ForgotPasswordHandler: forgotPasswordHandler,
		RefreshTokenHandler:   refreshTokenHandler,
	}
}
