package database

import (
	"context"
	"encoding/json"
	"fmt"

	"api/pkg/logging"
)

func NewDatabase(ctx context.Context, cfg *Config, logger *logging.Logger) (*PostgresDB, error) {
	switch cfg.Type {
	case "postgres":
		var opts PostgresOptions
		if err := json.Unmarshal(cfg.Options, &opts); err != nil {
			return nil, fmt.Errorf("failed to parse postgres options: %w", err)
		}
		opts.Password = cfg.Password
		return NewPostgresDB(ctx, &opts, logger)
	default:
		return nil, fmt.Errorf("unknown database type: %s", cfg.Type)
	}
}
