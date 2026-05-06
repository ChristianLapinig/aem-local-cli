package cmd

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/ChristianLapinig/aem-local-cli/models/config"
	"github.com/ChristianLapinig/aem-local-cli/models/environment"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func NewDeleteCmd() *cobra.Command {
	var purge bool
	var name string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Remove a configured AEM environment",
		Long: `Remove a configured AEM environment from the local CLI.

Removes the environment entry from the config file. Optionally deletes the
environment directory from the filesystem using the --purge flag.

If --name is not provided, an interactive prompt lets you select the environment
to delete. You will always be asked to confirm before any changes are made.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}
			environments := cfg.Environments
			if len(environments) == 0 {
				fmt.Println("No environments configured.")
				return nil
			}

			// User selects environment to delete if names flag is not passed
			if name == "" {
				names := make([]string, len(environments))
				for i, env := range environments {
					names[i] = env.Name
				}

				prompt := promptui.Select{
					Label: "Select environment to delete",
					Items: names,
				}
				_, name, err = prompt.Run()
				if err != nil {
					return err
				}
			}

			var env environment.Environment
			var found bool
			for _, e := range environments {
				if e.Name == name {
					env = e
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("Environment %s not found. Aborting.", name)
			}

			reader := bufio.NewReader(cmd.InOrStdin())

			fmt.Printf("Are you sure you want to delete %s? [y/N]: ", name)
			confirmDelete, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			if strings.ToLower(strings.TrimSpace(confirmDelete)) != "y" {
				fmt.Println("Aborting.")
				return nil
			}

			environments = slices.DeleteFunc(environments, func(e environment.Environment) bool {
				return e.Name == name
			})

			if !purge {
				fmt.Printf("Do you want to delete the environment folder at %s? [y/N]: ", env.Path)
				answer, err := reader.ReadString('\n')
				if err != nil {
					return err
				}
				purge = strings.ToLower(strings.TrimSpace(answer)) == "y"
			}

			if purge {
				if env.Path == "" {
					return fmt.Errorf("No path configured for %s. Aborting.", name)
				}
				if err := os.RemoveAll(env.Path); err != nil {
					return err
				}
			}

			configPath, err := config.GetConfigPath()
			if err != nil {
				return err
			}

			cfg.Environments = environments

			if err := config.UpdateConfig(configPath, cfg); err != nil {
				return err
			}

			fmt.Printf("Successfully deleted %s.\n", env.Name)
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Name of the environment to remove.")
	cmd.Flags().BoolVar(&purge, "purge", false, "Remove the environment from the filesystem.")

	return cmd
}
