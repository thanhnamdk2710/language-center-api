package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/errors"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/response"
)

type BaseHandler struct{}

func (h *BaseHandler) HandleError(c *gin.Context, err error) {
	httpErr := errors.MapError(err)

	if httpErr.Code == "VALIDATION_ERROR" && httpErr.Fields != nil {
		response.ValidationError(c, httpErr.Fields)
		return
	}
	response.Error(c, httpErr.Status, httpErr.Code, httpErr.Message)
}
