package main

import (
	"flag"
	"log"

	"app/internal/bootstrap"
	"app/internal/config"
	"app/internal/database/migrations"
	"app/internal/database/seeders"
	"app/pkg/database"
)

func main() {
	// Define command line flags
	seedFlag := flag.Bool("seed", false, "Run database seeding")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	if err := bootstrap.SetupDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run migrations
	db := migrations.InitMigrations(database.GetDB())
	if err := db.RunPending(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully")

	// Run seeding if flag is set
	if *seedFlag {
		if err := seeders.Seed(); err != nil {
			log.Fatalf("Failed to run seeding: %v", err)
		}
		log.Println("Seeding completed successfully")
	}
}
