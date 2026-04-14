package middleware

import (
	"api/internal/api/interfaces"
	"api/internal/shared/constants"
	"api/internal/shared/responses"
	"api/pkg/oauth"
	"context"
	"net/http"
	"strings"
)

type AuthConfig struct {
	DeviceService       interfaces.DeviceService
	AdminService        interfaces.AdminService
	OAuthVerifier       oauth.TokenVerifier
	PublicRoutes        []string
	AuthRoutes          []string
	AdminRoutes         []string
	DeviceLenientRoutes []string
}

func AuthMiddleware(cfg *AuthConfig) func(http.Handler) http.Handler {
	publicExact, publicPrefix := parseRoutes(cfg.PublicRoutes)
	authExact, authPrefix := parseRoutes(cfg.AuthRoutes)
	adminExact, adminPrefix := parseRoutes(cfg.AdminRoutes)
	lenientExact, lenientPrefix := parseRoutes(cfg.DeviceLenientRoutes)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			routeKey := r.Method + " " + r.URL.Path

			// 1. public routes — no auth
			if matchesRoute(routeKey, publicExact, publicPrefix) {
				next.ServeHTTP(w, r)
				return
			}

			// 2. authenticated routes — valid oauth token, no admin check
			if matchesRoute(routeKey, authExact, authPrefix) {
				claims, ok := verifyOAuthToken(r, cfg.OAuthVerifier, w)
				if !ok {
					return
				}

				ctx := context.WithValue(r.Context(), constants.OAuthClaimsKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// 3. device lenient routes — api key auth, allows inactive devices
			// checked before admin so exact matches here win over admin wildcards
			if matchesRoute(routeKey, lenientExact, lenientPrefix) {
				apiKey := r.Header.Get("x-api-key")
				if apiKey == "" {
					responses.ChallengeApiKey(w)
					return
				}

				device, err := cfg.DeviceService.AuthenticateAllowInactive(r.Context(), apiKey)
				if err != nil {
					responses.ApiKeyError(w)
					return
				}

				ctx := context.WithValue(r.Context(), constants.DeviceIDKey, device.ID)
				ctx = context.WithValue(ctx, constants.DeviceActiveKey, device.Active)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// 4. admin routes — valid oauth token + admin check
			if matchesRoute(routeKey, adminExact, adminPrefix) {
				claims, ok := verifyOAuthToken(r, cfg.OAuthVerifier, w)
				if !ok {
					return
				}

				_, err := cfg.AdminService.ValidateAdmin(r.Context(), claims.Email)
				if err != nil {
					responses.WriteError(w, "forbidden", http.StatusForbidden)
					return
				}

				ctx := context.WithValue(r.Context(), constants.OAuthClaimsKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// 5. everything else — device api key (rejects inactive)
			apiKey := r.Header.Get("x-api-key")
			if apiKey == "" {
				responses.ChallengeApiKey(w)
				return
			}

			device, err := cfg.DeviceService.Authenticate(r.Context(), apiKey)
			if err != nil {
				responses.ApiKeyError(w)
				return
			}

			ctx := context.WithValue(r.Context(), constants.DeviceIDKey, device.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func verifyOAuthToken(r *http.Request, verifier oauth.TokenVerifier, w http.ResponseWriter) (*oauth.Claims, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		responses.ChallengeBearer(w)
		return nil, false
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		responses.TokenError(w, "invalid authorization header format")
		return nil, false
	}

	claims, err := verifier.Verify(r.Context(), parts[1])
	if err != nil {
		responses.TokenError(w, "invalid or expired token")
		return nil, false
	}

	return claims, true
}

func parseRoutes(routes []string) (exact map[string]bool, prefixes []string) {
	exact = make(map[string]bool)
	for _, route := range routes {
		if strings.HasSuffix(route, "*") {
			prefixes = append(prefixes, strings.TrimSuffix(route, "*"))
		} else {
			exact[route] = true
		}
	}
	return
}

func matchesRoute(routeKey string, exact map[string]bool, prefixes []string) bool {
	if exact[routeKey] {
		return true
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(routeKey, prefix) {
			return true
		}
	}
	return false
}
