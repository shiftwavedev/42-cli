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
		fmt.Println(`42-cli, is in v0.0.1.
Cobra (version use), is in v1.8.0.
go-keyring (version use), is in v0.2.6.`)
	},
}
