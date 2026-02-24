package aristoteles

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v9/esapi"
)

type ElasticError struct {
	Operation  string
	StatusCode int
	Status     string
	Took       *int64
	Detail     interface{}
	Body       string
}

func (e *ElasticError) Error() string {
	status := e.Status
	if strings.TrimSpace(status) == "" {
		if txt := http.StatusText(e.StatusCode); txt != "" {
			status = fmt.Sprintf("%d %s", e.StatusCode, txt)
		} else {
			status = fmt.Sprintf("%d", e.StatusCode)
		}
	}

	msg := fmt.Sprintf("%s: %s: elasticsearch request failed (status=%s)", errorMessage, e.Operation, status)
	if e.Took != nil {
		msg += fmt.Sprintf(" took=%dms", *e.Took)
	}
	if e.Detail != nil {
		msg += fmt.Sprintf(" error=%v", e.Detail)
	} else if e.Body != "" {
		msg += fmt.Sprintf(" body=%s", e.Body)
	}

	return msg
}

func newElasticError(operation string, res *esapi.Response) error {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("%s: failed reading error response body: %w", operation, err)
	}
	return newElasticErrorFromBody(operation, res, body)
}

func newElasticErrorFromBody(operation string, res *esapi.Response, body []byte) error {
	errResp := &ElasticError{
		Operation:  operation,
		StatusCode: res.StatusCode,
		Status:     strings.TrimSpace(res.Status()),
		Body:       strings.TrimSpace(string(body)),
	}

	var parsed map[string]interface{}
	if json.Unmarshal(body, &parsed) == nil {
		if rawStatus, ok := parsed["status"]; ok {
			if status, ok := toInt(rawStatus); ok {
				errResp.StatusCode = status
			}
		}
		if rawTook, ok := parsed["took"]; ok {
			if took, ok := toInt64(rawTook); ok {
				errResp.Took = &took
			}
		}
		if rawDetail, ok := parsed["error"]; ok {
			errResp.Detail = rawDetail
		}
	}

	return errResp
}

func toInt(raw interface{}) (int, bool) {
	switch v := raw.(type) {
	case float64:
		return int(v), true
	case int:
		return v, true
	case int64:
		return int(v), true
	default:
		return 0, false
	}
}

func toInt64(raw interface{}) (int64, bool) {
	switch v := raw.(type) {
	case float64:
		return int64(v), true
	case int:
		return int64(v), true
	case int64:
		return v, true
	default:
		return 0, false
	}
}
