package cmd

import (
	"fmt"
	"regexp"

	"github.com/zalando/go-keyring"
)

func dataCheck(login42 string, clientUid string, clientSecret string) bool {
	checkLogin42, _ := regexp.Match("^[a-z-]+$", []byte(login42))
	checkClientUID, _ := regexp.Match("^u-s4t2ud-[a-z0-9]+$", []byte(clientUid))
	checkSecretUID, _ := regexp.Match("^s-s4t2ud-[a-z0-9]+$", []byte(clientSecret))
	return checkLogin42 && checkClientUID && checkSecretUID
}

func registerKeyring(login42 string, clientUid string, clientSecret string) error {
	if err := keyring.Set(service, "login42", login42); err != nil {
		return fmt.Errorf("failed to store login42 in keyring: %w", err)
	}

	if err := keyring.Set(service, "client_uid", clientUid); err != nil {
		return fmt.Errorf("failed to store client_uid in keyring: %w", err)
	}

	if err := keyring.Set(service, "client_secret", clientSecret); err != nil {
		return fmt.Errorf("failed to store client_secret in keyring: %w", err)
	}

	return nil
}

func unregisterKeyring() error {
	if err := keyring.DeleteAll(service); err != nil {
		return fmt.Errorf("failed to clear keyring: %w", err)
	}
	return nil
}
