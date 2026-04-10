package database

import (
	"encoding/json"
)

type Config struct {
	Type     string          `json:"type"`
	Password string          `json:"-"`
	Options  json.RawMessage `json:"options"`
}
