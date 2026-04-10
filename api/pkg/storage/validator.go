package storage

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
)

var (
	ErrInvalidFileType = errors.New("invalid file type")
	ErrFileTooLarge    = errors.New("file too large")
	ErrInvalidFileName = errors.New("invalid file name")
)

type FileValidator struct {
	AllowedMimeTypes  []string
	AllowedExtensions []string
	MaxFileSize       int64
}

func NewImageValidator() *FileValidator {
	return &FileValidator{
		AllowedMimeTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
		},
		AllowedExtensions: []string{
			".jpg", ".jpeg", ".png", ".gif", ".webp",
		},
		MaxFileSize: 10 * 1024 * 1024, // 10MB
	}
}

func (v *FileValidator) ValidateFile(file multipart.File, header *multipart.FileHeader) error {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to reset file position: %w", err)
	}

	// Check file size
	if header.Size > v.MaxFileSize {
		return fmt.Errorf("%w: file size %d exceeds maximum %d", ErrFileTooLarge, header.Size, v.MaxFileSize)
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !slices.Contains(v.AllowedExtensions, ext) {
		return fmt.Errorf("%w: extension %s not allowed", ErrInvalidFileType, ext)
	}

	// Validate MIME type from header
	headerMimeType := header.Header.Get("Content-Type")
	if !slices.Contains(v.AllowedMimeTypes, headerMimeType) {
		return fmt.Errorf("%w: MIME type %s not allowed", ErrInvalidFileType, headerMimeType)
	}

	// Detect actual MIME type from file content
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file for MIME type detection: %w", err)
	}

	detectedMimeType := http.DetectContentType(buffer)
	if !slices.Contains(v.AllowedMimeTypes, detectedMimeType) {
		return fmt.Errorf("%w: detected MIME type %s doesn't match allowed types", ErrInvalidFileType, detectedMimeType)
	}

	// Validate filename (basic security checks)
	if err := v.validateFileName(header.Filename); err != nil {
		return err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to reset file position: %w", err)
	}

	return nil
}

func (v *FileValidator) validateFileName(filename string) error {
	// Check for dangerous characters
	dangerous := []string{"..", "/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range dangerous {
		if strings.Contains(filename, char) {
			return fmt.Errorf("%w: contains dangerous character: %s", ErrInvalidFileName, char)
		}
	}

	// Check filename length
	if len(filename) > 255 {
		return fmt.Errorf("%w: filename too long", ErrInvalidFileName)
	}

	// Check for hidden files
	if strings.HasPrefix(filename, ".") {
		return fmt.Errorf("%w: hidden files not allowed", ErrInvalidFileName)
	}

	return nil
}

func (v *FileValidator) ValidateFileName(filename string) (string, string, error) {
	// Check for dangerous characters
	err := v.validateFileName(filename)
	if err != nil {
		return "", "", fmt.Errorf("invalid filename: %w", err)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(filename))
	if !slices.Contains(v.AllowedExtensions, ext) {
		return "", "", fmt.Errorf("%w: extension %s not allowed", ErrInvalidFileType, ext)
	}
	// Check MIME type
	mimeType := mime.TypeByExtension(ext)
	if !slices.Contains(v.AllowedMimeTypes, mimeType) {
		return "", "", fmt.Errorf("%w: MIME type %s not allowed", ErrInvalidFileType, mimeType)
	}

	return ext, mimeType, nil
}

func getMimeType(filePath string) string {
	ext := filepath.Ext(filePath) // Get file extension
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream" // Fallback
	}
	return mimeType
}
