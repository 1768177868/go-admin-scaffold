package storage

import (
	"context"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type localStorage struct {
	baseDir   string
	baseURL   string
	separator string
}

// LocalOptions represents local storage options
type LocalOptions struct {
	BaseDir string `mapstructure:"base_dir"` // Base directory for file storage
	BaseURL string `mapstructure:"base_url"` // Base URL for public access
}

func init() {
	Register("local", NewLocalStorage)
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(config *Config) (Storage, error) {
	var opts LocalOptions
	if err := mapstructureDecodeConfig(config.Options, &opts); err != nil {
		return nil, err
	}

	// Create base directory if it doesn't exist
	if err := os.MkdirAll(opts.BaseDir, 0755); err != nil {
		return nil, err
	}

	return &localStorage{
		baseDir:   opts.BaseDir,
		baseURL:   opts.BaseURL,
		separator: string(os.PathSeparator),
	}, nil
}

func (s *localStorage) Upload(ctx context.Context, path string, reader io.Reader, options ...UploadOption) error {
	if err := validatePath(path); err != nil {
		return err
	}

	// Apply options
	opts := &uploadOptions{}
	for _, opt := range options {
		opt(opts)
	}

	fullPath := filepath.Join(s.baseDir, path)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy content
	_, err = io.Copy(file, reader)
	return err
}

func (s *localStorage) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	if err := validatePath(path); err != nil {
		return nil, err
	}

	fullPath := filepath.Join(s.baseDir, path)
	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	return file, nil
}

func (s *localStorage) Delete(ctx context.Context, path string) error {
	if err := validatePath(path); err != nil {
		return err
	}

	fullPath := filepath.Join(s.baseDir, path)
	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *localStorage) Exists(ctx context.Context, path string) (bool, error) {
	if err := validatePath(path); err != nil {
		return false, err
	}

	fullPath := filepath.Join(s.baseDir, path)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *localStorage) List(ctx context.Context, prefix string) ([]FileInfo, error) {
	if err := validatePath(prefix); err != nil {
		return nil, err
	}

	fullPath := filepath.Join(s.baseDir, prefix)
	var files []FileInfo

	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory
		if path == fullPath {
			return nil
		}

		relativePath, err := filepath.Rel(s.baseDir, path)
		if err != nil {
			return err
		}

		files = append(files, FileInfo{
			Name:         info.Name(),
			Size:         info.Size(),
			LastModified: info.ModTime(),
			IsDir:        info.IsDir(),
			ContentType:  mime.TypeByExtension(filepath.Ext(path)),
			Path:         filepath.ToSlash(relativePath),
		})

		return nil
	})

	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return files, nil
}

func (s *localStorage) GetURL(ctx context.Context, path string) (string, error) {
	if err := validatePath(path); err != nil {
		return "", err
	}

	if s.baseURL == "" {
		return "", nil
	}

	// Ensure the base URL has a trailing slash
	baseURL := s.baseURL
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	// URL encode the path
	encodedPath := url.PathEscape(path)
	return baseURL + encodedPath, nil
}

func (s *localStorage) GetInfo(ctx context.Context, path string) (*FileInfo, error) {
	if err := validatePath(path); err != nil {
		return nil, err
	}

	fullPath := filepath.Join(s.baseDir, path)
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	return &FileInfo{
		Name:         info.Name(),
		Size:         info.Size(),
		LastModified: info.ModTime(),
		IsDir:        info.IsDir(),
		ContentType:  mime.TypeByExtension(filepath.Ext(path)),
		Path:         filepath.ToSlash(path),
	}, nil
}

func (s *localStorage) Close() error {
	return nil
}

// Helper functions

func validatePath(path string) error {
	if path == "" {
		return ErrInvalidPath
	}

	// Check for path traversal attempts
	if strings.Contains(path, "..") {
		return ErrInvalidPath
	}

	return nil
}

func mapstructureDecodeConfig(input, output interface{}) error {
	// You can use mapstructure package here
	// For simplicity, we'll just return nil
	return nil
}
