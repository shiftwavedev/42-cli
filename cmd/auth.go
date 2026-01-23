package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"

	"github.com/shiftwavedev/42-cli/internal/display"
)

const service string = "42-cli"

func maskToken(token string) string {
	if len(token) <= 16 {
		return "****"
	}
	return token[:8] + "..." + token[len(token)-8:]
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

		if err := registerKeyring(login42, clientUid, clientSecret); err != nil {
			log.Fatalf("Failed to store credentials: %v", err)
		}
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
		if err := unregisterKeyring(); err != nil {
			log.Fatalf("Failed to remove credentials: %v", err)
		}
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
		if err := registerKeyring(login42, clientUid, clientSecret); err != nil {
			log.Fatalf("Failed to update credentials: %v", err)
		}
		fmt.Println("Authentication credentials updated successfully")
	},
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Display stored authentication credentials and current access tokens",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=== Stored Authentication Credentials ===")
		ids := [3]string{"login42", "client_uid", "client_secret"}
		for index := range 3 {
			data, err := keyring.Get(service, ids[index])
			if err != nil {
				fmt.Printf("%v: <not found>\n", ids[index])
				continue
			}
			if ids[index] == "client_secret" {
				fmt.Printf("%v: %v\n", ids[index], maskToken(data))
			} else {
				fmt.Printf("%v: %v\n", ids[index], data)
			}
		}

		fmt.Println("\n=== Application Token (Client Credentials) ===")
		appToken, err := GetApplicationToken()
		if err != nil {
			fmt.Printf("Error getting application token: %v\n", err)
		} else {
			fmt.Printf("access_token: %s\n", maskToken(appToken))
			expiryStr, _ := keyring.Get(service, "token_expiry")
			display.TokenExpiryInfo(expiryStr)
		}

		fmt.Println("\n=== User Token (OAuth) ===")
		userToken, err := GetUserToken()
		if err != nil {
			fmt.Printf("No user token available: %v\n", err)
			fmt.Println("Run '42-cli auth login' to authenticate with OAuth")
		} else {
			fmt.Printf("access_token: %s\n", maskToken(userToken))
			expiryStr, _ := keyring.Get(service, "user_token_expiry")
			display.TokenExpiryInfo(expiryStr)

			refreshToken, err := keyring.Get(service, "user_refresh_token")
			if err == nil && refreshToken != "" {
				fmt.Printf("refresh_token: %s\n", maskToken(refreshToken))
			}
		}
	},
}
