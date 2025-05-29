package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"app/pkg/queue"

	"github.com/spf13/viper"
)

// QueueService 队列服务
type QueueService struct {
	manager *queue.Manager
	workers map[string]*queue.Worker
	mu      sync.RWMutex
}

// NewQueueService 创建队列服务
func NewQueueService() (*QueueService, error) {
	// 加载配置
	driver := viper.GetString("queue.driver")
	queueName := viper.GetString("queue.queue")

	config := queue.Config{
		Driver:  driver,
		Options: make(map[string]interface{}),
	}

	// 根据驱动类型设置选项
	switch driver {
	case "redis":
		// 构建Redis连接字符串
		redisHost := viper.GetString("redis.host")
		redisPort := viper.GetInt("redis.port")
		redisDB := viper.GetInt("redis.db")
		redisPassword := viper.GetString("redis.password")

		connectionStr := fmt.Sprintf("redis://%s:%d/%d", redisHost, redisPort, redisDB)
		if redisPassword != "" {
			connectionStr = fmt.Sprintf("redis://:%s@%s:%d/%d", redisPassword, redisHost, redisPort, redisDB)
		}

		config.Options["connection"] = connectionStr
		config.Options["queue"] = queueName

	case "database", "mysql":
		// 这里需要传入数据库连接实例
		// 暂时返回错误，提示需要在外部传入数据库连接
		return nil, fmt.Errorf("database driver requires external database connection setup")

	default:
		return nil, fmt.Errorf("unsupported queue driver: %s", driver)
	}

	// 创建队列管理器
	manager, err := queue.NewManager(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create queue manager: %v", err)
	}

	return &QueueService{
		manager: manager,
		workers: make(map[string]*queue.Worker),
	}, nil
}

// Start 启动队列服务
func (s *QueueService) Start() error {
	// 获取队列配置
	queues := viper.GetStringMap("queue.queues")
	if len(queues) == 0 {
		return fmt.Errorf("no queues configured")
	}

	// 启动每个队列的工作进程
	for name, config := range queues {
		queueConfig := config.(map[string]interface{})
		processes := int(queueConfig["processes"].(int))

		// 创建队列选项
		options := queue.WorkerOptions{
			Sleep:   time.Duration(viper.GetInt("queue.worker.sleep")) * time.Second,
			MaxJobs: viper.GetInt64("queue.worker.max_jobs"),
			MaxTime: time.Duration(viper.GetInt("queue.worker.max_time")) * time.Second,
			Rest:    time.Duration(viper.GetInt("queue.worker.rest")) * time.Second,
			Memory:  viper.GetInt64("queue.worker.memory"),
			Tries:   viper.GetInt("queue.worker.tries"),
			Timeout: time.Duration(viper.GetInt("queue.worker.timeout")) * time.Second,
		}

		// 启动指定数量的工作进程
		for i := 0; i < processes; i++ {
			worker := queue.NewWorker(s.manager, []string{name}, options)
			workerName := fmt.Sprintf("%s-%d", name, i+1)

			s.mu.Lock()
			s.workers[workerName] = worker
			s.mu.Unlock()

			// 启动工作进程
			go func(w *queue.Worker, name string) {
				log.Printf("Starting queue worker: %s", name)
				w.Start()
				log.Printf("Queue worker stopped: %s", name)

				s.mu.Lock()
				delete(s.workers, name)
				s.mu.Unlock()
			}(worker, workerName)
		}
	}

	return nil
}

// Stop 停止队列服务
func (s *QueueService) Stop() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 停止所有工作进程
	for name, worker := range s.workers {
		log.Printf("Stopping queue worker: %s", name)
		worker.Stop()
	}
}

// Push 推送任务到队列
func (s *QueueService) Push(ctx context.Context, job queue.JobInterface) error {
	return s.manager.Push(ctx, job)
}

// PushRaw 推送原始数据到队列
func (s *QueueService) PushRaw(ctx context.Context, queue string, payload []byte, options map[string]interface{}) error {
	return s.manager.PushRaw(ctx, queue, payload, options)
}

// Later 延迟推送任务
func (s *QueueService) Later(ctx context.Context, job queue.JobInterface, delay time.Duration) error {
	return s.manager.Later(ctx, job, delay)
}

// Pop 从队列中取出任务
func (s *QueueService) Pop(ctx context.Context, queue string) (queue.JobInterface, error) {
	return s.manager.Pop(ctx, queue)
}

// Size 获取队列大小
func (s *QueueService) Size(ctx context.Context, queue string) (int64, error) {
	return s.manager.Size(ctx, queue)
}

// Delete 删除任务
func (s *QueueService) Delete(ctx context.Context, queue string, job queue.JobInterface) error {
	return s.manager.Delete(ctx, queue, job)
}

// Release 释放任务回队列
func (s *QueueService) Release(ctx context.Context, queue string, job queue.JobInterface, delay time.Duration) error {
	return s.manager.Release(ctx, queue, job, delay)
}

// Clear 清空队列
func (s *QueueService) Clear(ctx context.Context, queue string) error {
	return s.manager.Clear(ctx, queue)
}

// GetWorkerCount 获取工作进程数量
func (s *QueueService) GetWorkerCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.workers)
}

// GetActiveQueues 获取活动队列列表
func (s *QueueService) GetActiveQueues() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	queues := make(map[string]struct{})
	for name := range s.workers {
		queues[name] = struct{}{}
	}

	result := make([]string, 0, len(queues))
	for name := range queues {
		result = append(result, name)
	}

	return result
}
