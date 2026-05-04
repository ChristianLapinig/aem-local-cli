package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ChristianLapinig/aem-local-cli/constants"
	"github.com/ChristianLapinig/aem-local-cli/internal/utils"
	"github.com/ChristianLapinig/aem-local-cli/models/config"
	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initializes the aemlocal configuration directory.",
		Long: `The init command initializes the aemlocal configuration directory, .aemlocal.

The .aemlocal directory includes a temp folder where 'create' jobs are temporarily
stored and are deleted if something fails.
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			home := utils.GetHomePath()
			if path == "" {
				path = home
			}

			if !utils.PathExists(path) {
				return errors.New("Path does not exist: " + path)
			}

			configPath := filepath.Join(path, constants.AemLocalFolder)
			if err := os.Mkdir(configPath, 0o755); err != nil {
				return utils.ErrorAndCleanup(configPath, err)
			}

			// Make config path discoverable to other commands via marker file
			markerPath := filepath.Join(home, constants.MarkerFile)
			if err := os.WriteFile(markerPath, []byte(configPath), 0o644); err != nil {
				return err
			}

			tempFolderPath := filepath.Join(configPath, "temp")
			if err := os.Mkdir(tempFolderPath, 0o755); err != nil {
				return utils.ErrorAndCleanup(configPath, err)
			}

			if err := config.CreateConfigFile(configPath); err != nil {
				return utils.ErrorAndCleanup(configPath, err)
			}

			fmt.Printf("Initialization completed. Configuration folder created at %s\n", configPath)

			return nil
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", "", "Path where to create the configuration directory.")

	return cmd
}
