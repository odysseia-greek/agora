package service

import (
	"crypto/tls"
	"github.com/odysseia-greek/agora/plato/models"
	"net/http"
)

type OdysseiaClient interface {
	Solon() Solon
	Herodotos() Herodotos
	Alexandros() Alexandros
	Dionysios() Dionysios
}

type Odysseia struct {
	solon      *SolonImpl
	herodotos  *HerodotosImpl
	alexandros *AlexandrosImpl
	dionysios  *DionysiosImpl
}

type Solon interface {
	Health(uuid string) (*http.Response, error)
	OneTimeToken(uuid string) (*http.Response, error)
	Register(requestBody models.SolonCreationRequest, uuid string) (*http.Response, error)
}

type Herodotos interface {
	Health(uuid string) (*http.Response, error)
	Create(body []byte, uuid string) (*http.Response, error)
	Analyze(body []byte, uuid string) (*http.Response, error)
	Check(body []byte, uuid string) (*http.Response, error)
	Options(uuid string) (*http.Response, error)
}

type Alexandros interface {
	Health(uuid string) (*http.Response, error)
	Search(word, language, mode, textSearch, uuid string) (*http.Response, error)
}

type Dionysios interface {
	Health(uuid string) (*http.Response, error)
	Grammar(word string, uuid string) (*http.Response, error)
}

type ClientConfig struct {
	Ca         []byte
	Solon      OdysseiaApi
	Ptolemaios OdysseiaApi
	Herodotos  OdysseiaApi
	Dionysios  OdysseiaApi
	Alexandros OdysseiaApi
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

	herodotosImpl, err := NewHerodotosConfig(config.Herodotos, config.Ca)
	if err != nil {
		return nil, err
	}

	alexandrosImpl, err := NewAlexnadrosConfig(config.Alexandros, config.Ca)
	if err != nil {
		return nil, err
	}

	dionysiosImpl, err := NewDionysiosConfig(config.Dionysios, config.Ca)
	if err != nil {
		return nil, err
	}

	return &Odysseia{
		solon:      solonImpl,
		herodotos:  herodotosImpl,
		alexandros: alexandrosImpl,
		dionysios:  dionysiosImpl,
	}, nil
}

func NewFakeClient(config ClientConfig, codes []int, responses []string) (OdysseiaClient, error) {
	client := NewFakeHttpClient(responses, codes)

	solonImpl, err := NewFakeSolonImpl(config.Solon.Scheme, config.Solon.Url, client)
	if err != nil {
		return nil, err
	}

	herodotosImpl, err := NewFakeHerodotosConfig(config.Herodotos.Scheme, config.Herodotos.Url, client)
	if err != nil {
		return nil, err
	}

	alexandrosImpl, err := NewFakeAlexandrosConfig(config.Alexandros.Scheme, config.Alexandros.Url, client)
	if err != nil {
		return nil, err
	}

	dionysiosImpl, err := NewFakeDionysiosConfig(config.Dionysios.Scheme, config.Dionysios.Url, client)
	if err != nil {
		return nil, err
	}

	return &Odysseia{
		solon:      solonImpl,
		herodotos:  herodotosImpl,
		alexandros: alexandrosImpl,
		dionysios:  dionysiosImpl,
	}, nil
}

func (o *Odysseia) Solon() Solon {
	if o == nil {
		return nil
	}
	return o.solon
}

func (o *Odysseia) Herodotos() Herodotos {
	if o == nil {
		return nil
	}
	return o.herodotos
}

func (o *Odysseia) Alexandros() Alexandros {
	if o == nil {
		return nil
	}
	return o.alexandros
}

func (o *Odysseia) Dionysios() Dionysios {
	if o == nil {
		return nil
	}
	return o.dionysios
}
