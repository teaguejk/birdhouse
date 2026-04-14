package command

import (
	"api/internal/api/interfaces"
	"api/internal/shared/constants"
	"api/internal/shared/models"
	"api/internal/shared/responses"
	"api/pkg/logging"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const Name = "commands"

type Handler struct {
	logger  *logging.Logger
	service interfaces.CommandService
}

func NewHandler(logger *logging.Logger, service interfaces.CommandService) interfaces.CommandHandler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) Name() string {
	return Name
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// admin creates commands for a device
	mux.HandleFunc("POST /device/{device_id}", h.Create)

	// device polls for its pending commands
	mux.HandleFunc("GET /pending", h.ListPending)

	// device acknowledges a command
	mux.HandleFunc("POST /{id}/ack", h.Acknowledge)
}

func (h *Handler) Initialize(ctx context.Context) error {
	h.logger.Info(fmt.Sprintf("initializing: %s", h.Name()))
	return nil
}

func (h *Handler) Shutdown(ctx context.Context) error {
	h.logger.Info(fmt.Sprintf("shutting down: %s", h.Name()))
	return nil
}

// Create is called by the admin UI to queue a command for a device.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("device_id")
	if deviceID == "" {
		responses.WriteError(w, "device ID is required", http.StatusBadRequest)
		return
	}

	var req models.CreateCommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Action == "" {
		responses.WriteError(w, "action is required", http.StatusBadRequest)
		return
	}

	cmd, err := h.service.Create(r.Context(), deviceID, &req)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to create command: %v", err))
		responses.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	responses.WriteJSON(w, cmd, http.StatusCreated)
}

// ListPending is called by the device to poll for commands to execute.
func (h *Handler) ListPending(w http.ResponseWriter, r *http.Request) {
	deviceID, ok := r.Context().Value(constants.DeviceIDKey).(string)
	if !ok || deviceID == "" {
		responses.WriteError(w, "device not authenticated", http.StatusUnauthorized)
		return
	}

	commands, err := h.service.ListPending(r.Context(), deviceID)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to list pending commands: %v", err))
		responses.WriteError(w, "failed to list commands", http.StatusInternalServerError)
		return
	}

	if commands == nil {
		commands = []models.Command{}
	}

	responses.WriteJSON(w, commands, http.StatusOK)
}

// Acknowledge is called by the device after executing a command.
func (h *Handler) Acknowledge(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		responses.WriteError(w, "command ID is required", http.StatusBadRequest)
		return
	}

	if err := h.service.Acknowledge(r.Context(), id); err != nil {
		h.logger.Error(fmt.Sprintf("failed to acknowledge command: %v", err))
		responses.WriteError(w, "failed to acknowledge command", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
