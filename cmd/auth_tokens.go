package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zalando/go-keyring"

	"github.com/shiftwavedev/42-cli/internal/api"
)

func GetAccessToken() (string, error) {
	return GetUserToken()
}

func GetApplicationToken() (string, error) {
	token, err := getValidToken()
	if err == nil && token != "" {
		return token, nil
	}

	return generateNewToken()
}

func GetUserToken() (string, error) {
	return getUserAccessToken()
}

func getValidToken() (string, error) {
	token, err := keyring.Get(service, "access_token")
	if err != nil {
		return "", err
	}

	expiryStr, err := keyring.Get(service, "token_expiry")
	if err != nil {
		return "", err
	}

	expiry, err := strconv.ParseInt(expiryStr, 10, 64)
	if err != nil {
		return "", err
	}

	if time.Now().Unix() >= expiry {
		return "", fmt.Errorf("token expired")
	}

	return token, nil
}

func generateNewToken() (string, error) {
	clientID, err := keyring.Get(service, "client_uid")
	if err != nil {
		return "", fmt.Errorf("client ID not found in keyring")
	}

	clientSecret, err := keyring.Get(service, "client_secret")
	if err != nil {
		return "", fmt.Errorf("client secret not found in keyring")
	}

	tokenResp, err := api.DefaultClient.GetClientCredentialsToken(clientID, clientSecret)
	if err != nil {
		return "", err
	}

	expiry := time.Now().Unix() + int64(tokenResp.ExpiresIn) - 300
	err = keyring.Set(service, "access_token", tokenResp.AccessToken)
	if err != nil {
		return "", err
	}

	err = keyring.Set(service, "token_expiry", strconv.FormatInt(expiry, 10))
	if err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func getUserAccessToken() (string, error) {
	token, err := getValidUserToken()
	if err == nil && token != "" {
		return token, nil
	}

	return refreshUserToken()
}

func getValidUserToken() (string, error) {
	token, err := keyring.Get(service, "user_access_token")
	if err != nil {
		return "", err
	}

	expiryStr, err := keyring.Get(service, "user_token_expiry")
	if err != nil {
		return "", err
	}

	expiry, err := strconv.ParseInt(expiryStr, 10, 64)
	if err != nil {
		return "", err
	}

	if time.Now().Unix() >= expiry {
		return "", fmt.Errorf("user token expired")
	}

	return token, nil
}

func refreshUserToken() (string, error) {
	refreshToken, err := keyring.Get(service, "user_refresh_token")
	if err != nil {
		return "", fmt.Errorf("no refresh token available, please login again")
	}

	clientID, err := keyring.Get(service, "client_uid")
	if err != nil {
		return "", fmt.Errorf("client ID not found in keyring")
	}

	clientSecret, err := keyring.Get(service, "client_secret")
	if err != nil {
		return "", fmt.Errorf("client secret not found in keyring")
	}

	tokenResp, err := refreshAccessToken(refreshToken, clientID, clientSecret)
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %v", err)
	}

	expiry := time.Now().Unix() + int64(tokenResp.ExpiresIn) - 300

	err = keyring.Set(service, "user_access_token", tokenResp.AccessToken)
	if err != nil {
		return "", fmt.Errorf("failed to store refreshed user access token: %v", err)
	}

	err = keyring.Set(service, "user_token_expiry", strconv.FormatInt(expiry, 10))
	if err != nil {
		return "", fmt.Errorf("failed to store refreshed user token expiry: %v", err)
	}

	if tokenResp.RefreshToken != "" {
		err = keyring.Set(service, "user_refresh_token", tokenResp.RefreshToken)
		if err != nil {
			return "", fmt.Errorf("failed to store new refresh token: %v", err)
		}
	}

	return tokenResp.AccessToken, nil
}
