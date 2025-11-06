package token

type Service interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	// ValidateToken(token string) (*TokenClaims, error)
}
