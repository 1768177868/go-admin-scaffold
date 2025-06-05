package commands

import (
	"context"
	"fmt"
	"time"

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
	if len(status) == 0 {
		fmt.Printf("ℹ️  \033[33mNo migrations found\033[0m\n")
		return
	}

	executed := 0
	pending := 0
	for _, s := range status {
		if s["status"].(string) == "Executed" {
			executed++
		} else {
			pending++
		}
	}

	fmt.Printf("\n📊 \033[34mMigration Status Overview\033[0m\n")
	fmt.Println("════════════════════════════════════════════════════════")
	fmt.Printf("✅ Executed: \033[32m%d\033[0m  ⏳ Pending: \033[33m%d\033[0m  📝 Total: %d\n", executed, pending, len(status))
	fmt.Println("════════════════════════════════════════════════════════")
	fmt.Printf("%-40s %-12s %-8s %-20s\n", "Migration", "Status", "Batch", "Executed At")
	fmt.Println("────────────────────────────────────────────────────────")

	for _, s := range status {
		name := s["name"].(string)
		status := s["status"].(string)
		batch := "─"
		executedAt := "─"

		if b, ok := s["batch"].(int); ok && b > 0 {
			batch = fmt.Sprintf("%d", b)
		}

		if createdAt, ok := s["created_at"]; ok && createdAt != nil {
			if t, ok := createdAt.(time.Time); ok {
				executedAt = t.Format("2006-01-02 15:04:05")
			}
		}

		// Color code the status
		statusColor := ""
		if status == "Executed" {
			statusColor = "\033[32m✓ " + status + "\033[0m"
		} else {
			statusColor = "\033[33m⏳ " + status + "\033[0m"
		}

		// Truncate name if too long
		displayName := name
		if len(name) > 40 {
			displayName = name[:37] + "..."
		}

		fmt.Printf("%-40s %-20s %-8s %-20s\n", displayName, statusColor, batch, executedAt)
	}
	fmt.Println("════════════════════════════════════════════════════════")
}
