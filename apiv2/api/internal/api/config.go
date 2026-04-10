package api

import (
	"api/pkg/ai"
	"api/pkg/database"
	"api/pkg/storage"
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
	Port             string        `env:"PORT"`
	ReadTimeout      time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout     time.Duration `env:"WRITE_TIMEOUT"`
	IdleTimeout      time.Duration `env:"IDLE_TIMEOUT"`
	RateLimitEnabled bool          `env:"RATE_LIMIT_ENABLED"`
	AllowedOrigins   string        `env:"ALLOWED_ORIGINS"`
	PublicRoutes     []string      `env:"PUBLIC_ROUTES"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Server:   NewServerDefaultConfig(),
		Storage:  NewStorageDefaultConfig(),
		Database: NewDatabaseDefaultConfig(),
		AI:       NewAIDefaultConfig(),
	}
}

func NewAIDefaultConfig() *ai.Config {
	return &ai.Config{
		APIKey:  os.Getenv("ANTHROPIC_API_KEY"),
		APIURL:  "https://api.anthropic.com/v1/messages",
		Model:   "claude-sonnet-4-20250514",
		Timeout: 30 * time.Second,
	}
}

func NewServerDefaultConfig() *ServerConfig {
	return &ServerConfig{
		Port:             "8080",
		ReadTimeout:      5 * time.Second,
		WriteTimeout:     10 * time.Second,
		IdleTimeout:      120 * time.Second,
		RateLimitEnabled: false,
		AllowedOrigins:   "*",
		PublicRoutes: []string{
			"GET /health",
			"GET /uploads/{filename}",
		},
	}
}

func NewStorageDefaultConfig() *storage.Config {
	// cfg := &storage.Config{
	// 	Type:    "local",
	// 	BaseURL: "http://localhost:8080/uploads",
	// 	Options: map[string]string{
	// 		"base_path": "/Users/jteague/dev/apps/birdhouse/shared/uploads",
	// 	},
	// }

	// cfg := &storage.Config{
	// 	Type: "s3",
	// 	Options: map[string]string{
	// 		"bucket": "birdhouse-uploads",
	// 		"region": "us-east-1",
	// 	},
	// }

	cfg := &storage.Config{
		Type: "gcs",
		Options: map[string]string{
			"bucket": "birdhouse-uploads",
		},
	}

	return cfg
}

func NewDatabaseDefaultConfig() *database.Config {
	cfg := &database.Config{
		Name:               "postgres",
		User:               "postgres",
		Host:               "127.0.0.1",
		Port:               "5432",
		SSLMode:            "disable",
		ConnectionTimeout:  30,
		Password:           "password",
		SSLCertPath:        "",
		SSLKeyPath:         "",
		SSLRootCertPath:    "",
		PoolMinConnections: "1",
		PoolMaxConnections: "10",
		PoolMaxConnLife:    5 * time.Minute,
		PoolMaxConnIdle:    1 * time.Minute,
		PoolHealthCheck:    1 * time.Minute,
	}

	return cfg
}
