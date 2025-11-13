package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/presenter"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/response"
)

type BaseHandler struct {
	errorPresenter *presenter.ErrorPresenter
}

func NewBaseHandler() BaseHandler {
	return BaseHandler{
		errorPresenter: presenter.NewErrorPresenter(),
	}
}

// HandleError converts domain errors to HTTP responses
func (h *BaseHandler) HandleError(c *gin.Context, err error) {
	httpErr := h.errorPresenter.PresentError(err)

	if httpErr.Code == "VALIDATION_ERROR" && httpErr.Fields != nil {
		response.ValidationError(c, httpErr.Fields)
		return
	}
	response.Error(c, httpErr.Status, httpErr.Code, httpErr.Message)
}
