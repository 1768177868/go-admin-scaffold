package migrations

import (
	"time"

	"gorm.io/gorm"
)

func init() {
	Register("20240304_create_permissions_table", &MigrationDefinition{
		Up: func(tx *gorm.DB) error {
			type Permission struct {
				ID          uint           `gorm:"primarykey"`
				Name        string         `gorm:"size:100;not null;unique;comment:'权限名称，如user:create'"`
				DisplayName string         `gorm:"size:100;not null;comment:'权限显示名称'"`
				Description string         `gorm:"size:255;comment:'权限描述'"`
				Module      string         `gorm:"size:50;not null;comment:'所属模块'"`
				Action      string         `gorm:"size:50;not null;comment:'操作类型：view,create,edit,delete'"`
				Resource    string         `gorm:"size:50;not null;comment:'资源类型：user,role,permission等'"`
				Status      int            `gorm:"default:1;comment:'状态：0-禁用，1-启用'"`
				CreatedAt   time.Time      `gorm:"type:timestamp"`
				UpdatedAt   time.Time      `gorm:"type:timestamp"`
				DeletedAt   gorm.DeletedAt `gorm:"index;type:timestamp"`
			}

			// Create permissions table
			if err := tx.AutoMigrate(&Permission{}); err != nil {
				return err
			}

			// Add indexes (check if they exist first)
			var count int64

			// Check and create idx_permissions_module
			tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'permissions' AND index_name = 'idx_permissions_module'").Scan(&count)
			if count == 0 {
				if err := tx.Exec("CREATE INDEX idx_permissions_module ON permissions(module)").Error; err != nil {
					return err
				}
			}

			// Check and create idx_permissions_resource_action
			tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'permissions' AND index_name = 'idx_permissions_resource_action'").Scan(&count)
			if count == 0 {
				if err := tx.Exec("CREATE INDEX idx_permissions_resource_action ON permissions(resource, action)").Error; err != nil {
					return err
				}
			}

			// Check and create idx_permissions_status
			tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'permissions' AND index_name = 'idx_permissions_status'").Scan(&count)
			if count == 0 {
				if err := tx.Exec("CREATE INDEX idx_permissions_status ON permissions(status)").Error; err != nil {
					return err
				}
			}

			return nil
		},
		Down: func(tx *gorm.DB) error {
			// Drop indexes first (MySQL compatible syntax)
			if err := tx.Exec("ALTER TABLE permissions DROP INDEX idx_permissions_module").Error; err != nil {
				// Ignore error if index doesn't exist
			}
			if err := tx.Exec("ALTER TABLE permissions DROP INDEX idx_permissions_resource_action").Error; err != nil {
				// Ignore error if index doesn't exist
			}
			if err := tx.Exec("ALTER TABLE permissions DROP INDEX idx_permissions_status").Error; err != nil {
				// Ignore error if index doesn't exist
			}

			// Drop table
			return tx.Migrator().DropTable("permissions")
		},
	})
}
