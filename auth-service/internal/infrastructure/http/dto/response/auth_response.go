package response

// AuthResponse represents authentication response
type AuthResponse struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Message      string `json:"message"`
}

// RegisterResponse represents registration response
type RegisterResponse struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

// RefreshTokenResponse represents token refresh response
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
