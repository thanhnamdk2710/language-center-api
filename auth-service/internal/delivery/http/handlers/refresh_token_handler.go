package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/shared/response"
	"github.com/thanhnamdk2710/auth-service/internal/usecase"
)

type RefreshTokenHandler struct {
	usecase usecase.RefreshTokenUsecase
}

func NewRefreshTokenHandler(usecase usecase.RefreshTokenUsecase) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		usecase: usecase,
	}
}

func (h RefreshTokenHandler) RefreshToken(c *gin.Context) {

	response.Success(c, http.StatusCreated, gin.H{
		"message": "RefreshToken",
	})
}
