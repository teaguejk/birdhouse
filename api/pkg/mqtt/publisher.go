package mqtt

import "api/pkg/logging"

type Config struct {
	BrokerURL string `json:"broker_url"`
	ClientID  string `json:"client_id"`
	Username  string `json:"username"`
	Password  string `json:"-"`
}

// Publisher is the interface services use to publish MQTT messages.
// This keeps services decoupled from the concrete MQTT client.
type Publisher interface {
	Publish(topic string, qos byte, retained bool, payload any) error
}

// NoopPublisher silently drops messages. Used when MQTT is not configured.
type NoopPublisher struct{}

func (n *NoopPublisher) Publish(string, byte, bool, any) error {
	return nil
}

// NewFromConfig returns a connected Client if MQTT is configured and reachable,
// or a NoopPublisher otherwise. Never fails — degrades gracefully.
func NewFromConfig(cfg *Config, logger *logging.Logger) (Publisher, *Client) {
	if cfg == nil || cfg.BrokerURL == "" {
		logger.Warn("mqtt was not configured, mqtt communication will be unavailable")
		return &NoopPublisher{}, nil
	}

	client, err := NewClient(cfg.BrokerURL, cfg.ClientID, cfg.Username, cfg.Password, logger)
	if err != nil {
		logger.Error("failed to connect to mqtt broker: " + err.Error())
		logger.Warn("mqtt failed to connect, mqtt communication will be unavailable")
		return &NoopPublisher{}, nil
	}

	logger.Info("mqtt connected to " + cfg.BrokerURL)
	return client, client
}
