package service

import (
	"io"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFakeHttpClientResponse(t *testing.T) {
	client := NewFakeHttpClient([]string{`{"result":"ok"}`}, []int{http.StatusCreated})
	endpoint := &url.URL{Scheme: "https", Host: "solon.example", Path: "/solon/v1/register"}

	response, err := client.Post(endpoint, []byte(`{"name":"test"}`), "trace-id")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, response.Body.Close()) })

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, "201 Created", response.Status)
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	require.NotNil(t, response.Request)
	assert.Equal(t, http.MethodPost, response.Request.Method)
	assert.Equal(t, endpoint.String(), response.Request.URL.String())
	assert.Equal(t, "trace-id", response.Request.Header.Get(HeaderKey))

	requestBody, err := io.ReadAll(response.Request.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"name":"test"}`, string(requestBody))
}

func TestFakeHttpClientRejectsInvalidFixtures(t *testing.T) {
	client := NewFakeHttpClient([]string{"one", "two"}, []int{http.StatusOK})

	response, err := client.Get(&url.URL{Scheme: "https", Host: "solon.example"}, "")

	assert.Nil(t, response)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "1 status codes and 2 response bodies")
}

func TestFakeHttpClientReturnsQueueExhausted(t *testing.T) {
	client := NewFakeHttpClient([]string{"ok"}, []int{http.StatusOK})
	endpoint := &url.URL{Scheme: "https", Host: "solon.example"}

	first, err := client.Get(endpoint, "")
	require.NoError(t, err)
	require.NoError(t, first.Body.Close())

	second, err := client.Get(endpoint, "")
	assert.Nil(t, second)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "queue exhausted after 1 requests")
}

func TestFakeHttpClientIsSafeForConcurrentUse(t *testing.T) {
	const requests = 20
	bodies := make([]string, requests)
	codes := make([]int, requests)
	for i := range requests {
		bodies[i] = "ok"
		codes[i] = http.StatusOK
	}

	client := NewFakeHttpClient(bodies, codes)
	endpoint := &url.URL{Scheme: "https", Host: "solon.example"}
	errors := make(chan error, requests)
	var waitGroup sync.WaitGroup

	for range requests {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			response, err := client.Get(endpoint, "")
			if err == nil {
				err = response.Body.Close()
			}
			errors <- err
		}()
	}

	waitGroup.Wait()
	close(errors)
	for err := range errors {
		assert.NoError(t, err)
	}
}
