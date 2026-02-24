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

type DocumentImpl struct {
	es *elasticsearch.Client
}

func NewDocumentImpl(suppliedClient *elasticsearch.Client) (*DocumentImpl, error) {
	return &DocumentImpl{es: suppliedClient}, nil
}

func (d *DocumentImpl) Create(index string, body []byte) (*models.CreateResult, error) {
	return d.CreateWithContext(context.Background(), index, body)
}

func (d *DocumentImpl) CreateWithContext(ctx context.Context, index string, body []byte) (*models.CreateResult, error) {
	res, err := esapi.CreateRequest{
		Index: index,
		Body:  bytes.NewReader(body),
	}.Do(ctx, d.es)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	jsonBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, newElasticErrorFromBody("create document", res, jsonBody)
	}

	elasticResult, err := models.UnmarshalCreateResult(jsonBody)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}

func (d *DocumentImpl) CreateWithId(index, documentId string, body []byte) (*models.CreateResult, error) {
	return d.CreateWithIdWithContext(context.Background(), index, documentId, body)
}

func (d *DocumentImpl) CreateWithIdWithContext(ctx context.Context, index, documentId string, body []byte) (*models.CreateResult, error) {
	res, err := esapi.CreateRequest{
		Index:      index,
		DocumentID: documentId,
		Body:       bytes.NewReader(body),
	}.Do(ctx, d.es)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	jsonBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, newElasticErrorFromBody("create document with id", res, jsonBody)
	}

	elasticResult, err := models.UnmarshalCreateResult(jsonBody)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}

func (d *DocumentImpl) CreateWithIdAndFirstItem(index, documentId, body, paramName string) (*models.CreateResult, error) {
	return d.CreateWithIdAndFirstItemWithContext(context.Background(), index, documentId, body, paramName)
}

func (d *DocumentImpl) CreateWithIdAndFirstItemWithContext(ctx context.Context, index, documentId, body, paramName string) (*models.CreateResult, error) {
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

	jsonBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, newElasticErrorFromBody("create document with first item", res, jsonBody)
	}

	elasticResult, err := models.UnmarshalCreateResult(jsonBody)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}

func (d *DocumentImpl) Update(index, id string, body []byte) (*models.CreateResult, error) {
	return d.UpdateWithContext(context.Background(), index, id, body)
}

func (d *DocumentImpl) UpdateWithContext(ctx context.Context, index, id string, body []byte) (*models.CreateResult, error) {
	res, err := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, body))),
	}.Do(ctx, d.es)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	jsonBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, newElasticErrorFromBody("update document", res, jsonBody)
	}

	elasticResult, err := models.UnmarshalCreateResult(jsonBody)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}

func (d *DocumentImpl) AddItemToDocument(index, id, body, paramName string) (*models.CreateResult, error) {
	return d.AddItemToDocumentWithContext(context.Background(), index, id, body, paramName)
}

func (d *DocumentImpl) AddItemToDocumentWithContext(ctx context.Context, index, id, body, paramName string) (*models.CreateResult, error) {
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

	jsonBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, newElasticErrorFromBody("add item to document", res, jsonBody)
	}

	elasticResult, err := models.UnmarshalCreateResult(jsonBody)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}
