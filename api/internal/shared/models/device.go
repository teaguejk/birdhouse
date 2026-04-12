package models

import "time"

type Device struct {
	ID         string     `json:"id" db:"id"`
	Name       string     `json:"name" db:"name"`
	APIKeyHash string     `json:"-" db:"api_key_hash"`
	Location   string     `json:"location" db:"location"`
	Active     bool       `json:"active" db:"active"`
	LastSeenAt *time.Time `json:"last_seen_at" db:"last_seen_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

type DeviceStatus struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Location   string     `json:"location"`
	Active     bool       `json:"active"`
	Online     bool       `json:"online"`
	LastSeenAt *time.Time `json:"last_seen_at"`
}

type CreateDeviceRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type UpdateDeviceRequest struct {
	Name     *string `json:"name,omitempty"`
	Location *string `json:"location,omitempty"`
	Active   *bool   `json:"active,omitempty"`
}

type CreateDeviceResponse struct {
	Device *Device `json:"device"`
	APIKey string  `json:"api_key"`
}

type RotateKeyResponse struct {
	APIKey string `json:"api_key"`
}
