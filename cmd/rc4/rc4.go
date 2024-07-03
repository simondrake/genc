package rc4

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rc4",
		Short: "RC4 releated commands",
	}

	cmd.AddCommand(newGenerateSecretCommand())
	cmd.AddCommand(newDecryptCommand())
	cmd.AddCommand(newEncryptCommand())

	return cmd
}
