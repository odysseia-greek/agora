package diogenes

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/api"
)

const (
	defaultMountPath = "kubernetes"
)

func (v *Vault) KubernetesAuthMethod(ctx context.Context, role, serviceAccountName, namespace, kubeHost string) error {
	return v.kubernetesAuthMethod(ctx, role, serviceAccountName, namespace, kubeHost, "")
}

func (v *Vault) KubernetesAuthMethodWithAudience(ctx context.Context, role, serviceAccountName, namespace, kubeHost, audience string) error {
	if audience == "" {
		return fmt.Errorf("Kubernetes auth audience cannot be empty")
	}

	return v.kubernetesAuthMethod(ctx, role, serviceAccountName, namespace, kubeHost, audience)
}

func (v *Vault) kubernetesAuthMethod(ctx context.Context, role, serviceAccountName, namespace, kubeHost, audience string) error {
	path := fmt.Sprintf("auth/%s", defaultMountPath)

	err := v.Connection.Sys().EnableAuthWithOptionsWithContext(ctx, defaultMountPath, &api.EnableAuthOptions{
		Type: defaultMountPath,
		Options: map[string]string{
			"mount_path":        defaultMountPath,
			"default_lease_ttl": "0",
			"max_lease_ttl":     "0",
		},
		Description: "Kubernetes authentication",
	})
	if err != nil {
		return err
	}

	configPath := fmt.Sprintf("%s/config", path)
	configData := map[string]interface{}{
		"kubernetes_host":        kubeHost,
		"disable_iss_validation": "true",
	}

	_, err = v.Connection.Logical().WriteWithContext(ctx, configPath, configData)
	if err != nil {
		return err
	}

	// Configure a role for the Kubernetes auth method
	rolePath := fmt.Sprintf("%s/role/%s", path, role)
	roleData := map[string]interface{}{
		"bound_service_account_names":      fmt.Sprintf("vault,%s", serviceAccountName),
		"bound_service_account_namespaces": namespace, // Replace with the namespace
		"policies":                         role,      // Replace with the desired policies
	}
	if audience != "" {
		roleData["audience"] = audience
	}

	_, err = v.Connection.Logical().WriteWithContext(ctx, rolePath, roleData)
	if err != nil {
		return err
	}

	return nil
}
