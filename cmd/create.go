package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ChristianLapinig/aem-local-cli/constants"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	var name string
	var path string
	var authorPort int
	var publishPort int

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Generates a local AEM environment with author and publish folders.",
		Long: `Generates a local AEM environment with author and publish folders. The
command assumes and requires that you have a valid license.properties and AEM Quickstart
JAR files.

By default, environments are stored under the 'envsPath' value set in .aemlocal/config.json.

Example: $ aemlocal create license.properties cq-quickstart.jar -p cloud-service

The new environment will be stored in /{envsPath}/aem/cloud-service.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := filepath.Join(os.Getenv(constants.AemLocalPathEnvVar), "config.json")
			fmt.Println(configPath)

			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "aem", "Name of the local AEM environment.")
	cmd.Flags().StringVarP(&path, "path", "p", "", "Where the environment should be stored. This will be created as a relative path inside envsPath.")
	cmd.Flags().IntVar(&authorPort, "author-port", constants.DefaultAuthorPort, "Author port.")
	cmd.Flags().IntVar(&publishPort, "publish-port", constants.DefaultPublishPort, "Publish port.")

	return cmd
}
