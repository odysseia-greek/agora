package service

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
)

type DionysiosImpl struct {
	Scheme  string
	BaseUrl string
	Client  HttpClient
}

func NewDionysiosConfig(schema OdysseiaApi, ca []byte) (*DionysiosImpl, error) {
	client := NewHttpClient(ca, schema.Cert)
	return &DionysiosImpl{Scheme: schema.Scheme, BaseUrl: schema.Url, Client: client}, nil
}

func NewFakeDionysiosConfig(scheme, baseUrl string, client HttpClient) (*DionysiosImpl, error) {
	return &DionysiosImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}

func (d *DionysiosImpl) Health(uuid string) (*http.Response, error) {
	healthPath := url.URL{
		Scheme: d.Scheme,
		Host:   d.BaseUrl,
		Path:   path.Join(dionysiosService, version, healthEndPoint),
	}

	return Health(healthPath, d.Client, uuid)
}

func (d *DionysiosImpl) Grammar(word string, uuid string) (*http.Response, error) {
	query := fmt.Sprintf("%s=%s", searchWord, word)

	urlPath := url.URL{
		Scheme:   d.Scheme,
		Host:     d.BaseUrl,
		Path:     path.Join(dionysiosService, version, grammar),
		RawQuery: query,
	}

	response, err := d.Client.Get(&urlPath, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return response, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, urlPath)
	}

	return response, nil
}
