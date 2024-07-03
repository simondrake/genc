package rc4

import (
	"crypto/rc4"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newDecryptCommand() *cobra.Command {
	var (
		cipherStr  string
		secret     string
		secretPath string
	)

	encryptCmd := &cobra.Command{
		Use:   "decrypt",
		Short: "decrypt RC4 cipher",
		Run: func(cmd *cobra.Command, args []string) {
			s, err := getSecret(secret, secretPath)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error getting secret: %w", err))
				os.Exit(1)
			}

			s, err = base64.StdEncoding.DecodeString(string(s))
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error decoding secret: %w", err))
				os.Exit(1)
			}

			b, err := decryptRC4(s, cipherStr)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error decrypting string: %w", err))
				os.Exit(1)
			}

			fmt.Fprintln(os.Stdout, string(b))
		},
	}

	encryptCmd.Flags().StringVar(&cipherStr, "cipher", "", "the cipher to decrypt")
	encryptCmd.Flags().StringVar(&secret, "secret", "", "the secret that can decrypt the cipher")
	encryptCmd.Flags().StringVar(&secretPath, "secret-path", "", "the location of the secret on disk")

	if err := encryptCmd.MarkFlagRequired("cipher"); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("internal error marking flag 'string' as required: %w", err))
	}

	encryptCmd.MarkFlagsOneRequired("secret", "secret-path")
	encryptCmd.MarkFlagsMutuallyExclusive("secret", "secret-path")

	return encryptCmd
}

func decryptRC4(key []byte, data string) ([]byte, error) {
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating new cipher: %w", err)
	}

	dc, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("error decoding cipher: %w", err)
	}

	plaintext := make([]byte, len(data))
	cipher.XORKeyStream(plaintext, dc)

	return plaintext, nil
}

func getSecret(secret, secretPath string) ([]byte, error) {
	if secret != "" {
		return []byte(secret), nil
	}

	return os.ReadFile(secretPath)
}
