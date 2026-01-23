package cmd

import (
	"github.com/shiftwavedev/42-cli/internal/api"
	"github.com/shiftwavedev/42-cli/pkg/auth"
)

var tokenManager *auth.TokenManager

func init() {
	storage := auth.NewKeyringTokenStorage(service)
	tokenManager = auth.NewTokenManager(storage, api.DefaultClient)
}

// GetAccessToken returns the user access token (for backward compatibility)
// TODO: Remove in future versions
func GetAccessToken() (string, error) {
	return GetUserToken()
}

// GetApplicationToken returns a valid application token
// TODO: Remove in future versions
func GetApplicationToken() (string, error) {
	return tokenManager.GetValidToken(auth.ScopeApplication)
}

// GetUserToken returns a valid user token
// TODO: Remove in future versions
func GetUserToken() (string, error) {
	return tokenManager.GetValidToken(auth.ScopeUser)
}

// StoreOAuthToken stores an OAuth token response
// TODO: Remove in future versions
func StoreOAuthToken(tokenResp *api.TokenResponse) error {
	return tokenManager.StoreOAuthToken(tokenResp)
}
