package storage

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Provider struct {
	client    *s3.Client
	presigner *s3.PresignClient
	uploader  *manager.Uploader
	bucket    string
	region    string
}

func NewS3Provider(ctx context.Context, bucket, region string) (*S3Provider, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)
	presigner := s3.NewPresignClient(client)
	uploader := manager.NewUploader(client)

	return &S3Provider{
		client:    client,
		presigner: presigner,
		uploader:  uploader,
		bucket:    bucket,
		region:    region,
	}, nil
}

func (p *S3Provider) Upload(ctx context.Context, file io.Reader, path string) (string, error) {
	result, err := p.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(path),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return result.Location, nil
}

func (p *S3Provider) Delete(ctx context.Context, path string) error {
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}
	return nil
}

func (p *S3Provider) GetURL(path string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", p.bucket, p.region, path)
}

func (p *S3Provider) GetFile(ctx context.Context, path string) (io.ReadCloser, error) {
	result, err := p.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file from S3: %w", err)
	}
	return result.Body, nil
}

func (p *S3Provider) Exists(ctx context.Context, path string) (bool, error) {
	_, err := p.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "NoSuchKey") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (p *S3Provider) Copy(ctx context.Context, srcPath, destPath string) error {
	copySource := fmt.Sprintf("%s/%s", p.bucket, srcPath)
	_, err := p.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(p.bucket),
		Key:        aws.String(destPath),
		CopySource: aws.String(copySource),
	})
	if err != nil {
		return fmt.Errorf("failed to copy file in S3: %w", err)
	}
	return nil
}

func (p *S3Provider) GetMetadata(ctx context.Context, path string) (map[string]string, error) {
	result, err := p.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata from S3: %w", err)
	}

	metadata := map[string]string{
		"size":          fmt.Sprintf("%d", *result.ContentLength),
		"content_type":  aws.ToString(result.ContentType),
		"etag":          strings.Trim(aws.ToString(result.ETag), "\""),
		"last_modified": result.LastModified.Format(time.RFC3339),
	}

	for key, value := range result.Metadata {
		metadata[key] = value
	}

	return metadata, nil
}

func (p *S3Provider) GetSignedURL(ctx context.Context, path string, expiration time.Duration, method string) (string, error) {
	switch strings.ToUpper(method) {
	case "PUT":
		presignReq, err := p.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(p.bucket),
			Key:    aws.String(path),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = expiration
		})
		if err != nil {
			return "", fmt.Errorf("failed to create signed PUT URL: %w", err)
		}
		return presignReq.URL, nil
	case "GET":
		presignReq, err := p.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(p.bucket),
			Key:    aws.String(path),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = expiration
		})
		if err != nil {
			return "", fmt.Errorf("failed to create signed GET URL: %w", err)
		}
		return presignReq.URL, nil
	case "DELETE":
		presignReq, err := p.presigner.PresignDeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(p.bucket),
			Key:    aws.String(path),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = expiration
		})
		if err != nil {
			return "", fmt.Errorf("failed to create signed DELETE URL: %w", err)
		}
		return presignReq.URL, nil
	default:
		return "", fmt.Errorf("unsupported method: %s", method)
	}
}
