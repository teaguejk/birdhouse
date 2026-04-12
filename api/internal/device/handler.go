package device

import (
	"api/internal/api/interfaces"
	"api/internal/shared/models"
	"api/internal/shared/responses"
	"api/pkg/logging"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const Name = "devices"

type Handler struct {
	logger  *logging.Logger
	service interfaces.DeviceService
}

func NewHandler(logger *logging.Logger, service interfaces.DeviceService) interfaces.DeviceHandler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) Name() string {
	return Name
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /", h.Create)
	mux.HandleFunc("GET /", h.List)
	mux.HandleFunc("GET /status", h.Status)
	mux.HandleFunc("GET /{id}", h.Get)
	mux.HandleFunc("PUT /{id}", h.Update)
	mux.HandleFunc("DELETE /{id}", h.Delete)
	mux.HandleFunc("POST /{id}/rotate-key", h.RotateKey)
}

func (h *Handler) Initialize(ctx context.Context) error {
	h.logger.Info(fmt.Sprintf("initializing: %s", h.Name()))
	return nil
}

func (h *Handler) Shutdown(ctx context.Context) error {
	h.logger.Info(fmt.Sprintf("shutting down: %s", h.Name()))
	return nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		responses.WriteError(w, "device name is required", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Create(r.Context(), &req)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to create device: %v", err))
		responses.WriteError(w, "failed to create device", http.StatusInternalServerError)
		return
	}

	responses.WriteJSON(w, resp, http.StatusCreated)
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	statuses, err := h.service.ListStatus(r.Context())
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to list device status: %v", err))
		responses.WriteError(w, "failed to list device status", http.StatusInternalServerError)
		return
	}

	responses.WriteJSON(w, statuses, http.StatusOK)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		responses.WriteError(w, "device ID is required", http.StatusBadRequest)
		return
	}

	device, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to get device: %v", err))
		responses.WriteError(w, "device not found", http.StatusNotFound)
		return
	}

	responses.WriteJSON(w, device, http.StatusOK)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	devices, err := h.service.List(r.Context())
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to list devices: %v", err))
		responses.WriteError(w, "failed to list devices", http.StatusInternalServerError)
		return
	}

	responses.WriteJSON(w, devices, http.StatusOK)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		responses.WriteError(w, "device ID is required", http.StatusBadRequest)
		return
	}

	var req models.UpdateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	device, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to update device: %v", err))
		responses.WriteError(w, "failed to update device", http.StatusInternalServerError)
		return
	}

	responses.WriteJSON(w, device, http.StatusOK)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		responses.WriteError(w, "device ID is required", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.logger.Error(fmt.Sprintf("failed to delete device: %v", err))
		responses.WriteError(w, "failed to delete device", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RotateKey(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		responses.WriteError(w, "device ID is required", http.StatusBadRequest)
		return
	}

	resp, err := h.service.RotateKey(r.Context(), id)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to rotate key: %v", err))
		responses.WriteError(w, "failed to rotate key", http.StatusInternalServerError)
		return
	}

	responses.WriteJSON(w, resp, http.StatusOK)
}
