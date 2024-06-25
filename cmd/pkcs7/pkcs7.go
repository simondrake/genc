package pkcs7

import (
	"github.com/spf13/cobra"
)

func NewPKCS7Command() *cobra.Command {
	pkcs7Cmd := &cobra.Command{
		Use:   "pkcs7",
		Short: "pkcs7 related commands",
	}

	pkcs7Cmd.AddCommand(newDecryptCommand())

	return pkcs7Cmd
}
