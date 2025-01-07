package diogenes

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"log"
)

func (v *Vault) WritePolicy(policyName string, policyContent []byte) error {
	// Check if the policy already exists
	existingPolicy, err := v.ReadPolicy(policyName)
	if err != nil {
		log.Printf("policy %s does not exist, creating...", policyName)
	} else {
		// If the policy already exists, compare the content
		if string(policyContent) == existingPolicy {
			log.Printf("policy %s already exists with the same content, skipping creation", policyName)
			return nil
		}
	}

	path := fmt.Sprintf("sys/policies/acl/%s", policyName)

	data := map[string]interface{}{
		"policy": string(policyContent),
	}

	_, err = v.Connection.Logical().Write(path, data)
	if err != nil {
		return fmt.Errorf("failed to write Vault policy: %w", err)
	}

	return nil
}

func (v *Vault) DeletePolicy(policyName string) (*api.Secret, error) {
	path := fmt.Sprintf("sys/policies/acl/%s", policyName)
	return v.Connection.Logical().Delete(path)
}

func (v *Vault) ListPolicies() ([]string, error) {
	path := "sys/policies/acl"

	secret, err := v.Connection.Logical().List(path)
	if err != nil {
		return nil, fmt.Errorf("unable to list policies in Vault: %w", err)
	}

	// Ensure the returned secret contains data
	if secret == nil || secret.Data == nil {
		return nil, nil // No policies found
	}

	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format when listing policies")
	}

	// Convert keys to a slice of strings
	policies := make([]string, len(keys))
	for i, key := range keys {
		policies[i], ok = key.(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert key to string")
		}
	}

	return policies, nil
}

func (v *Vault) ReadPolicy(policyName string) (string, error) {
	path := fmt.Sprintf("sys/policies/acl/%s", policyName)

	policy, err := v.Connection.Logical().Read(path)
	if err != nil {
		return "", fmt.Errorf("failed to read Vault policy: %w", err)
	}

	if policy == nil || policy.Data == nil {
		return "", fmt.Errorf("policy not found")
	}

	return policy.Data["policy"].(string), nil
}
