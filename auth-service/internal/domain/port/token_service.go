package port

type TokenService interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	// ValidateToken(token string) (*TokenClaims, error)
}
