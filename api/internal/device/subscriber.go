package device

import (
	"context"
	"fmt"
	"strings"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"api/internal/api/interfaces"
	"api/pkg/logging"
	"api/pkg/mqtt"
)

type StatusSubscriber struct {
	repo   interfaces.DeviceRepository
	logger *logging.Logger
}

func NewStatusSubscriber(logger *logging.Logger, repo interfaces.DeviceRepository) *StatusSubscriber {
	return &StatusSubscriber{
		repo:   repo,
		logger: logger,
	}
}

// Subscribe listens on birdhouse/+/status for all devices.
func (s *StatusSubscriber) Subscribe(client *mqtt.Client) error {
	return client.Subscribe("birdhouse/+/status", 0, s.handleStatus)
}

func (s *StatusSubscriber) handleStatus(_ pahomqtt.Client, msg pahomqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	if len(parts) != 3 {
		s.logger.Warn(fmt.Sprintf("unexpected status topic format: %s", msg.Topic()))
		return
	}
	deviceID := parts[1]

	if err := s.repo.TouchLastSeen(context.Background(), deviceID); err != nil {
		s.logger.Error(fmt.Sprintf("failed to update last_seen for device %s: %v", deviceID, err))
	}
}
