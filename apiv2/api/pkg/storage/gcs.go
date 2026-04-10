package storage

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCSProvider struct {
	client     *storage.Client
	bucket     *storage.BucketHandle
	bucketName string
}

func NewGCSProvider(ctx context.Context, bucketName string, opts ...option.ClientOption) (*GCSProvider, error) {
	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	bucket := client.Bucket(bucketName)

	return &GCSProvider{
		client:     client,
		bucket:     bucket,
		bucketName: bucketName,
	}, nil
}

func (p *GCSProvider) Upload(ctx context.Context, file io.Reader, path string) (string, error) {
	obj := p.bucket.Object(path)
	writer := obj.NewWriter(ctx)
	defer writer.Close()

	if _, err := io.Copy(writer, file); err != nil {
		return "", fmt.Errorf("failed to upload file to GCS: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close GCS writer: %w", err)
	}

	return p.GetURL(path), nil
}

func (p *GCSProvider) Delete(ctx context.Context, path string) error {
	obj := p.bucket.Object(path)
	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete file from GCS: %w", err)
	}
	return nil
}

func (p *GCSProvider) GetURL(path string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", p.bucketName, path)

}

func (p *GCSProvider) GetFile(ctx context.Context, path string) (io.ReadCloser, error) {
	obj := p.bucket.Object(path)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get file from GCS: %w", err)
	}
	return reader, nil
}

func (p *GCSProvider) Exists(ctx context.Context, path string) (bool, error) {
	obj := p.bucket.Object(path)
	_, err := obj.Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (p *GCSProvider) Copy(ctx context.Context, srcPath, destPath string) error {
	srcObj := p.bucket.Object(srcPath)
	destObj := p.bucket.Object(destPath)

	_, err := destObj.CopierFrom(srcObj).Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to copy file in GCS: %w", err)
	}
	return nil
}

func (p *GCSProvider) GetMetadata(ctx context.Context, path string) (map[string]string, error) {
	obj := p.bucket.Object(path)
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata from GCS: %w", err)
	}

	metadata := map[string]string{
		"size":          fmt.Sprintf("%d", attrs.Size),
		"content_type":  attrs.ContentType,
		"etag":          attrs.Etag,
		"last_modified": attrs.Updated.Format(time.RFC3339),
	}

	for key, value := range attrs.Metadata {
		metadata[key] = value
	}

	return metadata, nil
}

func (p *GCSProvider) GetSignedURL(ctx context.Context, filename string, expiration time.Duration, method string) (string, error) {
	var httpMethod string
	switch strings.ToUpper(method) {
	case "PUT":
		httpMethod = "PUT"
	case "GET":
		httpMethod = "GET"
	case "DELETE":
		httpMethod = "DELETE"
	default:
		return "", fmt.Errorf("unsupported method: %s", method)
	}

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  httpMethod,
		Expires: time.Now().Add(expiration),
	}

	signedURL, err := p.bucket.SignedURL(filename, opts)
	if err != nil {
		return "", fmt.Errorf("failed to create signed URL for method: %s: %w", method, err)
	}

	return signedURL, nil
}
