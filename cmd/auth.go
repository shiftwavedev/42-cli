package cmd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

const service string = "42-cli"

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
		for index := range 3 {
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
