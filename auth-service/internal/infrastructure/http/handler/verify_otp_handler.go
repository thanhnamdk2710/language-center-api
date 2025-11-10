package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/dto"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/verify_otp"
)

type VerifyOTPHandler struct {
	usecase verify_otp.Usecase
}

func NewVerifyOTPHandler(u verify_otp.Usecase) *VerifyOTPHandler {
	return &VerifyOTPHandler{
		usecase: u,
	}
}

func (h *VerifyOTPHandler) VerifyOTP(c *gin.Context) {
	req := c.MustGet("requestBody").(dto.VerifyOTPRequest)

	input := verify_otp.Input{
		Email: req.Email,
		OTP:   req.OTP,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		// Map domain errors to HTTP errors
		switch {
		case errors.Is(err, valueobject.ErrOTPExpired):
			response.Error(c, http.StatusBadRequest, "OTP_EXPIRED", "OTP has expired, please request a new one")
			return
		case errors.Is(err, valueobject.ErrOTPInvalid):
			response.Error(c, http.StatusBadRequest, "OTP_INVALID", "Invalid OTP code")
			return
		case errors.Is(err, valueobject.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
			return
		default:
			response.Error(c, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
			return
		}
	}

	response.Success(c, http.StatusOK, gin.H{
		"user_id":       output.UserID,
		"access_token":  output.AccessToken,
		"refresh_token": output.RefreshToken,
		"message":       "Email verified successfully",
	})
}
