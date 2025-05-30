package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// Migration represents a registered migration
type Migration struct {
	Name       string
	Definition *MigrationDefinition
}

// 使用切片替代map，保持注册顺序
var migrations []Migration

// Register registers a migration
func Register(name string, migration *MigrationDefinition) {
	// 检查是否已存在同名迁移
	for _, m := range migrations {
		if m.Name == name {
			panic(fmt.Sprintf("migration %s already exists", name))
		}
	}

	migrations = append(migrations, Migration{
		Name:       name,
		Definition: migration,
	})
}

// GetMigrations returns all registered migrations in registration order
func GetMigrations() []Migration {
	return migrations
}

// InitMigrations initializes all migrations in registration order
func InitMigrations(db *gorm.DB) *Migrator {
	migrator := NewMigrator(db)

	// 直接按注册顺序执行迁移
	for _, migration := range migrations {
		migrator.Register(migration.Name, migration.Definition)
	}

	return migrator
}
