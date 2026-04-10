package middleware

import (
	"api/internal/shared/responses"
	"api/pkg/logging"
	"fmt"
	"net/http"
	"time"
)

func LoggingMiddleware(logger *logging.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// create a response writer wrapper to capture status code
		wrapped := &responses.ResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		logger.Info(fmt.Sprintf("req: %s %s %d %v",
			r.Method, r.URL.Path, wrapped.StatusCode, duration))
	})
}
