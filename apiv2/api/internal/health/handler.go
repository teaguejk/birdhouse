package health

import (
	"api/internal/shared/responses"
	"api/pkg/database"
	"api/pkg/logging"
	"context"
	"fmt"
	"net/http"
)

const Name = "health"

type Handler struct {
	logger *logging.Logger
	db     database.Database
}

func NewHandler(logger *logging.Logger, db database.Database) *Handler {
	return &Handler{
		logger: logger,
		db:     db,
	}
}

func (h *Handler) Name() string {
	return Name
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthcheck", h.HealthCheck)
}

func (h *Handler) Initialize(ctx context.Context) error {
	h.logger.Info(fmt.Sprintf("initializing: %s", h.Name()))
	return nil
}

func (h *Handler) Shutdown(ctx context.Context) error {
	h.logger.Info(fmt.Sprintf("shutting down: %s", h.Name()))
	return nil
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	err := h.db.Ping(r.Context())
	if err != nil {
		message := fmt.Sprintf("database was unreachable %s", err)
		responses.WriteError(w, message, http.StatusServiceUnavailable)
		return
	}

	response := make(map[string]any)
	response["server_status"] = "server is active and accepting requests"
	response["db_status"] = "database is active and accepting connections"

	responses.WriteJSON(w, response, http.StatusOK)
}
