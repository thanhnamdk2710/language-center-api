package modules

import (
	"database/sql"

	"github.com/thanhnamdk2710/auth-service/internal/delivery/http/handlers"
	"github.com/thanhnamdk2710/auth-service/internal/infra/repository/postgres"
	"github.com/thanhnamdk2710/auth-service/internal/usecase"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/register"
)

type AuthModule struct {
	RegisterHandler       *handlers.RegisterHandler
	LoginHandler          *handlers.LoginHandler
	ForgotPasswordHandler *handlers.ForgotPasswordHandler
	RefreshTokenHandler   *handlers.RefreshTokenHandler
}

func NewAuthModule(db *sql.DB) *AuthModule {
	userRepo := postgres.NewUserRepository(db)

	registerUC := register.NewRegisterUsecase(userRepo)
	loginUC := usecase.NewLoginUsecase()
	forgotPasswordUC := usecase.NewForgotPasswordUsecase()
	refreshTokenUC := usecase.NewRefreshTokenUsecase()

	return &AuthModule{
		RegisterHandler:       handlers.NewRegisterHandler(registerUC),
		LoginHandler:          handlers.NewLoginHandler(loginUC),
		ForgotPasswordHandler: handlers.NewForgotPasswordHandler(forgotPasswordUC),
		RefreshTokenHandler:   handlers.NewRefreshTokenHandler(refreshTokenUC),
	}
}
