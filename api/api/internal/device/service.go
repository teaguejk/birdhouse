package device

import (
	"api/internal/api/interfaces"
	"api/internal/shared/models"
	"api/pkg/logging"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Service struct {
	repo   interfaces.DeviceRepository
	logger *logging.Logger
}

func NewService(logger *logging.Logger, repo interfaces.DeviceRepository) interfaces.DeviceService {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) Create(ctx context.Context, req *models.CreateDeviceRequest) (*models.CreateDeviceResponse, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("device name is required")
	}

	rawKey, hash := generateAPIKey()

	device := &models.Device{
		Name:       req.Name,
		APIKeyHash: hash,
		Location:   req.Location,
		Active:     true,
	}

	if err := s.repo.Create(ctx, device); err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	return &models.CreateDeviceResponse{
		Device: device,
		APIKey: rawKey,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*models.Device, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Authenticate(ctx context.Context, apiKey string) (*models.Device, error) {
	hash := hashKey(apiKey)

	device, err := s.repo.GetByAPIKeyHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	if device == nil {
		return nil, fmt.Errorf("invalid API key")
	}

	if !device.Active {
		return nil, fmt.Errorf("device is inactive")
	}

	return device, nil
}

func (s *Service) List(ctx context.Context) ([]models.Device, error) {
	return s.repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, id string, req *models.UpdateDeviceRequest) (*models.Device, error) {
	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		device.Name = *req.Name
	}
	if req.Location != nil {
		device.Location = *req.Location
	}
	if req.Active != nil {
		device.Active = *req.Active
	}

	if err := s.repo.Update(ctx, device); err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}

	return device, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) RotateKey(ctx context.Context, id string) (*models.RotateKeyResponse, error) {
	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	rawKey, hash := generateAPIKey()
	device.APIKeyHash = hash

	if err := s.repo.Update(ctx, device); err != nil {
		return nil, fmt.Errorf("failed to rotate key: %w", err)
	}

	return &models.RotateKeyResponse{APIKey: rawKey}, nil
}

func generateAPIKey() (raw string, hash string) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("failed to generate random bytes: %v", err))
	}
	raw = "bh_" + hex.EncodeToString(b)
	hash = hashKey(raw)
	return
}

func hashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}
