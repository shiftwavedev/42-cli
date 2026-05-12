package cmd

import (
	"github.com/shiftwavedev/42-cli/internal/api"
	"github.com/shiftwavedev/42-cli/pkg/auth"
)

var tokenManager *auth.TokenManager

func init() {
	storage := auth.NewKeyringTokenStorage(service)
	tokenManager = auth.NewTokenManager(storage, api.DefaultClient)
}
