package server

import (
	"api/internal/middleware"
	"net/http"
)

func (s *Server) setupMiddleware(handler http.Handler) http.Handler {
	h := middleware.AuthMiddleware(s.services.Device, s.env.Config.PublicRoutes)(handler)
	h = middleware.CorsMiddleware(s.env.Config.CORS, h)
	h = middleware.LoggingMiddleware(s.env.Logger, h)
	h = middleware.RecoveryMiddleware(s.env.Logger, h)

	if s.env.Config.RateLimitEnabled {
		h = middleware.RateLimitMiddleware(h)
	}

	return h
}
