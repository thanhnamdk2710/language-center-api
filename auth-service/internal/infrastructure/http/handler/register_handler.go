package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/dto/request"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/presenter"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/register"
)

type RegisterHandler struct {
	BaseHandler
	usecase   register.Usecase
	presenter *presenter.AuthPresenter
}

func NewRegisterHandler(u register.Usecase) *RegisterHandler {
	return &RegisterHandler{
		BaseHandler: NewBaseHandler(),
		usecase:     u,
		presenter:   presenter.NewAuthPresenter(),
	}
}

func (h *RegisterHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.HandleError(c, err)
		return
	}

	// Map to use case input
	input := register.Input{
		Email:    req.Email,
		Password: req.Password,
	}

	// Execute use case
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Present response
	resp := h.presenter.PresentRegister(output)
	response.Success(c, http.StatusCreated, resp)
}
