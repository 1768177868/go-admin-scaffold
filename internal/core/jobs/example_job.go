package jobs

import (
	"encoding/json"
	"fmt"
	"time"

	"app/pkg/queue"
)

// ExampleJob 示例任务
type ExampleJob struct {
	queue.BaseJob
	Message string `json:"message"`
}

// NewExampleJob 创建示例任务
func NewExampleJob(message string) *ExampleJob {
	return &ExampleJob{
		BaseJob: queue.BaseJob{
			Queue:       "default",
			Attempts:    0,
			MaxAttempts: 3,
			Delay:       0,
			Timeout:     60 * time.Second,
			RetryAfter:  60 * time.Second,
			Backoff:     []time.Duration{60 * time.Second, 300 * time.Second, 900 * time.Second},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Message: message,
	}
}

// Handle 处理任务
func (j *ExampleJob) Handle() error {
	// 模拟任务处理
	fmt.Printf("Processing example job: %s\n", j.Message)
	time.Sleep(2 * time.Second)
	fmt.Printf("Example job completed: %s\n", j.Message)
	return nil
}

// SendWelcomeEmailJob 发送欢迎邮件任务
type SendWelcomeEmailJob struct {
	queue.BaseJob
	Email   string `json:"email"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
}

// NewSendWelcomeEmailJob 创建发送欢迎邮件任务
func NewSendWelcomeEmailJob(email, name string) *SendWelcomeEmailJob {
	return &SendWelcomeEmailJob{
		BaseJob: queue.BaseJob{
			Queue:       "high",
			Attempts:    0,
			MaxAttempts: 5,
			Delay:       0,
			Timeout:     30 * time.Second,
			RetryAfter:  30 * time.Second,
			Backoff:     []time.Duration{30 * time.Second, 60 * time.Second, 180 * time.Second},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Email:   email,
		Name:    name,
		Subject: "Welcome to Go Admin",
	}
}

// Handle 处理任务
func (j *SendWelcomeEmailJob) Handle() error {
	// 模拟发送邮件
	fmt.Printf("Sending welcome email to %s (%s)\n", j.Name, j.Email)
	time.Sleep(1 * time.Second)
	fmt.Printf("Welcome email sent to %s\n", j.Email)
	return nil
}

// ProcessUploadJob 处理上传文件任务
type ProcessUploadJob struct {
	queue.BaseJob
	FileID   string `json:"file_id"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
}

// NewProcessUploadJob 创建处理上传文件任务
func NewProcessUploadJob(fileID, fileName string, fileSize int64) *ProcessUploadJob {
	return &ProcessUploadJob{
		BaseJob: queue.BaseJob{
			Queue:       "low",
			Attempts:    0,
			MaxAttempts: 2,
			Delay:       0,
			Timeout:     120 * time.Second,
			RetryAfter:  120 * time.Second,
			Backoff:     []time.Duration{120 * time.Second, 300 * time.Second},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		FileID:   fileID,
		FileName: fileName,
		FileSize: fileSize,
	}
}

// Handle 处理任务
func (j *ProcessUploadJob) Handle() error {
	// 模拟处理上传文件
	fmt.Printf("Processing uploaded file: %s (ID: %s, Size: %d bytes)\n", j.FileName, j.FileID, j.FileSize)
	time.Sleep(5 * time.Second)
	fmt.Printf("File processing completed: %s\n", j.FileName)
	return nil
}

// CleanupJob 清理任务
type CleanupJob struct {
	queue.BaseJob
	Target    string          `json:"target"`
	Options   json.RawMessage `json:"options"`
	StartTime time.Time       `json:"start_time"`
	EndTime   time.Time       `json:"end_time"`
}

// NewCleanupJob 创建清理任务
func NewCleanupJob(target string, options json.RawMessage, startTime, endTime time.Time) *CleanupJob {
	return &CleanupJob{
		BaseJob: queue.BaseJob{
			Queue:       "default",
			Attempts:    0,
			MaxAttempts: 3,
			Delay:       0,
			Timeout:     60 * time.Second,
			RetryAfter:  60 * time.Second,
			Backoff:     []time.Duration{60 * time.Second, 300 * time.Second, 900 * time.Second},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Target:    target,
		Options:   options,
		StartTime: startTime,
		EndTime:   endTime,
	}
}

// Handle 处理任务
func (j *CleanupJob) Handle() error {
	// 模拟清理操作
	fmt.Printf("Cleaning up %s from %s to %s\n", j.Target, j.StartTime.Format(time.RFC3339), j.EndTime.Format(time.RFC3339))
	time.Sleep(3 * time.Second)
	fmt.Printf("Cleanup completed for %s\n", j.Target)
	return nil
}
