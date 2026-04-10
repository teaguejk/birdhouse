package storage

import (
	"context"
	"io"
	"time"
)

type Provider interface {
	Upload(ctx context.Context, file io.Reader, path string) (string, error)
	Delete(ctx context.Context, path string) error
	GetURL(path string) string
	GetFile(ctx context.Context, path string) (io.ReadCloser, error)
	Exists(ctx context.Context, path string) (bool, error)
	Copy(ctx context.Context, srcPath, destPath string) error // For migrations
	GetMetadata(ctx context.Context, path string) (map[string]string, error)
	GetSignedURL(ctx context.Context, path string, expiration time.Duration, method string) (string, error)
}

type Config struct {
	Type    string            `json:"type"`     // "local", "s3", "gcs"
	BaseURL string            `json:"base_url"` // For serving files
	Options map[string]string `json:"options"`  // Provider-specific config
}
