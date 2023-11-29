package diogenes

import (
	"fmt"
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

func (v *Vault) ReadPolicy(policyName string) (string, error) {
	path := fmt.Sprintf("sys/policies/acl/%s", policyName)

	secret, err := v.Connection.Logical().Read(path)
	if err != nil {
		return "", fmt.Errorf("failed to read Vault policy: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return "", fmt.Errorf("policy not found")
	}

	return secret.Data["policy"].(string), nil
}
