package server

import (
	"api/internal/api/interfaces"
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (s *Server) RegisterHandler(handler interfaces.Handler) error {
	if err := handler.Initialize(context.Background()); err != nil {
		return fmt.Errorf("failed to initialize service: %w", err)
	}

	handlerMux := http.NewServeMux()

	handler.RegisterRoutes(handlerMux)

	s.router.Handle(fmt.Sprintf("/%s/", handler.Name()), http.StripPrefix(fmt.Sprintf("/%s", handler.Name()), handlerMux))
	s.router.Handle(fmt.Sprintf("/%s", handler.Name()), AddSuffix("/", http.StripPrefix(fmt.Sprintf("/%s", handler.Name()), handlerMux)))

	s.handlers = append(s.handlers, handler)

	s.env.Logger.Info(fmt.Sprintf("registered handler: %s", handler.Name()))

	return nil
}

func AddSuffix(suffix string, h http.Handler) http.Handler {
	if suffix == "" {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = r.URL.Path + suffix
		if r.URL.RawPath != "" {
			r2.URL.RawPath = r.URL.RawPath + suffix
		}
		h.ServeHTTP(w, r2)
	})
}
