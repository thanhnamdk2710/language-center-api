package presenter

import (
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/dto/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/login"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/register"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/verify_otp"
)

// AuthPresenter converts use case outputs to HTTP responses
type AuthPresenter struct{}

func NewAuthPresenter() *AuthPresenter {
	return &AuthPresenter{}
}

// PresentLogin converts login output to HTTP response
func (p *AuthPresenter) PresentLogin(output *login.Output) *response.AuthResponse {
	return &response.AuthResponse{
		UserID:       output.UserID,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		Message:      "Login successfully",
	}
}

// PresentRegister converts register output to HTTP response
func (p *AuthPresenter) PresentRegister(output *register.Output) *response.RegisterResponse {
	return &response.RegisterResponse{
		UserID:  output.UserID,
		Message: output.Message,
	}
}

// PresentVerifyOTP converts verify OTP output to HTTP response
func (p *AuthPresenter) PresentVerifyOTP(output *verify_otp.Output) *response.AuthResponse {
	return &response.AuthResponse{
		UserID:       output.UserID,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		Message:      "Email verified successfully",
	}
}
