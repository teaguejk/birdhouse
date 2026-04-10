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
	publicSet := make(map[string]bool, len(publicRoutes))
	for _, route := range publicRoutes {
		publicSet[route] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			routeKey := r.Method + " " + r.URL.Path
			if publicSet[routeKey] {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				responses.Challenge(w)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				responses.TokenError(w, "invalid authorization header format")
				return
			}

			apiKey := parts[1]

			device, err := deviceService.Authenticate(r.Context(), apiKey)
			if err != nil {
				responses.TokenError(w, "invalid or inactive API key")
				return
			}

			ctx := context.WithValue(r.Context(), constants.DeviceIDKey, device.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
