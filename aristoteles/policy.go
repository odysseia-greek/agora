package aristoteles

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/odysseia-greek/agora/aristoteles/models"
)

type PolicyImpl struct {
	es *elasticsearch.Client
}

func NewPolicyImpl(suppliedClient *elasticsearch.Client) (*PolicyImpl, error) {
	return &PolicyImpl{es: suppliedClient}, nil
}

func (p *PolicyImpl) CreatePolicyWithRollOver(name, maxAge, phase string) (*models.IndexCreateResult, error) {
	return p.CreatePolicyWithRollOverWithContext(context.Background(), name, maxAge, phase)
}

func (p *PolicyImpl) CreatePolicyWithRollOverWithContext(ctx context.Context, name, maxAge, phase string) (*models.IndexCreateResult, error) {
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

	return p.create(ctx, name, policyDefinition)
}

func (p *PolicyImpl) CreatePolicy(name, phase string) (*models.IndexCreateResult, error) {
	return p.CreatePolicyWithContext(context.Background(), name, phase)
}

func (p *PolicyImpl) CreatePolicyWithContext(ctx context.Context, name, phase string) (*models.IndexCreateResult, error) {
	policyDefinition := fmt.Sprintf(`{
	"policy": {
		"phases": {
			"%s": {
				"actions": {}
			}
		}
	}
}`, phase)

	return p.create(ctx, name, policyDefinition)
}

func (p *PolicyImpl) create(ctx context.Context, name, policyDefinition string) (*models.IndexCreateResult, error) {
	req := esapi.ILMPutLifecycleRequest{
		Policy: name,
		Body:   bytes.NewReader([]byte(policyDefinition)),
	}

	res, err := req.Do(ctx, p.es)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	jsonBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, newElasticErrorFromBody("create ilm policy", res, jsonBody)
	}

	elasticResult, err := models.UnmarshalIndexCreateResult(jsonBody)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}
