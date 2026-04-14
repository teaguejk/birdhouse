package command

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"api/internal/api/interfaces"
	"api/pkg/logging"
	"api/pkg/mqtt"
)

type AckSubscriber struct {
	repo   interfaces.CommandRepository
	logger *logging.Logger
}

func NewAckSubscriber(logger *logging.Logger, repo interfaces.CommandRepository) *AckSubscriber {
	return &AckSubscriber{
		repo:   repo,
		logger: logger,
	}
}

// Subscribe listens on birdhouse/+/ack for all devices.
func (s *AckSubscriber) Subscribe(client *mqtt.Client) error {
	return client.Subscribe("birdhouse/+/ack", 1, s.handleAck)
}

func (s *AckSubscriber) handleAck(_ pahomqtt.Client, msg pahomqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	if len(parts) != 3 {
		s.logger.Warn(fmt.Sprintf("unexpected ack topic format: %s", msg.Topic()))
		return
	}
	deviceID := parts[1]

	var payload struct {
		CommandID string `json:"command_id"`
	}
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		s.logger.Error(fmt.Sprintf("failed to unmarshal ack from device %s: %v", deviceID, err))
		return
	}

	if payload.CommandID == "" {
		s.logger.Warn(fmt.Sprintf("empty command_id in ack from device %s", deviceID))
		return
	}

	if err := s.repo.Acknowledge(context.Background(), payload.CommandID); err != nil {
		s.logger.Error(fmt.Sprintf("failed to acknowledge command %s: %v", payload.CommandID, err))
		return
	}

	s.logger.Info(fmt.Sprintf("command %s acknowledged via mqtt (device %s)", payload.CommandID, deviceID))
}
