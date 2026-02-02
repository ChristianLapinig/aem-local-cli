package main

import (
	"os"

	"github.com/ChristianLapinig/aem-local-cli/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	rootCmd.AddCommand(cmd.NewInitCmd())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
