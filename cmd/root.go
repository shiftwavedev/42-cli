package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(projectsCmd)
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
	Run:   func(cmd *cobra.Command, args []string) {},
}
