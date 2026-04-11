package middleware

import (
	"api/internal/api/interfaces"
	"api/internal/shared/constants"
	"api/internal/shared/responses"
	"context"
	"net/http"
	"strings"
)

func AuthMiddleware(deviceService interfaces.DeviceService, publicRoutes []string) func(http.Handler) http.Handler {
	exactRoutes := make(map[string]bool)
	var prefixRoutes []string

	for _, route := range publicRoutes {
		if strings.HasSuffix(route, "*") {
			prefixRoutes = append(prefixRoutes, strings.TrimSuffix(route, "*"))
		} else {
			exactRoutes[route] = true
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			routeKey := r.Method + " " + r.URL.Path
			if exactRoutes[routeKey] {
				next.ServeHTTP(w, r)
				return
			}

			for _, prefix := range prefixRoutes {
				if strings.HasPrefix(routeKey, prefix) {
					next.ServeHTTP(w, r)
					return
				}
			}

			apiKey := r.Header.Get("x-api-key")
			if apiKey == "" {
				responses.ChallengeApiKey(w)
				return
			}

			device, err := deviceService.Authenticate(r.Context(), apiKey)
			if err != nil {
				responses.ApiKeyError(w)
				return
			}

			ctx := context.WithValue(r.Context(), constants.DeviceIDKey, device.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
