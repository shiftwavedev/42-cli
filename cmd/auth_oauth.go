package cmd

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/zalando/go-keyring"

	"github.com/shiftwavedev/42-cli/internal/api"
)

var oauthState string

func generateSecureState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random state: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func startCallbackServer() (*CallbackResponse, error) {
	callbackChan := make(chan *CallbackResponse, 1)
	server := &http.Server{Addr: ":8080"}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		error := r.URL.Query().Get("error")

		if state != oauthState {
			fmt.Fprintf(w, "<html><body><h2>Security Error</h2><p>Invalid state parameter. Possible CSRF attack detected.</p><p>You can close this window.</p></body></html>")
			callbackChan <- &CallbackResponse{Error: "invalid_state"}
			return
		}

		if error != "" {
			fmt.Fprintf(w, "<html><body><h2>Authentication Error</h2><p>%s</p><p>You can close this window.</p></body></html>", error)
			callbackChan <- &CallbackResponse{Error: error}
			return
		}

		fmt.Fprintf(w, "<html><body><h2>Authentication Successful!</h2><p>You can close this window and return to the terminal.</p></body></html>")
		callbackChan <- &CallbackResponse{Code: code}
	})

	go func() {
		server.ListenAndServe()
	}()

	callback := <-callbackChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	return callback, nil
}

func exchangeCodeForToken(code, clientID, clientSecret, redirectURI string) (*api.TokenResponse, error) {
	return api.DefaultClient.ExchangeAuthorizationCode(code, clientID, clientSecret, redirectURI)
}

func refreshAccessToken(refreshToken, clientID, clientSecret string) (*api.TokenResponse, error) {
	return api.DefaultClient.RefreshToken(refreshToken, clientID, clientSecret)
}

func performOAuthLogin(clientID, clientSecret string) error {
	redirectURI := "http://localhost:8080/callback"

	state, err := generateSecureState()
	if err != nil {
		return fmt.Errorf("failed to generate state parameter: %w", err)
	}
	oauthState = state

	authURL := fmt.Sprintf("https://api.intra.42.fr/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=public%%20profile&state=%s",
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
		url.QueryEscape(state))

	fmt.Println("Opening browser for authentication...")
	fmt.Printf("If the browser doesn't open automatically, please visit: %s\n", authURL)

	err = openBrowser(authURL)
	if err != nil {
		fmt.Printf("Failed to open browser automatically: %v\n", err)
		fmt.Printf("Please manually visit: %s\n", authURL)
	}

	fmt.Println("Waiting for authentication callback...")
	callback, err := startCallbackServer()
	if err != nil {
		return fmt.Errorf("callback server error: %v", err)
	}

	if callback.Error != "" {
		return fmt.Errorf("authentication error: %s", callback.Error)
	}

	fmt.Println("Exchanging authorization code for access token...")
	tokenResp, err := exchangeCodeForToken(callback.Code, clientID, clientSecret, redirectURI)
	if err != nil {
		return fmt.Errorf("token exchange error: %v", err)
	}

	expiry := time.Now().Unix() + int64(tokenResp.ExpiresIn) - 300

	err = keyring.Set(service, "user_access_token", tokenResp.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to store user access token: %v", err)
	}

	err = keyring.Set(service, "user_token_expiry", strconv.FormatInt(expiry, 10))
	if err != nil {
		return fmt.Errorf("failed to store user token expiry: %v", err)
	}

	if tokenResp.RefreshToken != "" {
		err = keyring.Set(service, "user_refresh_token", tokenResp.RefreshToken)
		if err != nil {
			return fmt.Errorf("failed to store user refresh token: %v", err)
		}
	}

	return nil
}
