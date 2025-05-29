package commands

import (
	"context"

	"app/internal/database/seeders"
	"app/pkg/console"
)

type SeedCommand struct {
	*console.BaseCommand
}

func NewSeedCommand() *SeedCommand {
	cmd := &SeedCommand{
		BaseCommand: console.NewCommand("db:seed", "Seed the database with records"),
	}
	return cmd
}

func (c *SeedCommand) Configure(config *console.CommandConfig) {
	config.Name = "db:seed"
	config.Description = "Seed the database with records"
}

func (c *SeedCommand) Handle(ctx context.Context) error {
	return seeders.Seed()
}
