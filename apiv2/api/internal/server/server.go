package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"api/internal/api"
	"api/internal/api/interfaces"
)

type Server struct {
	env        *api.ServerEnv
	httpServer *http.Server
	router     *http.ServeMux
	handlers   []interfaces.Handler
	services   *api.Services
}

func New(env *api.ServerEnv, services *api.Services) *Server {
	mux := http.NewServeMux()

	server := &Server{
		env:      env,
		router:   mux,
		handlers: make([]interfaces.Handler, 0),
		services: services,
	}

	handler := server.setupMiddleware(mux)

	// setup server
	server.httpServer = &http.Server{
		Addr:         ":" + env.Config.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(env.Config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(env.Config.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(env.Config.IdleTimeout) * time.Second,
	}

	return server
}

// starts the HTTP server
func (s *Server) ListenAndServe() error {
	s.env.Logger.Info(fmt.Sprintf("starting server on port %s", s.env.Config.Port))
	return s.httpServer.ListenAndServe()
}

// gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.env.Logger.Info("server shutting down...")

	for _, handler := range s.handlers {
		if err := handler.Shutdown(ctx); err != nil {
			s.env.Logger.Info(fmt.Sprintf("error shutting down: %s: %v", handler.Name(), err))
		}
	}

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Handler() http.Handler {
	return s.router
}
