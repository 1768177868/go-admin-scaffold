package commands

import (
	"fmt"

	"app/internal/database/seeder"
	"app/internal/database/seeders"
	"app/pkg/database"
)

type SeedCommand struct {
	Name        string
	Description string
	Args        []string
}

func NewSeedCommand() *SeedCommand {
	return &SeedCommand{
		Name:        "seed",
		Description: "Database seeding commands",
		Args:        []string{"action", "seeder?"},
	}
}

func (c *SeedCommand) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("action required: run, reset, or status")
	}

	db := database.GetDB()
	seederManager := seeder.NewSeederManager(db)
	seeders.InitSeeders(seederManager)

	action := args[0]
	switch action {
	case "run":
		// If specific seeders are provided, run only those
		if len(args) > 1 {
			return seederManager.Run(args[1:]...)
		}
		// Otherwise run all seeders
		return seederManager.Run()

	case "reset":
		return seederManager.Reset()

	case "status":
		status, err := seederManager.Status()
		if err != nil {
			return err
		}
		printSeederStatus(status)
		return nil

	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}

func printSeederStatus(status []map[string]interface{}) {
	fmt.Println("\nSeeder Status:")
	fmt.Println("------------------------------------------")
	fmt.Printf("%-30s %-10s %-30s\n", "Seeder", "Status", "Description")
	fmt.Println("------------------------------------------")

	for _, s := range status {
		name := s["name"].(string)
		status := s["status"].(string)
		description := s["description"].(string)
		fmt.Printf("%-30s %-10s %-30s\n", name, status, description)
	}
	fmt.Println("------------------------------------------")
}
