package aristoteles

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/odysseia-greek/agora/aristoteles/models"
)

type IndexImpl struct {
	es *elasticsearch.Client
}

func NewIndexImpl(suppliedClient *elasticsearch.Client) (*IndexImpl, error) {
	return &IndexImpl{es: suppliedClient}, nil
}

func (i *IndexImpl) CreateDocument(index string, body []byte) (*models.CreateResult, error) {
	return i.CreateDocumentWithContext(context.Background(), index, body)
}

func (i *IndexImpl) CreateDocumentWithContext(ctx context.Context, index string, body []byte) (*models.CreateResult, error) {
	bodyString := strings.NewReader(string(body))

	esRequest := esapi.IndexRequest{
		Body:       bodyString,
		Refresh:    "true",
		Index:      index,
		DocumentID: "",
	}

	res, err := esRequest.Do(ctx, i.es)
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

func (i *IndexImpl) Create(index string, request map[string]interface{}) (*models.IndexCreateResult, error) {
	return i.CreateWithContext(context.Background(), index, request)
}

func (i *IndexImpl) CreateWithContext(ctx context.Context, index string, request map[string]interface{}) (*models.IndexCreateResult, error) {
	query, err := toBuffer(request)
	if err != nil {
		return nil, err
	}
	var elasticResult models.IndexCreateResult

	indexRequest := esapi.IndicesCreateRequest{
		Index: index,
		Body:  &query,
	}

	res, err := indexRequest.Do(ctx, i.es)
	if err != nil {
		return &elasticResult, err
	}
	defer res.Body.Close()

	jsonBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, newElasticErrorFromBody("create index", res, jsonBody)
	}

	elasticResult, err = models.UnmarshalIndexCreateResult(jsonBody)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}

func (i *IndexImpl) CreateWithAlias(indexName string, request map[string]interface{}) (*models.IndexCreateResult, error) {
	return i.CreateWithAliasWithContext(context.Background(), indexName, request)
}

func (i *IndexImpl) CreateWithAliasWithContext(ctx context.Context, indexName string, request map[string]interface{}) (*models.IndexCreateResult, error) {
	today := time.Now().Format("2006.01.02")
	indexNameWithDate := fmt.Sprintf("%s-%s", indexName, today)
	var elasticResult models.IndexCreateResult

	query, err := toBuffer(request)
	if err != nil {
		return nil, err
	}

	indexRequest := esapi.IndicesCreateRequest{
		Index: indexNameWithDate,
		Body:  &query,
	}

	res, err := indexRequest.Do(ctx, i.es)
	if err != nil {
		return &elasticResult, err
	}
	defer res.Body.Close()

	jsonBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, newElasticErrorFromBody("create index with alias", res, jsonBody)
	}

	aliasRequest := esapi.IndicesPutAliasRequest{
		Index: []string{indexNameWithDate},
		Name:  indexName,
	}
	aliasRes, err := aliasRequest.Do(ctx, i.es)
	if err != nil {
		return nil, err
	}
	defer aliasRes.Body.Close()
	if aliasRes.IsError() {
		return nil, newElasticError("put index alias", aliasRes)
	}

	elasticResult, err = models.UnmarshalIndexCreateResult(jsonBody)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}

func (i *IndexImpl) Delete(index string) (bool, error) {
	return i.DeleteWithContext(context.Background(), index)
}

func (i *IndexImpl) DeleteWithContext(ctx context.Context, index string) (bool, error) {
	res, err := i.es.Indices.Delete([]string{index}, i.es.Indices.Delete.WithContext(ctx))
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	jsonBody, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	if res.IsError() {
		return false, newElasticErrorFromBody("delete index", res, jsonBody)
	}

	var r map[string]interface{}
	if err := json.Unmarshal(jsonBody, &r); err != nil {
		return false, fmt.Errorf("error parsing successful response: %w", err)
	}

	acknowledged, ok := r["acknowledged"].(bool)
	if !ok {
		return false, fmt.Errorf("'acknowledged' field missing or not boolean in response: %v", r)
	}

	return acknowledged, nil
}

func (i *IndexImpl) IndexExists(index string) (bool, *models.IndexInfo, error) {
	return i.IndexExistsWithContext(context.Background(), index)
}

func (i *IndexImpl) IndexExistsWithContext(ctx context.Context, index string) (bool, *models.IndexInfo, error) {
	getRequest := esapi.IndicesGetRequest{
		Index: []string{index},
	}

	res, err := getRequest.Do(ctx, i.es)
	if err != nil {
		return false, nil, err
	}
	defer res.Body.Close()

	jsonBody, err := io.ReadAll(res.Body)
	if err != nil {
		return false, nil, err
	}

	if res.StatusCode == 404 {
		return false, nil, nil
	}

	if res.IsError() {
		return false, nil, newElasticErrorFromBody("check index existence", res, jsonBody)
	}

	var indexInfo map[string]interface{}
	if err := json.Unmarshal(jsonBody, &indexInfo); err != nil {
		return false, nil, err
	}

	if data, exists := indexInfo[index]; exists {
		indexData, ok := data.(map[string]interface{})
		if !ok {
			return false, nil, fmt.Errorf("invalid index data format for %q", index)
		}

		info := &models.IndexInfo{IndexName: index}

		if settings, ok := indexData["settings"].(map[string]interface{}); ok {
			info.Settings = settings
		}
		if mappings, ok := indexData["mappings"].(map[string]interface{}); ok {
			info.Mappings = mappings
		}
		if total, ok := indexData["total"].(map[string]interface{}); ok {
			if docs, ok := total["docs"].(map[string]interface{}); ok {
				if count, ok := docs["count"].(float64); ok {
					info.TotalDocuments = int64(count)
				}
			}
		}

		return true, info, nil
	}

	return false, nil, nil
}
