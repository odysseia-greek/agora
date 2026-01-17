package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
)

type Adapter func(http.HandlerFunc) http.HandlerFunc

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (rec *StatusRecorder) WriteHeader(code int) {
	rec.Status = code
	rec.ResponseWriter.WriteHeader(code)
}

// Adapt Iterate over adapters and run them one by one
func Adapt(h http.HandlerFunc, adapters ...Adapter) http.HandlerFunc {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

// LogRequestDetails logs the request details including the request method, URL, header keys and now also includes logging the latency
func LogRequestDetails() Adapter {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			blockList := []string{"health", "ping"}
			requestId := r.Header.Get(service.HeaderKey)
			for _, block := range blockList {
				if strings.Contains(r.URL.Path, block) {
					f(w, r)
					return
				}
			}

			f(w, r)

			duration := time.Since(startTime)
			clientIp := r.RemoteAddr
			method := r.Method
			path := r.URL.Path
			var statusCode int
			responseWriter, ok := w.(*StatusRecorder)
			if ok {
				statusCode = responseWriter.Status
			} else {
				// if w is not our wrapped response writer, we cannot get the status
				// so, let's set the status to StatusOK for this case
				statusCode = http.StatusOK
			}
			logging.Api(statusCode, method, requestId, clientIp, path, duration)
		}
	}
}

// ValidateRestMethod middleware to validate proper methods
func ValidateRestMethod(method string) Adapter {

	return func(f http.HandlerFunc) http.HandlerFunc {

		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				startTime := time.Now()
				var err models.MethodError
				e := models.MethodMessages{method, "Method " + r.Method + " not allowed at this endpoint"}
				err = models.MethodError{models.ErrorModel{UniqueCode: "methodError"}, append(err.Messages, e)}
				requestId := r.Header.Get(service.HeaderKey)
				duration := time.Since(startTime)
				logging.Api(http.StatusMethodNotAllowed, method, requestId, r.RemoteAddr, r.URL.Path, duration)
				ResponseWithJson(w, err)
				return
			}
			f(w, r)
		}
	}
}
func SetCorsHeaders() Adapter {
	// "open" CORS helper: useful for quick dev / tools, but do NOT use with credentials.
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Boule")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next(w, r)
		}
	}
}

func SetCorsHeadersLocal() Adapter {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				// Non-browser / same-origin / curl etc.
				next(w, r)
				return
			}

			if isAllowedLocalOrigin(origin) {
				logging.Debug(fmt.Sprintf("setting CORS header for origin: %s", origin))

				// Echo origin so credentials are possible
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Credentials", "true")

				// If you ever send custom headers, list them here
				w.Header().Set(
					"Access-Control-Allow-Headers",
					"Origin, X-Requested-With, Content-Type, Accept, Authorization, Boule",
				)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}

			next(w, r)
		}
	}
}

func isAllowedLocalOrigin(origin string) bool {
	// Allow common local dev origins with ports
	allowedPrefixes := []string{
		"http://localhost",
		"http://127.0.0.1",
		"http://0.0.0.0",
		"https://localhost",
		"https://127.0.0.1",
	}
	for _, p := range allowedPrefixes {
		if strings.HasPrefix(origin, p) {
			return true
		}
	}
	return false
}

func Env() string {
	if v := os.Getenv("ENV"); v != "" {
		return v
	}
	return "unknown"
}

// ResponseWithJson returns formed JSON
func ResponseWithJson(w http.ResponseWriter, payload interface{}) {
	code := 500

	switch payload.(type) {
	case models.SolonResponse:
		code = 200
	case models.ResultModel:
		code = 200
	case models.Word:
		code = 200
	case []models.Meros:
		code = 200
	case models.Health:
		code = 200
	case models.DeclensionTranslationResults:
		code = 200
	case models.TokenResponse:
		code = 200
	case map[string]interface{}:
		code = 200
	case models.ValidationError:
		code = 400
	case models.NotFoundError:
		code = 404
	case models.MethodError:
		code = 405
	case models.ElasticSearchError:
		code = 502
	default:
		code = 500
	}

	response, _ := json.Marshal(payload)
	resp := string(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(resp))
}

// ResponseWithJson returns formed JSON
func ResponseWithCustomCode(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	resp := string(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(resp))
}
