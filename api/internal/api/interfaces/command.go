package interfaces

import (
	"api/internal/shared/models"
	"context"
	"net/http"
)

type CommandRepository interface {
	Create(ctx context.Context, cmd *models.Command) error
	ListPending(ctx context.Context, deviceID string) ([]models.Command, error)
	Acknowledge(ctx context.Context, id string) error
}

type CommandService interface {
	Create(ctx context.Context, deviceID string, req *models.CreateCommandRequest) (*models.Command, error)
	ListPending(ctx context.Context, deviceID string) ([]models.Command, error)
	Acknowledge(ctx context.Context, id string) error
}

type CommandHandler interface {
	Handler
	Create(w http.ResponseWriter, r *http.Request)
	ListPending(w http.ResponseWriter, r *http.Request)
	Acknowledge(w http.ResponseWriter, r *http.Request)
}
