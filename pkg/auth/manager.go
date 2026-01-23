package auth

import (
	"fmt"
	"time"

	"github.com/shiftwavedev/42-cli/internal/api"
)

const (
	// TokenExpiryBuffer is the safety buffer before token expiry (5 minutes)
	TokenExpiryBuffer = 300
)

// TokenManager handles token operations
type TokenManager struct {
	storage TokenStorage
	client  *api.Client
}

// NewTokenManager creates a new TokenManager
func NewTokenManager(storage TokenStorage, client *api.Client) *TokenManager {
	return &TokenManager{
		storage: storage,
		client:  client,
	}
}

// GetValidToken retrieves a valid token for the given scope
// Automatically refreshes if expired
func (tm *TokenManager) GetValidToken(scope TokenScope) (string, error) {
	token, err := tm.storage.GetToken(scope)
	if err == nil && tm.isTokenValid(token) {
		return token.AccessToken, nil
	}

	// Token expired or not found, refresh/generate new one
	return tm.refreshToken(scope)
}

// isTokenValid checks if a token is still valid
func (tm *TokenManager) isTokenValid(token *Token) bool {
	return token != nil && time.Now().Unix() < token.ExpiresAt
}

// refreshToken refreshes or generates a new token based on scope
func (tm *TokenManager) refreshToken(scope TokenScope) (string, error) {
	switch scope {
	case ScopeApplication:
		return tm.generateApplicationToken()
	case ScopeUser:
		return tm.refreshUserToken()
	default:
		return "", fmt.Errorf("invalid token scope: %v", scope)
	}
}

// generateApplicationToken generates a new application token using client credentials
func (tm *TokenManager) generateApplicationToken() (string, error) {
	clientID, clientSecret, err := tm.storage.GetCredentials()
	if err != nil {
		return "", err
	}

	tokenResp, err := tm.client.GetClientCredentialsToken(clientID, clientSecret)
	if err != nil {
		return "", err
	}

	token := &Token{
		AccessToken: tokenResp.AccessToken,
		ExpiresAt:   time.Now().Unix() + int64(tokenResp.ExpiresIn) - TokenExpiryBuffer,
	}

	if err := tm.storage.StoreToken(ScopeApplication, token); err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

// refreshUserToken refreshes the user token using the refresh token
func (tm *TokenManager) refreshUserToken() (string, error) {
	refreshToken, err := tm.storage.GetRefreshToken(ScopeUser)
	if err != nil {
		return "", fmt.Errorf("no refresh token available, please login again")
	}

	clientID, clientSecret, err := tm.storage.GetCredentials()
	if err != nil {
		return "", err
	}

	tokenResp, err := tm.client.RefreshToken(refreshToken, clientID, clientSecret)
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}

	token := &Token{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Unix() + int64(tokenResp.ExpiresIn) - TokenExpiryBuffer,
	}

	if err := tm.storage.StoreToken(ScopeUser, token); err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

// StoreOAuthToken stores an OAuth token response
func (tm *TokenManager) StoreOAuthToken(tokenResp *api.TokenResponse) error {
	token := &Token{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Unix() + int64(tokenResp.ExpiresIn) - TokenExpiryBuffer,
	}

	return tm.storage.StoreToken(ScopeUser, token)
}
