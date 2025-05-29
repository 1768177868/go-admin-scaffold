package migrations

import "gorm.io/gorm"

var migrations = make(map[string]*MigrationDefinition)

// Register registers a migration
func Register(name string, migration *MigrationDefinition) {
	migrations[name] = migration
}

// GetMigrations returns all registered migrations
func GetMigrations() map[string]*MigrationDefinition {
	return migrations
}

// InitMigrations initializes all migrations
func InitMigrations(db *gorm.DB) *Migrator {
	migrator := NewMigrator(db)
	for name, migration := range migrations {
		migrator.Register(name, migration)
	}
	return migrator
}
