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

func (t *TheofratosHandler) Create() {
	health := t.Elastic.Health().Info()
	logging.Info(fmt.Sprintf("status: %v", health.Healthy))

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
	for userName, user := range config.Users {
		if err := t.createElasticUser(user, userName); err != nil {
			return fmt.Errorf("failed to create user %s: %v", userName, err)
		}
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
