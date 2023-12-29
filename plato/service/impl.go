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
	Sokrates() Sokrates
	Dionysios() Dionysios
}

type Odysseia struct {
	solon      *SolonImpl
	herodotos  *HerodotosImpl
	alexandros *AlexandrosImpl
	sokrates   *SokratesImpl
	dionysios  *DionysiosImpl
}

type Solon interface {
	Health(uuid string) (*http.Response, error)
	OneTimeToken(uuid string) (*http.Response, error)
	Register(requestBody models.SolonCreationRequest, uuid string) (*http.Response, error)
}

type Herodotos interface {
	Health(uuid string) (*http.Response, error)
	GetAuthors(uuid string) (*http.Response, error)
	GetBooks(authorId, uuid string) (*http.Response, error)
	CreateQuestion(author, book, uuid string) (*http.Response, error)
	AnalyseText(rootword string, uuid string) (*http.Response, error)
	CheckSentence(checkSentenceRequest models.CheckSentenceRequest, uuid string) (*http.Response, error)
}

type Alexandros interface {
	Health(uuid string) (*http.Response, error)
	Search(word, language, mode, uuid string) (*http.Response, error)
}

type Sokrates interface {
	Health(uuid string) (*http.Response, error)
	GetMethods(uuid string) (*http.Response, error)
	GetCategories(method, uuid string) (*http.Response, error)
	GetChapters(method, category, uuid string) (*http.Response, error)
	CreateQuestion(method, category, chapter, uuid string) (*http.Response, error)
	CheckAnswer(body []byte, uuid string) (*http.Response, error)
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
	Sokrates   OdysseiaApi
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

	sokratesImpl, err := NewSokratesConfig(config.Sokrates, config.Ca)
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
		sokrates:   sokratesImpl,
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

	sokratesImpl, err := NewFakeSokratesConfig(config.Sokrates.Scheme, config.Sokrates.Url, client)
	if err != nil {
		return nil, err
	}

	dionysiosImpl, err := NewFakeDionysiosConfig(config.Dionysios.Scheme, config.Sokrates.Url, client)
	if err != nil {
		return nil, err
	}

	return &Odysseia{
		solon:      solonImpl,
		herodotos:  herodotosImpl,
		alexandros: alexandrosImpl,
		sokrates:   sokratesImpl,
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

func (o *Odysseia) Sokrates() Sokrates {
	if o == nil {
		return nil
	}
	return o.sokrates
}
func (o *Odysseia) Dionysios() Dionysios {
	if o == nil {
		return nil
	}
	return o.dionysios
}
