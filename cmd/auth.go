package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"

	"github.com/shiftwavedev/42-cli/pkg/credentials"
	"github.com/shiftwavedev/42-cli/pkg/display"
)

const service string = "42-cli"

var credManager *credentials.Manager

func init() {
	storage := credentials.NewKeyringStorage(service)
	credManager = credentials.NewManager(storage)
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
		creds := &credentials.Credentials{
			Login42:      args[0],
			ClientID:     args[1],
			ClientSecret: args[2],
		}

		if err := credManager.Store(creds); err != nil {
			log.Fatalf("Failed to store credentials: %v", err)
		}
		fmt.Println("Credentials stored successfully")
		fmt.Println("Starting OAuth authentication...")

		err := performOAuthLogin(creds.ClientID, creds.ClientSecret)
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
		if err := credManager.Clear(); err != nil {
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
		creds := &credentials.Credentials{
			Login42:      args[0],
			ClientID:     args[1],
			ClientSecret: args[2],
		}

		if err := credManager.Store(creds); err != nil {
			log.Fatalf("Failed to update credentials: %v", err)
		}
		fmt.Println("Authentication credentials updated successfully")
	},
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Display stored authentication credentials and current access tokens",
	Run: func(cmd *cobra.Command, args []string) {
		// Gather credentials
		var creds *display.CredentialInfo
		login, loginErr := keyring.Get(service, "login42")
		clientID, clientErr := keyring.Get(service, "client_uid")
		clientSecret, secretErr := keyring.Get(service, "client_secret")

		if loginErr == nil || clientErr == nil || secretErr == nil {
			creds = &display.CredentialInfo{
				Login:        login,
				ClientID:     clientID,
				ClientSecret: clientSecret,
			}
		}

		// Gather app token info
		var appTokenInfo *display.TokenInfo
		appToken, err := GetApplicationToken()
		if err == nil {
			expiryStr, _ := keyring.Get(service, "token_expiry")
			expiry, _ := strconv.ParseInt(expiryStr, 10, 64)
			appTokenInfo = &display.TokenInfo{
				Token:      appToken,
				ExpiryUnix: expiry,
			}
		}

		// Gather user token info
		var userTokenInfo *display.TokenInfo
		userToken, err := GetUserToken()
		if err == nil {
			expiryStr, _ := keyring.Get(service, "user_token_expiry")
			expiry, _ := strconv.ParseInt(expiryStr, 10, 64)
			refreshToken, _ := keyring.Get(service, "user_refresh_token")
			userTokenInfo = &display.TokenInfo{
				Token:        userToken,
				ExpiryUnix:   expiry,
				RefreshToken: refreshToken,
			}
		}

		fmt.Print(display.RenderAuthStatus(creds, appTokenInfo, userTokenInfo))
	},
}
