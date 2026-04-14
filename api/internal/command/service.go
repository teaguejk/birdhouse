package command

import (
	"api/internal/api/interfaces"
	"api/internal/shared/models"
	"api/pkg/logging"
	"api/pkg/mqtt"
	"context"
	"encoding/json"
	"fmt"
)

var validActions = map[string]bool{
	"start_detection": true,
	"stop_detection":  true,
	"capture":         true,
	"update_config":   true,
}

type Service struct {
	repo      interfaces.CommandRepository
	logger    *logging.Logger
	publisher mqtt.Publisher
}

func NewService(logger *logging.Logger, repo interfaces.CommandRepository, publisher mqtt.Publisher) interfaces.CommandService {
	return &Service{
		repo:      repo,
		logger:    logger,
		publisher: publisher,
	}
}

func (s *Service) Create(ctx context.Context, deviceID string, req *models.CreateCommandRequest) (*models.Command, error) {
	if !validActions[req.Action] {
		return nil, fmt.Errorf("invalid action: %s", req.Action)
	}

	payload := req.Payload
	if payload == nil {
		payload = json.RawMessage(`{}`)
	}

	cmd := &models.Command{
		DeviceID: deviceID,
		Action:   req.Action,
		Payload:  payload,
	}

	if err := s.repo.Create(ctx, cmd); err != nil {
		return nil, fmt.Errorf("failed to create command: %w", err)
	}

	// publish to MQTT
	topic := fmt.Sprintf("birdhouse/%s/commands", deviceID)
	if err := s.publisher.Publish(topic, 1, false, cmd); err != nil {
		s.logger.Error(fmt.Sprintf("failed to publish command to mqtt: %v", err))
	}

	return cmd, nil
}

func (s *Service) ListPending(ctx context.Context, deviceID string) ([]models.Command, error) {
	return s.repo.ListPending(ctx, deviceID)
}

func (s *Service) Acknowledge(ctx context.Context, id string) error {
	return s.repo.Acknowledge(ctx, id)
}
