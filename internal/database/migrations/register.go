package migrations

import (
	"sort"

	"gorm.io/gorm"
)

var migrations = make(map[string]*MigrationDefinition)

// Register registers a migration
func Register(name string, migration *MigrationDefinition) {
	migrations[name] = migration
}

// GetMigrations returns all registered migrations
func GetMigrations() map[string]*MigrationDefinition {
	return migrations
}

// InitMigrations initializes all migrations in correct order
func InitMigrations(db *gorm.DB) *Migrator {
	migrator := NewMigrator(db)

	// Get sorted migration names
	var names []string
	for name := range migrations {
		names = append(names, name)
	}
	sort.Strings(names) // This ensures chronological order due to timestamp prefix

	// Register migrations in sorted order
	for _, name := range names {
		migrator.Register(name, migrations[name])
	}

	return migrator
}
