package api

import (
	"api/pkg/ai"
	"api/pkg/config"
	"api/pkg/database"
	"api/pkg/storage"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Server   *ServerConfig    `json:"server"`
	Storage  *storage.Config  `json:"storage"`
	Database *database.Config `json:"database"`
	AI       *ai.Config       `json:"ai"`
}

type ServerConfig struct {
	Port             string          `json:"port"`
	ReadTimeout      config.Duration `json:"read_timeout"`
	WriteTimeout     config.Duration `json:"write_timeout"`
	IdleTimeout      config.Duration `json:"idle_timeout"`
	RateLimitEnabled bool            `json:"rate_limit_enabled"`
	CORS             *CORSConfig     `json:"cors"`
	PublicRoutes     []string        `json:"public_routes"`
}

type CORSConfig struct {
	AllowedOrigins []string `json:"allowed_origins"`
	AllowedMethods []string `json:"allowed_methods"`
	AllowedHeaders []string `json:"allowed_headers"`
	MaxAge         int      `json:"max_age"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{
		Server: defaultServerConfig(),
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	applyEnvOverrides(cfg)

	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func validate(cfg *Config) error {
	if cfg.Database == nil {
		return fmt.Errorf("database config is required")
	}
	if cfg.Database.Type == "" {
		return fmt.Errorf("database.type is required")
	}

	if cfg.Storage == nil {
		return fmt.Errorf("storage config is required")
	}
	if cfg.Storage.Type == "" {
		return fmt.Errorf("storage.type is required")
	}

	if cfg.AI == nil {
		return fmt.Errorf("ai config is required")
	}
	if cfg.AI.Type == "" {
		return fmt.Errorf("ai.type is required")
	}

	return nil
}

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("PORT"); v != "" {
		cfg.Server.Port = v
	}
	if cfg.Database != nil {
		if v := os.Getenv("DB_PASSWORD"); v != "" {
			cfg.Database.Password = v
		}
	}
	if cfg.AI != nil {
		if v := os.Getenv("ANTHROPIC_API_KEY"); v != "" {
			cfg.AI.APIKey = v
		}
	}
}

func defaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:             "8080",
		ReadTimeout:      config.Duration{Duration: 5 * time.Second},
		WriteTimeout:     config.Duration{Duration: 10 * time.Second},
		IdleTimeout:      config.Duration{Duration: 120 * time.Second},
		RateLimitEnabled: false,
		CORS: &CORSConfig{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
			MaxAge:         86400,
		},
		PublicRoutes: []string{
			"GET /health/healthcheck",
		},
	}
}
