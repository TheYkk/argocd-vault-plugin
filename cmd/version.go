package cmd

import (
	"fmt"

	"github.com/IBM/argocd-vault-plugin/version"
	"github.com/spf13/cobra"
)

// NewVersionCommand returns a new instance of the version command
func NewVersionCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "version",
		Short: "Print argocd-vault-plugin version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Version)
		},
	}

	return command
}
