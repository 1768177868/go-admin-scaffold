package commands

import (
	"context"

	"app/internal/database/migrations"
	"app/pkg/console"
	"app/pkg/database"
)

type MigrateCommand struct {
	*console.BaseCommand
}

func NewMigrateCommand() *MigrateCommand {
	cmd := &MigrateCommand{
		BaseCommand: console.NewCommand("migrate", "Run database migrations"),
	}
	return cmd
}

func (c *MigrateCommand) Configure(config *console.CommandConfig) {
	config.Name = "migrate"
	config.Description = "Run database migrations"
}

func (c *MigrateCommand) Handle(ctx context.Context) error {
	db := database.GetDB()
	migrator := migrations.InitMigrations(db)
	return migrator.RunPending()
}
