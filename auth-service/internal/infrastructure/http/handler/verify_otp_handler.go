package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/dto"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/verify_otp"
)

type VerifyOTPHandler struct {
	BaseHandler
	usecase verify_otp.Usecase
}

func NewVerifyOTPHandler(u verify_otp.Usecase) *VerifyOTPHandler {
	return &VerifyOTPHandler{
		usecase: u,
	}
}

func (h *VerifyOTPHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.HandleError(c, err)
		return
	}

	input := verify_otp.Input{
		Email: req.Email,
		OTP:   req.OTP,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"user_id":       output.UserID,
		"access_token":  output.AccessToken,
		"refresh_token": output.RefreshToken,
		"message":       "Email verified successfully",
	})
}
