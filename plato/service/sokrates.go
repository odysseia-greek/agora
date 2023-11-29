package service

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
)

type SokratesImpl struct {
	Scheme  string
	BaseUrl string
	Client  HttpClient
}

func NewSokratesConfig(schema OdysseiaApi, ca []byte) (*SokratesImpl, error) {
	client := NewHttpClient(ca, schema.Cert)
	return &SokratesImpl{Scheme: schema.Scheme, BaseUrl: schema.Url, Client: client}, nil
}

func NewFakeSokratesConfig(scheme, baseUrl string, client HttpClient) (*SokratesImpl, error) {
	return &SokratesImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}
func (s *SokratesImpl) Health(uuid string) (*http.Response, error) {
	healthPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   path.Join(sokratesService, version, healthEndPoint),
	}

	return Health(healthPath, s.Client, uuid)
}

func (s *SokratesImpl) GetMethods(uuid string) (*http.Response, error) {
	methodPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s", sokratesService, version, methods),
	}

	response, err := s.Client.Get(&methodPath, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, methodPath)
	}

	return response, nil
}

func (s *SokratesImpl) GetCategories(method, uuid string) (*http.Response, error) {
	categoryPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s/%s/%s", sokratesService, version, methods, method, categories),
	}

	response, err := s.Client.Get(&categoryPath, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, categoryPath)
	}

	return response, nil
}

func (s *SokratesImpl) GetChapters(method, category, uuid string) (*http.Response, error) {
	chapterPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s", sokratesService, version, methods, method, categories, category, chapters),
	}

	response, err := s.Client.Get(&chapterPath, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, chapterPath)
	}

	return response, nil
}

func (s *SokratesImpl) CreateQuestion(method, category, chapter, uuid string) (*http.Response, error) {
	query := fmt.Sprintf("method=%s&category=%s&chapter=%s", method, category, chapter)
	questionPath := url.URL{
		Scheme:   s.Scheme,
		Host:     s.BaseUrl,
		Path:     fmt.Sprintf("%s/%s/%s", sokratesService, version, question),
		RawQuery: query,
	}

	response, err := s.Client.Get(&questionPath, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, questionPath)
	}

	return response, nil
}

func (s *SokratesImpl) CheckAnswer(body []byte, uuid string) (*http.Response, error) {
	answerPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s", sokratesService, version, answer),
	}

	response, err := s.Client.Post(&answerPath, body, uuid)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, answerPath)
	}

	return response, nil
}
