package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"log"
	"regexp"
)

const service string = "42-cli"

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate 42-cli with 42 intranet.",
}

var loginCmd = &cobra.Command{
	Use:   "login [login42] [client_id] [secret_id]",
	Short: "Authenticate with API 42 keys.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		login42, clientUid, clientSecret := args[0], args[1], args[2]
		isValidData := dataCheck(login42, clientUid, clientSecret)
		if !isValidData {
			log.Fatal("Error: The data supplied are not valid.")
		}
		registerKeyring(login42, clientUid, clientSecret)
		fmt.Println("Authentication token register.")
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Removing API 42 keys.",
	Run: func(cmd *cobra.Command, args []string) {
		unregisterKeyring()
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updating your 42 keys.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		login42, clientUid, clientSecret := args[0], args[1], args[2]
		isValidData := dataCheck(login42, clientUid, clientSecret)
		if !isValidData {
			log.Fatal("Error: The data supplied are not valid.")
		}
		registerKeyring(login42, clientUid, clientSecret)
		fmt.Println("Authentication token updated.")
	},
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Get your 42 keys.",
	Run: func(cmd *cobra.Command, args []string) {
		ids := [3]string{"login42", "client_uid", "client_secret"}
		for index := 0; index < 3; index++ {
			data, err := keyring.Get(service, ids[index])
			checkError(err)
			fmt.Printf("%v: %v\n", ids[index], data)
		}
	},
}

func checkError(err error) {
	if err != nil {
		log.Fatal("Error: Unsuccessful action.")
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
