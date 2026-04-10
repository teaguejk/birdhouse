package upload

import (
	"api/internal/api/interfaces"
	"api/internal/shared/models"
	"api/pkg/logging"
	"api/pkg/storage"
	"api/pkg/utils"
	"context"
	"errors"
	"fmt"
	"io"
	"time"
)

var (
	ErrTokenNotFound = errors.New("bearer token not found")
	ErrInvalidToken  = errors.New("bearer token is invalid")
	ErrExpiredToken  = errors.New("bearer token is expired")

	ErrInvalidFile     = errors.New("invalid file")
	ErrInvalidFileType = errors.New("invalid file type")
	ErrFileNotFound    = errors.New("file not found")
)

type Service struct {
	repo    interfaces.UploadRepository
	storage storage.Provider
	logger  *logging.Logger
}

func NewService(logger *logging.Logger, repo interfaces.UploadRepository, storage storage.Provider) interfaces.UploadService {
	return &Service{
		repo:    repo,
		storage: storage,
		logger:  logger,
	}
}

func (s *Service) Assign(ctx context.Context, resourceType, resourceID string, filenames []string) error {
	if len(filenames) == 0 {
		return fmt.Errorf("no filenames provided")
	}

	err := s.repo.Assign(ctx, resourceType, resourceID, filenames)
	if err != nil {
		return fmt.Errorf("failed to assign images: %w", err)
	}

	return nil
}

func (s *Service) UpdateSortOrder(ctx context.Context, resourceType, resourceID string, filenames []string) error {
	if len(filenames) == 0 {
		return nil
	}

	err := s.repo.UpdateSortOrder(ctx, resourceType, resourceID, filenames)
	if err != nil {
		return fmt.Errorf("failed to update sort order: %w", err)
	}

	return nil
}

func (s *Service) Complete(ctx context.Context, uploadKey string) error {
	if uploadKey == "" {
		return fmt.Errorf("upload key is required")
	}

	return s.repo.Complete(ctx, uploadKey)
}

func (s *Service) Create(ctx context.Context, image *models.File) error {
	if image == nil {
		return fmt.Errorf("image cannot be nil")
	}

	if image.Filename == "" || image.MimeType == "" || image.Size < 0 {
		return fmt.Errorf("invalid image data")
	}

	image.URL = s.storage.GetURL(image.Filename)

	err := s.repo.Create(ctx, image)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, userID, id string) error {
	image, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to retrieve image: %w", err)
	}

	if image.UserID != userID {
		return fmt.Errorf("unauthorized: user does not own the image")
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteByResource(ctx context.Context, userID, resourceType, resourceID string) error {
	images, err := s.repo.GetByResource(ctx, resourceType, resourceID, false)
	if err != nil {
		return fmt.Errorf("failed to retrieve images: %w", err)
	}

	for _, image := range images {
		if image.UserID != userID {
			return fmt.Errorf("unauthorized: user does not own all images")
		}

		err = s.repo.Delete(ctx, image.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) DeleteByFilename(ctx context.Context, userID, filename string) error {
	image, err := s.repo.GetByFilename(ctx, filename)
	if err != nil {
		return fmt.Errorf("failed to retrieve images: %w", err)
	}

	if image.UserID != userID {
		return fmt.Errorf("unauthorized: user does not own the image")
	}

	err = s.repo.Delete(ctx, image.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetByResource(ctx context.Context, resourceType, resourceID string, assignedOnly bool) ([]models.File, error) {
	images, err := s.repo.GetByResource(ctx, resourceType, resourceID, assignedOnly)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*models.File, error) {
	image, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (s *Service) GetByFilename(ctx context.Context, filename string) (*models.File, error) {
	image, err := s.repo.GetByFilename(ctx, filename)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (s *Service) GetExpiredPending(ctx context.Context) ([]models.File, error) {
	images, err := s.repo.GetExpiredPending(ctx)
	if err != nil {
		return nil, err
	}

	for i, img := range images {
		images[i].URL = s.storage.GetURL(img.Filename)
	}

	return images, nil
}

// storage

func (s *Service) GenerateUploadURL(ctx context.Context, filename string, ext string) (string, string, error) {
	uploadKey := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), utils.RandomHash(filename), ext)

	url, err := s.storage.GetSignedURL(ctx, uploadKey, time.Second*60, "PUT")
	if err != nil {
		return "", "", err
	}

	return url, uploadKey, nil
}

func (s *Service) GetFile(ctx context.Context, path string) (io.ReadCloser, error) {
	file, err := s.storage.GetFile(ctx, path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *Service) CopyFile(ctx context.Context, srcPath, destPath string) error {
	return s.storage.Copy(ctx, srcPath, destPath)
}
