package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "42-cli",
	Short: "42-cli, is a cli that provides access to information on the intra of 42 School.",
	Long:  `42-cli, is a cli that provides access to information on the intra of 42 School via API. This tool is a alternative way of accessing public intranet data.`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

