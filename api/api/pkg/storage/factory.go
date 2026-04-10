package storage

import (
	"context"
	"fmt"
)

func NewProvider(ctx context.Context, cfg *Config) (Provider, error) {
	switch cfg.Type {
	case "local":
		basePath := cfg.Options["base_path"]
		if basePath == "" {
			return nil, fmt.Errorf("base_path must be set for local storage provider")
		}

		baseURL := cfg.BaseURL
		if baseURL == "" {
			return nil, fmt.Errorf("base_url must be set for local storage provider")
		}

		return NewLocalProvider(basePath, baseURL)

	case "s3":
		// TODO: Implement S3 storage
		bucket := cfg.Options["bucket"]
		if bucket == "" {
			return nil, fmt.Errorf("bucket must be set for S3 storage provider")
		}
		region := cfg.Options["region"]
		if region == "" {
			return nil, fmt.Errorf("region must be set for S3 storage provider")
		}

		return NewS3Provider(ctx, bucket, region)

	case "gcs":
		bucket := cfg.Options["bucket"]
		if bucket == "" {
			return nil, fmt.Errorf("bucket must be set for GCS storage provider")
		}

		return NewGCSProvider(ctx, bucket)

	default:
		return nil, fmt.Errorf("unknown storage provider type: %s", cfg.Type)
	}
}
