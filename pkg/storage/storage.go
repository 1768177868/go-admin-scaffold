package storage

import (
	"context"
	"errors"
	"io"
	"time"
)

var (
	ErrDriverNotFound = errors.New("storage driver not found")
	ErrFileNotFound   = errors.New("file not found")
	ErrInvalidPath    = errors.New("invalid path")
)

// FileInfo represents information about a file
type FileInfo struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	IsDir        bool      `json:"is_dir"`
	ContentType  string    `json:"content_type"`
	Path         string    `json:"path"`
}

// Storage defines the interface that storage drivers must implement
type Storage interface {
	// Upload uploads a file to storage
	Upload(ctx context.Context, path string, reader io.Reader, options ...UploadOption) error

	// Download downloads a file from storage
	Download(ctx context.Context, path string) (io.ReadCloser, error)

	// Delete deletes a file from storage
	Delete(ctx context.Context, path string) error

	// Exists checks if a file exists in storage
	Exists(ctx context.Context, path string) (bool, error)

	// List lists files in a directory
	List(ctx context.Context, prefix string) ([]FileInfo, error)

	// GetURL gets the public URL for a file
	GetURL(ctx context.Context, path string) (string, error)

	// GetInfo gets file information
	GetInfo(ctx context.Context, path string) (*FileInfo, error)

	// Close closes the storage connection
	Close() error
}

// Config represents storage configuration
type Config struct {
	Driver  string                 `mapstructure:"driver"` // local or s3
	Options map[string]interface{} `mapstructure:"options"`
}

// UploadOption represents upload options
type UploadOption func(*uploadOptions)

type uploadOptions struct {
	ContentType string
	Public      bool
	MaxAge      time.Duration
}

// WithContentType sets the content type for upload
func WithContentType(contentType string) UploadOption {
	return func(o *uploadOptions) {
		o.ContentType = contentType
	}
}

// WithPublic sets whether the file should be publicly accessible
func WithPublic(public bool) UploadOption {
	return func(o *uploadOptions) {
		o.Public = public
	}
}

// WithMaxAge sets the cache control max-age
func WithMaxAge(maxAge time.Duration) UploadOption {
	return func(o *uploadOptions) {
		o.MaxAge = maxAge
	}
}

var (
	drivers = make(map[string]func(config *Config) (Storage, error))
)

// Register registers a storage driver
func Register(name string, driver func(config *Config) (Storage, error)) {
	drivers[name] = driver
}

// New creates a new storage instance
func New(config *Config) (Storage, error) {
	driver, ok := drivers[config.Driver]
	if !ok {
		return nil, ErrDriverNotFound
	}
	return driver(config)
}
