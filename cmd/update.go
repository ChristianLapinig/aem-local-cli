package cmd

import (
	"fmt"

	"github.com/ChristianLapinig/aem-local-cli/internal/updater"
	"github.com/spf13/cobra"
)

func NewUpdateCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update aemlocal to the latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(version)
		},
	}
}

func runUpdate(currentVersion string) error {
	fmt.Println("Checking for updates...")

	latest, err := updater.FetchLatestVersion()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !updater.IsNewer(latest, currentVersion) {
		fmt.Printf("Already on the latest version (%s)\n", currentVersion)
		return nil
	}

	fmt.Printf("Updating %s → %s\n", currentVersion, latest)

	if err := updater.SelfUpdate(latest); err != nil {
		return err
	}

	fmt.Printf("Successfully updated to %s\n", latest)
	return nil
}
