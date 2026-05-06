package main

import (
	"os"

	"github.com/ChristianLapinig/aem-local-cli/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	rootCmd := cmd.NewRootCmd(version)
	rootCmd.AddCommand(cmd.NewInitCmd())
	rootCmd.AddCommand(cmd.NewCreateCommand())
	rootCmd.AddCommand(cmd.NewListCommand())
	rootCmd.AddCommand(cmd.NewAddCommand())
	rootCmd.AddCommand(cmd.NewDeleteCmd())
	rootCmd.AddCommand(cmd.NewUpdateCmd(version))
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
