package diogenes

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"
	"log"
	"time"
)

type Client interface {
	CheckHealthyStatus(ticks, tick time.Duration) bool
	Health() (bool, error)
	CreateOneTimeToken(policy []string) (string, error)
	CreateNewSecret(name string, payload []byte) (bool, error)
	GetSecret(name string) (*api.Secret, error)
	DeleteSecret(name string) error
	ListSecrets() ([]string, error)
	SetOnetimeToken(token string)
	LoginWithRootToken(rootToken string) error
	GetCurrentToken() string
	Unseal(keys []string) (bool, error)
	AutoUnsealGCP(keyRing, cryptoKey, location string, keys []string) (bool, error)
	Status() (*api.SealStatusResponse, error)
	Initialize(shares, threshold int) (*api.InitResponse, error)
	InitializeAutoUnseal(shares, threshold int) (*api.InitResponse, error)
	EnableKVSecretsEngine(namespace, configName string) error
	WritePolicy(policyName string, policyContent []byte) error
	DeletePolicy(policyName string) (*api.Secret, error)
	ReadPolicy(policyName string) (string, error)
	KubernetesAuthMethod(role, serviceAccountName, namespace, kubeHost string) error
	RaftJoin(leaderAddress string, cert, key, ca []byte) (*api.RaftJoinResponse, error)
	Leader() (*api.LeaderResponse, error)
}

type Vault struct {
	SecretPath   string
	KVSecretPath string
	Connection   *api.Client
}

const (
	defaultKVSecretPath string = "configs"
	defaultKVSecretData string = "configs/data"
	fixtureSecretName   string = "isitsecretisitsafe"
)

func NewVaultClient(address, token string, tlsConfig *api.TLSConfig) (Client, error) {
	config := api.Config{
		Address: address,
	}

	log.Print(tlsConfig)

	if tlsConfig != nil {
		err := config.ConfigureTLS(tlsConfig)
		if err != nil {
			return nil, err
		}
	}

	client, err := api.NewClient(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	if token != "" {
		client.SetToken(token)
	}

	return &Vault{Connection: client, SecretPath: defaultKVSecretData, KVSecretPath: defaultKVSecretPath}, nil
}

func CreateVaultClientKubernetes(address, vaultRole, jwt string, tlsConfig *api.TLSConfig) (Client, error) {
	config := api.Config{
		Address: address,
	}

	if tlsConfig != nil {
		err := config.ConfigureTLS(tlsConfig)
		if err != nil {
			return nil, err
		}
	}

	client, err := api.NewClient(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	k8sAuth, err := auth.NewKubernetesAuth(
		vaultRole,
		auth.WithServiceAccountToken(jwt),
	)

	// log in to Vault's Kubernetes auth method
	resp, err := client.Auth().Login(context.Background(), k8sAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with Kubernetes auth: %w", err)
	}
	if resp == nil || resp.Auth == nil || resp.Auth.ClientToken == "" {
		return nil, fmt.Errorf("login response did not return client token")
	}

	client.SetToken(resp.Auth.ClientToken)

	return &Vault{Connection: client, SecretPath: defaultKVSecretData, KVSecretPath: defaultKVSecretPath}, nil
}

func CreateTLSConfig(ca, cert, key, caPath string) *api.TLSConfig {
	return &api.TLSConfig{
		CAPath:     caPath,
		CACert:     ca,
		ClientCert: cert,
		ClientKey:  key,
		Insecure:   false,
	}
}
