package aesgcm

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aesgcm",
		Short: "AES-GCM releated commands",
	}

	cmd.AddCommand(newEncryptCommand())
	cmd.AddCommand(newDecryptCommand())
	cmd.AddCommand(newGenerateSecretCommand())

	return cmd
}
