package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/adapter/http/dto"
	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
	"github.com/thanhnamdk2710/auth-service/internal/shared/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/login"
)

type LoginHandler struct {
	usecase login.Usecase
}

func NewLoginHandler(usecase login.Usecase) *LoginHandler {
	return &LoginHandler{
		usecase: usecase,
	}
}

func (h *LoginHandler) Login(c *gin.Context) {
	req := c.MustGet("requestBody").(dto.LoginRequest)

	input := login.Input{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, valueobject.ErrInvalidCredentials):
			response.Error(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
			return
		case errors.Is(err, valueobject.ErrEmailNotVerified):
			response.Error(c, http.StatusUnauthorized, "EMAIL_NOT_VERIFIED", "Please verify your email before logging in")
			return
		case errors.Is(err, valueobject.ErrAccountDisabled):
			response.Error(c, http.StatusUnauthorized, "ACCOUNT_DISABLED", "Your account has been disabled")
			return
		case errors.Is(err, valueobject.ErrAccountLocked):
			response.Error(c, http.StatusUnauthorized, "ACCOUNT_LOCKED", "Account locked due to too many failed attempts")
			return
		default:
			response.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "An error occurred during login")
			return
		}
	}

	response.Success(c, http.StatusOK, gin.H{
		"user_id":       output.UserID,
		"access_token":  output.AccessToken,
		"refresh_token": output.RefreshToken,
		"message":       "Login successfully",
	})
}
