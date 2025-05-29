package migrations

import (
	"gorm.io/gorm"
)

func init() {
	Register("20240305_create_role_permissions_table", &MigrationDefinition{
		Up: func(tx *gorm.DB) error {
			type RolePermission struct {
				ID           uint `gorm:"primarykey"`
				RoleID       uint `gorm:"not null;comment:'角色ID'"`
				PermissionID uint `gorm:"not null;comment:'权限ID'"`
			}

			// Create role_permissions table
			if err := tx.AutoMigrate(&RolePermission{}); err != nil {
				return err
			}

			// Add indexes (check if they exist first)
			var count int64

			// Check and create idx_role_permissions_role_id
			tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'role_permissions' AND index_name = 'idx_role_permissions_role_id'").Scan(&count)
			if count == 0 {
				if err := tx.Exec("CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id)").Error; err != nil {
					return err
				}
			}

			// Check and create idx_role_permissions_permission_id
			tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'role_permissions' AND index_name = 'idx_role_permissions_permission_id'").Scan(&count)
			if count == 0 {
				if err := tx.Exec("CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id)").Error; err != nil {
					return err
				}
			}

			// Check and create unique index for role_id + permission_id
			tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'role_permissions' AND index_name = 'idx_role_permissions_unique'").Scan(&count)
			if count == 0 {
				if err := tx.Exec("CREATE UNIQUE INDEX idx_role_permissions_unique ON role_permissions(role_id, permission_id)").Error; err != nil {
					return err
				}
			}

			// Add foreign key constraints (check if they exist first)
			tx.Raw("SELECT COUNT(*) FROM information_schema.key_column_usage WHERE table_schema = DATABASE() AND table_name = 'role_permissions' AND constraint_name = 'fk_role_permissions_role_id'").Scan(&count)
			if count == 0 {
				if err := tx.Exec("ALTER TABLE role_permissions ADD CONSTRAINT fk_role_permissions_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE").Error; err != nil {
					return err
				}
			}

			tx.Raw("SELECT COUNT(*) FROM information_schema.key_column_usage WHERE table_schema = DATABASE() AND table_name = 'role_permissions' AND constraint_name = 'fk_role_permissions_permission_id'").Scan(&count)
			if count == 0 {
				if err := tx.Exec("ALTER TABLE role_permissions ADD CONSTRAINT fk_role_permissions_permission_id FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE").Error; err != nil {
					return err
				}
			}

			return nil
		},
		Down: func(tx *gorm.DB) error {
			// Drop foreign keys first (MySQL compatible syntax)
			if err := tx.Exec("ALTER TABLE role_permissions DROP FOREIGN KEY fk_role_permissions_role_id").Error; err != nil {
				// Ignore error if foreign key doesn't exist
			}
			if err := tx.Exec("ALTER TABLE role_permissions DROP FOREIGN KEY fk_role_permissions_permission_id").Error; err != nil {
				// Ignore error if foreign key doesn't exist
			}

			// Drop indexes (MySQL compatible syntax)
			if err := tx.Exec("ALTER TABLE role_permissions DROP INDEX idx_role_permissions_role_id").Error; err != nil {
				// Ignore error if index doesn't exist
			}
			if err := tx.Exec("ALTER TABLE role_permissions DROP INDEX idx_role_permissions_permission_id").Error; err != nil {
				// Ignore error if index doesn't exist
			}
			if err := tx.Exec("ALTER TABLE role_permissions DROP INDEX idx_role_permissions_unique").Error; err != nil {
				// Ignore error if index doesn't exist
			}

			// Drop table
			return tx.Migrator().DropTable("role_permissions")
		},
	})
}
