package service

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
)

type HerodotosImpl struct {
	Scheme  string
	BaseUrl string
	Client  HttpClient
}

func NewHerodotosConfig(schema OdysseiaApi, ca []byte) (*HerodotosImpl, error) {
	client := NewHttpClient(ca, schema.Cert)
	return &HerodotosImpl{Scheme: schema.Scheme, BaseUrl: schema.Url, Client: client}, nil
}

func NewFakeHerodotosConfig(scheme, baseUrl string, client HttpClient) (*HerodotosImpl, error) {
	return &HerodotosImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}

func (h *HerodotosImpl) Options(uuid string) (*http.Response, error) {
	optionsPath := url.URL{
		Scheme: h.Scheme,
		Host:   h.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s/%s", herodotosService, version, texts, options),
	}

	response, err := h.Client.Get(&optionsPath, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return response, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, optionsPath)
	}

	return response, nil
}

func (h *HerodotosImpl) Analyze(body []byte, uuid string) (*http.Response, error) {
	textPath := url.URL{
		Scheme: h.Scheme,
		Host:   h.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s/%s", herodotosService, version, texts, analyzeVerb),
	}

	response, err := h.Client.Post(&textPath, body, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return response, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, textPath)
	}

	return response, nil
}

func (h *HerodotosImpl) Create(body []byte, uuid string) (*http.Response, error) {
	createPath := url.URL{
		Scheme: h.Scheme,
		Host:   h.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s/%s", herodotosService, version, texts, createVerb),
	}

	response, err := h.Client.Post(&createPath, body, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return response, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, createPath)
	}

	return response, nil
}

func (h *HerodotosImpl) Check(body []byte, uuid string) (*http.Response, error) {
	checkPath := url.URL{
		Scheme: h.Scheme,
		Host:   h.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s/%s", herodotosService, version, texts, checkVerb),
	}

	response, err := h.Client.Post(&checkPath, body, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return response, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, checkPath)
	}

	return response, nil
}

func (h *HerodotosImpl) Health(uuid string) (*http.Response, error) {
	healthPath := url.URL{
		Scheme: h.Scheme,
		Host:   h.BaseUrl,
		Path:   path.Join(herodotosService, version, healthEndPoint),
	}

	return Health(healthPath, h.Client, uuid)
}
