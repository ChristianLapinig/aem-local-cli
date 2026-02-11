package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ChristianLapinig/aem-local-cli/constants"
	"github.com/ChristianLapinig/aem-local-cli/internal/utils"
	"github.com/ChristianLapinig/aem-local-cli/models/config"
	"github.com/ChristianLapinig/aem-local-cli/models/environment"
	"github.com/spf13/cobra"
)

var (
	path     string
	envsPath string
)

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initializes the aemlocal configuration directory.",
		Long: `The init command initializes the aemlocal configuration directory, .aemlocal.

The .aemlocal directory includes a temp folder where 'create' jobs are temporarily
stored and are deleted if something fails, and a config.json file with the following
structure:

{
	"envsPaths": /path/to/aem/environments, // Program will exit if this path doesn't exist
	"environments: [] // List of local AEM environments
}
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			home := utils.GetHomePath()
			if path == "" {
				path = home
			}

			if !utils.PathExists(envsPath) && envsPath != "" {
				return errors.New("Environments path does not exist: " + envsPath)
			} else if envsPath == "" {
				envsPath = home
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

			config := &config.Config{
				EnvsPath:     envsPath,
				Environments: []environment.Environment{},
			}
			jsonData, err := json.MarshalIndent(config, "", " ")
			if err != nil {
				return utils.ErrorAndCleanup(configPath, err)
			}
			jsonPath := filepath.Join(configPath, constants.ConfigJSON)
			if err := os.WriteFile(jsonPath, jsonData, 0o644); err != nil {
				return utils.ErrorAndCleanup(configPath, err)
			}

			fmt.Printf("Initialization completed. Configuration folder created at %s\n", configPath)

			return nil
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", "", "Path where to create the configuration directory.")
	cmd.Flags().StringVarP(&envsPath, "envsPath", "e", "", "Path where local AEM environments are stored.")

	return cmd
}
