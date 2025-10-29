package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	response.Success(c, http.StatusCreated, gin.H{
		"message": "Login",
	})
}
