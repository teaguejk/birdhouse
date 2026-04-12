package admin

import (
	"api/internal/api/interfaces"
	"api/internal/shared/models"
	"api/pkg/logging"
	"context"
	"fmt"
)

type Service struct {
	repo   interfaces.AdminRepository
	logger *logging.Logger
}

func NewService(logger *logging.Logger, repo interfaces.AdminRepository) interfaces.AdminService {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) ValidateAdmin(ctx context.Context, email string) (*models.AdminUser, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to validate admin: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("unauthorized: %s is not an authorized admin", email)
	}

	return user, nil
}
