package device

import (
	"context"
	"fmt"
	"strings"

	"api/internal/api/interfaces"
	"api/pkg/logging"
	"api/pkg/mqtt"
	pahomqtt "github.com/eclipse/paho.mqtt.golang"
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

	if err := s.repo.UpdateStatus(context.Background(), deviceID, msg.Payload()); err != nil {
		s.logger.Error(fmt.Sprintf("failed to update status for device %s: %v", deviceID, err))
	}
}
