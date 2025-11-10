package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/forgot_password"
)

type ForgotPasswordHandler struct {
	usecase forgot_password.Usecase
}

func NewForgotPasswordHandler(u forgot_password.Usecase) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{
		usecase: u,
	}
}

func (h *ForgotPasswordHandler) ForgotPassword(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{
		"message": "ForgotPassword",
	})
}
