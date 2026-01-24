package main

import "github.com/shiftwavedev/42-cli/cmd"

// Version is set via ldflags during build
var Version = "dev"

func main() {
	cmd.SetVersion(Version)
	cmd.Execute()
}
