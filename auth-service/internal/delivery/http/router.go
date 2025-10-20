package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/delivery/http/handlers"
	usecase "github.com/thanhnamdk2710/auth-service/internal/usecases/register"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			registerUsecase := usecase.NewRegisterUsecase()
			registerHandler := handlers.NewAuthHandler(registerUsecase)
			auth.POST("/register", registerHandler.Register)
		}
	}

	return r
}
