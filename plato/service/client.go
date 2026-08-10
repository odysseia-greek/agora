package service

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type HttpClient interface {
	Get(u *url.URL, uuid string) (*http.Response, error)
	Post(u *url.URL, body []byte, uuid string) (*http.Response, error)
}

type ClientImpl struct {
	client *http.Client
}

type FakeClientImpl struct {
	responseBodies []string
	codes          []int
	index          int
	mu             sync.Mutex
}

func NewHttpClient(caCert []byte, certs []tls.Certificate) HttpClient {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	if caCert != nil {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		client = http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      caCertPool,
					Certificates: certs,
				},
			},
			Timeout: 10 * time.Second,
		}
	}

	return &ClientImpl{
		client: &client,
	}
}

func NewFakeHttpClient(responseBodies []string, codes []int) HttpClient {
	return &FakeClientImpl{
		responseBodies: responseBodies,
		codes:          codes,
		index:          0,
	}
}

func (c *ClientImpl) Get(u *url.URL, uuid string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(HeaderKey, uuid)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *ClientImpl) Post(u *url.URL, body []byte, uuid string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set(HeaderKey, uuid)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (f *FakeClientImpl) Get(u *url.URL, uuid string) (*http.Response, error) {
	return f.do(http.MethodGet, u, nil, uuid)
}

func (f *FakeClientImpl) Post(u *url.URL, body []byte, uuid string) (*http.Response, error) {
	return f.do(http.MethodPost, u, bytes.NewReader(body), uuid)
}

func (f *FakeClientImpl) do(method string, u *url.URL, body io.Reader, uuid string) (*http.Response, error) {
	if u == nil {
		return nil, fmt.Errorf("create fake %s request: URL is nil", method)
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("create fake %s request: %w", method, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(HeaderKey, uuid)

	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.codes) != len(f.responseBodies) {
		return nil, fmt.Errorf(
			"invalid fake responses: got %d status codes and %d response bodies",
			len(f.codes),
			len(f.responseBodies),
		)
	}
	if f.index >= len(f.codes) {
		return nil, fmt.Errorf("fake response queue exhausted after %d requests", f.index)
	}

	statusCode := f.codes[f.index]
	responseBody := f.responseBodies[f.index]
	f.index++

	return &http.Response{
		StatusCode: statusCode,
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(responseBody)),
		Request:    req,
	}, nil
}
