package interfaces

import (
	"api/internal/shared/models"
	"context"
)

type AdminRepository interface {
	GetByEmail(ctx context.Context, email string) (*models.AdminUser, error)
}

type AdminService interface {
	ValidateAdmin(ctx context.Context, email string) (*models.AdminUser, error)
}
