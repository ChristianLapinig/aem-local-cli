package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRootCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:     "aemlocal",
		Version: version,
		Short: "A CLI that helps manage local AEM environments and SDKs.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("aemlocal")
			return nil
		},
	}
}
