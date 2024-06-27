package pkcs7

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/fullsailor/pkcs7"
)

func TestDecryptPKCS7(t *testing.T) {
	plaintext := "thisissupersecret!@$%#"

	t.Run("PKCS1", func(t *testing.T) {
		// Generate new Private Key
		privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			t.Errorf("GenerateKey returned an error when one wasn't expected: %+v", err)
		}

		// Generate a new X509 Certificate Template
		tmpl, err := certTemplate()
		if err != nil {
			t.Errorf("certTemplate returned an error when one wasn't expected: %+v", err)
		}

		// Create Certificate DER
		der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &privateKey.PublicKey, privateKey)
		if err != nil {
			t.Errorf("CreateCertificate returned an error when one wasn't expected: %+v", err)
		}

		// Convert Public Key to X509
		x509PublicCert, err := x509.ParseCertificates(der)
		if err != nil {
			t.Errorf("ParseCertificate returned an error when one wasn't expected: %+v", err)
		}

		// Encrypt our plaintext value
		enc, err := pkcs7.Encrypt([]byte(plaintext), x509PublicCert)
		if err != nil {
			t.Errorf("Encrypt returned an error when one wasn't expected: %+v", err)
		}

		// The actual test
		out, err := decryptPKCS7(
			pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}),
			pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}),
			true,
			base64.StdEncoding.EncodeToString(enc),
		)
		if err != nil {
			t.Errorf("decryptPKCS7 returned an error when one wasn't expected: %+v", err)
		}

		if string(out) != plaintext {
			t.Errorf("result of decryptPKCS7 was expected to be '%s' but was '%s'", plaintext, string(out))
		}
	})
	t.Run("PKCS8", func(t *testing.T) {
		// Generate new Private Key
		privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			t.Errorf("GenerateKey returned an error when one wasn't expected: %+v", err)
		}

		// Generate a new X509 Certificate Template
		tmpl, err := certTemplate()
		if err != nil {
			t.Errorf("certTemplate returned an error when one wasn't expected: %+v", err)
		}

		// Create Certificate DER
		der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &privateKey.PublicKey, privateKey)
		if err != nil {
			t.Errorf("CreateCertificate returned an error when one wasn't expected: %+v", err)
		}

		// Convert Public Key to X509
		x509PublicCert, err := x509.ParseCertificates(der)
		if err != nil {
			t.Errorf("ParseCertificate returned an error when one wasn't expected: %+v", err)
		}

		// Encrypt our plaintext value
		enc, err := pkcs7.Encrypt([]byte(plaintext), x509PublicCert)
		if err != nil {
			t.Errorf("Encrypt returned an error when one wasn't expected: %+v", err)
		}

		pkcs8pk, err := x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			t.Errorf("MarshalPKCS8PrivateKey returned an error when one wasn't expected: %+v", err)
		}

		// The actual test
		out, err := decryptPKCS7(
			pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8pk}),
			pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}),
			true,
			base64.StdEncoding.EncodeToString(enc),
		)
		if err != nil {
			t.Errorf("decryptPKCS7 returned an error when one wasn't expected: %+v", err)
		}

		if string(out) != plaintext {
			t.Errorf("result of decryptPKCS7 was expected to be '%s' but was '%s'", plaintext, string(out))
		}
	})
}

func certTemplate() (*x509.Certificate, error) {
	limit := new(big.Int).Lsh(big.NewInt(1), 128)

	sn, err := rand.Int(rand.Reader, limit)
	if err != nil {
		return nil, err
	}

	return &x509.Certificate{
		SerialNumber:          sn,
		Subject:               pkix.Name{Organization: []string{"Wibble Wobble, Inc."}},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // valid for a year
		BasicConstraintsValid: true,
	}, nil
}
