package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/shiftwavedev/42-cli/pkg/display"
)

var noColorFlag bool

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVar(&noColorFlag, "no-color", false, "Disable colored output")

	// Register commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(projectsCmd)
	rootCmd.AddCommand(correctionsCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(updateCmd)
	authCmd.AddCommand(tokenCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "42-cli",
	Short: "42-cli brings 42's intranet to your terminal.",
	Long:  `42-cli brings 42's intranet to your terminal. This tool is a alternative way of accessing public intranet data.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize display configuration before any command runs
		display.InitConfig(noColorFlag)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
