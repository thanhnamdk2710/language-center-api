package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/delivery/http/dto"
	"github.com/thanhnamdk2710/auth-service/internal/shared/password"
	"github.com/thanhnamdk2710/auth-service/internal/shared/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/register"
)

type RegisterHandler struct {
	usecase register.Usecase
}

func NewRegisterHandler(u register.Usecase) *RegisterHandler {
	return &RegisterHandler{
		usecase: u,
	}
}

func (h RegisterHandler) Register(c *gin.Context) {
	req := c.MustGet("requestBody").(dto.RegisterRequest)

	passwordHash, err := password.HashPassword(req.Password)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "PASSWORD_INVALID", err.Error())
		return
	}

	input := register.Input{
		Email:        req.Email,
		PasswordHash: passwordHash,
	}

	err = h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}
