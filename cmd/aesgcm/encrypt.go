package aesgcm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
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
		Long:  "encrypt plaintext value, using AES-GCM encryption, returning the cipher text base64 encoded",
		Example: `
    
    `,
		Run: func(cmd *cobra.Command, args []string) {
			s, err := getSecret(secret, secretPath)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error retrieving secret: %w", err))
				os.Exit(1)
			}

			bytes, err := encryptAESGCM(plaintext, s)
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

func encryptAESGCM(str string, secret []byte) ([]byte, error) {
	secret, err := base64.StdEncoding.DecodeString(string(secret))
	if err != nil {
		return nil, fmt.Errorf("error decoding secret: %w", err)
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, fmt.Errorf("error creating new cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating new gcm: %w", err)
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("error creating nonce: %w", err)
	}

	return aesgcm.Seal(nonce, nonce, []byte(str), nil), nil
}
