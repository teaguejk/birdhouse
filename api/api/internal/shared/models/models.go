package models

import (
	"time"
)

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// PaginatedResponse wraps paginated results
type PaginatedResponse struct {
	Data       interface{}     `json:"data"`
	Pagination *PaginationMeta `json:"pagination"`
}

// APIResponse is a standard API response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ValidationError represents field validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse represents error response format
type ErrorResponse struct {
	Success    bool              `json:"success"`
	Message    string            `json:"message"`
	Errors     []ValidationError `json:"errors,omitempty"`
	StatusCode int               `json:"status_code"`
}

// HealthCheck represents health check response
type HealthCheck struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Database  string    `json:"database"`
}
