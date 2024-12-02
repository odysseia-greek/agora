package tlsmanager

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
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
	caPath := filepath.Join(tempDir, "/tls.pem")

	err := generateSelfSignedCert(certPath, keyPath, caPath, "localhost")
	assert.NoError(t, err, "Certificate generation should not fail")

	// Initialize the TLSManager
	tlsManager := NewTLSManager("", tempDir, 1*time.Second)
	err = tlsManager.LoadCertificates()
	assert.NoError(t, err, "Initial certificate load should not fail")

	// Assert the initial certificate is loaded
	assert.Len(t, tlsManager.GetTLSConfig().Certificates, 1, "TLS Config should have one certificate")

	err = generateSelfSignedCert(certPath, keyPath, caPath, "localhost")
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
	caPath := filepath.Join(tempDir, "/tls.pem")

	err := generateSelfSignedCert(certPath, keyPath, caPath, "localhost")
	assert.NoError(t, err, "Certificate generation should not fail")

	// Initialize the TLSManager
	tlsManager := NewTLSManager("", tempDir, 1000*time.Millisecond)
	err = tlsManager.LoadCertificates()
	assert.NoError(t, err, "Initial certificate load should not fail")

	// Start watching for changes
	tlsManager.WatchCertificates(200 * time.Millisecond)

	err = generateSelfSignedCert(certPath, keyPath, caPath, "localhost")
	assert.NoError(t, err, "Certificate generation should not fail")

	// Wait for watch to detect change
	time.Sleep(500 * time.Millisecond)

	// Assert the updated certificate is added
	fmt.Printf("%v", len(tlsManager.GetTLSConfig().Certificates))
	assert.True(t, len(tlsManager.GetTLSConfig().Certificates) >= 2)

	// Wait for grace period and assert old certificate removal
	time.Sleep(2 * time.Second)
	assert.Len(t, tlsManager.GetTLSConfig().Certificates, 1, "Old certificate should be removed after grace period")
}

// GenerateSelfSignedCert generates a self-signed CA certificate, signs a leaf certificate, and writes them to the specified paths.
func generateSelfSignedCert(certPath, keyPath, caPath string, commonName string) error {
	// Generate CA private key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	// Create the CA certificate
	caCert := &x509.Certificate{
		SerialNumber:          big.NewInt(2024),
		Subject:               pkix.Name{CommonName: commonName + " CA"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // CA valid for 1 year
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	// Self-sign the CA certificate
	caCertBytes, err := x509.CreateCertificate(rand.Reader, caCert, caCert, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return err
	}

	// Write the CA certificate to the specified file
	caFile, err := os.Create(caPath)
	if err != nil {
		return err
	}
	defer caFile.Close()

	if err := pem.Encode(caFile, &pem.Block{Type: "CERTIFICATE", Bytes: caCertBytes}); err != nil {
		return err
	}

	// Generate leaf certificate private key
	leafPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	// Create the leaf certificate
	leafCert := &x509.Certificate{
		SerialNumber: big.NewInt(2025),
		Subject:      pkix.Name{CommonName: commonName},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(30 * 24 * time.Hour), // Leaf cert valid for 30 days
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	// Sign the leaf certificate with the CA private key
	leafCertBytes, err := x509.CreateCertificate(rand.Reader, leafCert, caCert, &leafPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return err
	}

	// Write the leaf certificate to the specified file
	leafCertFile, err := os.Create(certPath)
	if err != nil {
		return err
	}
	defer leafCertFile.Close()

	if err := pem.Encode(leafCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: leafCertBytes}); err != nil {
		return err
	}

	// Write the leaf private key to the specified file
	leafKeyFile, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer leafKeyFile.Close()

	privBytes, err := x509.MarshalECPrivateKey(leafPrivKey)
	if err != nil {
		return err
	}

	if err := pem.Encode(leafKeyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}); err != nil {
		return err
	}

	return nil
}
