package migrations

import (
	"gorm.io/gorm"
)

func init() {
	up := func(tx *gorm.DB) error {
		type RoleMenu struct {
			RoleID uint `gorm:"primaryKey;column:role_id"`
			MenuID uint `gorm:"primaryKey;column:menu_id"`
		}

		// Create role_menus table
		if err := tx.AutoMigrate(&RoleMenu{}); err != nil {
			return err
		}

		// Add indexes and foreign keys
		var count int64

		// Check and create idx_role_menus_role_id
		tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'role_menus' AND index_name = 'idx_role_menus_role_id'").Scan(&count)
		if count == 0 {
			if err := tx.Exec("CREATE INDEX idx_role_menus_role_id ON role_menus(role_id)").Error; err != nil {
				return err
			}
		}

		// Check and create idx_role_menus_menu_id
		tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'role_menus' AND index_name = 'idx_role_menus_menu_id'").Scan(&count)
		if count == 0 {
			if err := tx.Exec("CREATE INDEX idx_role_menus_menu_id ON role_menus(menu_id)").Error; err != nil {
				return err
			}
		}

		// Check and create unique index
		tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'role_menus' AND index_name = 'idx_role_menus_unique'").Scan(&count)
		if count == 0 {
			if err := tx.Exec("CREATE UNIQUE INDEX idx_role_menus_unique ON role_menus(role_id, menu_id)").Error; err != nil {
				return err
			}
		}

		// Add foreign key constraints
		tx.Raw("SELECT COUNT(*) FROM information_schema.key_column_usage WHERE table_schema = DATABASE() AND table_name = 'role_menus' AND constraint_name = 'fk_role_menus_role_id'").Scan(&count)
		if count == 0 {
			if err := tx.Exec("ALTER TABLE role_menus ADD CONSTRAINT fk_role_menus_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE").Error; err != nil {
				return err
			}
		}

		tx.Raw("SELECT COUNT(*) FROM information_schema.key_column_usage WHERE table_schema = DATABASE() AND table_name = 'role_menus' AND constraint_name = 'fk_role_menus_menu_id'").Scan(&count)
		if count == 0 {
			if err := tx.Exec("ALTER TABLE role_menus ADD CONSTRAINT fk_role_menus_menu_id FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE").Error; err != nil {
				return err
			}
		}

		return nil
	}

	down := func(tx *gorm.DB) error {
		// Drop foreign keys first
		if err := tx.Exec("ALTER TABLE role_menus DROP FOREIGN KEY fk_role_menus_role_id").Error; err != nil {
			// Ignore error if foreign key doesn't exist
		}
		if err := tx.Exec("ALTER TABLE role_menus DROP FOREIGN KEY fk_role_menus_menu_id").Error; err != nil {
			// Ignore error if foreign key doesn't exist
		}

		// Drop indexes
		indexes := []string{
			"idx_role_menus_role_id",
			"idx_role_menus_menu_id",
			"idx_role_menus_unique",
		}

		for _, idx := range indexes {
			if err := tx.Exec("ALTER TABLE role_menus DROP INDEX " + idx).Error; err != nil {
				// Ignore error if index doesn't exist
			}
		}

		// Drop table
		return tx.Migrator().DropTable("role_menus")
	}

	Register("create_role_menus_table", NewMigration("20240311_create_role_menus_table.go", up, down))
}
