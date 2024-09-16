package jwt

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jwt",
		Short: "jwt related commands",
	}

	cmd.AddCommand(newCreateCommand())
	cmd.AddCommand(newParseCommand())

	return cmd
}
