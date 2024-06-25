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

func newDecryptCommand() *cobra.Command {
	var (
		encString  string
		publicKey  string
		privateKey string
		base64     bool
	)

	decryptCmd := &cobra.Command{
		Use:   "decrypt",
		Short: "decrypt pkcs7 secret",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := decryptPKCS7(encString, publicKey, privateKey, base64)
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
	decryptCmd.Flags().BoolVar(&base64, "base64", true, "whether the encrypted string is base64 encoded")

	decryptCmd.MarkFlagRequired("string")
	decryptCmd.MarkFlagRequired("public-key")
	decryptCmd.MarkFlagRequired("private-key")

	return decryptCmd
}

func decryptPKCS7(enc string, pubKey string, privKey string, b64 bool) ([]byte, error) {
	bPriv, err := os.ReadFile(privKey)
	if err != nil {
		return nil, fmt.Errorf("error reading private key: %w", err)
	}

	bPub, err := os.ReadFile(pubKey)
	if err != nil {
		return nil, fmt.Errorf("error reading public key: %w", err)
	}

	pemPriv, _ := pem.Decode(bPriv)
	if pemPriv == nil {
		return nil, errors.New("unable to decode private key")
	}

	pemPub, _ := pem.Decode(bPub)
	if pemPub == nil {
		return nil, errors.New("unable to decode public key")
	}

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

	x509Priv, err := x509.ParsePKCS1PrivateKey(pemPriv.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %w", err)
	}

	x509Pub, err := x509.ParseCertificate(pemPub.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %w", err)
	}

	return p7.Decrypt(x509Pub, x509Priv)
}
