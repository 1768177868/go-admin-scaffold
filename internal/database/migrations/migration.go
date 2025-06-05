package migrations

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// MigrationRecord represents a database migration record
type MigrationRecord struct {
	ID        uint      `gorm:"primarykey"`
	Name      string    `gorm:"size:255;not null;unique"`
	Batch     int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

// Migration interface defines the contract for database migrations
type Migration interface {
	Up(db *gorm.DB) error
	Down(db *gorm.DB) error
	File() string
}

// MigrationDefinition provides a convenient way to create migrations using functions
type MigrationDefinition struct {
	up   func(db *gorm.DB) error
	down func(db *gorm.DB) error
	file string
}

// Up implements the Migration interface
func (m *MigrationDefinition) Up(db *gorm.DB) error {
	return m.up(db)
}

// Down implements the Migration interface
func (m *MigrationDefinition) Down(db *gorm.DB) error {
	return m.down(db)
}

// File implements the Migration interface
func (m *MigrationDefinition) File() string {
	return m.file
}

// NewMigration creates a new migration using the function-based approach
func NewMigration(file string, up, down func(db *gorm.DB) error) *MigrationDefinition {
	return &MigrationDefinition{
		up:   up,
		down: down,
		file: file,
	}
}

// RegisteredMigration represents a registered migration
type RegisteredMigration struct {
	Name       string
	Definition Migration
}

// Migrator handles database migrations
type Migrator struct {
	db         *gorm.DB
	migrations []RegisteredMigration
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make([]RegisteredMigration, 0),
	}
}

// Register registers a new migration
func (m *Migrator) Register(name string, migration Migration) {
	m.migrations = append(m.migrations, RegisteredMigration{
		Name:       name,
		Definition: migration,
	})
}

// CreateMigrationsTable creates the migrations table if it doesn't exist
func (m *Migrator) CreateMigrationsTable() error {
	return m.db.AutoMigrate(&MigrationRecord{})
}

// GetLastBatch gets the last batch number
func (m *Migrator) GetLastBatch() (int, error) {
	var lastBatch int
	err := m.db.Model(&MigrationRecord{}).Select("COALESCE(MAX(batch), 0)").Scan(&lastBatch).Error
	return lastBatch, err
}

// RunPending runs all pending migrations
func (m *Migrator) RunPending() error {
	// Create migrations table if not exists
	if err := m.CreateMigrationsTable(); err != nil {
		return err
	}

	// Get executed migrations
	var executed []MigrationRecord
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

	// Count pending migrations
	pendingCount := 0
	for _, migration := range m.migrations {
		if !executedNames[migration.Name] {
			pendingCount++
		}
	}

	if pendingCount == 0 {
		fmt.Printf("✅ \033[32mNo pending migrations found\033[0m\n")
		return nil
	}

	fmt.Printf("🚀 \033[34mRunning %d pending migration(s)...\033[0m\n", pendingCount)
	fmt.Println("────────────────────────────────────────")

	// Run pending migrations in registration order
	currentIndex := 0
	for _, migration := range m.migrations {
		if !executedNames[migration.Name] {
			currentIndex++
			fmt.Printf("📄 \033[33m[%d/%d]\033[0m %s", currentIndex, pendingCount, migration.Name)

			// Begin transaction
			err := m.db.Transaction(func(tx *gorm.DB) error {
				// Run migration
				if err := migration.Definition.Up(tx); err != nil {
					return err
				}

				// Record migration
				return tx.Create(&MigrationRecord{
					Name:      migration.Name,
					Batch:     lastBatch + 1,
					CreatedAt: time.Now(),
				}).Error
			})

			if err != nil {
				fmt.Printf(" \033[31m✗ FAILED\033[0m\n")
				return fmt.Errorf("failed to run migration %s: %w", migration.Name, err)
			}

			fmt.Printf(" \033[32m✓ COMPLETED\033[0m\n")
		}
	}

	fmt.Println("────────────────────────────────────────")
	fmt.Printf("🎉 \033[32mAll migrations completed successfully!\033[0m\n")
	return nil
}

// Rollback rolls back the last batch of migrations
func (m *Migrator) Rollback() error {
	lastBatch, err := m.GetLastBatch()
	if err != nil {
		return err
	}

	if lastBatch == 0 {
		fmt.Printf("ℹ️  \033[33mNo migrations to rollback\033[0m\n")
		return nil
	}

	var migrations []MigrationRecord
	if err := m.db.Where("batch = ?", lastBatch).Order("id DESC").Find(&migrations).Error; err != nil {
		return err
	}

	if len(migrations) == 0 {
		fmt.Printf("ℹ️  \033[33mNo migrations found for batch %d\033[0m\n", lastBatch)
		return nil
	}

	fmt.Printf("🔄 \033[34mRolling back %d migration(s) from batch %d...\033[0m\n", len(migrations), lastBatch)
	fmt.Println("────────────────────────────────────────")

	for i, migration := range migrations {
		fmt.Printf("📄 \033[33m[%d/%d]\033[0m %s", i+1, len(migrations), migration.Name)

		// Find the migration definition
		var def Migration
		for _, m := range m.migrations {
			if m.Name == migration.Name {
				def = m.Definition
				break
			}
		}

		if def != nil {
			err := m.db.Transaction(func(tx *gorm.DB) error {
				// Run down migration
				if err := def.Down(tx); err != nil {
					return err
				}

				// Remove migration record
				return tx.Delete(&migration).Error
			})

			if err != nil {
				fmt.Printf(" \033[31m✗ FAILED\033[0m\n")
				return fmt.Errorf("failed to rollback migration %s: %w", migration.Name, err)
			}

			fmt.Printf(" \033[32m✓ ROLLED BACK\033[0m\n")
		} else {
			fmt.Printf(" \033[33m⚠ SKIPPED (definition not found)\033[0m\n")
		}
	}

	fmt.Println("────────────────────────────────────────")
	fmt.Printf("🎉 \033[32mRollback completed successfully!\033[0m\n")
	return nil
}

// Reset rolls back all migrations
func (m *Migrator) Reset() error {
	var migrations []MigrationRecord
	if err := m.db.Order("id DESC").Find(&migrations).Error; err != nil {
		return err
	}

	if len(migrations) == 0 {
		fmt.Printf("ℹ️  \033[33mNo migrations to reset\033[0m\n")
		return nil
	}

	fmt.Printf("🔥 \033[34mResetting %d migration(s)...\033[0m\n", len(migrations))
	fmt.Println("────────────────────────────────────────")

	for i, migration := range migrations {
		fmt.Printf("📄 \033[33m[%d/%d]\033[0m %s", i+1, len(migrations), migration.Name)

		// Find the migration definition
		var def Migration
		for _, m := range m.migrations {
			if m.Name == migration.Name {
				def = m.Definition
				break
			}
		}

		if def != nil {
			err := m.db.Transaction(func(tx *gorm.DB) error {
				if err := def.Down(tx); err != nil {
					return err
				}
				return tx.Delete(&migration).Error
			})

			if err != nil {
				fmt.Printf(" \033[31m✗ FAILED\033[0m\n")
				return fmt.Errorf("failed to reset migration %s: %w", migration.Name, err)
			}

			fmt.Printf(" \033[32m✓ RESET\033[0m\n")
		} else {
			fmt.Printf(" \033[33m⚠ SKIPPED (definition not found)\033[0m\n")
		}
	}

	fmt.Println("────────────────────────────────────────")
	fmt.Printf("🎉 \033[32mReset completed successfully!\033[0m\n")
	return nil
}

// Refresh resets and reruns all migrations
func (m *Migrator) Refresh() error {
	fmt.Printf("🔄 \033[34mRefreshing migrations (Reset + Run)...\033[0m\n")
	fmt.Println("════════════════════════════════════════")

	fmt.Printf("📋 \033[33mStep 1: Resetting all migrations\033[0m\n")
	if err := m.Reset(); err != nil {
		return err
	}

	fmt.Printf("\n📋 \033[33mStep 2: Running all migrations\033[0m\n")
	return m.RunPending()
}

// Status returns the status of all migrations
func (m *Migrator) Status() ([]map[string]interface{}, error) {
	var executed []MigrationRecord
	if err := m.db.Find(&executed).Error; err != nil {
		return nil, err
	}

	executedMap := make(map[string]MigrationRecord)
	for _, migration := range executed {
		executedMap[migration.Name] = migration
	}

	var status []map[string]interface{}
	for _, migration := range m.migrations {
		if record, ok := executedMap[migration.Name]; ok {
			status = append(status, map[string]interface{}{
				"name":       migration.Name,
				"batch":      record.Batch,
				"created_at": record.CreatedAt,
				"status":     "Executed",
			})
		} else {
			status = append(status, map[string]interface{}{
				"name":       migration.Name,
				"batch":      0,
				"created_at": nil,
				"status":     "Pending",
			})
		}
	}

	return status, nil
}
