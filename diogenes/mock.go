package diogenes

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/vault/api"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	fixtures = make(map[string]io.ReadCloser)
)

func init() {
	_, callingFile, _, _ := runtime.Caller(0)
	callingDir := filepath.Dir(callingFile)
	dirParts := strings.Split(callingDir, string(os.PathSeparator))
	var aristotelesElasticPath []string
	for i, part := range dirParts {
		if strings.Contains(part, "diogenes") {
			aristotelesElasticPath = dirParts[0 : i+1]
		}
	}
	l := "/"
	for _, path := range aristotelesElasticPath {
		l = filepath.Join(l, path)
	}
	eratosthenesDir := filepath.Join(l, "eratosthenes", "*.json")
	fixtureFiles, err := filepath.Glob(eratosthenesDir)
	if err != nil {
		panic(fmt.Sprintf("Cannot glob fixture files: %s", err))
	}

	for _, fpath := range fixtureFiles {
		f, err := ioutil.ReadFile(fpath)
		if err != nil {
			panic(fmt.Sprintf("Cannot read fixture file: %s", err))
		}
		fixtures[filepath.Base(fpath)] = ioutil.NopCloser(bytes.NewReader(f))
	}
}

func fixture(fname string) io.ReadCloser {
	out := new(bytes.Buffer)
	b1 := bytes.NewBuffer([]byte{})
	b2 := bytes.NewBuffer([]byte{})
	tr := io.TeeReader(fixtures[fname], b1)

	defer func() { fixtures[fname] = ioutil.NopCloser(b1) }()
	io.Copy(b2, tr)
	out.ReadFrom(b2)

	return ioutil.NopCloser(out)
}

type MockVaultTransport struct {
	Responses   []*http.Response
	ResponseIdx int
}

func (t *MockVaultTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	response := t.Responses[t.ResponseIdx]
	t.ResponseIdx = (t.ResponseIdx + 1) % len(t.Responses)
	return response, nil
}

func CreateMockVaultClient(fixtureFiles []string, statusCode int) (Client, error) {
	mockCode := 500
	switch statusCode {
	case 200:
		mockCode = http.StatusOK
	case 404:
		mockCode = http.StatusNotFound
	case 500:
		mockCode = http.StatusInternalServerError
	case 502:
		mockCode = http.StatusBadGateway
	default:
		mockCode = 200
	}

	var responses []*http.Response
	for _, fix := range fixtureFiles {
		body := fixture(fmt.Sprintf("%s.json", fix))
		response := &http.Response{
			StatusCode: mockCode,
			Body:       body,
		}
		responses = append(responses, response)
	}

	mockTrans := &MockVaultTransport{
		Responses:   responses,
		ResponseIdx: 0,
	}

	config := api.DefaultConfig()
	config.Address = "http://example.com" // Replace with your Vault server address
	config.HttpClient.Transport = mockTrans

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Vault{Connection: client, SecretPath: defaultPath}, nil
}
