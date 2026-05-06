package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/ChristianLapinig/aem-local-cli/internal/updater"
	"github.com/spf13/cobra"
)

func NewRootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "aemlocal",
		Version: version,
		Short:   "A CLI that helps manage local AEM environments and SDKs.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("aemlocal")
			return nil
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Name() == "update" {
				return
			}
			entry, err := updater.ReadCache()
			if err != nil || time.Since(entry.CheckedAt) > updater.CacheTTL {
				go updater.RefreshCache()
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if cmd.Name() == "update" {
				return
			}
			if latest := updater.CheckForUpdate(version); latest != "" {
				fmt.Fprintf(os.Stderr, "\nA new version of aemlocal is available: %s (you have %s)\n", latest, version)
				fmt.Fprintln(os.Stderr, "Run 'aemlocal update' to upgrade.")
			}
		},
	}
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	return cmd
}
