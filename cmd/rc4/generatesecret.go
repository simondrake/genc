package rc4

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newGenerateSecretCommand() *cobra.Command {
	var size int16

	generateSecretCmd := &cobra.Command{
		Use:   "generate-secret",
		Short: "generate a new secret",
		Long:  "generate a new secret, base64 encoded, to be used to encrypt/decrypt with RC4, with a variable size",
		Run: func(cmd *cobra.Command, args []string) {
			if size < 1 || size > 256 {
				fmt.Fprintln(os.Stderr, fmt.Errorf("secret must be between 1 and 256"))
				os.Exit(1)
			}

			// Generate a random key, of n size
			bytes := make([]byte, size)
			if _, err := rand.Read(bytes); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error generating random secret: %w", err))
				os.Exit(1)
			}

			fmt.Fprintln(os.Stdout, base64.StdEncoding.EncodeToString(bytes))
		},
	}

	generateSecretCmd.Flags().Int16Var(&size, "size", 1, "the size of the secret to generate")

	return generateSecretCmd
}
