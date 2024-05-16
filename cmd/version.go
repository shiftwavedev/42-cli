package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of 42-CLI.",
	Long:  `Print the version of 42-CLI and third party.`,
	Run:   func(cmd *cobra.Command, args []string) {
		fmt.Println(`42-CLI, is in v0.0.1.
Cobra (version use), is in v1.8.0.`)
	},
}
