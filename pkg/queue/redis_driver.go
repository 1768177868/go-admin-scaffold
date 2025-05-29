package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
)

type redisQueue struct {
	client *redis.Client
}

// RedisOptions represents Redis connection options
type RedisOptions struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func init() {
	Register("redis", NewRedisQueue)
}

// NewRedisQueue creates a new Redis queue instance
func NewRedisQueue(config *Config) (Queue, error) {
	var opts RedisOptions
	if err := mapstructureDecodeConfig(config.Options, &opts); err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Password: opts.Password,
		DB:       opts.DB,
	})

	// Test connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &redisQueue{
		client: client,
	}, nil
}

func (q *redisQueue) Push(ctx context.Context, queueName string, payload interface{}, opts ...JobOption) error {
	// Apply options
	options := &jobOptions{
		MaxRetry: 3, // default max retry
	}
	for _, opt := range opts {
		opt(options)
	}

	// Create job
	job := &Job{
		ID:        uuid.New().String(),
		Queue:     queueName,
		MaxRetry:  options.MaxRetry,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Marshal payload
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	job.Payload = json.RawMessage(data)

	// Marshal job
	jobData, err := json.Marshal(job)
	if err != nil {
		return err
	}

	pipe := q.client.Pipeline()

	// Add job to queue
	if options.Delay > 0 {
		pipe.ZAdd(ctx, getDelayedKey(queueName), redis.Z{
			Score:  float64(time.Now().Add(options.Delay).Unix()),
			Member: job.ID,
		})
	} else {
		pipe.LPush(ctx, getQueueKey(queueName), job.ID)
	}

	// Store job data
	pipe.HSet(ctx, getJobKey(queueName), job.ID, jobData)

	_, err = pipe.Exec(ctx)
	return err
}

func (q *redisQueue) Pop(ctx context.Context, queueName string) (*Job, error) {
	// Move delayed jobs to main queue if they're ready
	q.moveDelayedJobs(ctx, queueName)

	// Get job ID from queue
	jobID, err := q.client.RPop(ctx, getQueueKey(queueName)).Result()
	if err == redis.Nil {
		return nil, ErrQueueEmpty
	}
	if err != nil {
		return nil, err
	}

	// Get job data
	jobData, err := q.client.HGet(ctx, getJobKey(queueName), jobID).Result()
	if err == redis.Nil {
		return nil, ErrJobNotFound
	}
	if err != nil {
		return nil, err
	}

	var job Job
	if err := json.Unmarshal([]byte(jobData), &job); err != nil {
		return nil, err
	}

	// Update job status
	job.Status = "processing"
	job.Attempts++
	job.UpdatedAt = time.Now()

	// Save updated job
	data, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}
	if err := q.client.HSet(ctx, getJobKey(queueName), jobID, string(data)).Err(); err != nil {
		return nil, err
	}

	return &job, nil
}

func (q *redisQueue) Delete(ctx context.Context, queueName, jobID string) error {
	return q.client.HDel(ctx, getJobKey(queueName), jobID).Err()
}

func (q *redisQueue) Release(ctx context.Context, queueName, jobID string, delay time.Duration) error {
	// Get job data
	jobData, err := q.client.HGet(ctx, getJobKey(queueName), jobID).Result()
	if err == redis.Nil {
		return ErrJobNotFound
	}
	if err != nil {
		return err
	}

	var job Job
	if err := json.Unmarshal([]byte(jobData), &job); err != nil {
		return err
	}

	pipe := q.client.Pipeline()

	if delay > 0 {
		pipe.ZAdd(ctx, getDelayedKey(queueName), redis.Z{
			Score:  float64(time.Now().Add(delay).Unix()),
			Member: jobID,
		})
	} else {
		pipe.LPush(ctx, getQueueKey(queueName), jobID)
	}

	_, err = pipe.Exec(ctx)
	return err
}

func (q *redisQueue) Clear(ctx context.Context, queueName string) error {
	pipe := q.client.Pipeline()
	pipe.Del(ctx, getQueueKey(queueName))
	pipe.Del(ctx, getDelayedKey(queueName))
	pipe.Del(ctx, getJobKey(queueName))
	_, err := pipe.Exec(ctx)
	return err
}

func (q *redisQueue) Size(ctx context.Context, queueName string) (int64, error) {
	pipe := q.client.Pipeline()
	pending := pipe.LLen(ctx, getQueueKey(queueName))
	delayed := pipe.ZCard(ctx, getDelayedKey(queueName))
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return pending.Val() + delayed.Val(), nil
}

func (q *redisQueue) Failed(ctx context.Context, queueName string) ([]*Job, error) {
	// Get all jobs
	jobs, err := q.client.HGetAll(ctx, getJobKey(queueName)).Result()
	if err != nil {
		return nil, err
	}

	var failed []*Job
	for _, data := range jobs {
		var job Job
		if err := json.Unmarshal([]byte(data), &job); err != nil {
			continue
		}
		if job.Status == "failed" {
			failed = append(failed, &job)
		}
	}

	return failed, nil
}

func (q *redisQueue) Retry(ctx context.Context, queueName, jobID string) error {
	// Get job data
	jobData, err := q.client.HGet(ctx, getJobKey(queueName), jobID).Result()
	if err == redis.Nil {
		return ErrJobNotFound
	}
	if err != nil {
		return err
	}

	var job Job
	if err := json.Unmarshal([]byte(jobData), &job); err != nil {
		return err
	}

	if job.Attempts >= job.MaxRetry {
		return errors.New("max retry attempts exceeded")
	}

	job.Status = "pending"
	job.UpdatedAt = time.Now()

	// Save updated job
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	pipe := q.client.Pipeline()
	pipe.HSet(ctx, getJobKey(queueName), jobID, string(data))
	pipe.LPush(ctx, getQueueKey(queueName), jobID)
	_, err = pipe.Exec(ctx)
	return err
}

func (q *redisQueue) Close() error {
	return q.client.Close()
}

func (q *redisQueue) moveDelayedJobs(ctx context.Context, queueName string) error {
	now := float64(time.Now().Unix())

	// Get delayed jobs that are ready
	jobs, err := q.client.ZRangeByScore(ctx, getDelayedKey(queueName), &redis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("%f", now),
	}).Result()
	if err != nil {
		return err
	}

	if len(jobs) == 0 {
		return nil
	}

	pipe := q.client.Pipeline()

	// Move jobs to main queue
	for _, jobID := range jobs {
		pipe.LPush(ctx, getQueueKey(queueName), jobID)
		pipe.ZRem(ctx, getDelayedKey(queueName), jobID)
	}

	_, err = pipe.Exec(ctx)
	return err
}

// Helper functions
func getQueueKey(name string) string {
	return fmt.Sprintf("queues:%s", name)
}

func getDelayedKey(name string) string {
	return fmt.Sprintf("queues:%s:delayed", name)
}

func getJobKey(name string) string {
	return fmt.Sprintf("queues:%s:jobs", name)
}

func mapstructureDecodeConfig(input, output interface{}) error {
	return mapstructure.Decode(input, output)
}
