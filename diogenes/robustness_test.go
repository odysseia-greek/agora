package diogenes

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateVaultClientRejectsUnknownAuthMethod(t *testing.T) {
	t.Setenv(EnvAuthMethod, "kubernets")

	client, err := CreateVaultClient(context.Background(), false)

	assert.Nil(t, client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `unsupported AUTH_METHOD "kubernets"`)
}

func TestVaultKubernetesTokenPath(t *testing.T) {
	t.Run("defaults to the standard Kubernetes token", func(t *testing.T) {
		t.Setenv(EnvVaultKubernetesTokenPath, "")

		path := getStringFromEnv(EnvVaultKubernetesTokenPath, defaultServiceAccountTokenPath)

		assert.Equal(t, defaultServiceAccountTokenPath, path)
	})

	t.Run("uses a dedicated projected Vault token", func(t *testing.T) {
		const projectedTokenPath = "/var/run/secrets/vault/token"
		t.Setenv(EnvVaultKubernetesTokenPath, projectedTokenPath)

		path := getStringFromEnv(EnvVaultKubernetesTokenPath, defaultServiceAccountTokenPath)

		assert.Equal(t, projectedTokenPath, path)
	})
}

func TestCheckHealthyStatusValidatesDurations(t *testing.T) {
	client, err := CreateMockVaultClient([]string{"health"}, http.StatusOK)
	require.NoError(t, err)

	assert.False(t, client.CheckHealthyStatus(context.Background(), 0, time.Millisecond))
	assert.False(t, client.CheckHealthyStatus(context.Background(), time.Millisecond, 0))
}

func TestCheckHealthyStatusRetriesUntilHealthy(t *testing.T) {
	response := func(status int, body string) *http.Response {
		return &http.Response{
			StatusCode: status,
			Status:     http.StatusText(status),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(bytes.NewBufferString(body)),
		}
	}
	transport := &MockVaultTransport{Responses: []*http.Response{
		response(http.StatusServiceUnavailable, `{"errors":["not ready"]}`),
		response(http.StatusOK, `{}`),
	}}
	config := api.DefaultConfig()
	config.Address = "http://example.com"
	config.HttpClient.Transport = transport
	connection, err := api.NewClient(config)
	require.NoError(t, err)
	client := &Vault{Connection: connection}

	assert.True(t, client.CheckHealthyStatus(context.Background(), 250*time.Millisecond, time.Millisecond))
	assert.Equal(t, 0, transport.ResponseIdx)
}

func TestCheckHealthyStatusStopsWhenContextIsCancelled(t *testing.T) {
	client, err := CreateMockVaultClient([]string{"health"}, http.StatusOK)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	started := time.Now()
	healthy := client.CheckHealthyStatus(ctx, time.Minute, time.Second)

	assert.False(t, healthy)
	assert.Less(t, time.Since(started), 100*time.Millisecond)
}

func TestMockVaultClientValidation(t *testing.T) {
	t.Run("requires fixtures", func(t *testing.T) {
		client, err := CreateMockVaultClient(nil, http.StatusOK)
		assert.Nil(t, client)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "at least one")
	})

	t.Run("validates status code", func(t *testing.T) {
		client, err := CreateMockVaultClient([]string{"health"}, 42)
		assert.Nil(t, client)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid mock Vault HTTP status")
	})

	t.Run("reports missing fixture", func(t *testing.T) {
		client, err := CreateMockVaultClient([]string{"does-not-exist"}, http.StatusOK)
		assert.Nil(t, client)
		require.Error(t, err)
		assert.Contains(t, err.Error(), `fixture "does-not-exist"`)
	})
}

func TestMockVaultClientIsSafeForConcurrentUse(t *testing.T) {
	client, err := CreateMockVaultClient([]string{"health"}, http.StatusOK)
	require.NoError(t, err)

	const requests = 20
	errors := make(chan error, requests)
	var waitGroup sync.WaitGroup
	for range requests {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			_, err := client.Health(context.Background())
			errors <- err
		}()
	}

	waitGroup.Wait()
	close(errors)
	for err := range errors {
		assert.NoError(t, err)
	}
}

func TestReadPolicyValidatesResponse(t *testing.T) {
	tests := []struct {
		name       string
		fixture    string
		errorMatch string
	}{
		{name: "missing policy", fixture: "policy-missing", errorMatch: "does not contain a policy field"},
		{name: "invalid policy", fixture: "policy-invalid", errorMatch: "non-string policy field"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := CreateMockVaultClient([]string{test.fixture}, http.StatusOK)
			require.NoError(t, err)

			policy, err := client.ReadPolicy(context.Background(), "test")
			assert.Empty(t, policy)
			require.Error(t, err)
			assert.Contains(t, err.Error(), test.errorMatch)
		})
	}
}

func TestCreateOneTimeTokenValidatesResponse(t *testing.T) {
	client, err := CreateMockVaultClient([]string{"empty"}, http.StatusOK)
	require.NoError(t, err)

	token, err := client.CreateOneTimeToken(context.Background(), []string{"test"})

	assert.Empty(t, token)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "did not contain a client token")
}
