package aristoteles

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/odysseia-greek/aristoteles/models"
	"io"
)

type PolicyImpl struct {
	es *elasticsearch.Client
}

func NewPolicyImpl(suppliedClient *elasticsearch.Client) (*PolicyImpl, error) {
	return &PolicyImpl{es: suppliedClient}, nil
}

const (
	HOT  = "hot"
	WARM = "warm"
	COLD = "cold"
)

func (p *PolicyImpl) CreatePolicyWithRollOver(name, maxAge, phase string) (*models.IndexCreateResult, error) {
	policyDefinition := fmt.Sprintf(`{
		"policy": {
			"phases": {
				"%s": {
					"actions": {
						"rollover": {
							"max_age": "%sd"
						}
					}
				}
			}
		}
	}`, phase, maxAge)

	return p.create(name, policyDefinition)
}

func (p *PolicyImpl) CreateHotPolicy(name string) (*models.IndexCreateResult, error) {
	policyDefinition := fmt.Sprintf(`{
	"policy": {
		"phases": {
			"%s": {
				"actions": {}
			}
		}
	}
}`, HOT)

	return p.create(name, policyDefinition)
}

func (p *PolicyImpl) CreateWarmPolicy(name string) (*models.IndexCreateResult, error) {
	policyDefinition := fmt.Sprintf(`{
	"policy": {
		"phases": {
			"%s": {
				"actions": {}
			}
		}
	}
}`, WARM)

	return p.create(name, policyDefinition)
}

func (p *PolicyImpl) CreateColdPolicy(name string) (*models.IndexCreateResult, error) {
	policyDefinition := fmt.Sprintf(`{
	"policy": {
		"phases": {
			"%s": {
				"actions": {}
			}
		}
	}
}`, COLD)

	return p.create(name, policyDefinition)
}

func (p *PolicyImpl) create(name, policyDefinition string) (*models.IndexCreateResult, error) {
	var elasticResult models.IndexCreateResult

	req := esapi.ILMPutLifecycleRequest{
		Policy: name,
		Body:   bytes.NewReader([]byte(policyDefinition)),
	}

	res, err := req.Do(context.Background(), p.es)
	if err != nil {
		return &elasticResult, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("%s: %s", errorMessage, res.Status())
	}

	jsonBody, _ := io.ReadAll(res.Body)
	elasticResult, err = models.UnmarshalIndexCreateResult(jsonBody)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}
