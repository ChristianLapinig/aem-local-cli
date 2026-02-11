package cmd

import (
	"os"

	"github.com/ChristianLapinig/aem-local-cli/models/config"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all local AEM environments.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}

			table := tablewriter.NewTable(os.Stdout)
			table.Header("Name", "Path")
			for _, environment := range cfg.Environments {
				if err := table.Append([]string{environment.Name, environment.Path}); err != nil {
					return err
				}
			}
			if err := table.Render(); err != nil {
				return err
			}
			return nil
		},
	}
}
