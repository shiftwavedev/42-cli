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
	"time"

	"github.com/shiftwavedev/42-cli/internal/api"
	"github.com/shiftwavedev/42-cli/pkg/display"
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

func performOAuthLogin(clientID, clientSecret string) error {
	redirectURI := "http://localhost:8080/callback"

	fmt.Println(display.Header("OAuth Authentication", ""))
	fmt.Println()

	// Step 1: Generate state
	state, err := generateSecureState()
	if err != nil {
		return fmt.Errorf("failed to generate state parameter: %w", err)
	}
	oauthState = state

	authURL := fmt.Sprintf("https://api.intra.42.fr/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=public%%20profile&state=%s",
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
		url.QueryEscape(state))

	// Step 2: Open browser
	fmt.Printf("%s Opening browser...\n", display.Badge("1", display.Primary))
	fmt.Printf("%sIf the browser doesn't open, visit:\n%s%s\n\n", display.Indent, display.Indent, authURL)

	err = openBrowser(authURL)
	if err != nil {
		fmt.Printf("%sFailed to open browser automatically\n", display.Indent)
	}

	// Step 3: Wait for callback
	fmt.Printf("%s Waiting for authorization...\n", display.Badge("2", display.Primary))
	spinner := display.NewSimpleSpinner("Listening for callback on localhost:8080...")
	spinner.Start()

	callback, err := startCallbackServer()
	if err != nil {
		spinner.StopWithError("Callback server error")
		return fmt.Errorf("callback server error: %v", err)
	}

	if callback.Error != "" {
		spinner.StopWithError("Authorization denied")
		return fmt.Errorf("authentication error: %s", callback.Error)
	}
	spinner.StopWithSuccess("Authorization received")

	// Step 4: Exchange token
	fmt.Printf("%s Exchanging tokens...\n", display.Badge("3", display.Primary))
	exchangeSpinner := display.NewSimpleSpinner("Exchanging authorization code...")
	exchangeSpinner.Start()

	tokenResp, err := exchangeCodeForToken(callback.Code, clientID, clientSecret, redirectURI)
	if err != nil {
		exchangeSpinner.StopWithError("Token exchange failed")
		return fmt.Errorf("token exchange error: %v", err)
	}
	exchangeSpinner.StopWithSuccess("Tokens received")

	// Step 5: Store tokens
	fmt.Printf("%s Storing credentials...\n", display.Badge("4", display.Primary))
	if err := StoreOAuthToken(tokenResp); err != nil {
		return fmt.Errorf("failed to store OAuth token: %v", err)
	}

	fmt.Printf("\n%s Authentication complete!\n", display.Badge("✓", display.Success))

	return nil
}
