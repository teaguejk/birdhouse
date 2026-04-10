package upload

import (
	"api/internal/api/interfaces"
	"api/internal/shared/constants"
	"api/internal/shared/models"
	"api/internal/shared/responses"
	"api/pkg/logging"
	"api/pkg/storage"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const Name = "uploads"

type Handler struct {
	logger  *logging.Logger
	service interfaces.UploadService
}

func NewHandler(logger *logging.Logger, service interfaces.UploadService) interfaces.UploadHandler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) Name() string {
	return Name
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// mux.HandleFunc("POST /", h.Create)
	mux.HandleFunc("POST /init", h.GenerateUploadURL)
	mux.HandleFunc("POST /complete", h.Complete)

	mux.HandleFunc("GET /meta/{id}", h.Get)
	mux.HandleFunc("GET /meta/{resource_type}/{resource_id}", h.GetByResource)

	mux.HandleFunc("DELETE /{id}", h.Delete)
}

func (h *Handler) Initialize(ctx context.Context) error {
	h.logger.Info(fmt.Sprintf("initializing: %s", h.Name()))
	return nil
}

func (h *Handler) Shutdown(ctx context.Context) error {
	h.logger.Info(fmt.Sprintf("shutting down: %s", h.Name()))
	return nil
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	imageID := r.PathValue("id")
	if imageID == "" {
		responses.WriteError(w, "file ID is required", http.StatusBadRequest)
		return
	}

	image, err := h.service.GetByID(r.Context(), imageID)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to get file: %v", err))
		responses.WriteError(w, "failed to get file", http.StatusInternalServerError)
		return
	}

	if image == nil {
		responses.WriteError(w, "file not found", http.StatusNotFound)
		return
	}

	responses.WriteJSON(w, image, http.StatusOK)
}

func (h *Handler) GetByResource(w http.ResponseWriter, r *http.Request) {
	resourceType := r.PathValue("resource_type")
	resourceID := r.PathValue("resource_id")
	if resourceType == "" || resourceID == "" {
		responses.WriteError(w, "resource type and ID are required", http.StatusBadRequest)
		return
	}

	images, err := h.service.GetByResource(r.Context(), resourceType, resourceID, true)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to get images by resource ID: %v", err))
		responses.WriteError(w, "failed to get images", http.StatusInternalServerError)
		return
	}

	responses.WriteJSON(w, images, http.StatusOK)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(constants.UserIDKey).(string)
	if uid == "" {
		responses.WriteError(w, "can not determine requesting user", http.StatusBadRequest)
		return
	}

	imageID := r.PathValue("id")
	if imageID == "" {
		responses.WriteError(w, "image ID is required", http.StatusBadRequest)
		return
	}

	err := h.service.Delete(r.Context(), uid, imageID)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to delete image: %v", err))
		responses.WriteError(w, "failed to delete image", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Complete(w http.ResponseWriter, r *http.Request) {
	var request models.UploadCompleteRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responses.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if request.Key == "" {
		responses.WriteError(w, "upload key is required", http.StatusBadRequest)
		return
	}

	err := h.service.Complete(r.Context(), request.Key)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to complete upload: %v", err))
		responses.WriteError(w, "failed to complete upload", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GenerateUploadURL(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(constants.UserIDKey).(string)
	if uid == "" {
		responses.WriteError(w, "can not determine requesting user", http.StatusBadRequest)
		return
	}

	var request *models.FileUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responses.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(request.Filenames) == 0 {
		responses.WriteError(w, "filenames required", http.StatusBadRequest)
		return
	}

	response := make(models.UploadDetailsResponse)

	for _, filename := range request.Filenames {
		validator := storage.NewImageValidator()
		ext, mimeType, err := validator.ValidateFileName(filename)
		if err != nil {
			responses.WriteError(w, fmt.Sprintf("failed to validate filename: %v", err), http.StatusBadRequest)
			return
		}

		url, uploadKey, err := h.service.GenerateUploadURL(r.Context(), filename, ext)
		if err != nil {
			responses.WriteError(w, fmt.Sprintf("failed to generate upload URL: %v", err), http.StatusInternalServerError)
			return
		}

		err = h.service.Create(r.Context(), &models.File{
			UserID:       uid,
			Filename:     uploadKey,
			OriginalName: filename,
			MimeType:     mimeType,
			Size:         0,
			Status:       "pending",
			URL:          "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
		if err != nil {
			responses.WriteError(w, fmt.Sprintf("failed to write upload metadata: %v", err), http.StatusInternalServerError)
			return
		}

		response[filename] = models.UploadDetails{
			URL: url,
			Key: uploadKey,
		}
	}

	responses.WriteJSON(w, response, http.StatusOK)
}
