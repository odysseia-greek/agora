package aristoteles

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"io"
)

type PolicyImpl struct {
	es *elasticsearch.Client
}

func NewPolicyImpl(suppliedClient *elasticsearch.Client) (*PolicyImpl, error) {
	return &PolicyImpl{es: suppliedClient}, nil
}

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

func (p *PolicyImpl) CreatePolicy(name, phase string) (*models.IndexCreateResult, error) {
	policyDefinition := fmt.Sprintf(`{
	"policy": {
		"phases": {
			"%s": {
				"actions": {}
			}
		}
	}
}`, phase)

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
