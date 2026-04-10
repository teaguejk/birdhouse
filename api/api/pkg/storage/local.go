package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LocalProvider struct {
	basePath string // upload location
	baseURL  string // access URL prefix
}

func NewLocalProvider(basePath, baseURL string) (*LocalProvider, error) {
	return &LocalProvider{
		basePath: basePath,
		baseURL:  strings.TrimSuffix(baseURL, "/"),
	}, nil
}

func (p *LocalProvider) Upload(ctx context.Context, file io.Reader, path string) (string, error) {
	// ensure subdirectories exist
	fullPath := filepath.Join(p.basePath, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// create file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// copy content
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(fullPath) // cleanup on failure
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// return public URL
	url := fmt.Sprintf("%s/%s", p.baseURL, strings.ReplaceAll(path, "\\", "/"))
	return url, nil
}

func (p *LocalProvider) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(p.basePath, path)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (p *LocalProvider) GetURL(path string) string {
	return fmt.Sprintf("%s/%s", p.baseURL, strings.ReplaceAll(path, "\\", "/"))
}

func (p *LocalProvider) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(p.basePath, path)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (p *LocalProvider) Copy(ctx context.Context, srcPath, destPath string) error {
	src := filepath.Join(p.basePath, srcPath)
	dest := filepath.Join(p.basePath, destPath)

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

func (p *LocalProvider) GetFile(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(p.basePath, path)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return file, nil
}

func (p *LocalProvider) GetMetadata(ctx context.Context, path string) (map[string]string, error) {
	fullPath := filepath.Join(p.basePath, path)

	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	metadata := map[string]string{
		"size":     fmt.Sprintf("%d", fileInfo.Size()),
		"mode":     fileInfo.Mode().String(),
		"mod_time": fileInfo.ModTime().Format("2006-01-02 15:04:05"),
	}

	return metadata, nil
}

func (p *LocalProvider) GetSignedURL(ctx context.Context, path string, expiration time.Duration, method string) (string, error) {
	return "", fmt.Errorf("signed URLs not supported for local provider")
}
