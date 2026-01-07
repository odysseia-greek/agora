package aristoteles

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"io"
	"log"
)

type DocumentImpl struct {
	es *elasticsearch.Client
}

func NewDocumentImpl(suppliedClient *elasticsearch.Client) (*DocumentImpl, error) {
	return &DocumentImpl{es: suppliedClient}, nil
}

func (d *DocumentImpl) Create(index string, body []byte) (*models.CreateResult, error) {
	var elasticResult models.CreateResult

	ctx := context.Background()
	res, err := esapi.CreateRequest{
		Index: index,
		Body:  bytes.NewReader(body),
	}.Do(ctx, d.es)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("%s: %s", errorMessage, res.Status())
	}

	jsonBody, _ := io.ReadAll(res.Body)
	elasticResult, err = models.UnmarshalCreateResult(jsonBody)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}

func (d *DocumentImpl) CreateWithId(index, documentId string, body []byte) (*models.CreateResult, error) {
	var elasticResult models.CreateResult

	ctx := context.Background()
	res, err := esapi.CreateRequest{
		Index:      index,
		DocumentID: documentId,
		Body:       bytes.NewReader(body),
	}.Do(ctx, d.es)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("%s: %s", errorMessage, res.Status())
	}

	jsonBody, _ := io.ReadAll(res.Body)
	elasticResult, err = models.UnmarshalCreateResult(jsonBody)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}

func (d *DocumentImpl) CreateWithIdAndFirstItem(index, documentId, body, paramName string) (*models.CreateResult, error) {
	var elasticResult models.CreateResult

	ctx := context.Background()
	res, err := esapi.CreateRequest{
		Index:      index,
		DocumentID: documentId,
		Body: bytes.NewReader([]byte(fmt.Sprintf(`{
			"script": {
				"source": "ctx._source.items.addAll(params.items)",
				"params": {
					"%s": %s
				}
			}
		}`, paramName, body))),
	}.Do(ctx, d.es)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("%s: %s", errorMessage, res.Status())
	}

	jsonBody, _ := io.ReadAll(res.Body)
	elasticResult, err = models.UnmarshalCreateResult(jsonBody)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}

func (d *DocumentImpl) Update(index, id string, body []byte) (*models.CreateResult, error) {
	var elasticResult models.CreateResult

	ctx := context.Background()
	res, err := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, body))),
	}.Do(ctx, d.es)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		jsonBody, _ := io.ReadAll(res.Body)
		log.Print(jsonBody)
		return nil, fmt.Errorf("%s: %s", errorMessage, res.Status())
	}

	jsonBody, _ := io.ReadAll(res.Body)
	elasticResult, err = models.UnmarshalCreateResult(jsonBody)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}

func (d *DocumentImpl) AddItemToDocument(index, id, body, paramName string) (*models.CreateResult, error) {
	var elasticResult models.CreateResult

	ctx := context.Background()

	// Build the update request
	res, err := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body: bytes.NewReader([]byte(fmt.Sprintf(`{
			"script": {
				"source": "ctx._source.%s.add(params.item)",
				"lang": "painless",
				"params": {
					"item": %s
				}
			}
		}`, paramName, body))),
	}.Do(ctx, d.es)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		jsonBody, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("%s: %s", errorMessage, string(jsonBody))
	}

	jsonBody, _ := io.ReadAll(res.Body)
	elasticResult, err = models.UnmarshalCreateResult(jsonBody)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}
