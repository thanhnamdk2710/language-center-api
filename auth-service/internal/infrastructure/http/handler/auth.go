package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/request"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/login"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/register"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/verify_otp"
)

type AuthHandler struct {
	BaseHandler
	registerUC  register.Usecase
	verifyOTPUC verify_otp.Usecase
	loginUC     login.Usecase
}

func NewAuthHandler(registerUC register.Usecase, verifyOTPUC verify_otp.Usecase, loginUC login.Usecase) *AuthHandler {
	return &AuthHandler{
		registerUC:  registerUC,
		verifyOTPUC: verifyOTPUC,
		loginUC:     loginUC,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.HandleError(c, err)
		return
	}

	// Map to use case input
	input := register.Input{
		Email:    req.Email,
		Password: req.Password,
	}

	// Execute use case
	output, err := h.registerUC.Execute(c.Request.Context(), input)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, &response.RegisterResponse{
		UserID:  output.UserID,
		Message: output.Message,
	})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req request.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.HandleError(c, err)
		return
	}

	// Map to use case input
	input := verify_otp.Input{
		Email: req.Email,
		OTP:   req.OTP,
	}

	// Execute use case
	output, err := h.verifyOTPUC.Execute(c.Request.Context(), input)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, &response.AuthResponse{
		UserID:       output.UserID,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		Message:      "Email verified successfully",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.HandleError(c, err)
		return
	}

	// Map to use case input
	input := login.Input{
		Email:    req.Email,
		Password: req.Password,
	}

	// Execute use case
	output, err := h.loginUC.Execute(c.Request.Context(), input)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, &response.AuthResponse{
		UserID:       output.UserID,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		Message:      "Login successfully",
	})
}
