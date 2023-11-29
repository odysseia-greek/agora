package service

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
)

type AlexandrosImpl struct {
	Scheme  string
	BaseUrl string
	Client  HttpClient
}

func NewAlexnadrosConfig(schema OdysseiaApi, ca []byte) (*AlexandrosImpl, error) {
	client := NewHttpClient(ca, schema.Cert)
	return &AlexandrosImpl{Scheme: schema.Scheme, BaseUrl: schema.Url, Client: client}, nil
}

func NewFakeAlexandrosConfig(scheme, baseUrl string, client HttpClient) (*AlexandrosImpl, error) {
	return &AlexandrosImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}

func (a *AlexandrosImpl) Search(word, language, mode, uuid string) (*http.Response, error) {
	values := url.Values{}
	values.Set(searchWord, word)
	if language != "" {
		values.Set(searchLanguage, language)
	}
	if mode != "" {
		values.Set(searchMode, mode)
	}

	query := values.Encode()

	urlPath := url.URL{
		Scheme:   a.Scheme,
		Host:     a.BaseUrl,
		Path:     fmt.Sprintf("%s/%s/%s", alexandrosService, version, search),
		RawQuery: query,
	}

	response, err := a.Client.Get(&urlPath, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, urlPath)
	}

	return response, nil
}

func (a *AlexandrosImpl) Health(uuid string) (*http.Response, error) {
	healthPath := url.URL{
		Scheme: a.Scheme,
		Host:   a.BaseUrl,
		Path:   path.Join(alexandrosService, version, healthEndPoint),
	}

	return Health(healthPath, a.Client, uuid)
}
