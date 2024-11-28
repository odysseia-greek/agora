package diogenes

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"log"
	"os"
	"path/filepath"
	"time"
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

	if debugMode {
		log.Printf("vaultAuthMethod set to %s", vaultAuthMethod)
		log.Printf("secretPath set to %s", secretPath)
		log.Printf("tlsEnabled set to %v", tlsEnabled)
	}

	var tlsConfig *api.TLSConfig

	if tlsEnabled {
		ca := fmt.Sprintf("%s/vault.ca", secretPath)
		cert := fmt.Sprintf("%s/vault.crt", secretPath)
		key := fmt.Sprintf("%s/vault.key", secretPath)

		if debugMode {
			log.Print(ca)
			log.Print(cert)
			log.Print(key)
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
