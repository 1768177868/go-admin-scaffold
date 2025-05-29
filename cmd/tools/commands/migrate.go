package commands

import (
	"context"
	"fmt"

	"app/internal/database/migrations"
	"app/pkg/console"
	"app/pkg/database"
)

type MigrateCommand struct {
	*console.BaseCommand
}

func NewMigrateCommand() *MigrateCommand {
	cmd := &MigrateCommand{
		BaseCommand: console.NewCommand("migrate", "Database migration commands"),
	}
	return cmd
}

func (c *MigrateCommand) Configure(config *console.CommandConfig) {
	config.Name = "migrate"
	config.Description = "Database migration commands"
	config.Usage = "migrate [action]"
	config.Arguments = []console.Argument{
		{
			Name:        "action",
			Description: "Action to perform (run, rollback, reset, refresh, or status)",
			Required:    true,
		},
	}
}

func (c *MigrateCommand) Handle(ctx context.Context) error {
	args := ctx.Value("args").([]string)
	if len(args) < 2 {
		return fmt.Errorf("action required: run, rollback, reset, refresh, or status")
	}

	db := database.GetDB()
	migrator := migrations.InitMigrations(db)

	action := args[1]
	switch action {
	case "run":
		return migrator.RunPending()
	case "rollback":
		return migrator.Rollback()
	case "reset":
		return migrator.Reset()
	case "refresh":
		return migrator.Refresh()
	case "status":
		status, err := migrator.Status()
		if err != nil {
			return err
		}
		printMigrationStatus(status)
		return nil
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}

func printMigrationStatus(status []map[string]interface{}) {
	fmt.Println("\nMigration Status:")
	fmt.Println("------------------------------------------")
	fmt.Printf("%-30s %-10s %-20s\n", "Migration", "Status", "Batch")
	fmt.Println("------------------------------------------")

	for _, s := range status {
		name := s["name"].(string)
		status := s["status"].(string)
		batch := "N/A"
		if b, ok := s["batch"].(int); ok && b > 0 {
			batch = fmt.Sprintf("%d", b)
		}
		fmt.Printf("%-30s %-10s %-20s\n", name, status, batch)
	}
	fmt.Println("------------------------------------------")
}
