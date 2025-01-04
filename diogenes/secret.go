package diogenes

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"log"
)

const SecretPrefix string = "secret"

func (v *Vault) CreateNewSecret(name string, payload []byte) (bool, error) {
	vaultPath := fmt.Sprintf("%s/%s", v.SecretPath, name)

	_, err := v.Connection.Logical().WriteBytes(vaultPath, payload)
	if err != nil {
		return false, fmt.Errorf("unable to connect to vault: %w", err)
	}

	return true, nil
}

func (v *Vault) GetSecret(name string) (*api.Secret, error) {
	vaultPath := fmt.Sprintf("%s/%s", v.SecretPath, name)

	secret, err := v.Connection.Logical().Read(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read data from vault: %w", err)
	}

	return secret, nil
}

func (v *Vault) ListSecrets() ([]string, error) {
	//curl -k -H "X-Vault-Token: $token" \
	//--request LIST \
	//https://vault:8200/v1/configs/metadata
	vaultPath := fmt.Sprintf("%s/metadata", v.KVSecretPath)
	secret, err := v.Connection.Logical().List(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("unable to list secrets in vault: %w", err)
	}

	// Ensure the returned secret contains data
	if secret == nil || secret.Data == nil {
		return nil, nil // No secrets found
	}

	// Extract the keys from the data
	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format when listing secrets")
	}

	// Convert keys to a slice of strings
	result := make([]string, len(keys))
	for i, key := range keys {
		result[i], ok = key.(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert key to string")
		}
	}

	return result, nil
}

func (v *Vault) DeleteSecret(name string) error {
	vaultPath := fmt.Sprintf("%s/%s", v.SecretPath, name)

	_, err := v.Connection.Logical().Delete(vaultPath)
	if err != nil {
		return fmt.Errorf("unable to delete secret in vault: %w", err)
	}

	return nil
}

// EnableKVSecretsEngine enables the Key-Value (KV) secrets engine in HashiCorp Vault.
//
// Parameters:
//
// namespace: In Vault, a namespace is a way to create a logical grouping or isolation of data within a Vault cluster.
// If you're not using namespaces, you can typically set this to an empty string or ignore it.
//
// configName: This is the name you want to give to your KV (Key-Value) secrets engine. For instance, to create a KV secrets
// engine at the path "configs," pass "configs" as configName.
//
// Returns:
//
// error: If there is an error during the process of enabling the KV secrets engine, an error is returned. Otherwise, nil
// is returned.
//
// Usage example:
//
// // Enable KV secrets engine without using namespaces
// err := EnableKVSecretsEngine("", "configs")
//
//	if err != nil {
//	    log.Printf("Error enabling KV secrets engine: %v", err)
//	}
func (v *Vault) EnableKVSecretsEngine(namespace, configName string) error {
	var path string

	// Set the path based on whether a namespace is provided
	if namespace == "" {
		path = fmt.Sprintf("%s", configName)
	} else {
		path = fmt.Sprintf("%s/%s", namespace, configName)
	}

	options := map[string]string{
		"version": "2", // Set the version to v2
	}

	// Enable the KV secrets engine
	err := v.Connection.Sys().Mount(path, &api.MountInput{
		Type:        "kv",
		Description: "KV secrets engine for odysseia-greek",
		Options:     options,
	})
	if err != nil {
		return fmt.Errorf("failed to enable KV secrets engine: %w", err)
	}

	log.Printf("KV secrets engine enabled at path %s\n", path)

	return nil
}
