package diogenes

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/hashicorp/vault/api"
)

//go:embed eratosthenes/*.json
var fixtureFiles embed.FS

type MockVaultTransport struct {
	Responses      []*http.Response
	ResponseIdx    int
	mu             sync.Mutex
	responseBodies [][]byte
}

func (t *MockVaultTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.Responses) == 0 {
		return nil, fmt.Errorf("mock Vault transport has no responses")
	}

	template := t.Responses[t.ResponseIdx]
	t.ResponseIdx = (t.ResponseIdx + 1) % len(t.Responses)
	if template == nil {
		return nil, fmt.Errorf("mock Vault response is nil")
	}
	if len(t.responseBodies) != len(t.Responses) {
		t.responseBodies = make([][]byte, len(t.Responses))
	}
	responseIndex := (t.ResponseIdx - 1 + len(t.Responses)) % len(t.Responses)
	if t.responseBodies[responseIndex] == nil {
		if template.Body == nil {
			return nil, fmt.Errorf("mock Vault response body is nil")
		}
		body, err := io.ReadAll(template.Body)
		if err != nil {
			return nil, fmt.Errorf("read mock Vault response body: %w", err)
		}
		t.responseBodies[responseIndex] = body
	}

	response := new(http.Response)
	*response = *template
	response.Body = io.NopCloser(bytes.NewReader(t.responseBodies[responseIndex]))
	response.Request = req
	return response, nil
}

func CreateMockVaultClient(fixtures []string, statusCode int) (Client, error) {
	if statusCode < 100 || statusCode > 599 {
		return nil, fmt.Errorf("invalid mock Vault HTTP status code %d", statusCode)
	}
	if len(fixtures) == 0 {
		return nil, fmt.Errorf("at least one mock Vault fixture is required")
	}

	responses := make([]*http.Response, 0, len(fixtures))
	for _, fixtureName := range fixtures {
		fixturePath := fmt.Sprintf("eratosthenes/%s.json", fixtureName)
		body, err := fixtureFiles.ReadFile(fixturePath)
		if err != nil {
			return nil, fmt.Errorf("read mock Vault fixture %q: %w", fixtureName, err)
		}

		responses = append(responses, &http.Response{
			StatusCode: statusCode,
			Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(bytes.NewReader(body)),
		})
	}

	config := api.DefaultConfig()
	config.Address = "http://example.com"
	config.HttpClient.Transport = &MockVaultTransport{Responses: responses}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create mock Vault client: %w", err)
	}

	return &Vault{
		Connection:   client,
		SecretPath:   defaultKVSecretData,
		KVSecretPath: defaultKVSecretPath,
	}, nil
}
