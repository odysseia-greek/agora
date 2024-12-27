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

func (s *SokratesImpl) Create(body []byte, requestID string) (*http.Response, error) {
	createPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s/%s", sokratesService, version, quiz, create),
	}

	response, err := s.Client.Post(&createPath, body, requestID)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return response, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, createPath)
	}

	return response, nil
}

func (s *SokratesImpl) Check(body []byte, requestID string) (*http.Response, error) {
	answerPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s/%s", sokratesService, version, quiz, answer),
	}

	response, err := s.Client.Post(&answerPath, body, requestID)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return response, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, answerPath)
	}

	return response, nil
}

func (s *SokratesImpl) Options(quizType string, requestID string) (*http.Response, error) {
	query := fmt.Sprintf("%s=%s", QuizType, quizType)

	answerPath := url.URL{
		Scheme:   s.Scheme,
		Host:     s.BaseUrl,
		Path:     fmt.Sprintf("%s/%s/%s/%s", sokratesService, version, quiz, options),
		RawQuery: query,
	}

	response, err := s.Client.Get(&answerPath, requestID)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return response, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, answerPath)
	}

	return response, nil
}
