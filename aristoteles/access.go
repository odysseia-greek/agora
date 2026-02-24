package aristoteles

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

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
	return a.CreateRoleWithContext(context.Background(), name, roleRequest)
}

func (a *AccessImpl) CreateRoleWithContext(ctx context.Context, name string, roleRequest models.CreateRoleRequest) (bool, error) {
	jsonRole, err := roleRequest.Marshal()
	if err != nil {
		return false, err
	}
	buffer := bytes.NewBuffer(jsonRole)
	res, err := a.es.Security.PutRole(name, buffer, a.es.Security.PutRole.WithContext(ctx))
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return false, newElasticError("create role", res)
	}

	return true, nil
}

func (a *AccessImpl) CreateUser(name string, userCreation models.CreateUserRequest) (bool, error) {
	return a.CreateUserWithContext(context.Background(), name, userCreation)
}

func (a *AccessImpl) CreateUserWithContext(ctx context.Context, name string, userCreation models.CreateUserRequest) (bool, error) {
	jsonUser, err := userCreation.Marshal()
	if err != nil {
		return false, err
	}
	buffer := bytes.NewBuffer(jsonUser)
	res, err := a.es.Security.PutUser(name, buffer, a.es.Security.PutUser.WithContext(ctx))
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return false, newElasticError("create user", res)
	}

	return true, nil
}

func (a *AccessImpl) ListUsers() ([]string, error) {
	return a.ListUsersWithContext(context.Background())
}

func (a *AccessImpl) ListUsersWithContext(ctx context.Context) ([]string, error) {
	res, err := a.es.Security.GetUser(a.es.Security.GetUser.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("error fetching users from Elasticsearch: %w", err)
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, newElasticErrorFromBody("list users", res, respBody)
	}

	var users map[string]interface{}
	if err := json.Unmarshal(respBody, &users); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	userList := make([]string, 0, len(users))
	for user := range users {
		userList = append(userList, user)
	}

	return userList, nil
}

func (a *AccessImpl) DeleteUser(name string) (bool, error) {
	return a.DeleteUserWithContext(context.Background(), name)
}

func (a *AccessImpl) DeleteUserWithContext(ctx context.Context, name string) (bool, error) {
	res, err := a.es.Security.DeleteUser(name, a.es.Security.DeleteUser.WithContext(ctx))
	if err != nil {
		return false, fmt.Errorf("error deleting user %s: %w", name, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return false, newElasticError("delete user", res)
	}

	return true, nil
}
