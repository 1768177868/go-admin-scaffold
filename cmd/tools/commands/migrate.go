package commands

import (
	"fmt"

	"app/internal/database/migrations"
	"app/pkg/database"
)

type MigrateCommand struct {
	Name        string
	Description string
	Args        []string
}

func NewMigrateCommand() *MigrateCommand {
	return &MigrateCommand{
		Name:        "migrate",
		Description: "Database migration commands",
		Args:        []string{"action"},
	}
}

func (c *MigrateCommand) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("action required: run, rollback, reset, refresh, or status")
	}

	db := database.GetDB()
	migrator := migrations.InitMigrations(db)

	action := args[0]
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
