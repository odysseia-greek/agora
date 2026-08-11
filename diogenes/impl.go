package diogenes

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"
)

type Client interface {
	CheckHealthyStatus(ctx context.Context, ticks, tick time.Duration) bool
	Health(ctx context.Context) (bool, error)
	CreateOneTimeToken(ctx context.Context, policy []string) (string, error)
	ValidateToken(ctx context.Context, oneTimeToken string) (map[string]interface{}, error)
	CreateNewSecret(ctx context.Context, name string, payload []byte) (bool, error)
	GetSecret(ctx context.Context, name string) (*api.Secret, error)
	RemoveSecret(ctx context.Context, name string) error
	DeleteSecret(ctx context.Context, name string) error
	ListSecrets(ctx context.Context) ([]string, error)
	SetOnetimeToken(token string)
	LoginWithRootToken(rootToken string) error
	GetCurrentToken() string
	Unseal(ctx context.Context, keys []string) (bool, error)
	AutoUnsealGCP(ctx context.Context, keyRing, cryptoKey, location string, keys []string) (bool, error)
	Status(ctx context.Context) (*api.SealStatusResponse, error)
	Initialize(ctx context.Context, shares, threshold int) (*api.InitResponse, error)
	InitializeAutoUnseal(ctx context.Context, shares, threshold int) (*api.InitResponse, error)
	EnableKVSecretsEngine(ctx context.Context, namespace, configName string) error
	WritePolicy(ctx context.Context, policyName string, policyContent []byte) error
	DeletePolicy(ctx context.Context, policyName string) (*api.Secret, error)
	ReadPolicy(ctx context.Context, policyName string) (string, error)
	ListPolicies(ctx context.Context) ([]string, error)
	KubernetesAuthMethod(ctx context.Context, role, serviceAccountName, namespace, kubeHost string) error
	KubernetesAuthMethodWithAudience(ctx context.Context, role, serviceAccountName, namespace, kubeHost, audience string) error
	RaftJoin(ctx context.Context, leaderAddress string, cert, key, ca []byte) (*api.RaftJoinResponse, error)
	Leader(ctx context.Context) (*api.LeaderResponse, error)
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

func CreateVaultClientKubernetes(ctx context.Context, address, vaultRole, jwt string, tlsConfig *api.TLSConfig) (Client, error) {
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
	if err != nil {
		return nil, fmt.Errorf("create Kubernetes auth method: %w", err)
	}

	// log in to Vault's Kubernetes auth method
	resp, err := client.Auth().Login(ctx, k8sAuth)
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
