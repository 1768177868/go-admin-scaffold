package queue

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrQueueEmpty    = errors.New("queue is empty")
	ErrQueueFull     = errors.New("queue is full")
	ErrQueueNotFound = errors.New("queue not found")
	ErrJobNotFound   = errors.New("job not found")
)

// Job represents a queue job
type Job struct {
	ID        string          `json:"id"`
	Queue     string          `json:"queue"`
	Payload   json.RawMessage `json:"payload"`
	Attempts  int             `json:"attempts"`
	MaxRetry  int             `json:"max_retry"`
	Status    string          `json:"status"` // pending, processing, completed, failed
	Error     string          `json:"error,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// Queue defines the interface that queue drivers must implement
type Queue interface {
	// Push adds a job to the queue
	Push(ctx context.Context, queueName string, payload interface{}, opts ...JobOption) error

	// Pop retrieves and reserves a job from the queue
	Pop(ctx context.Context, queueName string) (*Job, error)

	// Delete removes a job from the queue
	Delete(ctx context.Context, queueName, jobID string) error

	// Release puts a job back to the queue
	Release(ctx context.Context, queueName, jobID string, delay time.Duration) error

	// Clear removes all jobs from the queue
	Clear(ctx context.Context, queueName string) error

	// Size returns the number of jobs in the queue
	Size(ctx context.Context, queueName string) (int64, error)

	// Failed returns the list of failed jobs
	Failed(ctx context.Context, queueName string) ([]*Job, error)

	// Retry retries a failed job
	Retry(ctx context.Context, queueName, jobID string) error

	// Close closes the queue connection
	Close() error
}

// Config represents queue configuration
type Config struct {
	Driver  string                 `mapstructure:"driver"`  // redis or mysql
	Options map[string]interface{} `mapstructure:"options"` // driver-specific options
}

// JobOption represents job options
type JobOption func(*jobOptions)

type jobOptions struct {
	Delay    time.Duration
	MaxRetry int
}

// WithDelay sets the delay for the job
func WithDelay(delay time.Duration) JobOption {
	return func(o *jobOptions) {
		o.Delay = delay
	}
}

// WithMaxRetry sets the maximum retry attempts for the job
func WithMaxRetry(maxRetry int) JobOption {
	return func(o *jobOptions) {
		o.MaxRetry = maxRetry
	}
}

var (
	drivers = make(map[string]func(config *Config) (Queue, error))
)

// Register registers a queue driver
func Register(name string, driver func(config *Config) (Queue, error)) {
	drivers[name] = driver
}

// New creates a new queue instance
func New(config *Config) (Queue, error) {
	driver, ok := drivers[config.Driver]
	if !ok {
		return nil, errors.New("unknown queue driver: " + config.Driver)
	}
	return driver(config)
}
