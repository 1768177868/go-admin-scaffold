package schedule

import (
	"context"
	"log"
)

// Kernel manages the scheduler
type Kernel struct {
	scheduler *Scheduler
}

// NewKernel creates a new scheduler kernel
func NewKernel(scheduler *Scheduler) *Kernel {
	return &Kernel{
		scheduler: scheduler,
	}
}

// Schedule defines scheduled tasks
func (k *Kernel) Schedule() {
	// Add a test task that runs every minute
	k.scheduler.Command("hello:world").EveryMinute().Register()

	log.Println("Scheduled tasks initialized")
}

// Start starts the scheduler
func (k *Kernel) Start(ctx context.Context) error {
	log.Println("Starting scheduler...")
	k.Schedule()
	return k.scheduler.Start(ctx)
}

// Stop stops the scheduler
func (k *Kernel) Stop() {
	log.Println("Stopping scheduler...")
	k.scheduler.Stop()
}
