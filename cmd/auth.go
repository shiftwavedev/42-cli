package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

const service string = "42-cli"

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	CreatedAt    int64  `json:"created_at"`
	RefreshToken string `json:"refresh_token"`
}


type CallbackResponse struct {
	Code  string
	Error string
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands for 42 intranet",
	Long:  "Commands to manage authentication credentials for 42 intranet API access",
}

var loginCmd = &cobra.Command{
	Use:   "login [login42] [client_id] [secret_id]",
	Short: "Store authentication credentials and authenticate with 42 API using OAuth",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		login42, clientUid, clientSecret := args[0], args[1], args[2]
		isValidData := dataCheck(login42, clientUid, clientSecret)
		if !isValidData {
			log.Fatal("Error: Invalid credentials format")
		}
		
		registerKeyring(login42, clientUid, clientSecret)
		fmt.Println("Credentials stored successfully")
		fmt.Println("Starting OAuth authentication...")
		
		err := performOAuthLogin(clientUid, clientSecret)
		if err != nil {
			log.Fatalf("OAuth authentication failed: %v", err)
		}
		
		fmt.Println("OAuth authentication successful! You can now access user endpoints.")
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored authentication credentials",
	Run: func(cmd *cobra.Command, args []string) {
		unregisterKeyring()
		fmt.Println("Authentication credentials removed successfully")
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [login42] [client_id] [secret_id]",
	Short: "Update stored authentication credentials",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		login42, clientUid, clientSecret := args[0], args[1], args[2]
		isValidData := dataCheck(login42, clientUid, clientSecret)
		if !isValidData {
			log.Fatal("Error: Invalid credentials format")
		}
		registerKeyring(login42, clientUid, clientSecret)
		fmt.Println("Authentication credentials updated successfully")
	},
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Display stored authentication credentials and current access tokens",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=== Stored Authentication Credentials ===")
		ids := [3]string{"login42", "client_uid", "client_secret"}
		for index := 0; index < 3; index++ {
			data, err := keyring.Get(service, ids[index])
			checkError(err, "Keyring data not found")
			fmt.Printf("%v: %v\n", ids[index], data)
		}
		
		fmt.Println("\n=== Application Token (Client Credentials) ===")
		appToken, err := GetApplicationToken()
		if err != nil {
			fmt.Printf("Error getting application token: %v\n", err)
		} else {
			fmt.Printf("access_token: %s\n", appToken)
			
			expiryStr, err := keyring.Get(service, "token_expiry")
			if err == nil {
				expiry, err := strconv.ParseInt(expiryStr, 10, 64)
				if err == nil {
					expiryTime := time.Unix(expiry, 0)
					fmt.Printf("expires_at: %s\n", expiryTime.Format("2006-01-02 15:04:05"))
					
					timeLeft := time.Until(expiryTime)
					if timeLeft > 0 {
						fmt.Printf("time_left: %v\n", timeLeft.Round(time.Minute))
					} else {
						fmt.Printf("status: expired\n")
					}
				}
			}
		}

		fmt.Println("\n=== User Token (OAuth) ===")
		userToken, err := GetUserToken()
		if err != nil {
			fmt.Printf("No user token available: %v\n", err)
			fmt.Println("Run '42-cli auth login' to authenticate with OAuth")
		} else {
			fmt.Printf("access_token: %s\n", userToken)
			
			expiryStr, err := keyring.Get(service, "user_token_expiry")
			if err == nil {
				expiry, err := strconv.ParseInt(expiryStr, 10, 64)
				if err == nil {
					expiryTime := time.Unix(expiry, 0)
					fmt.Printf("expires_at: %s\n", expiryTime.Format("2006-01-02 15:04:05"))
					
					timeLeft := time.Until(expiryTime)
					if timeLeft > 0 {
						fmt.Printf("time_left: %v\n", timeLeft.Round(time.Minute))
					} else {
						fmt.Printf("status: expired\n")
					}
				}
			}

			refreshToken, err := keyring.Get(service, "user_refresh_token")
			if err == nil && refreshToken != "" {
				fmt.Printf("refresh_token: %s\n", refreshToken)
			}
		}
	},
}

func checkError(err error, message string) {
	if err != nil {
		log.Fatal("Error: " + message)
	}
}

func dataCheck(login42 string, clientUid string, clientSecret string) bool {
	checkLogin42, _ := regexp.Match("^[a-z-]+$", []byte(login42))
	checkClientUID, _ := regexp.Match("^u-s4t2ud-[a-z0-9]+$", []byte(clientUid))
	checkSecretUID, _ := regexp.Match("^s-s4t2ud-[a-z0-9]+$", []byte(clientSecret))
	return checkLogin42 && checkClientUID && checkSecretUID
}

func registerKeyring(login42 string, clientUid string, clientSecret string) {
	ids := [3]string{"login42", "client_uid", "client_secret"}
	var err error
	err = keyring.Set(service, ids[0], login42)
	checkError(err, "Authentication regiter operation failed (login 42)")

	err = keyring.Set(service, ids[1], clientUid)
	checkError(err, "Authentication regiter operation failed (client_id)")

	err = keyring.Set(service, ids[2], clientSecret)
	checkError(err, "Authentication regiter operation failed (client_secret)")
}

func unregisterKeyring() {
	err := keyring.DeleteAll(service)
	checkError(err, "Unregiter operation failed")
}

func GetAccessToken() (string, error) {
	return GetApplicationToken()
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

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("scope", "public profile")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.PostForm("https://api.intra.42.fr/oauth/token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
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
		error := r.URL.Query().Get("error")

		if error != "" {
			fmt.Fprintf(w, "<html><body><h2>Erreur d'authentification</h2><p>%s</p><p>Vous pouvez fermer cette fenêtre.</p></body></html>", error)
			callbackChan <- &CallbackResponse{Error: error}
			return
		}

		fmt.Fprintf(w, "<html><body><h2>Authentification réussie!</h2><p>Vous pouvez fermer cette fenêtre et retourner au terminal.</p></body></html>")
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

func exchangeCodeForToken(code, clientID, clientSecret, redirectURI string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.PostForm("https://api.intra.42.fr/oauth/token", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func refreshAccessToken(refreshToken, clientID, clientSecret string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.PostForm("https://api.intra.42.fr/oauth/token", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token refresh failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func performOAuthLogin(clientID, clientSecret string) error {
	redirectURI := "http://localhost:8080/callback"
	
	authURL := fmt.Sprintf("https://api.intra.42.fr/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=public%%20profile",
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI))

	fmt.Println("Opening browser for authentication...")
	fmt.Printf("If the browser doesn't open automatically, please visit: %s\n", authURL)
	
	err := openBrowser(authURL)
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
