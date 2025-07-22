package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of 42-cli.",
	Long:  `Print the version of 42-cli and third party.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=== Version ===\n42-cli v0.1.0")
	},
}
