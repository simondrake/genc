package cidr

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cidr",
		Short: "cidr related commands",
	}

	cmd.AddCommand(newOverlapCommand())

	return cmd
}
