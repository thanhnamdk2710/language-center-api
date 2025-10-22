package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/shared/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase"
)

type LoginHandler struct {
	usecase usecase.LoginUsecase
}

func NewLoginHandler(usecase usecase.LoginUsecase) *LoginHandler {
	return &LoginHandler{
		usecase: usecase,
	}
}

func (h LoginHandler) Login(c *gin.Context) {
	response.Success(c, http.StatusCreated, gin.H{
		"message": "Login",
	})
}
