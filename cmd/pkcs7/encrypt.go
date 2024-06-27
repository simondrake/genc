package pkcs7

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/fullsailor/pkcs7"
	"github.com/spf13/cobra"
)

func newEncryptCommand() *cobra.Command {
	var (
		str       string
		publicKey string
		b64       bool
	)

	encryptCmd := &cobra.Command{
		Use:   "encrypt",
		Short: "encrypt plaintext with pkcs7",
		Run: func(cmd *cobra.Command, args []string) {
			pk, err := os.ReadFile(publicKey)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error reading public key: %w", err))
				os.Exit(1)
			}

			b, err := encryptPKCS7(str, pk)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error encrypting string: %w", err))
				os.Exit(1)
			}

			if b64 {
				fmt.Fprintln(os.Stdout, base64.StdEncoding.EncodeToString(b))
			} else {
				fmt.Fprintln(os.Stdout, string(b))
			}
		},
	}

	encryptCmd.Flags().StringVar(&str, "string", "", "the string to encrypt")
	encryptCmd.Flags().StringVar(&publicKey, "public-key", "", "the location of the public key on disk")
	encryptCmd.Flags().BoolVar(&b64, "base64", true, "whether the string should be base64 encoded after encryption")

	if err := encryptCmd.MarkFlagRequired("string"); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("internal error marking flag 'string' as required: %w", err))
	}
	if err := encryptCmd.MarkFlagRequired("public-key"); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("internal error marking flag 'public-key' as required: %w", err))
	}

	return encryptCmd
}

func encryptPKCS7(str string, pubKey []byte) ([]byte, error) {
	pemPub, _ := pem.Decode(pubKey)
	if pemPub == nil {
		return nil, errors.New("unable to decode public key")
	}

	certs, err := x509.ParseCertificates(pemPub.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificates: %w", err)
	}

	return pkcs7.Encrypt([]byte(str), certs)
}
