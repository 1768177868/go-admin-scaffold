package migrations

import (
	"time"

	"gorm.io/gorm"
)

func init() {
	Register("20240301_create_users_table", &MigrationDefinition{
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

			// Add indexes (check if they exist first)
			var count int64

			// Check and create idx_users_email
			tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'users' AND index_name = 'idx_users_email'").Scan(&count)
			if count == 0 {
				if err := tx.Exec("CREATE INDEX idx_users_email ON users(email)").Error; err != nil {
					return err
				}
			}

			// Check and create idx_users_status
			tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'users' AND index_name = 'idx_users_status'").Scan(&count)
			if count == 0 {
				if err := tx.Exec("CREATE INDEX idx_users_status ON users(status)").Error; err != nil {
					return err
				}
			}

			return nil
		},
		Down: func(tx *gorm.DB) error {
			// Drop indexes first (MySQL compatible syntax)
			if err := tx.Exec("ALTER TABLE users DROP INDEX idx_users_email").Error; err != nil {
				// Ignore error if index doesn't exist
			}
			if err := tx.Exec("ALTER TABLE users DROP INDEX idx_users_status").Error; err != nil {
				// Ignore error if index doesn't exist
			}

			// Drop table
			return tx.Migrator().DropTable("users")
		},
	})
}
