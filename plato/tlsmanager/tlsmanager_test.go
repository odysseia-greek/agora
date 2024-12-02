package tlsmanager

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTLSManager_LoadCertificates(t *testing.T) {
	// Create temporary files for cert, key, and CA
	tempDir := t.TempDir()
	certPath := filepath.Join(tempDir, "/tls.crt")
	keyPath := filepath.Join(tempDir, "/tls.key")

	err := generateSelfSignedCert(certPath, keyPath, "localhost")
	assert.NoError(t, err, "Certificate generation should not fail")

	// Initialize the TLSManager
	tlsManager := NewTLSManager("", tempDir, 1*time.Second)
	err = tlsManager.LoadCertificates()
	assert.NoError(t, err, "Initial certificate load should not fail")

	// Assert the initial certificate is loaded
	assert.Len(t, tlsManager.GetTLSConfig().Certificates, 1, "TLS Config should have one certificate")

	err = generateSelfSignedCert(certPath, keyPath, "localhost")
	assert.NoError(t, err, "Certificate generation should not fail")

	// Reload certificates
	err = tlsManager.LoadCertificates()
	assert.NoError(t, err, "Certificate reload should not fail")

	// Assert the updated certificate is added
	assert.Len(t, tlsManager.GetTLSConfig().Certificates, 2, "TLS Config should have two certificates")

	// Wait for grace period and assert old certificate removal
	time.Sleep(2 * time.Second)
	assert.Len(t, tlsManager.GetTLSConfig().Certificates, 1, "Old certificate should be removed after grace period")
}

func TestTLSManager_WatchCertificates(t *testing.T) {
	// Create temporary files for cert, key, and CA
	tempDir := t.TempDir()
	certPath := filepath.Join(tempDir, "/tls.crt")
	keyPath := filepath.Join(tempDir, "/tls.key")

	err := generateSelfSignedCert(certPath, keyPath, "localhost")
	assert.NoError(t, err, "Certificate generation should not fail")

	// Initialize the TLSManager
	tlsManager := NewTLSManager("", tempDir, 1*time.Second)
	err = tlsManager.LoadCertificates()
	assert.NoError(t, err, "Initial certificate load should not fail")

	// Start watching for changes
	tlsManager.WatchCertificates(500 * time.Millisecond)

	err = generateSelfSignedCert(certPath, keyPath, "localhost")
	assert.NoError(t, err, "Certificate generation should not fail")

	// Wait for watch to detect change
	time.Sleep(1 * time.Second)

	// Assert the updated certificate is added
	assert.Len(t, tlsManager.GetTLSConfig().Certificates, 2, "TLS Config should have two certificates after file change")

	// Wait for grace period and assert old certificate removal
	time.Sleep(2 * time.Second)
	assert.Len(t, tlsManager.GetTLSConfig().Certificates, 1, "Old certificate should be removed after grace period")
}

// Helper to write content to a file
func writeToFile(path, content string) {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}

// GenerateSelfSignedCert generates a self-signed certificate and writes it to the specified paths.
func generateSelfSignedCert(certPath, keyPath string, commonName string) error {
	// Generate a private key
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	// Create a certificate template
	template := x509.Certificate{
		SerialNumber:          big.NewInt(2019),
		Subject:               pkix.Name{CommonName: commonName},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour), // Valid for 1 day
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Self-sign the certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	// Encode and write the certificate to a file
	certFile, err := os.Create(certPath)
	if err != nil {
		return err
	}
	defer certFile.Close()

	err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	if err != nil {
		return err
	}

	// Encode and write the private key to a file
	keyFile, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return err
	}

	err = pem.Encode(keyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	if err != nil {
		return err
	}

	return nil
}
