package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func NewCACertificate(caCert x509.Certificate, caPrivKey rsa.PrivateKey) (*x509.Certificate, error) {
	return NewPeerCertificate(caCert, caPrivKey, caCert, caPrivKey)
}

func NewPeerCertificate(serverCert x509.Certificate, serverPrivKey rsa.PrivateKey, caCertificate x509.Certificate, caPrivKey rsa.PrivateKey) (*x509.Certificate, error) {
	certBytes, err := x509.CreateCertificate(rand.Reader, &serverCert, &caCertificate, &serverPrivKey.PublicKey, &caPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cert: %w", err)
	}

	return x509.ParseCertificate(certBytes)
}

func SaveCertificateAsPEM(cert x509.Certificate, path string) error {
	// Save cert and key to disk
	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	if err := os.WriteFile(path, certPem, 0o600); err != nil {
		return fmt.Errorf("failed to write cert: %w", err)
	}
	return nil
}

func SavePrivateKeyAsPEM(priv rsa.PrivateKey, path string) error {
	serverPrivKeyPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(&priv)})
	if err := os.WriteFile(path, serverPrivKeyPem, 0o600); err != nil {
		return fmt.Errorf("failed to write key: %w", err)
	}
	return nil
}

func LoadPEMCertificateFromPath(path string) (*x509.Certificate, error) {
	certBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read cert: %w", err)
	}
	certPem, _ := pem.Decode(certBytes)
	return x509.ParseCertificate(certPem.Bytes)
}

func LoadPEMPrivateKeyFromPath(path string) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key: %w", err)
	}
	keyPem, _ := pem.Decode(keyBytes)
	return x509.ParsePKCS1PrivateKey(keyPem.Bytes)
}
