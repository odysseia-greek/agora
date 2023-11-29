package aristoteles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	var elasticResult BulkResponse
	res, err := d.es.Bulk(bytes.NewReader(buf.Bytes()), d.es.Bulk.WithIndex(index))
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
	err = json.Unmarshal(jsonBody, &elasticResult)
	if err != nil {
		return nil, err
	}

	return &elasticResult, nil
}
