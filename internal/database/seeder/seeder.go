package seeder

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Seeder represents a database seeder
type Seeder struct {
	Name         string
	Description  string
	Run          func(tx *gorm.DB) error
	Dependencies []string
}

// SeederManager manages database seeders
type SeederManager struct {
	db      *gorm.DB
	seeders map[string]*Seeder
}

// NewSeederManager creates a new seeder manager
func NewSeederManager(db *gorm.DB) *SeederManager {
	return &SeederManager{
		db:      db,
		seeders: make(map[string]*Seeder),
	}
}

// Register registers a new seeder
func (m *SeederManager) Register(name string, seeder *Seeder) {
	m.seeders[name] = seeder
}

// CreateSeedersTable creates the seeders table if it doesn't exist
func (m *SeederManager) CreateSeedersTable() error {
	type SeederHistory struct {
		ID         uint      `gorm:"primarykey"`
		Name       string    `gorm:"size:255;not null;unique"`
		ExecutedAt time.Time `gorm:"not null"`
	}
	return m.db.AutoMigrate(&SeederHistory{})
}

// Run executes specified seeders
func (m *SeederManager) Run(names ...string) error {
	if err := m.CreateSeedersTable(); err != nil {
		return err
	}

	// If no specific seeders are specified, run all
	if len(names) == 0 {
		for name := range m.seeders {
			names = append(names, name)
		}
	}

	// Check which seeders have already been executed
	var executed []struct {
		Name string
	}
	if err := m.db.Table("seeder_histories").Select("name").Find(&executed).Error; err != nil {
		return err
	}

	executedMap := make(map[string]bool)
	for _, e := range executed {
		executedMap[e.Name] = true
	}

	// Build dependency graph for only non-executed seeders
	graph := make(map[string][]string)
	for _, name := range names {
		if seeder, ok := m.seeders[name]; ok {
			if !executedMap[name] {
				graph[name] = seeder.Dependencies
			}
		} else {
			return fmt.Errorf("seeder not found: %s", name)
		}
	}

	// Resolve dependencies
	executedInThisRun := make(map[string]bool)
	var execute func(name string) error

	execute = func(name string) error {
		if executedInThisRun[name] || executedMap[name] {
			return nil
		}

		seeder := m.seeders[name]
		for _, dep := range seeder.Dependencies {
			if err := execute(dep); err != nil {
				return err
			}
		}

		// Execute seeder in transaction
		err := m.db.Transaction(func(tx *gorm.DB) error {
			if err := seeder.Run(tx); err != nil {
				return err
			}

			// Record execution
			return tx.Exec(
				"INSERT INTO seeder_histories (name, executed_at) VALUES (?, ?)",
				name,
				time.Now(),
			).Error
		})

		if err != nil {
			return fmt.Errorf("failed to run seeder %s: %w", name, err)
		}

		executedInThisRun[name] = true
		fmt.Printf("Seeded: %s\n", name)
		return nil
	}

	// Execute seeders
	for _, name := range names {
		if !executedMap[name] {
			if err := execute(name); err != nil {
				return err
			}
		} else {
			fmt.Printf("Skipped: %s (already executed)\n", name)
		}
	}

	return nil
}

// Reset removes all seeded data
func (m *SeederManager) Reset() error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// Get all executed seeders in reverse dependency order
		var executed []struct {
			Name string
		}
		if err := tx.Table("seeder_histories").Select("name").Find(&executed).Error; err != nil {
			return err
		}

		// Clear data for each executed seeder in reverse order
		executedMap := make(map[string]bool)
		for _, e := range executed {
			executedMap[e.Name] = true
		}

		// Define table clearing order (reverse of dependency order)
		clearOrder := []string{"user_roles", "users", "roles"}

		for _, seederName := range clearOrder {
			if executedMap[seederName] {
				switch seederName {
				case "user_roles":
					if err := tx.Exec("DELETE FROM user_roles").Error; err != nil {
						return err
					}
				case "users":
					if err := tx.Exec("DELETE FROM users").Error; err != nil {
						return err
					}
				case "roles":
					if err := tx.Exec("DELETE FROM roles").Error; err != nil {
						return err
					}
				}
			}
		}

		// Clear seeder history
		if err := tx.Exec("DELETE FROM seeder_histories").Error; err != nil {
			return err
		}

		return nil
	})
}

// Status returns the status of all seeders
func (m *SeederManager) Status() ([]map[string]interface{}, error) {
	var executed []struct {
		Name       string
		ExecutedAt time.Time
	}

	if err := m.db.Table("seeder_histories").Find(&executed).Error; err != nil {
		return nil, err
	}

	executedMap := make(map[string]time.Time)
	for _, e := range executed {
		executedMap[e.Name] = e.ExecutedAt
	}

	var status []map[string]interface{}
	for name, seeder := range m.seeders {
		if executedAt, ok := executedMap[name]; ok {
			status = append(status, map[string]interface{}{
				"name":        name,
				"description": seeder.Description,
				"executed_at": executedAt,
				"status":      "Executed",
			})
		} else {
			status = append(status, map[string]interface{}{
				"name":        name,
				"description": seeder.Description,
				"executed_at": nil,
				"status":      "Pending",
			})
		}
	}

	return status, nil
}
