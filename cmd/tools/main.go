package main

import (
	"log"

	"app/internal/commands"
	"app/internal/config"
	"app/pkg/console"
	"app/pkg/database"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize database
	if err := database.Init(cfg); err != nil {
		log.Fatal(err)
	}

	// Create command manager
	manager := console.NewManager()

	// Register database commands
	manager.Register(commands.NewMigrateCommand())
	manager.Register(commands.NewSeedCommand())

	// Run command from arguments
	if err := manager.RunFromArgs(); err != nil {
		log.Fatal(err)
	}
}
