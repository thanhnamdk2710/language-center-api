package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/dto/request"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/presenter"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/verify_otp"
)

type VerifyOTPHandler struct {
	BaseHandler
	usecase   verify_otp.Usecase
	presenter *presenter.AuthPresenter
}

func NewVerifyOTPHandler(u verify_otp.Usecase) *VerifyOTPHandler {
	return &VerifyOTPHandler{
		BaseHandler: NewBaseHandler(),
		usecase:     u,
		presenter:   presenter.NewAuthPresenter(),
	}
}

func (h *VerifyOTPHandler) VerifyOTP(c *gin.Context) {
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
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Present response
	resp := h.presenter.PresentVerifyOTP(output)
	response.Success(c, http.StatusOK, resp)
}
