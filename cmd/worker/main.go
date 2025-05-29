package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"app/internal/config"
	"app/pkg/queue"
)

type Worker struct {
	queue    queue.QueueInterface
	handlers map[string]JobHandler
	stop     chan struct{}
	wg       sync.WaitGroup
}

type JobHandler func(ctx context.Context, payload []byte) error

func NewWorker(q queue.QueueInterface) *Worker {
	return &Worker{
		queue:    q,
		handlers: make(map[string]JobHandler),
		stop:     make(chan struct{}),
	}
}

func (w *Worker) RegisterHandler(queueName string, handler JobHandler) {
	w.handlers[queueName] = handler
}

func (w *Worker) Start(concurrency int) {
	for i := 0; i < concurrency; i++ {
		w.wg.Add(1)
		go w.process()
	}
}

func (w *Worker) Stop() {
	close(w.stop)
	w.wg.Wait()
}

func (w *Worker) process() {
	defer w.wg.Done()

	for {
		select {
		case <-w.stop:
			return
		default:
			for queueName := range w.handlers {
				ctx := context.Background()
				job, err := w.queue.Pop(ctx, queueName)
				if err == queue.ErrQueueEmpty {
					time.Sleep(time.Second)
					continue
				}
				if err != nil {
					log.Printf("Error popping job from queue %s: %v", queueName, err)
					continue
				}

				handler := w.handlers[queueName]
				if err := handler(ctx, job.GetPayload()); err != nil {
					log.Printf("Error processing job %s: %v", job.GetID(), err)
					// Retry with exponential backoff if under max attempts
					if job.GetAttempts() < job.GetMaxAttempts() {
						delay := time.Duration(job.GetAttempts()*job.GetAttempts()) * time.Second
						if err := w.queue.Release(ctx, queueName, job, delay); err != nil {
							log.Printf("Error releasing job %s: %v", job.GetID(), err)
						}
					} else {
						// Delete job if max attempts exceeded
						if err := w.queue.Delete(ctx, queueName, job); err != nil {
							log.Printf("Error deleting job %s: %v", job.GetID(), err)
						}
					}
				} else {
					// Delete successful job
					if err := w.queue.Delete(ctx, queueName, job); err != nil {
						log.Printf("Error deleting job %s: %v", job.GetID(), err)
					}
				}
			}
		}
	}
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create queue instance
	q, err := queue.New(queue.Config{
		Driver:  cfg.Queue.Driver,
		Options: cfg.Queue.Options,
	})
	if err != nil {
		log.Fatalf("Failed to create queue: %v", err)
	}
	defer q.Close()

	// Create worker
	worker := NewWorker(q)

	// Register job handlers
	worker.RegisterHandler("emails", func(ctx context.Context, payload []byte) error {
		// Handle email sending
		return nil
	})

	worker.RegisterHandler("notifications", func(ctx context.Context, payload []byte) error {
		// Handle notification sending
		return nil
	})

	// Start worker with concurrency
	worker.Start(5)

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down worker...")
	worker.Stop()
	log.Println("Worker stopped")
}
