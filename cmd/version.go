package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Version holds the application version
var Version = "dev"

// SetVersion sets the version from main package
func SetVersion(v string) {
	Version = v
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of 42-cli.",
	Long:  `Print the version of 42-cli and third party.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("42-cli %s\n", Version)
	},
}
