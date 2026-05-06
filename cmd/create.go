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
	"github.com/manifoldco/promptui"
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

By default, the environment is created in the current working directory.

Example: $ aemlocal create /path/to/license.properties /path/to/cq-quickstart.jar -n cloud-service`,
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

			// Resolve base directory
			var base string
			if path != "" {
				base = path
			} else {
				base, err = os.Getwd()
				if err != nil {
					return err
				}
			}
			if err := utils.PathExistsWithError(base); err != nil {
				return err
			}

			// Prompt for name if not provided
			if name == "" {
				prompt := promptui.Prompt{Label: "Environment name"}
				name, err = prompt.Run()
				if err != nil {
					return err
				}
			}

			// Check for duplicate name in config
			for _, e := range cfg.Environments {
				if e.Name == name {
					return fmt.Errorf("environment %q already exists", name)
				}
			}

			// Destination is always a named subdirectory of base
			dest := filepath.Join(base, name)
			if utils.PathExists(dest) {
				return fmt.Errorf("directory %s already exists", dest)
			}

			// Temp location is deleted in-case something goes wrong
			srcPath := filepath.Join(tempFolderPath, name)
			logf("creating temp directory at %s", srcPath)
			if err := os.Mkdir(srcPath, 0o755); err != nil {
				return err
			}

			paths := &paths.Paths{
				Name:              srcPath,
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
			logf("creating author instance (port %d)", authorPort)
			if err := authorInstance.Create(paths); err != nil {
				return utils.ErrorAndCleanup(srcPath, err)
			}

			logf("creating publish instance (port %d)", publishPort)
			if err := publishInstance.Create(paths); err != nil {
				return utils.ErrorAndCleanup(srcPath, err)
			}

			// dest is guaranteed not to exist; move temp folder into place
			logf("moving environment to %s", dest)
			if err := os.Rename(srcPath, dest); err != nil {
				return utils.ErrorAndCleanup(srcPath, err)
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
			logf("updating config at %s", configPath)
			if err := config.UpdateConfig(configPath, cfg); err != nil {
				return err
			}

			fmt.Printf("Successfully created AEM environment at %s\n", dest)

			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Name of the local AEM environment.")
	cmd.Flags().StringVarP(&path, "path", "p", "", "Path where environment should be created.")
	cmd.Flags().IntVar(&authorPort, "author-port", constants.DefaultAuthorPort, "Author port.")
	cmd.Flags().IntVar(&publishPort, "publish-port", constants.DefaultPublishPort, "Publish port.")

	return cmd
}
