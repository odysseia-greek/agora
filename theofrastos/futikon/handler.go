package futikon

import (
	"fmt"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/thales"
	"gopkg.in/yaml.v3"
)

type TheofratosHandler struct {
	Elastic       elastic.Client
	Namespace     string
	ConfigMapName string
	Kube          *thales.KubeClient
}

func (t *TheofratosHandler) handle(configMap map[string]string) error {
	var config Config

	// Unmarshal roles
	rolesYaml, ok := configMap["roles"]
	if !ok {
		return fmt.Errorf("missing 'roles' key in configMap")
	}
	err := yaml.Unmarshal([]byte(rolesYaml), &config.Roles)
	if err != nil {
		return fmt.Errorf("failed to unmarshal 'roles': %v", err)
	}

	// Unmarshal users
	usersYaml, ok := configMap["users"]
	if !ok {
		return fmt.Errorf("missing 'users' key in configMap")
	}
	err = yaml.Unmarshal([]byte(usersYaml), &config.Users)
	if err != nil {
		return fmt.Errorf("failed to unmarshal 'users': %v", err)
	}

	// Unmarshal policies
	policiesYaml, ok := configMap["policies"]
	if !ok {
		return fmt.Errorf("missing 'policies' key in configMap")
	}
	err = yaml.Unmarshal([]byte(policiesYaml), &config.Policies)
	if err != nil {
		return fmt.Errorf("failed to unmarshal 'policies': %v", err)
	}

	// Create roles
	for roleName, role := range config.Roles {
		if err := t.createElasticRoles(role, roleName); err != nil {
			return fmt.Errorf("failed to create role %s: %v", roleName, err)
		}
	}

	// Create users
	existingUsers, err := t.Elastic.Access().ListUsers()
	if err != nil {
		return fmt.Errorf("failed to list existing users: %w", err)
	}

	userSet := make(map[string]struct{}, len(existingUsers))
	for _, username := range existingUsers {
		userSet[username] = struct{}{}
	}

	for name, user := range config.Users {
		if _, exists := userSet[name]; exists {
			logging.Debug(fmt.Sprintf("user: %s already exists — skipping", name))
			continue
		}

		if err := t.createElasticUser(user, name); err != nil {
			return fmt.Errorf("failed to create user %s: %w", name, err)
		}

		logging.Info(fmt.Sprintf("user: %s created", name))
	}

	// Create policies
	for phaseName, policyPhase := range config.Policies {
		for _, policy := range policyPhase.Rollover {
			if err := t.createILMPolicy(policy.Name, phaseName, policy.MaxAge); err != nil {
				return fmt.Errorf("failed to create ILM policy %s: %v", policy.Name, err)
			}
		}
	}

	return nil
}
