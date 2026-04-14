package interfaces

import (
	"api/internal/shared/models"
	"context"
	"net/http"
)

type DeviceRepository interface {
	Create(ctx context.Context, device *models.Device) error
	GetByID(ctx context.Context, id string) (*models.Device, error)
	GetByAPIKeyHash(ctx context.Context, hash string) (*models.Device, error)
	List(ctx context.Context) ([]models.Device, error)
	ListStatus(ctx context.Context) ([]models.DeviceStatus, error)
	TouchLastSeen(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status []byte) error
	Update(ctx context.Context, device *models.Device) error
	Delete(ctx context.Context, id string) error
}

type DeviceService interface {
	Create(ctx context.Context, req *models.CreateDeviceRequest) (*models.CreateDeviceResponse, error)
	GetByID(ctx context.Context, id string) (*models.Device, error)
	Authenticate(ctx context.Context, apiKey string) (*models.Device, error)
	AuthenticateAllowInactive(ctx context.Context, apiKey string) (*models.Device, error)
	List(ctx context.Context) ([]models.Device, error)
	ListStatus(ctx context.Context) ([]models.DeviceStatus, error)
	Update(ctx context.Context, id string, req *models.UpdateDeviceRequest) (*models.Device, error)
	Delete(ctx context.Context, id string) error
	RotateKey(ctx context.Context, id string) (*models.RotateKeyResponse, error)
}

type DeviceHandler interface {
	Handler
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	GetConfig(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Status(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	RotateKey(w http.ResponseWriter, r *http.Request)
}
