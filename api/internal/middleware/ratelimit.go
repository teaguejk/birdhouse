package middleware

import "net/http"

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: rate limit logic if needed
		next.ServeHTTP(w, r)
	})
}
