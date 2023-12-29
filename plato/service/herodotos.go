package service

import (
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/models"
	"net/http"
	"net/url"
	"path"
	//"path"
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

func (h *HerodotosImpl) GetAuthors(uuid string) (*http.Response, error) {
	authorPath := url.URL{
		Scheme: h.Scheme,
		Host:   h.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s", herodotosService, version, authors),
	}

	response, err := h.Client.Get(&authorPath, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, authorPath)
	}

	return response, nil
}

func (h *HerodotosImpl) GetBooks(authorId string, uuid string) (*http.Response, error) {
	bookPath := url.URL{
		Scheme: h.Scheme,
		Host:   h.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s/%s/%s", herodotosService, version, authors, authorId, books),
	}

	response, err := h.Client.Get(&bookPath, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, bookPath)
	}

	return response, nil
}

func (h *HerodotosImpl) AnalyseText(rootword string, uuid string) (*http.Response, error) {
	query := fmt.Sprintf("author=%s", rootword)
	textPath := url.URL{
		Scheme:   h.Scheme,
		Host:     h.BaseUrl,
		Path:     fmt.Sprintf("%s/%s/%s", herodotosService, version, text),
		RawQuery: query,
	}

	response, err := h.Client.Get(&textPath, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, textPath)
	}

	return response, nil
}

func (h *HerodotosImpl) CreateQuestion(author, book, uuid string) (*http.Response, error) {
	query := fmt.Sprintf("author=%s&book=%s", author, book)
	questionPath := url.URL{
		Scheme:   h.Scheme,
		Host:     h.BaseUrl,
		Path:     fmt.Sprintf("%s/%s/%s", herodotosService, version, question),
		RawQuery: query,
	}

	response, err := h.Client.Get(&questionPath, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, questionPath)
	}

	return response, nil
}

func (h *HerodotosImpl) CheckSentence(checkSentenceRequest models.CheckSentenceRequest, uuid string) (*http.Response, error) {
	sentencePath := url.URL{
		Scheme: h.Scheme,
		Host:   h.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s", herodotosService, version, sentence),
	}

	body, err := json.Marshal(checkSentenceRequest)
	if err != nil {
		return nil, err
	}

	response, err := h.Client.Post(&sentencePath, body, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, sentencePath)
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
