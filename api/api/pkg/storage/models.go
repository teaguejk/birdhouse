package storage

import "time"

type ImageMetadata struct {
	ID           string    `json:"id"`
	OriginalName string    `json:"original_name"`
	ContentType  string    `json:"content_type"`
	Size         int64     `json:"size"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	URL          string    `json:"url"`
	CreatedAt    time.Time `json:"created_at"`
}

type UploadOptions struct {
	MaxSize       int64    // Maximum file size in bytes
	AllowedTypes  []string // Allowed MIME types
	GenerateThumb bool     // Whether to generate thumbnails
	Quality       int      // JPEG quality (1-100)
	MaxWidth      int      // Maximum width for resizing
	MaxHeight     int      // Maximum height for resizing
}
