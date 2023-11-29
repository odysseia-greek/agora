package aristoteles

import (
	"bytes"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"io"
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
		fixtures[filepath.Base(fpath)] = io.NopCloser(bytes.NewReader(f))
	}
}

func fixture(fname string) io.ReadCloser {
	out := new(bytes.Buffer)
	b1 := bytes.NewBuffer([]byte{})
	b2 := bytes.NewBuffer([]byte{})
	tr := io.TeeReader(fixtures[fname], b1)

	defer func() { fixtures[fname] = io.NopCloser(b1) }()
	io.Copy(b2, tr)
	out.ReadFrom(b2)

	return io.NopCloser(out)
}

type MockTransport struct {
	Responses   []*http.Response
	ResponseIdx int
	RoundTripFn func(req *http.Request) (*http.Response, error)
}

//func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
//	return t.RoundTripFn(req)
//}

func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	response := t.Responses[t.ResponseIdx]
	t.ResponseIdx = (t.ResponseIdx + 1) % len(t.Responses)
	return response, nil
}

func CreateMockClient(fixtureFiles []string, statusCode int) (*elasticsearch.Client, error) {
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
