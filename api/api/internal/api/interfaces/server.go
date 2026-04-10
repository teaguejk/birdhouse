package interfaces

import (
	"context"
	"net/http"
)

type RouteRegistrar interface {
	RegisterRoutes(mux *http.ServeMux)
}

type Handler interface {
	RouteRegistrar
	Name() string
	Initialize(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
