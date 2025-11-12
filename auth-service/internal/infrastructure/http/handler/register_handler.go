package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/dto"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/register"
)

type RegisterHandler struct {
	BaseHandler
	usecase register.Usecase
}

func NewRegisterHandler(u register.Usecase) *RegisterHandler {
	return &RegisterHandler{
		usecase: u,
	}
}

func (h *RegisterHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.HandleError(c, err)
		return
	}

	input := register.Input{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, gin.H{
		"user_id": output.UserID,
		"message": output.Message,
	})
}
