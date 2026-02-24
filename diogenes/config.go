package diogenes

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
)

const (
	defaultVault            = "https://vault:8200"
	VAULT                   = "vault"
	defaultRoleName         = "solon"
	EnvVaultService         = "VAULT_SERVICE"
	EnvAuthMethod           = "AUTH_METHOD"
	EnvTLSEnabled           = "VAULT_TLS"
	EnvVaultRole            = "VAULT_ROLE"
	EnvRootTlSDir           = "CERT_ROOT"
	AuthMethodKube          = "kubernetes"
	AuthMethodToken         = "token"
	defaultTLSFileLocation  = "/etc/certs"
	serviceAccountTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)

func getStringFromEnv(envName, defaultValue string) string {
	var value string
	value = os.Getenv(envName)
	if value == "" {
		value = defaultValue
	}

	return value
}

func getBoolFromEnv(envName string) bool {
	var value bool
	envValue := os.Getenv(envName)
	if envValue == "" || envValue == "no" || envValue == "false" {
		value = false
	} else {
		value = true
	}

	return value
}

func CreateVaultClient(healthCheck bool) (Client, error) {
	var vaultClient Client

	vaultAuthMethod := getStringFromEnv(EnvAuthMethod, AuthMethodToken)
	vaultService := getStringFromEnv(EnvVaultService, defaultVault)
	vaultRole := getStringFromEnv(EnvVaultRole, defaultRoleName)
	tlsEnabled := getBoolFromEnv(EnvTLSEnabled)
	rootPath := getStringFromEnv(EnvRootTlSDir, defaultTLSFileLocation)
	secretPath := filepath.Join(rootPath, VAULT)

	var tlsConfig *api.TLSConfig

	if tlsEnabled {
		ca, cert, key, err := resolveTLSFiles(secretPath)
		if err != nil {
			return nil, err
		}
		tlsConfig = CreateTLSConfig(ca, cert, key, secretPath)
	}

	if vaultAuthMethod == AuthMethodKube {
		jwt, err := os.ReadFile(serviceAccountTokenPath)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		vaultToken := string(jwt)

		client, err := CreateVaultClientKubernetes(vaultService, vaultRole, vaultToken, tlsConfig)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		if healthCheck {
			ticks := 120 * time.Second
			tick := 1 * time.Second
			healthy := client.CheckHealthyStatus(ticks, tick)
			if !healthy {
				return nil, fmt.Errorf("error getting healthy status from vault")
			}
		}

		vaultClient = client
	} else {
		client, err := NewVaultClient(vaultService, "", tlsConfig)
		if err != nil {
			return nil, err
		}

		vaultClient = client
	}

	return vaultClient, nil
}

// resolveTLSFiles returns file paths for CA, cert, and key in secretPath.
// It supports both legacy (vault.*) and cert-manager (ca.crt/tls.*) naming.
func resolveTLSFiles(secretPath string) (caPath, certPath, keyPath string, err error) {
	// Priority order: prefer cert-manager if present, else legacy
	caCandidates := []string{
		"ca.crt",   // cert-manager
		"vault.ca", // legacy
	}
	certCandidates := []string{
		"tls.crt",   // cert-manager
		"vault.crt", // legacy
	}
	keyCandidates := []string{
		"tls.key",   // cert-manager
		"vault.key", // legacy
	}

	caPath = firstExisting(secretPath, caCandidates)
	certPath = firstExisting(secretPath, certCandidates)
	keyPath = firstExisting(secretPath, keyCandidates)

	var missing []string
	if caPath == "" {
		missing = append(missing, "CA (ca.crt or vault.ca)")
	}
	if certPath == "" {
		missing = append(missing, "cert (tls.crt or vault.crt)")
	}
	if keyPath == "" {
		missing = append(missing, "key (tls.key or vault.key)")
	}

	if len(missing) > 0 {
		return "", "", "", fmt.Errorf(
			"TLS enabled but missing %s in %s",
			strings.Join(missing, ", "),
			secretPath,
		)
	}

	return caPath, certPath, keyPath, nil
}

func firstExisting(dir string, candidates []string) string {
	for _, name := range candidates {
		p := filepath.Join(dir, name)
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
			return p
		}
	}
	return ""
}
