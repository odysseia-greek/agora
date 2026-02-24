package aristoteles

import (
	"bytes"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	fixtures = make(map[string][]byte)
)

func AddRawFixture(name string, data []byte) {
	fixtures[name] = data
}

func init() {
	_, callingFile, _, _ := runtime.Caller(0)
	callingDir := filepath.Dir(callingFile)
	dirParts := strings.Split(callingDir, string(os.PathSeparator))
	var aristotelesElasticPath []string
	for i, part := range dirParts {
		if strings.Contains(part, "aristoteles") {
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
		f, err := os.ReadFile(fpath)
		if err != nil {
			panic(fmt.Sprintf("Cannot read fixture file: %s", err))
		}
		fixtures[filepath.Base(fpath)] = f
	}
}

func fixture(name string) io.ReadCloser {
	data, ok := fixtures[name]
	if !ok {
		panic(fmt.Sprintf("Fixture not found: %s", name))
	}
	return io.NopCloser(bytes.NewReader(data))
}

type MockTransport struct {
	Responses   []*http.Response
	ResponseIdx int
	RoundTripFn func(req *http.Request) (*http.Response, error)
}

func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	response := t.Responses[t.ResponseIdx]
	t.ResponseIdx = (t.ResponseIdx + 1) % len(t.Responses)
	return response, nil
}

func CreateMockClient(fixtureFiles []string, statusCode int) (*elasticsearch.Client, error) {
	mockCode := toHTTPStatusCode(statusCode)
	statusCodes := make([]int, len(fixtureFiles))
	for i := range fixtureFiles {
		statusCodes[i] = mockCode
	}
	return CreateMockClientWithStatusCodes(fixtureFiles, statusCodes)
}

func CreateMockClientWithStatusCodes(fixtureFiles []string, statusCodes []int) (*elasticsearch.Client, error) {
	if len(fixtureFiles) != len(statusCodes) {
		return nil, fmt.Errorf("fixtureFiles and statusCodes must have equal length")
	}

	var responses []*http.Response
	for i, fix := range fixtureFiles {
		if !strings.Contains(fix, ".json") {
			fix = fmt.Sprintf("%s.json", fix)
		}
		body := fixture(fix)
		response := &http.Response{
			StatusCode: toHTTPStatusCode(statusCodes[i]),
			Body:       body,
			Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
		}
		responses = append(responses, response)
	}

	mockTrans := MockTransport{
		Responses:   responses,
		ResponseIdx: 0,
	}
	mockTrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mockTrans.RoundTrip(req) }

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Transport: &mockTrans,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func toHTTPStatusCode(statusCode int) int {
	switch statusCode {
	case 200:
		return http.StatusOK
	case 401:
		return http.StatusUnauthorized
	case 403:
		return http.StatusForbidden
	case 404:
		return http.StatusNotFound
	case 500:
		return http.StatusInternalServerError
	case 502:
		return http.StatusBadGateway
	default:
		return http.StatusOK
	}
}
