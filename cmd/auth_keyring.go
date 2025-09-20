package cmd

import (
	"log"
	"regexp"

	"github.com/zalando/go-keyring"
)

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
