package migrations

import (
	"gorm.io/gorm"
)

func init() {
	Register("20240310_create_user_roles_table", &MigrationDefinition{
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

			// Add indexes and foreign keys
			return tx.Exec(`
				CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
				CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
				ALTER TABLE user_roles ADD CONSTRAINT fk_user_roles_user_id 
					FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
				ALTER TABLE user_roles ADD CONSTRAINT fk_user_roles_role_id 
					FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;
			`).Error
		},
		Down: func(tx *gorm.DB) error {
			// Drop foreign keys first
			if err := tx.Exec(`
				ALTER TABLE user_roles DROP FOREIGN KEY IF EXISTS fk_user_roles_user_id;
				ALTER TABLE user_roles DROP FOREIGN KEY IF EXISTS fk_user_roles_role_id;
				DROP INDEX IF EXISTS idx_user_roles_user_id ON user_roles;
				DROP INDEX IF EXISTS idx_user_roles_role_id ON user_roles;
			`).Error; err != nil {
				return err
			}

			// Drop table
			return tx.Migrator().DropTable("user_roles")
		},
	})
}
