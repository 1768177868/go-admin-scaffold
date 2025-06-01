package main

import (
	"log"

	"app/cmd/tools/commands"
	"app/internal/bootstrap"
	"app/internal/config"
	"app/pkg/console"
)

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

	// Create command manager
	manager := console.NewManager()

	// Register database commands
	manager.Register(commands.NewMigrateCommand())
	manager.Register(commands.NewSeedCommand())
	manager.Register(commands.NewMakeMigrationCommand())

	// Run command from arguments
	if err := manager.RunFromArgs(); err != nil {
		log.Fatal(err)
	}
}
