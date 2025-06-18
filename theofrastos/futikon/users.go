package futikon

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/generator"
	"github.com/odysseia-greek/agora/plato/logging"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"time"
)

func (t *TheofratosHandler) createElasticUser(user UserMapping, username string) error {
	logging.Debug(fmt.Sprintf("creating a user %s with role %s", username, user.Role))

	password, err := generator.RandomPassword(24)
	if err != nil {
		return err
	}

	secretName := fmt.Sprintf("%s-elastic", user)
	secretData := map[string][]byte{
		"user":     []byte(username),
		"password": []byte(password),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	secretExists := true
	_, err = t.Kube.CoreV1().Secrets(t.Namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			secretExists = false
		}
	}

	if secretExists {
		logging.Info(fmt.Sprintf("secret %s already exists", secretName))
		return nil
	}

	logging.Info(fmt.Sprintf("secret %s does not exist", secretName))

	scr := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		Immutable:  nil,
		Data:       secretData,
		StringData: nil,
		Type:       corev1.SecretTypeOpaque,
	}
	_, err = t.Kube.CoreV1().Secrets(t.Namespace).Create(ctx, scr, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	var index string
	roleName := user.Role
	switch username {
	case config.DefaultTracingName:
		index = config.TracingElasticIndex
		roleName = fmt.Sprintf("%s_%s", index, user.Role)
	case config.DefaultMetricsName:
		index = config.MetricsElasticIndex
		roleName = fmt.Sprintf("%s_%s", index, user.Role)
	}

	putUser := models.CreateUserRequest{
		Password: password,
		Roles:    []string{roleName},
		FullName: username,
		Email:    fmt.Sprintf("%s@odysseia-greek.com", username),
		Metadata: &models.Metadata{Version: 1},
	}

	userCreated, err := t.Elastic.Access().CreateUser(username, putUser)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("user %s created: %v in elasticSearch", username, userCreated))

	return nil
}
