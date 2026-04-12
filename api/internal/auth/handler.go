package auth

import (
	"api/internal/api/interfaces"
	"api/internal/shared/constants"
	"api/internal/shared/responses"
	"api/pkg/logging"
	"api/pkg/oauth"
	"context"
	"fmt"
	"net/http"
)

const Name = "auth"

type Handler struct {
	logger       *logging.Logger
	adminService interfaces.AdminService
}

func NewHandler(logger *logging.Logger, adminService interfaces.AdminService) interfaces.AuthHandler {
	return &Handler{
		logger:       logger,
		adminService: adminService,
	}
}

func (h *Handler) Name() string {
	return Name
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /me", h.Me)
}

func (h *Handler) Initialize(ctx context.Context) error {
	h.logger.Info(fmt.Sprintf("initializing: %s", h.Name()))
	return nil
}

func (h *Handler) Shutdown(ctx context.Context) error {
	h.logger.Info(fmt.Sprintf("shutting down: %s", h.Name()))
	return nil
}

type MeResponse struct {
	Subject string `json:"subject"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(constants.OAuthClaimsKey).(*oauth.Claims)
	if !ok || claims == nil {
		responses.WriteError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	isAdmin := false
	admin, err := h.adminService.ValidateAdmin(r.Context(), claims.Email)
	if err == nil && admin != nil {
		isAdmin = true
	}

	responses.WriteJSON(w, MeResponse{
		Subject: claims.Subject,
		Email:   claims.Email,
		Name:    claims.Name,
		IsAdmin: isAdmin,
	}, http.StatusOK)
}
