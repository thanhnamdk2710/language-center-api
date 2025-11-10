package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase/refresh_token"
)

type RefreshTokenHandler struct {
	usecase refresh_token.Usecase
}

func NewRefreshTokenHandler(usecase refresh_token.Usecase) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		usecase: usecase,
	}
}

func (h *RefreshTokenHandler) RefreshToken(c *gin.Context) {

	response.Success(c, http.StatusCreated, gin.H{
		"message": "RefreshToken",
	})
}
