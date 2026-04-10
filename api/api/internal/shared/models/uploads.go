package models

import "time"

type File struct {
	ID           string     `json:"id" db:"id"`
	DeviceID       string     `json:"device_id" db:"device_id"`
	ResourceType *string    `json:"resource_type" db:"resource_type"` // nullable
	ResourceID   *string    `json:"resource_id" db:"resource_id"`     // nullable
	Status       string     `json:"status" db:"status"`               // "pending", "assigned", "orphaned"
	Filename     string     `json:"filename" db:"filename"`
	OriginalName string     `json:"original_name" db:"original_name"`
	MimeType     string     `json:"mime_type" db:"mime_type"`
	Size         int64      `json:"size" db:"size"`
	URL          string     `json:"url" db:"url"`
	SortOrder    int        `json:"sort_order" db:"sort_order"`
	ExpiresAt    *time.Time `json:"expires_at" db:"expires_at"` // for cleanup
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type FileUploadRequest struct {
	Filenames []string `json:"filenames"`
}

type UploadDetails struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

type UploadCompleteRequest struct {
	Key string `json:"key"`
}

type UploadDetailsResponse map[string]UploadDetails
