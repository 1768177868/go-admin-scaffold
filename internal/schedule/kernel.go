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
	// Example: Run cleanup every day at midnight (unique)
	k.scheduler.Command("cleanup").Daily().At("00:00").Unique().Register()

	// Example: Send daily report (unique)
	k.scheduler.Command("report:daily").Daily().At("23:00").Unique().Register()

	// Example: Database backup every 6 hours (unique)
	k.scheduler.Command("db:backup").Cron("0 */6 * * *").Unique().Register()

	// Example: Cache cleanup every 30 minutes (can run on multiple servers)
	k.scheduler.Command("cache:clear").EveryThirtyMinutes().Register()

	// Log scheduled tasks
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
