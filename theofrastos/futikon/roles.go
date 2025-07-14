package futikon

import (
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/logging"
)

func (t *TheofratosHandler) createElasticRoles(role RoleMapping, roleName string) error {
	// Create the role for each index
	if len(role.Indices) == 0 || role.Indices[0] == "" {
		logging.Debug(fmt.Sprintf("creating a role for admin with role %s", roleName))
		clusterPerms := role.Role.Cluster

		putRole := models.CreateRoleRequest{
			Cluster:      clusterPerms,
			Indices:      make([]models.Indices, 0),
			Applications: []models.Application{},
			RunAs:        nil,
			Metadata:     models.Metadata{Version: 1},
		}

		roleCreated, err := t.Elastic.Access().CreateRole(roleName, putRole)
		if err != nil {
			return fmt.Errorf("failed to create role %s: %v", roleName, err)
		}

		logging.Info(fmt.Sprintf("role: %s - created: %v", roleName, roleCreated))

		return nil
	} else {
		for _, index := range role.Indices {
			logging.Debug(fmt.Sprintf("creating a role for index %s with role %s", index, roleName))

			// Prepare the Elasticsearch role creation request
			names := []string{index}

			// extra rule that should be set in the configmap
			if roleName == "alias" {
				names = []string{fmt.Sprintf("%s*", index)}
			}

			elasticIndices := []models.Indices{
				{
					Names:      names,
					Privileges: role.Role.Privileges,
					Query:      "",
				},
			}

			clusterPerms := role.Role.Cluster
			if clusterPerms == nil {
				clusterPerms = []string{}
			}

			putRole := models.CreateRoleRequest{
				Cluster:      clusterPerms,
				Indices:      elasticIndices,
				Applications: []models.Application{},
				RunAs:        nil,
				Metadata:     models.Metadata{Version: 1},
			}

			nameInElastic := fmt.Sprintf("%s_%s", index, roleName)
			roleCreated, err := t.Elastic.Access().CreateRole(nameInElastic, putRole)
			if err != nil {
				return fmt.Errorf("failed to create role %s for index %s: %v", roleName, index, err)
			}

			logging.Info(fmt.Sprintf("role: %s - created: %v", nameInElastic, roleCreated))
		}
	}

	return nil
}
