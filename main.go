package main

import (
	"os"

	"github.com/ChristianLapinig/aem-local-cli/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	rootCmd.AddCommand(cmd.NewInitCmd())
	rootCmd.AddCommand(cmd.NewCreateCommand())
	rootCmd.AddCommand(cmd.NewListCommand())
	rootCmd.AddCommand(cmd.NewAddCommand())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
