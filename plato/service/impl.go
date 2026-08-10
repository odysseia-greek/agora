package service

import (
	"crypto/tls"
	"github.com/odysseia-greek/agora/plato/models"
	"net/http"
)

type OdysseiaClient interface {
	Solon() Solon
}

type Odysseia struct {
	solon *SolonImpl
}

type Solon interface {
	Health(uuid string) (*http.Response, error)
	OneTimeToken(uuid string) (*http.Response, error)
	Register(requestBody models.SolonCreationRequest, uuid string) (*http.Response, error)
}

type ClientConfig struct {
	Ca    []byte
	Solon OdysseiaApi
}

type OdysseiaApi struct {
	Url    string
	Scheme string
	Cert   []tls.Certificate
}

func NewClient(config ClientConfig) (OdysseiaClient, error) {
	solonImpl, err := NewSolonImpl(config.Solon, config.Ca)
	if err != nil {
		return nil, err
	}

	return &Odysseia{
		solon: solonImpl,
	}, nil
}

func NewFakeClient(config ClientConfig, codes []int, responses []string) (OdysseiaClient, error) {
	client := NewFakeHttpClient(responses, codes)

	solonImpl, err := NewFakeSolonImpl(config.Solon.Scheme, config.Solon.Url, client)
	if err != nil {
		return nil, err
	}

	return &Odysseia{
		solon: solonImpl,
	}, nil
}

func (o *Odysseia) Solon() Solon {
	if o == nil {
		return nil
	}
	return o.solon
}
