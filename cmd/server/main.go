package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/cmd/server/setup"
	_ "app/docs" // 导入 swagger 文档
	"app/internal/bootstrap"
	"app/internal/commands"
	"app/internal/config"
	"app/internal/schedule"
	"app/pkg/console"
	"app/pkg/locker"

	"github.com/redis/go-redis/v9"
)

// @title Go Admin Scaffold API
// @version 1.0
// @description A modern Go admin scaffold API server.
// @host localhost:8080
// @BasePath /api
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize database
	if err := bootstrap.SetupDatabase(cfg); err != nil {
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
	manager.Register(commands.NewMigrateCommand())
	manager.Register(commands.NewSeedCommand())
	manager.Register(commands.NewSendEmailsCommand())
	manager.Register(commands.NewMakeCommand())

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

	// Initialize the HTTP server
	app, err := setup.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Start HTTP server in a goroutine
	srv := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: app.Engine(),
	}

	go func() {
		log.Printf("Server is running on %s", cfg.Server.Address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown gracefully
	log.Println("Shutting down server...")

	// Create a deadline for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Stop scheduler
	log.Println("Shutting down scheduler...")
	kernel.Stop()

	// Close Redis client
	if err := redisClient.Close(); err != nil {
		log.Printf("Error closing Redis client: %v", err)
	}

	log.Println("Server exited")
}
