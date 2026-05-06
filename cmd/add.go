package cmd

import (
	"fmt"

	"github.com/ChristianLapinig/aem-local-cli/constants"
	"github.com/ChristianLapinig/aem-local-cli/internal/utils"
	"github.com/ChristianLapinig/aem-local-cli/models/config"
	"github.com/ChristianLapinig/aem-local-cli/models/environment"
	"github.com/spf13/cobra"
)

func NewAddCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "add <name> <path-to-environment>",
		Args:  cobra.ExactArgs(2),
		Short: "Add an existing environment.",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			path := args[1]
			if !utils.PathExists(path) {
				return fmt.Errorf("%s %s", constants.PathDoesNotExist, path)
			}
			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}

			for _, e := range cfg.Environments {
				if e.Name == name {
					return fmt.Errorf("Environment %s already exists.", name)
				}
			}

			env := environment.Environment{
				Name: name,
				Path: path,
			}
			cfg.Environments = append(cfg.Environments, env)
			configPath, err := config.GetConfigPath()
			if err != nil {
				return err
			}
			if err := config.UpdateConfig(configPath, cfg); err != nil {
				return err
			}

			fmt.Printf("Successfully added environment %s at %s\n", name, path)
			return nil
		},
	}
}
