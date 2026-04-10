package ai

import "encoding/json"

type Config struct {
	Type    string          `json:"type"`
	APIKey  string          `json:"-"`
	Options json.RawMessage `json:"options"`
}
