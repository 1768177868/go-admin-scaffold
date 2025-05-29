package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mysqlQueue struct {
	db *gorm.DB
}

// MySQLJob represents a job in the MySQL database
type MySQLJob struct {
	ID        string `gorm:"primaryKey;type:varchar(36)"`
	Queue     string `gorm:"type:varchar(100);index"`
	Payload   string `gorm:"type:text"`
	Attempts  int    `gorm:"default:0"`
	MaxRetry  int    `gorm:"default:3"`
	Status    string `gorm:"type:varchar(20);index"`
	Error     string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName returns the table name for the job
func (MySQLJob) TableName() string {
	return "queue_jobs"
}

func init() {
	Register("mysql", NewMySQLQueue)
}

// NewMySQLQueue creates a new MySQL queue instance
func NewMySQLQueue(config *Config) (Queue, error) {
	db, ok := config.Options["db"].(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("mysql queue requires a *gorm.DB instance")
	}

	// Auto migrate the jobs table
	if err := db.AutoMigrate(&MySQLJob{}); err != nil {
		return nil, err
	}

	return &mysqlQueue{
		db: db,
	}, nil
}

func (q *mysqlQueue) Push(ctx context.Context, queueName string, payload interface{}, opts ...JobOption) error {
	// Apply options
	options := &jobOptions{
		MaxRetry: 3, // default max retry
	}
	for _, opt := range opts {
		opt(options)
	}

	// Marshal payload
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	job := &MySQLJob{
		ID:        uuid.New().String(),
		Queue:     queueName,
		Payload:   string(data),
		MaxRetry:  options.MaxRetry,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if options.Delay > 0 {
		job.CreatedAt = time.Now().Add(options.Delay)
	}

	return q.db.WithContext(ctx).Create(job).Error
}

func (q *mysqlQueue) Pop(ctx context.Context, queueName string) (*Job, error) {
	var mysqlJob MySQLJob

	err := q.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get the next available job
		if err := tx.Where("queue = ? AND status = ? AND created_at <= ?", queueName, "pending", time.Now()).
			Order("created_at ASC").
			First(&mysqlJob).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return ErrQueueEmpty
			}
			return err
		}

		// Update job status
		mysqlJob.Status = "processing"
		mysqlJob.Attempts++
		mysqlJob.UpdatedAt = time.Now()

		return tx.Save(&mysqlJob).Error
	})

	if err != nil {
		return nil, err
	}

	return &Job{
		ID:        mysqlJob.ID,
		Queue:     mysqlJob.Queue,
		Payload:   json.RawMessage(mysqlJob.Payload),
		Attempts:  mysqlJob.Attempts,
		MaxRetry:  mysqlJob.MaxRetry,
		Status:    mysqlJob.Status,
		Error:     mysqlJob.Error,
		CreatedAt: mysqlJob.CreatedAt,
		UpdatedAt: mysqlJob.UpdatedAt,
	}, nil
}

func (q *mysqlQueue) Delete(ctx context.Context, queueName, jobID string) error {
	return q.db.WithContext(ctx).Where("queue = ? AND id = ?", queueName, jobID).Delete(&MySQLJob{}).Error
}

func (q *mysqlQueue) Release(ctx context.Context, queueName, jobID string, delay time.Duration) error {
	updates := map[string]interface{}{
		"status":     "pending",
		"created_at": time.Now().Add(delay),
		"updated_at": time.Now(),
	}
	return q.db.WithContext(ctx).Model(&MySQLJob{}).Where("queue = ? AND id = ?", queueName, jobID).Updates(updates).Error
}

func (q *mysqlQueue) Clear(ctx context.Context, queueName string) error {
	return q.db.WithContext(ctx).Where("queue = ?", queueName).Delete(&MySQLJob{}).Error
}

func (q *mysqlQueue) Size(ctx context.Context, queueName string) (int64, error) {
	var count int64
	err := q.db.WithContext(ctx).Model(&MySQLJob{}).Where("queue = ? AND status = ?", queueName, "pending").Count(&count).Error
	return count, err
}

func (q *mysqlQueue) Failed(ctx context.Context, queueName string) ([]*Job, error) {
	var mysqlJobs []MySQLJob
	err := q.db.WithContext(ctx).Where("queue = ? AND status = ?", queueName, "failed").Find(&mysqlJobs).Error
	if err != nil {
		return nil, err
	}

	jobs := make([]*Job, len(mysqlJobs))
	for i, mysqlJob := range mysqlJobs {
		jobs[i] = &Job{
			ID:        mysqlJob.ID,
			Queue:     mysqlJob.Queue,
			Payload:   json.RawMessage(mysqlJob.Payload),
			Attempts:  mysqlJob.Attempts,
			MaxRetry:  mysqlJob.MaxRetry,
			Status:    mysqlJob.Status,
			Error:     mysqlJob.Error,
			CreatedAt: mysqlJob.CreatedAt,
			UpdatedAt: mysqlJob.UpdatedAt,
		}
	}

	return jobs, nil
}

func (q *mysqlQueue) Retry(ctx context.Context, queueName, jobID string) error {
	updates := map[string]interface{}{
		"status":     "pending",
		"updated_at": time.Now(),
	}
	return q.db.WithContext(ctx).Model(&MySQLJob{}).Where("queue = ? AND id = ?", queueName, jobID).Updates(updates).Error
}

func (q *mysqlQueue) Close() error {
	sqlDB, err := q.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
