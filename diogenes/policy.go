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
