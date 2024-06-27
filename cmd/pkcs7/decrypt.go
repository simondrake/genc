package pkcs7

import (
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/fullsailor/pkcs7"
	"github.com/spf13/cobra"
)

func newDecryptCommand() *cobra.Command {
	var (
		encString  string
		publicKey  string
		privateKey string
		b64        bool
	)

	decryptCmd := &cobra.Command{
		Use:   "decrypt",
		Short: "decrypt pkcs7 secret",
		Run: func(cmd *cobra.Command, args []string) {
			privKey, err := os.ReadFile(privateKey)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error reading private key: %w", err))
				os.Exit(1)
			}

			pubKey, err := os.ReadFile(publicKey)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error reading public key: %w", err))
				os.Exit(1)
			}

			b, err := decryptPKCS7(privKey, pubKey, b64, encString)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error decrypting string: %w", err))
				os.Exit(1)
			}

			fmt.Fprintln(os.Stdout, string(b))
		},
	}

	decryptCmd.Flags().StringVar(&encString, "string", "", "the encrypted string")
	decryptCmd.Flags().StringVar(&publicKey, "public-key", "", "the location of the public key on disk")
	decryptCmd.Flags().StringVar(&privateKey, "private-key", "", "the location of the private key on disk")
	decryptCmd.Flags().BoolVar(&b64, "base64", true, "whether the encrypted string is base64 encoded")

	if err := decryptCmd.MarkFlagRequired("string"); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("internal error marking flag 'string' as required: %w", err))
	}
	if err := decryptCmd.MarkFlagRequired("public-key"); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("internal error marking flag 'public-key' as required: %w", err))
	}
	if err := decryptCmd.MarkFlagRequired("private-key"); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("internal error marking flag 'private-key' as required: %w", err))
	}

	return decryptCmd
}

func decryptPKCS7(privKey []byte, pubKey []byte, b64 bool, enc string) ([]byte, error) {
	p7b := []byte(enc)
	if b64 {
		var err error
		p7b, err = base64.StdEncoding.DecodeString(enc)
		if err != nil {
			return nil, fmt.Errorf("error decoding encrypted string: %w", err)
		}
	}

	p7, err := pkcs7.Parse(p7b)
	if err != nil {
		return nil, fmt.Errorf("error parsing encrypted string: %w", err)
	}

	pemPriv, _ := pem.Decode(privKey)
	if pemPriv == nil {
		return nil, errors.New("unable to decode private key")
	}

	pemPub, _ := pem.Decode(pubKey)
	if pemPub == nil {
		return nil, errors.New("unable to decode public key")
	}

	cpk, err := parsePrivateKey(pemPriv)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %w", err)
	}

	x509Pub, err := x509.ParseCertificate(pemPub.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %w", err)
	}

	return p7.Decrypt(x509Pub, cpk)
}

func parsePrivateKey(b *pem.Block) (crypto.PrivateKey, error) {
	switch b.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(b.Bytes)
	case "PRIVATE KEY":
		return x509.ParsePKCS8PrivateKey(b.Bytes)
	}

	return nil, fmt.Errorf("unsupported private key type '%s'", b.Type)
}
