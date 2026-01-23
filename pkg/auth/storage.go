package auth

// TokenStorage defines the interface for token storage
type TokenStorage interface {
	GetToken(scope TokenScope) (*Token, error)
	StoreToken(scope TokenScope, token *Token) error
	GetCredentials() (clientID, clientSecret string, err error)
	GetRefreshToken(scope TokenScope) (string, error)
}
