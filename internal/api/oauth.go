package api

import (
	"net/url"
)

// TokenResponse represents the OAuth token response from 42 API
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
	CreatedAt    int64  `json:"created_at"`
}

// ExchangeAuthorizationCode exchanges an authorization code for an access token
func (c *Client) ExchangeAuthorizationCode(code, clientID, clientSecret, redirectURI string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	var tokenResp TokenResponse
	err := c.PostForm("/oauth/token", data, &tokenResp)
	if err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// RefreshToken refreshes an access token using a refresh token
func (c *Client) RefreshToken(refreshToken, clientID, clientSecret string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	var tokenResp TokenResponse
	err := c.PostForm("/oauth/token", data, &tokenResp)
	if err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// GetClientCredentialsToken gets an access token using client credentials
func (c *Client) GetClientCredentialsToken(clientID, clientSecret string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("scope", "public profile")

	var tokenResp TokenResponse
	err := c.PostForm("/oauth/token", data, &tokenResp)
	if err != nil {
		return nil, err
	}

	return &tokenResp, nil
}
