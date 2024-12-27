package service

import (
	"fmt"
	"net/http"
	"net/url"
)

func Health(path url.URL, client HttpClient, uuid string) (*http.Response, error) {
	response, err := client.Get(&path, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return response, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	return response, nil
}
