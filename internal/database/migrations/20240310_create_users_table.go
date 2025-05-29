package migrations

import (
	"time"

	"gorm.io/gorm"
)

func init() {
	migrations["20240310_create_users_table"] = &MigrationDefinition{
		Up: func(tx *gorm.DB) error {
			type User struct {
				ID        uint   `gorm:"primarykey"`
				Username  string `gorm:"size:50;not null;unique"`
				Password  string `gorm:"size:255;not null"`
				Email     string `gorm:"size:100;not null;unique"`
				Nickname  string `gorm:"size:50"`
				Avatar    string `gorm:"size:255"`
				Status    int    `gorm:"default:1;comment:'Status: 0-inactive, 1-active'"`
				LastLogin *time.Time
				CreatedAt time.Time
				UpdatedAt time.Time
				DeletedAt gorm.DeletedAt `gorm:"index"`
			}

			// Create users table
			if err := tx.AutoMigrate(&User{}); err != nil {
				return err
			}

			// Add indexes
			return tx.Exec(`
				CREATE INDEX idx_users_email ON users(email);
				CREATE INDEX idx_users_status ON users(status);
			`).Error
		},
		Down: func(tx *gorm.DB) error {
			// Drop indexes first
			if err := tx.Exec(`
				DROP INDEX IF EXISTS idx_users_email ON users;
				DROP INDEX IF EXISTS idx_users_status ON users;
			`).Error; err != nil {
				return err
			}

			// Drop table
			return tx.Migrator().DropTable("users")
		},
	}
}
