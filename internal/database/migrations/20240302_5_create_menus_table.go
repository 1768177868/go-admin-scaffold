package migrations

import (
	"time"

	"gorm.io/gorm"
)

func init() {
	up := func(tx *gorm.DB) error {
		type Menu struct {
			ID         uint           `gorm:"primarykey"`
			Name       string         `gorm:"size:50;not null;comment:'菜单名称'"`
			Title      string         `gorm:"size:50;not null;comment:'菜单标题'"`
			Icon       string         `gorm:"size:50;comment:'菜单图标'"`
			Path       string         `gorm:"size:200;comment:'菜单路径'"`
			Component  string         `gorm:"size:200;comment:'组件路径'"`
			ParentID   *uint          `gorm:"comment:'父菜单ID'"`
			Sort       int            `gorm:"default:0;comment:'排序值'"`
			Type       int            `gorm:"default:1;comment:'菜单类型：1-菜单，2-按钮'"`
			Visible    int            `gorm:"default:1;comment:'是否可见：0-隐藏，1-显示'"`
			Status     int            `gorm:"default:1;comment:'状态：0-禁用，1-启用'"`
			KeepAlive  bool           `gorm:"default:false;comment:'是否缓存'"`
			External   bool           `gorm:"default:false;comment:'是否外链'"`
			Permission string         `gorm:"size:100;comment:'权限标识'"`
			Meta       string         `gorm:"type:json;comment:'菜单元信息'"`
			CreatedAt  time.Time      `gorm:"type:timestamp"`
			UpdatedAt  time.Time      `gorm:"type:timestamp"`
			DeletedAt  gorm.DeletedAt `gorm:"index;type:timestamp"`
		}

		// Create menus table
		if err := tx.AutoMigrate(&Menu{}); err != nil {
			return err
		}

		// Add indexes
		var count int64

		// Check and create idx_menus_parent_id
		tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'menus' AND index_name = 'idx_menus_parent_id'").Scan(&count)
		if count == 0 {
			if err := tx.Exec("CREATE INDEX idx_menus_parent_id ON menus(parent_id)").Error; err != nil {
				return err
			}
		}

		// Check and create idx_menus_status
		tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'menus' AND index_name = 'idx_menus_status'").Scan(&count)
		if count == 0 {
			if err := tx.Exec("CREATE INDEX idx_menus_status ON menus(status)").Error; err != nil {
				return err
			}
		}

		// Check and create idx_menus_sort
		tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'menus' AND index_name = 'idx_menus_sort'").Scan(&count)
		if count == 0 {
			if err := tx.Exec("CREATE INDEX idx_menus_sort ON menus(sort)").Error; err != nil {
				return err
			}
		}

		// Check and create idx_menus_type
		tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'menus' AND index_name = 'idx_menus_type'").Scan(&count)
		if count == 0 {
			if err := tx.Exec("CREATE INDEX idx_menus_type ON menus(type)").Error; err != nil {
				return err
			}
		}

		// Check and create idx_menus_visible
		tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'menus' AND index_name = 'idx_menus_visible'").Scan(&count)
		if count == 0 {
			if err := tx.Exec("CREATE INDEX idx_menus_visible ON menus(visible)").Error; err != nil {
				return err
			}
		}

		return nil
	}

	down := func(tx *gorm.DB) error {
		// Drop indexes first
		indexes := []string{
			"idx_menus_parent_id",
			"idx_menus_status",
			"idx_menus_sort",
			"idx_menus_type",
			"idx_menus_visible",
		}

		for _, idx := range indexes {
			if err := tx.Exec("ALTER TABLE menus DROP INDEX " + idx).Error; err != nil {
				// Ignore error if index doesn't exist
			}
		}

		// Drop table
		return tx.Migrator().DropTable("menus")
	}

	Register("create_menus_table", NewMigration("20240302_5_create_menus_table.go", up, down))
}
