package migrations

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration struct {
	ID        uint      `gorm:"primarykey"`
	Name      string    `gorm:"size:255;not null;unique"`
	Batch     int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

// MigrationFunc defines a migration function
type MigrationFunc func(tx *gorm.DB) error

// MigrationDefinition defines a migration with up and down functions
type MigrationDefinition struct {
	Up   MigrationFunc
	Down MigrationFunc
}

// Migrator handles database migrations
type Migrator struct {
	db         *gorm.DB
	migrations map[string]*MigrationDefinition
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make(map[string]*MigrationDefinition),
	}
}

// Register registers a new migration
func (m *Migrator) Register(name string, migration *MigrationDefinition) {
	m.migrations[name] = migration
}

// CreateMigrationsTable creates the migrations table if it doesn't exist
func (m *Migrator) CreateMigrationsTable() error {
	return m.db.AutoMigrate(&Migration{})
}

// GetLastBatch gets the last batch number
func (m *Migrator) GetLastBatch() (int, error) {
	var lastBatch int
	err := m.db.Model(&Migration{}).Select("COALESCE(MAX(batch), 0)").Scan(&lastBatch).Error
	return lastBatch, err
}

// RunPending runs all pending migrations
func (m *Migrator) RunPending() error {
	// Create migrations table if not exists
	if err := m.CreateMigrationsTable(); err != nil {
		return err
	}

	// Get executed migrations
	var executed []Migration
	if err := m.db.Find(&executed).Error; err != nil {
		return err
	}

	// Get last batch number
	lastBatch, err := m.GetLastBatch()
	if err != nil {
		return err
	}

	// Track executed migration names
	executedNames := make(map[string]bool)
	for _, migration := range executed {
		executedNames[migration.Name] = true
	}

	// Run pending migrations
	for name, migration := range m.migrations {
		if !executedNames[name] {
			// Begin transaction
			err := m.db.Transaction(func(tx *gorm.DB) error {
				// Run migration
				if err := migration.Up(tx); err != nil {
					return err
				}

				// Record migration
				return tx.Create(&Migration{
					Name:      name,
					Batch:     lastBatch + 1,
					CreatedAt: time.Now(),
				}).Error
			})

			if err != nil {
				return fmt.Errorf("failed to run migration %s: %w", name, err)
			}

			fmt.Printf("Migrated: %s\n", name)
		}
	}

	return nil
}

// Rollback rolls back the last batch of migrations
func (m *Migrator) Rollback() error {
	lastBatch, err := m.GetLastBatch()
	if err != nil {
		return err
	}

	var migrations []Migration
	if err := m.db.Where("batch = ?", lastBatch).Order("id DESC").Find(&migrations).Error; err != nil {
		return err
	}

	for _, migration := range migrations {
		if def, ok := m.migrations[migration.Name]; ok {
			err := m.db.Transaction(func(tx *gorm.DB) error {
				// Run down migration
				if err := def.Down(tx); err != nil {
					return err
				}

				// Remove migration record
				return tx.Delete(&migration).Error
			})

			if err != nil {
				return fmt.Errorf("failed to rollback migration %s: %w", migration.Name, err)
			}

			fmt.Printf("Rolled back: %s\n", migration.Name)
		}
	}

	return nil
}

// Reset rolls back all migrations
func (m *Migrator) Reset() error {
	var migrations []Migration
	if err := m.db.Order("id DESC").Find(&migrations).Error; err != nil {
		return err
	}

	for _, migration := range migrations {
		if def, ok := m.migrations[migration.Name]; ok {
			err := m.db.Transaction(func(tx *gorm.DB) error {
				if err := def.Down(tx); err != nil {
					return err
				}
				return tx.Delete(&migration).Error
			})

			if err != nil {
				return fmt.Errorf("failed to reset migration %s: %w", migration.Name, err)
			}

			fmt.Printf("Reset: %s\n", migration.Name)
		}
	}

	return nil
}

// Refresh resets and reruns all migrations
func (m *Migrator) Refresh() error {
	if err := m.Reset(); err != nil {
		return err
	}
	return m.RunPending()
}

// Status returns the status of all migrations
func (m *Migrator) Status() ([]map[string]interface{}, error) {
	var executed []Migration
	if err := m.db.Find(&executed).Error; err != nil {
		return nil, err
	}

	executedMap := make(map[string]Migration)
	for _, migration := range executed {
		executedMap[migration.Name] = migration
	}

	var status []map[string]interface{}
	for name := range m.migrations {
		if migration, ok := executedMap[name]; ok {
			status = append(status, map[string]interface{}{
				"name":       name,
				"batch":      migration.Batch,
				"created_at": migration.CreatedAt,
				"status":     "Executed",
			})
		} else {
			status = append(status, map[string]interface{}{
				"name":       name,
				"batch":      0,
				"created_at": nil,
				"status":     "Pending",
			})
		}
	}

	return status, nil
}
