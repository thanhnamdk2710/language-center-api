package container

import (
	"database/sql"

	"github.com/thanhnamdk2710/auth-service/internal/config"
	"github.com/thanhnamdk2710/auth-service/internal/delivery/http/handlers"
	"github.com/thanhnamdk2710/auth-service/internal/domain"
	"github.com/thanhnamdk2710/auth-service/internal/infra/repository/postgres"
	"github.com/thanhnamdk2710/auth-service/internal/infra/service"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/forgot_password"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/login"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/refresh_password"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/register"
)

type Container struct {
	// Infrastructure - Repositories
	UserRepo domain.UserRepository

	// Infrastructure - Services
	PasswordService domain.PasswordService
	TokenService    domain.TokenService
	EmailService    domain.EmailService

	// Usecase
	RegisterUsecase       register.Usecase
	LoginUsecase          login.Usecase
	ForgotPasswordUsecase forgot_password.Usecase
	RefreshTokenUsecase   refresh_password.Usecase

	// Usecase
	RegisterHandler       *handlers.RegisterHandler
	LoginHandler          *handlers.LoginHandler
	ForgotPasswordHandler *handlers.ForgotPasswordHandler
	RefreshTokenHandler   *handlers.RefreshTokenHandler
}

func NewContainer(db *sql.DB, cfg *config.Config) *Container {
	// Infrastructure
	userRepo := postgres.NewUserRepository(db)
	passwordSvc := service.NewBcryptPasswordService()
	emailSvc := service.NewSMTPEmailService(cfg.SMTP)

	// Usecases
	registerUC := register.NewRegisterUsecase(userRepo, passwordSvc, emailSvc)
	loginUC := login.NewLoginUsecase()
	forgotPasswordUC := forgot_password.NewForgotPasswordUsecase()
	refreshTokenUC := refresh_password.NewRefreshTokenUsecase()

	// Handlers
	registerHandler := handlers.NewRegisterHandler(registerUC)
	loginHandler := handlers.NewLoginHandler(loginUC)
	forgotPasswordHandler := handlers.NewForgotPasswordHandler(forgotPasswordUC)
	refreshTokenHandler := handlers.NewRefreshTokenHandler(refreshTokenUC)

	return &Container{
		RegisterHandler:       registerHandler,
		LoginHandler:          loginHandler,
		ForgotPasswordHandler: forgotPasswordHandler,
		RefreshTokenHandler:   refreshTokenHandler,
	}
}
