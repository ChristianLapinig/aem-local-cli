package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "aemlocal",
		Short: "A CLI that helps manage local AEM environments and SDKs.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("aemlocal")
			return nil
		},
	}
}
