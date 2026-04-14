package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"api/pkg/logging"
)

type Client struct {
	client pahomqtt.Client
	logger *logging.Logger
}

func NewClient(brokerURL, clientID, username, password string, logger *logging.Logger) (*Client, error) {
	opts := pahomqtt.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID(clientID).
		SetAutoReconnect(true).
		SetConnectTimeout(10 * time.Second).
		SetOnConnectHandler(func(_ pahomqtt.Client) {
			logger.Info("mqtt connected")
		}).
		SetConnectionLostHandler(func(_ pahomqtt.Client, err error) {
			logger.Warn(fmt.Sprintf("mqtt connection lost: %v", err))
		}).
		SetReconnectingHandler(func(_ pahomqtt.Client, _ *pahomqtt.ClientOptions) {
			logger.Info("mqtt reconnecting...")
		})

	if username != "" {
		opts.SetUsername(username)
		opts.SetPassword(password)
	}

	pahoClient := pahomqtt.NewClient(opts)
	token := pahoClient.Connect()
	if !token.WaitTimeout(10 * time.Second) {
		return nil, fmt.Errorf("mqtt connect timed out")
	}
	if token.Error() != nil {
		return nil, fmt.Errorf("mqtt connect failed: %w", token.Error())
	}

	return &Client{
		client: pahoClient,
		logger: logger,
	}, nil
}

func (c *Client) Publish(topic string, qos byte, retained bool, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("mqtt marshal failed: %w", err)
	}

	token := c.client.Publish(topic, qos, retained, data)
	token.Wait()
	return token.Error()
}

func (c *Client) Subscribe(topic string, qos byte, handler pahomqtt.MessageHandler) error {
	token := c.client.Subscribe(topic, qos, handler)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("mqtt subscribe failed: %w", token.Error())
	}
	c.logger.Info(fmt.Sprintf("mqtt subscribed to %s", topic))
	return nil
}

func (c *Client) Disconnect() {
	c.client.Disconnect(250)
	c.logger.Info("mqtt disconnected")
}

func (c *Client) IsConnected() bool {
	return c.client.IsConnected()
}
