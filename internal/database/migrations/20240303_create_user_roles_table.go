package migrations

import (
	"gorm.io/gorm"
)

func init() {
	Register("20240303_create_user_roles_table", &MigrationDefinition{
		Up: func(tx *gorm.DB) error {
			type UserRole struct {
				ID     uint `gorm:"primarykey"`
				UserID uint `gorm:"not null"`
				RoleID uint `gorm:"not null"`
			}

			// Create user_roles table
			if err := tx.AutoMigrate(&UserRole{}); err != nil {
				return err
			}

			// Add indexes (check if they exist first)
			var count int64

			// Check and create idx_user_roles_user_id
			tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'user_roles' AND index_name = 'idx_user_roles_user_id'").Scan(&count)
			if count == 0 {
				if err := tx.Exec("CREATE INDEX idx_user_roles_user_id ON user_roles(user_id)").Error; err != nil {
					return err
				}
			}

			// Check and create idx_user_roles_role_id
			tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'user_roles' AND index_name = 'idx_user_roles_role_id'").Scan(&count)
			if count == 0 {
				if err := tx.Exec("CREATE INDEX idx_user_roles_role_id ON user_roles(role_id)").Error; err != nil {
					return err
				}
			}

			// Add foreign key constraints (check if they exist first)
			tx.Raw("SELECT COUNT(*) FROM information_schema.key_column_usage WHERE table_schema = DATABASE() AND table_name = 'user_roles' AND constraint_name = 'fk_user_roles_user_id'").Scan(&count)
			if count == 0 {
				if err := tx.Exec("ALTER TABLE user_roles ADD CONSTRAINT fk_user_roles_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE").Error; err != nil {
					return err
				}
			}

			tx.Raw("SELECT COUNT(*) FROM information_schema.key_column_usage WHERE table_schema = DATABASE() AND table_name = 'user_roles' AND constraint_name = 'fk_user_roles_role_id'").Scan(&count)
			if count == 0 {
				if err := tx.Exec("ALTER TABLE user_roles ADD CONSTRAINT fk_user_roles_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE").Error; err != nil {
					return err
				}
			}

			return nil
		},
		Down: func(tx *gorm.DB) error {
			// Drop foreign keys first (MySQL compatible syntax)
			if err := tx.Exec("ALTER TABLE user_roles DROP FOREIGN KEY fk_user_roles_user_id").Error; err != nil {
				// Ignore error if foreign key doesn't exist
			}
			if err := tx.Exec("ALTER TABLE user_roles DROP FOREIGN KEY fk_user_roles_role_id").Error; err != nil {
				// Ignore error if foreign key doesn't exist
			}

			// Drop indexes (MySQL compatible syntax)
			if err := tx.Exec("ALTER TABLE user_roles DROP INDEX idx_user_roles_user_id").Error; err != nil {
				// Ignore error if index doesn't exist
			}
			if err := tx.Exec("ALTER TABLE user_roles DROP INDEX idx_user_roles_role_id").Error; err != nil {
				// Ignore error if index doesn't exist
			}

			// Drop table
			return tx.Migrator().DropTable("user_roles")
		},
	})
}
