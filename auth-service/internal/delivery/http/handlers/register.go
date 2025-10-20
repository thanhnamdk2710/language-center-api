package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	response "github.com/thanhnamdk2710/auth-service/internal/delivery/http/responses"
	"github.com/thanhnamdk2710/auth-service/internal/domain"
	usecase "github.com/thanhnamdk2710/auth-service/internal/usecases/register"
	"github.com/thanhnamdk2710/auth-service/internal/utils/password"
)

type RegisterHandler struct {
	usecase usecase.RegisterUsecase
}

func NewAuthHandler(usecase usecase.RegisterUsecase) *RegisterHandler {
	return &RegisterHandler{
		usecase: usecase,
	}
}

type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8,max=64"`
	PasswordConfirm string `json:"password_confim" binding:"required,min=8,max=64,eqfield=Password"`
}

func (r *RegisterRequest) Validate() error {
	ps := password.PasswordService{}
	if err := ps.Validate(r.Password); err != nil {
		return err
	}
	return nil
}

func (h RegisterHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INPUT_INVALID", err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		response.Error(c, http.StatusBadRequest, "EMAIL_INVALID", err.Error())
		return
	}

	user := &domain.User{
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	err := h.usecase.Register(c.Request.Context(), user)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}
