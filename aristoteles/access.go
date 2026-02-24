package aristoteles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/odysseia-greek/agora/aristoteles/models"
)

type AccessImpl struct {
	es *elasticsearch.Client
}

func NewAccessImpl(suppliedClient *elasticsearch.Client) (*AccessImpl, error) {
	if suppliedClient == nil {
		return nil, fmt.Errorf("cannot create interface with empty client")
	}
	return &AccessImpl{es: suppliedClient}, nil
}

func (a *AccessImpl) CreateRole(name string, roleRequest models.CreateRoleRequest) (bool, error) {
	jsonRole, err := roleRequest.Marshal()
	if err != nil {
		return false, err
	}
	buffer := bytes.NewBuffer(jsonRole)
	res, err := a.es.Security.PutRole(name, buffer)
	if err != nil {
		return false, err
	}

	if res.IsError() {
		return false, fmt.Errorf("%s: %s", errorMessage, res.Status())
	}

	return true, nil
}

func (a *AccessImpl) CreateUser(name string, userCreation models.CreateUserRequest) (bool, error) {
	jsonUser, err := userCreation.Marshal()
	if err != nil {
		return false, err
	}
	buffer := bytes.NewBuffer(jsonUser)
	res, err := a.es.Security.PutUser(name, buffer)
	if err != nil {
		return false, err
	}

	if res.IsError() {
		return false, fmt.Errorf("%s: %s", errorMessage, res.Status())
	}

	return true, nil
}

func (a *AccessImpl) ListUsers() ([]string, error) {
	// Use the Security API to get all users
	res, err := a.es.Security.GetUser()
	if err != nil {
		return nil, fmt.Errorf("error fetching users from Elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("failed to list users: %s", res.Status())
	}

	// Parse the response into a map
	var users map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// Extract user names from the map
	var userList []string
	for user := range users {
		userList = append(userList, user)
	}

	return userList, nil
}

func (a *AccessImpl) DeleteUser(name string) (bool, error) {
	// Use the Security API to delete a user
	res, err := a.es.Security.DeleteUser(name)
	if err != nil {
		return false, fmt.Errorf("error deleting user %s: %w", name, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return false, fmt.Errorf("failed to delete user %s: %s", name, res.Status())
	}

	return true, nil
}
