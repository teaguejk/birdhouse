package oauth

import (
	"context"
	"fmt"
)

type Claims struct {
	Subject string
	Email   string
	Name    string
}

type TokenVerifier interface {
	Verify(ctx context.Context, token string) (*Claims, error)
}

type Config struct {
	Type     string `json:"type"`
	Audience string `json:"audience"`
}

func NewVerifier(cfg *Config) (TokenVerifier, error) {
	switch cfg.Type {
	case "google":
		return newGoogleVerifier(cfg.Audience), nil
	default:
		return nil, fmt.Errorf("unknown oauth provider type: %s", cfg.Type)
	}
}
