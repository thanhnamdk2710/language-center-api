package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/dto"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/login"
)

type LoginHandler struct {
	BaseHandler
	usecase login.Usecase
}

func NewLoginHandler(usecase login.Usecase) *LoginHandler {
	return &LoginHandler{
		usecase: usecase,
	}
}

func (h *LoginHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.HandleError(c, err)
		return
	}

	input := login.Input{
		Email:    req.Email,
		Password: req.Password,
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
		"message":       "Login successfully",
	})
}
