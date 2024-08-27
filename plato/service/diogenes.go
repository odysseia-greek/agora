package service

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
)

type DiogenesImpl struct {
	Scheme  string
	BaseUrl string
	Client  HttpClient
}

func NewDiogenesConfig(schema OdysseiaApi, ca []byte) (*DiogenesImpl, error) {
	client := NewHttpClient(ca, schema.Cert)
	return &DiogenesImpl{Scheme: schema.Scheme, BaseUrl: schema.Url, Client: client}, nil
}

func NewFakeDiogenesConfig(scheme, baseUrl string, client HttpClient) (*DiogenesImpl, error) {
	return &DiogenesImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}

func (d *DiogenesImpl) Convert(body []byte, uuid string) (*http.Response, error) {
	textPath := url.URL{
		Scheme: d.Scheme,
		Host:   d.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s/%s", diogenesService, version, words, convert),
	}

	response, err := d.Client.Post(&textPath, body, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, textPath)
	}

	return response, nil
}

func (d *DiogenesImpl) Health(uuid string) (*http.Response, error) {
	healthPath := url.URL{
		Scheme: d.Scheme,
		Host:   d.BaseUrl,
		Path:   path.Join(diogenesService, version, healthEndPoint),
	}

	return Health(healthPath, d.Client, uuid)
}
