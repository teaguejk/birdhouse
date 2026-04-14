package models

import (
	"encoding/json"
	"time"
)

type DeviceConfig struct {
	MinContourArea  int     `json:"min_contour_area"`
	Threshold       int     `json:"threshold"`
	CooldownSeconds float64 `json:"cooldown_seconds"`
}

var DefaultDeviceConfig = DeviceConfig{
	MinContourArea:  500,
	Threshold:       25,
	CooldownSeconds: 2.0,
}

type Device struct {
	ID         string          `json:"id" db:"id"`
	Name       string          `json:"name" db:"name"`
	APIKeyHash string          `json:"-" db:"api_key_hash"`
	Location   string          `json:"location" db:"location"`
	Active     bool            `json:"active" db:"active"`
	Config     json.RawMessage `json:"config" db:"config"`
	LastSeenAt *time.Time      `json:"last_seen_at" db:"last_seen_at"`
	LastStatus json.RawMessage `json:"last_status" db:"last_status"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at" db:"updated_at"`
}

type DeviceStatus struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Location   string          `json:"location"`
	Active     bool            `json:"active"`
	Online     bool            `json:"online"`
	Config     json.RawMessage `json:"config"`
	LastSeenAt *time.Time      `json:"last_seen_at"`
	LastStatus json.RawMessage `json:"last_status"`
}

type CreateDeviceRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type UpdateDeviceRequest struct {
	Name     *string       `json:"name,omitempty"`
	Location *string       `json:"location,omitempty"`
	Active   *bool         `json:"active,omitempty"`
	Config   *DeviceConfig `json:"config,omitempty"`
}

type CreateDeviceResponse struct {
	Device *Device `json:"device"`
	APIKey string  `json:"api_key"`
}

type RotateKeyResponse struct {
	APIKey string `json:"api_key"`
}
