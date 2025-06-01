package migrations

import (
	"time"

	"gorm.io/gorm"
)

func init() {
	up := func(tx *gorm.DB) error {
		type Role struct {
			ID          uint           `gorm:"primarykey"`
			Name        string         `gorm:"size:50;not null"`
			Code        string         `gorm:"size:50;not null;unique"`
			Description string         `gorm:"size:255"`
			Status      int            `gorm:"default:1;comment:'Status: 0-inactive, 1-active'"`
			PermList    []string       `gorm:"type:json"`
			CreatedAt   time.Time      `gorm:"type:timestamp"`
			UpdatedAt   time.Time      `gorm:"type:timestamp"`
			DeletedAt   gorm.DeletedAt `gorm:"index;type:timestamp"`
		}

		// Create roles table
		if err := tx.AutoMigrate(&Role{}); err != nil {
			return err
		}

		// Add indexes (check if they exist first)
		var count int64

		// Check and create idx_roles_code
		tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'roles' AND index_name = 'idx_roles_code'").Scan(&count)
		if count == 0 {
			if err := tx.Exec("CREATE INDEX idx_roles_code ON roles(code)").Error; err != nil {
				return err
			}
		}

		// Check and create idx_roles_status
		tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'roles' AND index_name = 'idx_roles_status'").Scan(&count)
		if count == 0 {
			if err := tx.Exec("CREATE INDEX idx_roles_status ON roles(status)").Error; err != nil {
				return err
			}
		}

		return nil
	}

	down := func(tx *gorm.DB) error {
		// Drop indexes first (MySQL compatible syntax)
		if err := tx.Exec("ALTER TABLE roles DROP INDEX idx_roles_code").Error; err != nil {
			// Ignore error if index doesn't exist
		}
		if err := tx.Exec("ALTER TABLE roles DROP INDEX idx_roles_status").Error; err != nil {
			// Ignore error if index doesn't exist
		}

		// Drop table
		return tx.Migrator().DropTable("roles")
	}

	Register("create_roles_table", NewMigration("20240302_create_roles_table.go", up, down))
}
