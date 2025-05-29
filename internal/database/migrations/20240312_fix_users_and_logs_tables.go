package migrations

import (
	"time"

	"gorm.io/gorm"
)

func init() {
	Register("20240312_fix_users_and_logs_tables", &MigrationDefinition{
		Up: func(tx *gorm.DB) error {
			// 1. 检查并重命名 last_login 列为 last_login_at
			var count int64
			tx.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 'users' AND column_name = 'last_login'").Scan(&count)
			if count > 0 {
				if err := tx.Exec("ALTER TABLE users CHANGE COLUMN last_login last_login_at datetime").Error; err != nil {
					return err
				}
			}

			// 2. 检查并创建 sys_login_logs 表
			tx.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'sys_login_logs'").Scan(&count)
			if count == 0 {
				type SysLoginLog struct {
					ID        uint      `gorm:"primarykey"`
					UserID    uint      `gorm:"index"`
					Username  string    `gorm:"size:50"`
					IP        string    `gorm:"size:50"`
					UserAgent string    `gorm:"size:255"`
					Status    int       `gorm:"default:1;comment:'Status: 0-failed, 1-success'"`
					Message   string    `gorm:"size:255"`
					LoginTime time.Time `gorm:"not null"`
				}

				// 创建表
				if err := tx.AutoMigrate(&SysLoginLog{}); err != nil {
					return err
				}

				// 添加索引
				indexQueries := []struct {
					name  string
					query string
				}{
					{"idx_sys_login_logs_login_time", "CREATE INDEX idx_sys_login_logs_login_time ON sys_login_logs(login_time)"},
					{"idx_sys_login_logs_status", "CREATE INDEX idx_sys_login_logs_status ON sys_login_logs(status)"},
				}

				for _, idx := range indexQueries {
					tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'sys_login_logs' AND index_name = ?", idx.name).Scan(&count)
					if count == 0 {
						if err := tx.Exec(idx.query).Error; err != nil {
							return err
						}
					}
				}
			}

			return nil
		},
		Down: func(tx *gorm.DB) error {
			// 1. 重命名 last_login_at 列回 last_login
			var count int64
			tx.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 'users' AND column_name = 'last_login_at'").Scan(&count)
			if count > 0 {
				if err := tx.Exec("ALTER TABLE users CHANGE COLUMN last_login_at last_login datetime").Error; err != nil {
					return err
				}
			}

			// 2. 删除 sys_login_logs 表的索引和表
			indexes := []string{
				"idx_sys_login_logs_login_time",
				"idx_sys_login_logs_status",
			}

			for _, idx := range indexes {
				if err := tx.Exec("ALTER TABLE sys_login_logs DROP INDEX " + idx).Error; err != nil {
					// 忽略索引不存在的错误
				}
			}

			return tx.Migrator().DropTable("sys_login_logs")
		},
	})
}
