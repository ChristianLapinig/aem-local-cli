package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ChristianLapinig/aem-local-cli/constants"
	"github.com/ChristianLapinig/aem-local-cli/internal/utils"
	"github.com/ChristianLapinig/aem-local-cli/models/config"
	"github.com/ChristianLapinig/aem-local-cli/models/environment"
	"github.com/ChristianLapinig/aem-local-cli/models/instance"
	"github.com/ChristianLapinig/aem-local-cli/models/paths"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	var name string
	var path string
	var authorPort int
	var publishPort int

	cmd := &cobra.Command{
		Use:   "create <license-properties-path> <aem-quickstart-jar-path>",
		Args:  cobra.ExactArgs(2),
		Short: "Generates a local AEM environment with author and publish folders.",
		Long: `Generates a local AEM environment with author and publish folders. The
command assumes and requires that you have a valid license.properties and AEM Quickstart
JAR files.

By default, environments are stored under the 'envsPath' value set in .aemlocal/config.json.

Example: $ aemlocal create /path/to/license.properties /path/to/cq-quickstart.jar -p cloud-service

The new environment will be stored in /{envsPath}/aem/cloud-service.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			tempFolderPath, err := config.GetTempFolderPath()
			if err != nil {
				return err
			}
			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}
			licensePropertiesPath := args[0]
			quickstartJarPath := args[1]
			if err := utils.PathExistsWithError(licensePropertiesPath); err != nil {
				return err
			}
			if err := utils.PathExistsWithError(quickstartJarPath); err != nil {
				return err
			}

			// Temp location is deleted in-case something goes wrong
			src := filepath.Join(tempFolderPath, name)
			if err := os.Mkdir(src, 0o755); err != nil {
				return err
			}

			paths := &paths.Paths{
				Name:              src,
				LicenseProperties: licensePropertiesPath,
				QuickstartJAR:     quickstartJarPath,
			}
			authorInstance := &instance.Instance{
				Name: constants.Author,
				Port: authorPort,
			}
			publishInstance := &instance.Instance{
				Name: constants.Publish,
				Port: publishPort,
			}

			// Create author and publish instance folders
			if err := authorInstance.Create(paths); err != nil {
				return utils.ErrorAndCleanup(src, err)
			}

			if err := publishInstance.Create(paths); err != nil {
				return utils.ErrorAndCleanup(src, err)
			}

			// Move environment from temp folder to final destination
			dest := cfg.EnvsPath
			if path != "" && utils.PathExists(filepath.Join(dest, path)) {
				dest = filepath.Join(dest, path, name)
			} else {
				dest = filepath.Join(dest, name)
			}

			if err := os.Rename(src, dest); err != nil {
				return utils.ErrorAndCleanup(dest, err)
			}

			environment := environment.Environment{
				Name: name,
				Path: dest,
			}
			cfg.Environments = append(cfg.Environments, environment)
			configPath, err := config.GetConfigPath()
			if err != nil {
				return err
			}
			if err := config.UpdateConfig(configPath, cfg); err != nil {
				return err
			}

			fmt.Printf("Successfully created AEM environment at %s\n", dest)

			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "aem", "Name of the local AEM environment.")
	cmd.Flags().StringVarP(&path, "path", "p", "", "Where the environment should be stored. This should be a relative path inside envsPath.")
	cmd.Flags().IntVar(&authorPort, "author-port", constants.DefaultAuthorPort, "Author port.")
	cmd.Flags().IntVar(&publishPort, "publish-port", constants.DefaultPublishPort, "Publish port.")

	return cmd
}
