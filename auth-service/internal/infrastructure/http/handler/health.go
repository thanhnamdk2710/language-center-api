package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/response"
)

type HealthHandler struct{}

func (h *HealthHandler) Check(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{
		"message": "OK",
	})
}
