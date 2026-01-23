package auth

import (
	"fmt"
	"strconv"

	"github.com/zalando/go-keyring"
)

// KeyringTokenStorage implements TokenStorage using system keyring
type KeyringTokenStorage struct {
	service string
}

// NewKeyringTokenStorage creates a new KeyringTokenStorage
func NewKeyringTokenStorage(service string) *KeyringTokenStorage {
	return &KeyringTokenStorage{
		service: service,
	}
}

// GetToken retrieves a token from keyring
func (k *KeyringTokenStorage) GetToken(scope TokenScope) (*Token, error) {
	var tokenKey, expiryKey string

	switch scope {
	case ScopeApplication:
		tokenKey = "access_token"
		expiryKey = "token_expiry"
	case ScopeUser:
		tokenKey = "user_access_token"
		expiryKey = "user_token_expiry"
	default:
		return nil, fmt.Errorf("invalid token scope")
	}

	accessToken, err := keyring.Get(k.service, tokenKey)
	if err != nil {
		return nil, err
	}

	expiryStr, err := keyring.Get(k.service, expiryKey)
	if err != nil {
		return nil, err
	}

	expiresAt, err := strconv.ParseInt(expiryStr, 10, 64)
	if err != nil {
		return nil, err
	}

	token := &Token{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}

	// Get refresh token if it exists (user scope only)
	if scope == ScopeUser {
		refreshToken, err := keyring.Get(k.service, "user_refresh_token")
		if err == nil {
			token.RefreshToken = refreshToken
		}
	}

	return token, nil
}

// StoreToken stores a token in keyring
func (k *KeyringTokenStorage) StoreToken(scope TokenScope, token *Token) error {
	var tokenKey, expiryKey string

	switch scope {
	case ScopeApplication:
		tokenKey = "access_token"
		expiryKey = "token_expiry"
	case ScopeUser:
		tokenKey = "user_access_token"
		expiryKey = "user_token_expiry"
	default:
		return fmt.Errorf("invalid token scope")
	}

	if err := keyring.Set(k.service, tokenKey, token.AccessToken); err != nil {
		return fmt.Errorf("failed to store access token: %w", err)
	}

	if err := keyring.Set(k.service, expiryKey, strconv.FormatInt(token.ExpiresAt, 10)); err != nil {
		return fmt.Errorf("failed to store token expiry: %w", err)
	}

	// Store refresh token if present (user scope only)
	if scope == ScopeUser && token.RefreshToken != "" {
		if err := keyring.Set(k.service, "user_refresh_token", token.RefreshToken); err != nil {
			return fmt.Errorf("failed to store refresh token: %w", err)
		}
	}

	return nil
}

// GetCredentials retrieves client credentials from keyring
func (k *KeyringTokenStorage) GetCredentials() (string, string, error) {
	clientID, err := keyring.Get(k.service, "client_uid")
	if err != nil {
		return "", "", fmt.Errorf("client ID not found: %w", err)
	}

	clientSecret, err := keyring.Get(k.service, "client_secret")
	if err != nil {
		return "", "", fmt.Errorf("client secret not found: %w", err)
	}

	return clientID, clientSecret, nil
}

// GetRefreshToken retrieves the refresh token from keyring
func (k *KeyringTokenStorage) GetRefreshToken(scope TokenScope) (string, error) {
	if scope != ScopeUser {
		return "", fmt.Errorf("refresh tokens only available for user scope")
	}

	refreshToken, err := keyring.Get(k.service, "user_refresh_token")
	if err != nil {
		return "", fmt.Errorf("no refresh token available: %w", err)
	}

	return refreshToken, nil
}
