package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/delivery/http/dto"
	"github.com/thanhnamdk2710/auth-service/internal/domain"
	"github.com/thanhnamdk2710/auth-service/internal/shared/password"
	"github.com/thanhnamdk2710/auth-service/internal/shared/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase"
)

type RegisterHandler struct {
	usecase usecase.RegisterUsecase
}

func NewRegisterHandler(usecase usecase.RegisterUsecase) *RegisterHandler {
	return &RegisterHandler{
		usecase: usecase,
	}
}

func (h RegisterHandler) Register(c *gin.Context) {
	req := c.MustGet("requestBody").(dto.RegisterRequest)

	passwordHash, err := password.HashPassword(req.Password)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "PASSWORD_INVALID", err.Error())
		return
	}

	user := &domain.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
	}

	err = h.usecase.Register(c.Request.Context(), user)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}
