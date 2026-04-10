package ai

import (
	"encoding/json"
	"fmt"
)

func NewClient(cfg *Config) (Client, error) {
	switch cfg.Type {
	case "anthropic":
		var opts AnthropicOptions
		if err := json.Unmarshal(cfg.Options, &opts); err != nil {
			return nil, fmt.Errorf("failed to parse anthropic options: %w", err)
		}
		opts.APIKey = cfg.APIKey
		return newAnthropicClient(&opts), nil
	default:
		return nil, fmt.Errorf("unknown ai provider type: %s", cfg.Type)
	}
}
