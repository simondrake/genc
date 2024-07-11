package rc4

import (
	"crypto/rc4"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newEncryptCommand() *cobra.Command {
	var (
		plaintext  string
		secret     string
		secretPath string
	)

	encryptCmd := &cobra.Command{
		Use:   "encrypt",
		Short: "encrypt plaintext",
		Long:  "encrypt plaintext value, using RC4 encryption, returning the cipher text base64 encoded",
		Run: func(cmd *cobra.Command, args []string) {
			s, err := getSecret(secret, secretPath)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error retrieving secret: %w", err))
				os.Exit(1)
			}

			bytes, err := encryptRC4([]byte(plaintext), s)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error encrypting string: %w", err))
				os.Exit(1)
			}

			fmt.Fprintln(os.Stdout, base64.StdEncoding.EncodeToString(bytes))
		},
	}

	encryptCmd.Flags().StringVar(&plaintext, "plaintext", "", "the plaintext to encrypt")
	encryptCmd.Flags().StringVar(&secret, "secret", "", "the secret to encrypt the string with")
	encryptCmd.Flags().StringVar(&secretPath, "secret-path", "", "the location of the secret on disk")

	if err := encryptCmd.MarkFlagRequired("plaintext"); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("internal error marking flag 'string' as required: %w", err))
	}

	encryptCmd.MarkFlagsOneRequired("secret", "secret-path")
	encryptCmd.MarkFlagsMutuallyExclusive("secret", "secret-path")

	return encryptCmd
}

func encryptRC4(data []byte, key []byte) ([]byte, error) {
	secret, err := base64.StdEncoding.DecodeString(string(key))
	if err != nil {
		return nil, fmt.Errorf("error decoding secret: %w", err)
	}

	cipher, err := rc4.NewCipher(secret)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher: %w", err)
	}

	cipherText := make([]byte, len(data))
	cipher.XORKeyStream(cipherText, data)

	return cipherText, nil
}
