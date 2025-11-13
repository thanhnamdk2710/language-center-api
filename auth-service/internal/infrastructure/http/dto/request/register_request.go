package request

type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8,max=64"`
	PasswordConfirm string `json:"password_confirm" binding:"eqfield=Password"`
}
