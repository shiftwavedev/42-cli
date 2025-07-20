package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

const service string = "42-cli"

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	CreatedAt   int64  `json:"created_at"`
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands for 42 intranet",
	Long:  "Commands to manage authentication credentials for 42 intranet API access",
}

var loginCmd = &cobra.Command{
	Use:   "login [login42] [client_id] [secret_id]",
	Short: "Store authentication credentials for 42 API access",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		login42, clientUid, clientSecret := args[0], args[1], args[2]
		isValidData := dataCheck(login42, clientUid, clientSecret)
		if !isValidData {
			log.Fatal("Error: Invalid credentials format")
		}
		registerKeyring(login42, clientUid, clientSecret)
		fmt.Println("Authentication credentials stored successfully")
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
	Short: "Display stored authentication credentials and current access token",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=== Stored Authentication Credentials ===")
		ids := [3]string{"login42", "client_uid", "client_secret"}
		for index := 0; index < 3; index++ {
			data, err := keyring.Get(service, ids[index])
			checkError(err)
			fmt.Printf("%v: %v\n", ids[index], data)
		}
		
		fmt.Println("\n=== Current Access Token ===")
		token, err := GetAccessToken()
		if err != nil {
			fmt.Printf("Error getting access token: %v\n", err)
		} else {
			fmt.Printf("access_token: %s\n", token)
			
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
	},
}

func checkError(err error) {
	if err != nil {
		log.Fatal("Error: Authentication operation failed")
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
	checkError(err)

	err = keyring.Set(service, ids[1], clientUid)
	checkError(err)

	err = keyring.Set(service, ids[2], clientSecret)
	checkError(err)
}

func unregisterKeyring() {
	err := keyring.DeleteAll(service)
	checkError(err)
}

func GetAccessToken() (string, error) {
	token, err := getValidToken()
	if err == nil && token != "" {
		return token, nil
	}

	return generateNewToken()
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
