package aesgcm

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newGenerateSecretCommand() *cobra.Command {
	size := secretSize24

	generateSecretCmd := &cobra.Command{
		Use:   "generate-secret",
		Short: "generate a new secret",
		Long:  "generate a new secret, base64 encoded, to be used to encrypt/decrypt with AES-GCM, with a variable size",
		Example: `
    # Create a new secret of size 32
    genc aesgcm generate-secret --size 32

    # Create a new secret of size 24
    genc aesgcm generate-secret --size 24`,
		Run: func(cmd *cobra.Command, args []string) {
			// Generate a random key, of n size
			bytes := make([]byte, size)
			if _, err := rand.Read(bytes); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error generating random secret: %w", err))
				os.Exit(1)
			}

			fmt.Fprintln(os.Stdout, base64.StdEncoding.EncodeToString(bytes))
		},
	}

	generateSecretCmd.Flags().Var(&size, "size", "the size of the secret to generate")

	return generateSecretCmd
}
