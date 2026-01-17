package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
)

type GraphqlAdapter func(http.Handler) http.Handler

// GraphqlAdapt Iterate over adapters and run them one by one, meant for graphql servers
func GraphqlAdapt(h http.Handler, adapters ...GraphqlAdapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func SetCorsHeadersLocal() GraphqlAdapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && isAllowedLocalOrigin(origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set(
					"Access-Control-Allow-Headers",
					"Origin, X-Requested-With, Content-Type, Accept, Authorization, Boule",
				)

				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
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

func LogGraphql() GraphqlAdapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			// Skip noisy endpoints
			blockList := []string{"health", "ping"}
			for _, block := range blockList {
				if strings.Contains(r.URL.Path, block) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Wrap writer to capture status
			rec := &StatusRecorder{
				ResponseWriter: w,
				Status:         http.StatusOK, // default if handler never calls WriteHeader
			}

			requestId := r.Header.Get(service.HeaderKey)

			// Call the wrapped handler
			next.ServeHTTP(rec, r)

			// Log
			duration := time.Since(startTime)
			clientIp := r.RemoteAddr
			method := r.Method
			path := r.URL.Path

			logging.Api(rec.Status, method, requestId, clientIp, path, duration)
		})
	}
}
