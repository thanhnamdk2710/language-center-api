package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thanhnamdk2710/auth-service/internal/domain"
	"github.com/thanhnamdk2710/auth-service/internal/dto"
	"github.com/thanhnamdk2710/auth-service/internal/shared/password"
	response "github.com/thanhnamdk2710/auth-service/internal/shared/responses"
	"github.com/thanhnamdk2710/auth-service/internal/shared/validation"
	usecase "github.com/thanhnamdk2710/auth-service/internal/usecase/register"
)

type RegisterHandler struct {
	usecase usecase.RegisterUsecase
}

func NewAuthHandler(usecase usecase.RegisterUsecase) *RegisterHandler {
	return &RegisterHandler{
		usecase: usecase,
	}
}

func (h RegisterHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			errors := validation.ConvertValidationErrors(verrs)
			response.Error(c, http.StatusBadRequest, "INPUT_INVALID", errors)
			return
		}

		response.Error(c, http.StatusBadRequest, "INPUT_INVALID", "Invalid request body format")
		return
	}

	passwordHash, err := password.HashPassword(req.Password)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "PASSWORD_INVALID", err.Error())
		return
	}

	user := &domain.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
	}

	err = h.usecase.Register(c.Request.Context(), user)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}
