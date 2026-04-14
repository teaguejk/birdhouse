package models

import (
	"encoding/json"
	"time"
)

type Command struct {
	ID        string          `json:"id" db:"id"`
	DeviceID  string          `json:"device_id" db:"device_id"`
	Action    string          `json:"action" db:"action"`
	Payload   json.RawMessage `json:"payload" db:"payload"`
	Status    string          `json:"status" db:"status"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

type CreateCommandRequest struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type PendingCommandsResponse struct {
	Commands []Command `json:"commands"`
	Active   bool      `json:"active"`
}
