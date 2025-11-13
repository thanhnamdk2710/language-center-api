package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/dto/request"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/presenter"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/login"
)

type LoginHandler struct {
	BaseHandler
	usecase   login.Usecase
	presenter *presenter.AuthPresenter
}

func NewLoginHandler(usecase login.Usecase) *LoginHandler {
	return &LoginHandler{
		BaseHandler: NewBaseHandler(),
		usecase:     usecase,
		presenter:   presenter.NewAuthPresenter(),
	}
}

func (h *LoginHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.HandleError(c, err)
		return
	}

	// Map to use case input
	input := login.Input{
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
	resp := h.presenter.PresentLogin(output)
	response.Success(c, http.StatusOK, resp)
}
