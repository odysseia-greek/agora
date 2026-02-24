package aristoteles

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
)

type BulkResponse struct {
	Errors bool `json:"errors"`
	Items  []struct {
		Index struct {
			ID     string `json:"_id"`
			Result string `json:"result"`
			Status int    `json:"status"`
			Error  struct {
				Type   string `json:"type"`
				Reason string `json:"reason"`
				Cause  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
				} `json:"caused_by"`
			} `json:"error"`
		} `json:"index"`
	} `json:"items"`
}

func (d *DocumentImpl) Bulk(buf bytes.Buffer, index string) (*BulkResponse, error) {
	return d.BulkWithContext(context.Background(), buf, index)
}

func (d *DocumentImpl) BulkWithContext(ctx context.Context, buf bytes.Buffer, index string) (*BulkResponse, error) {
	res, err := d.es.Bulk(
		bytes.NewReader(buf.Bytes()),
		d.es.Bulk.WithContext(ctx),
		d.es.Bulk.WithIndex(index),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	jsonBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, newElasticErrorFromBody("bulk create documents", res, jsonBody)
	}

	var elasticResult BulkResponse
	if err := json.Unmarshal(jsonBody, &elasticResult); err != nil {
		return nil, err
	}

	return &elasticResult, nil
}
