package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	response "github.com/thanhnamdk2710/auth-service/internal/delivery/http/responses"
)

func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Success(c, http.StatusCreated, gin.H{
			"message": "User registered successfully",
		})
	}
}
