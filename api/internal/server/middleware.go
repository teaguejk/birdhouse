package server

import (
	"api/internal/middleware"
	"net/http"
)

func (s *Server) setupMiddleware(handler http.Handler) http.Handler {
	h := middleware.AuthMiddleware(&middleware.AuthConfig{
		DeviceService:       s.services.Device,
		AdminService:        s.services.Admin,
		OAuthVerifier:       s.env.OAuthVerifier,
		PublicRoutes:        s.env.Config.PublicRoutes,
		AuthRoutes:          s.env.Config.AuthRoutes,
		AdminRoutes:         s.env.Config.AdminRoutes,
		DeviceLenientRoutes: s.env.Config.DeviceLenientRoutes,
	})(handler)
	h = middleware.CorsMiddleware(s.env.Config.CORS, h)
	h = middleware.LoggingMiddleware(s.env.Logger, h)
	h = middleware.RecoveryMiddleware(s.env.Logger, h)

	if s.env.Config.RateLimitEnabled {
		h = middleware.RateLimitMiddleware(h)
	}

	return h
}
