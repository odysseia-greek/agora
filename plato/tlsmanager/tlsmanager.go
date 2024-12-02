package tlsmanager

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"io"
	"os"
	"path/filepath"
	"time"
)

const DefaultCertRoot string = "/app/config/"

type TLSManager struct {
	currentTLSConfig *tls.Config
	certPath         string
	keyPath          string
	caPath           string
	caPool           *x509.CertPool
	gracePeriod      time.Duration
	lastCertHash     string
	lastKeyHash      string
	cipherSuites     []uint16
	curvePreferences []tls.CurveID
}

func NewTLSManager(hostName, rootPath string, gracePeriod time.Duration) *TLSManager {
	if rootPath == "" {
		rootPath = DefaultCertRoot
	}
	certPath := filepath.Join(rootPath, hostName, "tls.crt")
	keyPath := filepath.Join(rootPath, hostName, "tls.key")
	caPath := filepath.Join(rootPath, hostName, "tls.pem")

	// Load CA certificate
	caPool := x509.NewCertPool()
	caFromFile, err := os.ReadFile(caPath)
	if err != nil {
		logging.Error(fmt.Sprintf("Failed to read CA certificate: %v", err))
		return nil
	}
	if !caPool.AppendCertsFromPEM(caFromFile) {
		logging.Error("Failed to append CA certificate")
		return nil
	}

	// Define cipher suites and curve preferences
	cipherSuites := []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	}
	curvePreferences := []tls.CurveID{
		tls.CurveP521,
		tls.CurveP384,
		tls.CurveP256,
	}

	return &TLSManager{
		certPath:         certPath,
		keyPath:          keyPath,
		caPath:           caPath,
		caPool:           caPool,
		gracePeriod:      gracePeriod,
		cipherSuites:     cipherSuites,
		curvePreferences: curvePreferences,
	}
}

// GetTLSConfig returns the current TLS configuration for the server.
func (tm *TLSManager) GetTLSConfig() *tls.Config {
	return tm.currentTLSConfig
}

// LoadCertificates loads and updates the TLS configuration.
func (tm *TLSManager) LoadCertificates() error {
	// Load the new certificate
	cert, err := tls.LoadX509KeyPair(tm.certPath, tm.keyPath)
	if err != nil {
		return fmt.Errorf("failed to load certificate: %w", err)
	}

	// Initialize the TLS configuration if it's the first run
	if tm.currentTLSConfig == nil {
		tm.currentTLSConfig = &tls.Config{
			Certificates:     []tls.Certificate{cert},
			ClientCAs:        tm.caPool,
			MinVersion:       tls.VersionTLS12,
			CipherSuites:     tm.cipherSuites,
			CurvePreferences: tm.curvePreferences,
			GetCertificate: func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
				for _, cert := range tm.currentTLSConfig.Certificates {
					return &cert, nil
				}
				return nil, fmt.Errorf("no certificate available")
			},
		}
		logging.System("TLS configuration initialized with the first certificate.")
		return nil
	}

	// Clone the existing configuration
	oldTLSConfig := tm.currentTLSConfig
	newCertificates := append(oldTLSConfig.Certificates, cert)

	// Update the current TLS configuration
	tm.currentTLSConfig = &tls.Config{
		Certificates:     newCertificates,
		ClientCAs:        tm.caPool,
		MinVersion:       tls.VersionTLS12,
		CipherSuites:     tm.cipherSuites,
		CurvePreferences: tm.curvePreferences,
		GetCertificate: func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			for _, cert := range newCertificates {
				return &cert, nil
			}
			return nil, fmt.Errorf("no certificate available")
		},
	}

	// Schedule removal of old certificates after the grace period
	if len(oldTLSConfig.Certificates) > 0 {
		go tm.removeOldCertificate(oldTLSConfig.Certificates[0])
	}

	return nil
}

// WatchCertificates starts monitoring the certificate files for changes.
func (tm *TLSManager) WatchCertificates(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)

			certChanged := tm.checkFileContentChange(tm.certPath, &tm.lastCertHash)
			keyChanged := tm.checkFileContentChange(tm.keyPath, &tm.lastKeyHash)

			if certChanged || keyChanged {
				logging.Debug("Certificate content changed; reloading...")
				if err := tm.LoadCertificates(); err != nil {
					logging.Error(fmt.Sprintf("Error reloading certificates: %v\n", err))
				}
			}
		}
	}()
}

// Helper function to remove old certificates after grace period
func (tm *TLSManager) removeOldCertificate(oldCert tls.Certificate) {
	time.Sleep(tm.gracePeriod)

	var remainingCertificates []tls.Certificate
	for _, cert := range tm.currentTLSConfig.Certificates {
		if !certEqual(cert, oldCert) {
			remainingCertificates = append(remainingCertificates, cert)
		}
	}

	tm.currentTLSConfig = &tls.Config{
		Certificates: remainingCertificates,
		GetCertificate: func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			for _, cert := range remainingCertificates {
				return &cert, nil
			}
			return nil, fmt.Errorf("no certificate available")
		},
	}

	logging.Debug("Grace period ended; old certificate removed.")
}

// Check if the file content has changed
func (tm *TLSManager) checkFileContentChange(filePath string, lastHash *string) bool {
	currentHash, err := hashFile(filePath)
	if err != nil {
		logging.Error(fmt.Sprintf("Failed to hash file %s: %v", filePath, err))
		return false
	}

	if currentHash != *lastHash {
		*lastHash = currentHash
		return true
	}

	return false
}

// Hash a file's content
func hashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

// Helper to compare certificates
func certEqual(cert1, cert2 tls.Certificate) bool {
	return string(cert1.Certificate[0]) == string(cert2.Certificate[0])
}
