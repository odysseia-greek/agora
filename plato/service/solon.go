package service

import (
	"fmt"
	"github.com/odysseia-greek/agora/plato/models"
	"net/http"
	"net/url"
	"path"
)

type SolonImpl struct {
	Scheme  string
	BaseUrl string
	Client  HttpClient
}

func NewSolonImpl(schema OdysseiaApi, ca []byte) (*SolonImpl, error) {
	client := NewHttpClient(ca, schema.Cert)
	return &SolonImpl{Scheme: schema.Scheme, BaseUrl: schema.Url, Client: client}, nil
}

func NewFakeSolonImpl(scheme, baseUrl string, client HttpClient) (*SolonImpl, error) {
	return &SolonImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}

func (s *SolonImpl) OneTimeToken(uuid string) (*http.Response, error) {
	urlPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   path.Join(solonService, version, token),
	}

	response, err := s.Client.Get(&urlPath, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	return response, nil
}

func (s *SolonImpl) Register(requestBody models.SolonCreationRequest, uuid string) (*http.Response, error) {
	urlPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   path.Join(solonService, version, register),
	}

	body, err := requestBody.Marshal()
	if err != nil {
		return nil, err
	}

	response, err := s.Client.Post(&urlPath, body, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	return response, nil
}

func (s *SolonImpl) Health(uuid string) (*http.Response, error) {
	healthPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   path.Join(solonService, version, healthEndPoint),
	}

	return Health(healthPath, s.Client, uuid)
}
