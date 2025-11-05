package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/shared/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/refresh_password"
)

type RefreshTokenHandler struct {
	usecase refresh_password.Usecase
}

func NewRefreshTokenHandler(usecase refresh_password.Usecase) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		usecase: usecase,
	}
}

func (h *RefreshTokenHandler) RefreshToken(c *gin.Context) {

	response.Success(c, http.StatusCreated, gin.H{
		"message": "RefreshToken",
	})
}
