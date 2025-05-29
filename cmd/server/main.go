package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"app/internal/commands"
	"app/internal/config"
	"app/internal/schedule"
	"app/pkg/console"
	"app/pkg/database"
	"app/pkg/locker"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize database
	err = database.Setup(&database.DBConfig{
		Driver:   "mysql",
		Host:     cfg.MySQL.Host,
		Port:     cfg.MySQL.Port,
		Username: cfg.MySQL.Username,
		Password: cfg.MySQL.Password,
		Database: cfg.MySQL.Database,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Create Redis locker
	redisLocker := locker.NewRedisLocker(redisClient)

	// Create command manager for scheduled tasks
	manager := console.NewManager()

	// Register commands that can be scheduled
	manager.Register(commands.NewHelloWorldCommand())
	// Register other commands that need to be scheduled...

	// Create and start scheduler with Redis locker
	scheduler := schedule.NewScheduler(manager, redisLocker)
	kernel := schedule.NewKernel(scheduler)

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start scheduler in a goroutine
	go func() {
		if err := kernel.Start(ctx); err != nil {
			log.Printf("Scheduler error: %v", err)
		}
	}()

	// Start HTTP server in a goroutine
	go func() {
		// Your HTTP server setup and start code here
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown gracefully
	log.Println("Shutting down scheduler...")
	kernel.Stop()

	// Close Redis client
	if err := redisClient.Close(); err != nil {
		log.Printf("Error closing Redis client: %v", err)
	}

	log.Println("Server exiting")
}
