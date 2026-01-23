package helpers

import (
	"fmt"

	"github.com/shiftwavedev/42-cli/pkg/credentials"
)

// GetLoginOrDefault returns the login from args if provided, otherwise retrieves from credentials manager
func GetLoginOrDefault(args []string, credMgr *credentials.Manager) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	login, err := credMgr.GetLogin()
	if err != nil {
		return "", fmt.Errorf("you need to login first. Use '42-cli auth login'")
	}

	return login, nil
}
