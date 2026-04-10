package middleware

import (
	"fmt"
	"net/http"

	"api/internal/shared/responses"
	"api/pkg/logging"
)

func RecoveryMiddleware(logger *logging.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(fmt.Sprintf("panic recovered: %v", err))
				responses.WriteError(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
