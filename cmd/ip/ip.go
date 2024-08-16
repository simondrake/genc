package ip

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ip",
		Short: "ip related commands",
	}

	cmd.AddCommand(newInCIDRCommand())

	return cmd
}
