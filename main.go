/*
Copyright © 2026 Christian Lapinig <lapinig.a.christian@gmail.com>
*/
package main

import (
	"fmt"
	"os"

	"github.com/ChristianLapinig/aem-local-cli/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
