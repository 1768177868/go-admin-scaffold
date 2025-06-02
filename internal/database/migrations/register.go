package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// 使用切片替代map，保持注册顺序
var registeredMigrations []RegisteredMigration

// Register registers a migration
func Register(name string, migration Migration) {
	// 检查是否已存在同名迁移
	for _, m := range registeredMigrations {
		if m.Name == name {
			panic(fmt.Sprintf("migration %s already exists", name))
		}
	}

	registeredMigrations = append(registeredMigrations, RegisteredMigration{
		Name:       name,
		Definition: migration,
	})
}

// GetMigrations returns all registered migrations in registration order
func GetMigrations() []RegisteredMigration {
	return registeredMigrations
}

// InitMigrations initializes all migrations in registration order
func InitMigrations(db *gorm.DB) *Migrator {
	migrator := NewMigrator(db)

	// 直接按注册顺序执行迁移
	for _, migration := range registeredMigrations {
		migrator.Register(migration.Name, migration.Definition)
	}

	return migrator
}
