package pkcs7

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pkcs7",
		Short: "pkcs7 related commands",
	}

	cmd.AddCommand(newDecryptCommand())
	cmd.AddCommand(newEncryptCommand())

	return cmd
}
