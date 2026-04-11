package interfaces

import (
	"api/internal/shared/models"
	"context"
	"io"
	"net/http"
)

type UploadRepository interface {
	Assign(ctx context.Context, resourceType, resourceID string, filenames []string) error
	UpdateSortOrder(ctx context.Context, resourceType, resourceID string, filenames []string) error
	Complete(ctx context.Context, uploadKey string) error
	Create(ctx context.Context, image *models.File) error
	Delete(ctx context.Context, id string) error
	GetByResource(ctx context.Context, resourceType, resourceID string, assignedOnly bool) ([]models.File, error)
	GetByID(ctx context.Context, id string) (*models.File, error)
	GetByFilename(ctx context.Context, filename string) (*models.File, error)
	GetLatest(ctx context.Context) (*models.File, error)
	GetExpiredPending(ctx context.Context) ([]models.File, error)
}

type UploadService interface {
	Assign(ctx context.Context, resourceType, resourceID string, filenames []string) error
	UpdateSortOrder(ctx context.Context, resourceType, resourceID string, filenames []string) error
	Complete(ctx context.Context, uploadKey string) error
	Create(ctx context.Context, image *models.File) error
	Delete(ctx context.Context, deviceID, id string) error
	DeleteByResource(ctx context.Context, deviceID, resourceType, resourceID string) error
	DeleteByFilename(ctx context.Context, deviceID, filename string) error
	GetByResource(ctx context.Context, resourceType, resourceID string, assignedOnly bool) ([]models.File, error)
	GetByID(ctx context.Context, id string) (*models.File, error)
	GetByFilename(ctx context.Context, filename string) (*models.File, error)
	GetLatest(ctx context.Context) (*models.File, error)
	GetExpiredPending(ctx context.Context) ([]models.File, error)

	// storage
	GenerateUploadURL(ctx context.Context, filename, ext string) (string, string, error)
	GetFile(ctx context.Context, path string) (io.ReadCloser, error)
	CopyFile(ctx context.Context, srcPath, destPath string) error
}

type UploadHandler interface {
	Handler
	GenerateUploadURL(w http.ResponseWriter, r *http.Request)
	Complete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	GetByResource(w http.ResponseWriter, r *http.Request)
	GetLatest(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
