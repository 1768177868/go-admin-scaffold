package migrations

import (
	"time"

	"gorm.io/gorm"
)

func init() {
	Register("20240310_create_roles_table", &MigrationDefinition{
		Up: func(tx *gorm.DB) error {
			type Role struct {
				ID          uint     `gorm:"primarykey"`
				Name        string   `gorm:"size:50;not null"`
				Code        string   `gorm:"size:50;not null;unique"`
				Description string   `gorm:"size:255"`
				Status      int      `gorm:"default:1;comment:'Status: 0-inactive, 1-active'"`
				PermList    []string `gorm:"type:json"`
				CreatedAt   time.Time
				UpdatedAt   time.Time
				DeletedAt   gorm.DeletedAt `gorm:"index"`
			}

			// Create roles table
			if err := tx.AutoMigrate(&Role{}); err != nil {
				return err
			}

			// Add indexes
			if err := tx.Exec("CREATE INDEX idx_roles_code ON roles(code)").Error; err != nil {
				return err
			}
			if err := tx.Exec("CREATE INDEX idx_roles_status ON roles(status)").Error; err != nil {
				return err
			}

			return nil
		},
		Down: func(tx *gorm.DB) error {
			// Drop indexes first
			if err := tx.Exec("DROP INDEX IF EXISTS idx_roles_code ON roles").Error; err != nil {
				return err
			}
			if err := tx.Exec("DROP INDEX IF EXISTS idx_roles_status ON roles").Error; err != nil {
				return err
			}

			// Drop table
			return tx.Migrator().DropTable("roles")
		},
	})
}
