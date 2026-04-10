package database

import (
	"context"
	"encoding/json"
	"fmt"
)

func NewDatabase(ctx context.Context, cfg *Config) (*PostgresDB, error) {
	switch cfg.Type {
	case "postgres":
		var opts PostgresOptions
		if err := json.Unmarshal(cfg.Options, &opts); err != nil {
			return nil, fmt.Errorf("failed to parse postgres options: %w", err)
		}
		opts.Password = cfg.Password
		return NewPostgresDB(ctx, &opts)
	default:
		return nil, fmt.Errorf("unknown database type: %s", cfg.Type)
	}
}
