package aristoteles

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"io"
	"strings"
	"time"
)

type IndexImpl struct {
	es *elasticsearch.Client
}

func NewIndexImpl(suppliedClient *elasticsearch.Client) (*IndexImpl, error) {
	return &IndexImpl{es: suppliedClient}, nil
}

func (i *IndexImpl) CreateDocument(index string, body []byte) (*models.CreateResult, error) {
	var elasticResult models.CreateResult
	bodyString := strings.NewReader(string(body))

	esRequest := esapi.IndexRequest{
		Body:       bodyString,
		Refresh:    "true",
		Index:      index,
		DocumentID: "",
	}

	res, err := esRequest.Do(context.Background(), i.es)
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

func (i *IndexImpl) Create(index string, request map[string]interface{}) (*models.IndexCreateResult, error) {
	query, err := toBuffer(request)
	if err != nil {
		return nil, err
	}

	var elasticResult models.IndexCreateResult
	indexRequest := esapi.IndicesCreateRequest{
		Index: index,
		Body:  &query,
	}

	res, err := indexRequest.Do(context.Background(), i.es)
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

func (i *IndexImpl) CreateWithAlias(indexName string, request map[string]interface{}) (*models.IndexCreateResult, error) {
	// Generate index name with the current date
	today := time.Now().Format("2006.01.02") // Use dashes (-) instead of dots (.)
	indexNameWithDate := fmt.Sprintf("%s-%s", indexName, today)

	// Set up the alias name
	aliasName := indexName

	// Create the index and alias
	query, err := toBuffer(request)
	if err != nil {
		return nil, err
	}

	var elasticResult models.IndexCreateResult
	indexRequest := esapi.IndicesCreateRequest{
		Index: indexNameWithDate,
		Body:  &query,
	}

	res, err := indexRequest.Do(context.Background(), i.es)
	if err != nil {
		return &elasticResult, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("%s: %s", errorMessage, res.Status())
	}

	// Create the alias
	aliasRequest := esapi.IndicesPutAliasRequest{
		Index: []string{indexNameWithDate},
		Name:  aliasName,
	}
	_, err = aliasRequest.Do(context.Background(), i.es)
	if err != nil {
		return nil, err
	}

	jsonBody, _ := io.ReadAll(res.Body)
	elasticResult, err = models.UnmarshalIndexCreateResult(jsonBody)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}

func (i *IndexImpl) Update(index string, request map[string]interface{}) (*models.IndexCreateResult, error) {
	query, err := toBuffer(request)
	if err != nil {
		return nil, err
	}

	var elasticResult models.IndexCreateResult
	indexRequest := esapi.IndicesCreateRequest{
		Index: index,
		Body:  &query,
	}

	res, err := indexRequest.Do(context.Background(), i.es)
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

func (i *IndexImpl) Delete(index string) (bool, error) {
	res, err := i.es.Indices.Delete([]string{index})
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	// Check for API errors first before trying to parse JSON
	if res.IsError() {
		// Read the raw body for error details
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return false, fmt.Errorf("error reading error response: %v", err)
		}

		// Try to parse as JSON if possible for structured error
		var errorResponse map[string]interface{}
		if json.Unmarshal(bodyBytes, &errorResponse) == nil {
			// Successfully parsed JSON error
			errorJSON, _ := json.Marshal(errorResponse)
			return false, fmt.Errorf("%s", string(errorJSON))
		}

		// Couldn't parse as JSON, return raw error
		return false, fmt.Errorf("elasticsearch error: %s", string(bodyBytes))
	}

	// Only try to parse JSON for successful responses
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		bodyBytes, _ := io.ReadAll(res.Body)
		return false, fmt.Errorf("error parsing successful response: %v, body: %s", err, string(bodyBytes))
	}

	// Handle case where acknowledged might not be present
	acknowledged, ok := r["acknowledged"].(bool)
	if !ok {
		return false, fmt.Errorf("'acknowledged' field missing or not boolean in response: %v", r)
	}

	return acknowledged, nil

}

func (i *IndexImpl) IndexExists(index string) (bool, *models.IndexInfo, error) {
	// Send a request to check the index
	getRequest := esapi.IndicesGetRequest{
		Index: []string{index},
	}

	res, err := getRequest.Do(context.Background(), i.es)
	if err != nil {
		return false, nil, err
	}
	defer res.Body.Close()

	// If the index does not exist
	if res.StatusCode == 404 {
		return false, nil, nil
	}

	if res.IsError() {
		return false, nil, fmt.Errorf("error checking index existence: %s", res.Status())
	}

	// Parse the response body
	var indexInfo map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&indexInfo); err != nil {
		return false, nil, err
	}

	// Extract relevant data about the index
	if data, exists := indexInfo[index]; exists {
		indexData := data.(map[string]interface{})

		// Prepare the struct
		info := &models.IndexInfo{
			IndexName: index,
		}

		// Extract settings
		if settings, ok := indexData["settings"].(map[string]interface{}); ok {
			info.Settings = settings
		}

		// Extract mappings
		if mappings, ok := indexData["mappings"].(map[string]interface{}); ok {
			info.Mappings = mappings
		}

		// Extract stats
		if total, ok := indexData["total"].(map[string]interface{}); ok {
			if docs, ok := total["docs"].(map[string]interface{}); ok {
				if count, ok := docs["count"].(float64); ok {
					info.TotalDocuments = int64(count)
				}
			}
			if store, ok := total["store"].(map[string]interface{}); ok {
				if size, ok := store["size_in_bytes"].(float64); ok {
					info.SizeInBytes = int64(size)
				}
			}
		}

		return true, info, nil
	}

	return true, nil, fmt.Errorf("unexpected response format while checking index existence")
}
