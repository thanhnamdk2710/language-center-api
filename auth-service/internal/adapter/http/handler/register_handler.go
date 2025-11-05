package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/adapter/http/dto"
	"github.com/thanhnamdk2710/auth-service/internal/domain"
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

func (h *RegisterHandler) Register(c *gin.Context) {
	req := c.MustGet("requestBody").(dto.RegisterRequest)

	input := register.Input{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		// Map domain errors to HTTP errors
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			response.Error(c, http.StatusConflict, "USER_EXISTS", err.Error())
			return
		}
		if errors.Is(err, domain.ErrInvalidPassword) {
			response.Error(c, http.StatusConflict, "INVALID_PASSWORD", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, gin.H{
		"user_id": output.UserID,
		"message": "User registered successfully",
	})
}
