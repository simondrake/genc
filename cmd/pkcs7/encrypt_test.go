package pkcs7

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/fullsailor/pkcs7"
)

func TestEncryptPKCS7(t *testing.T) {
	plaintext := "thisissupersecret!@$%#"

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

	out, err := encryptPKCS7(plaintext, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
	if err != nil {
		t.Errorf("encryptPKCS7 returned an error when one wasn't expected: %+v", err)
	}

	// Now make sure the output can be decrypted
	p7, err := pkcs7.Parse(out)
	if err != nil {
		t.Errorf("Parse returned an error when one wasn't expected: %+v", err)
	}

	x509PubCert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Errorf("ParseCertificate returned an error when one wasn't expected: %+v", err)
	}

	dec, err := p7.Decrypt(x509PubCert, privateKey)
	if err != nil {
		t.Errorf("Decrypt returned an error when one wasn't expected: %+v", err)
	}

	if string(dec) != plaintext {
		t.Errorf("result of Decrypt was expected to be '%s' but was '%s'", plaintext, string(out))
	}
}
