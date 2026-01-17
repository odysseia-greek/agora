package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	router := Router()
	body := map[string]string{"answerSentence": "The Foenicians. ;came to Argos,,.;:'' afd set out some cargo",
		"sentenceId": "GmBFYHkBkbwXxxT5S6F_",
		"author":     "herodotos"}

	jsonBody, _ := json.Marshal(body)
	bodyInBytes := bytes.NewReader(jsonBody)
	code := "sometestcode"

	t.Run("HappyFlowPost", func(t *testing.T) {
		response := PerformPostRequest(router, "/test/v1/postping", bodyInBytes)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("HappyFlowPost", func(t *testing.T) {
		response := PerformGetRequest(router, "/test/v1/getping")
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("CannotUsePostOnGetMethod", func(t *testing.T) {
		response := PerformPostRequest(router, "/test/v1/getping", bodyInBytes)
		assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
	})

	t.Run("CannotUseGetOnPostMethod", func(t *testing.T) {
		response := PerformGetRequest(router, "/test/v1/postping")
		assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
	})

	t.Run("LoggingAischylosCodeGet", func(t *testing.T) {
		var stringBuffer bytes.Buffer
		log.SetOutput(&stringBuffer)
		PerformGetRequestWithInjectedCode(router, "/test/v1/logger", code)
		assert.Contains(t, stringBuffer.String(), code)
	})

	t.Run("LoggingAischylosCodePost", func(t *testing.T) {
		var stringBuffer bytes.Buffer
		log.SetOutput(&stringBuffer)
		PerformPostRequestWithInjectedCode(router, "/test/v1/logger", code, bodyInBytes)
		assert.Contains(t, stringBuffer.String(), code)
	})
}

func Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/test/v1/getping", Adapt(pingPong, ValidateRestMethod("GET"), LogRequestDetails(), SetCorsHeaders()))
	mux.HandleFunc("/test/v1/postping", Adapt(pingPong, ValidateRestMethod("POST"), LogRequestDetails(), SetCorsHeaders()))
	mux.HandleFunc("/test/v1/logger", Adapt(pingPong, ValidateRestMethod("POST"), LogRequestDetails(), SetCorsHeaders()))
	mux.HandleFunc("/test/v1/health", Adapt(pingPong, ValidateRestMethod("GET"), LogRequestDetails(), SetCorsHeaders()))

	return mux
}

// PingPong pongs the ping
func pingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	ResponseWithJson(w, pingPong)
}

func PerformGetRequestWithInjectedCode(r http.Handler, path, code string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Set(service.HeaderKey, code)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func PerformGetRequest(r http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func PerformPostRequestWithInjectedCode(r http.Handler, path, code string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, body)
	req.Header.Set(service.HeaderKey, code)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func PerformPostRequest(r http.Handler, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
