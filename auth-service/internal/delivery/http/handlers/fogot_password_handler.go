package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/shared/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase"
)

type ForgotPasswordHandler struct {
	usecase usecase.ForgotPasswordUsecase
}

func NewForgotPasswordHandler(usecase usecase.ForgotPasswordUsecase) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{
		usecase: usecase,
	}
}

func (h ForgotPasswordHandler) ForgotPassword(c *gin.Context) {
	response.Success(c, http.StatusCreated, gin.H{
		"message": "ForgotPassword",
	})
}
