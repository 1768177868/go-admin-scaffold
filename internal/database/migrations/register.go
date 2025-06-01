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

// AddPhoneToUsers migration adds phone column to users table
type AddPhoneToUsers struct{}

func (m *AddPhoneToUsers) Up(db *gorm.DB) error {
	return db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(20) DEFAULT NULL COMMENT '手机号' AFTER email").Error
}

func (m *AddPhoneToUsers) Down(db *gorm.DB) error {
	return db.Exec("ALTER TABLE users DROP COLUMN IF EXISTS phone").Error
}

func (m *AddPhoneToUsers) File() string {
	return "2025_06_02_000300_add_phone_to_users.go"
}
