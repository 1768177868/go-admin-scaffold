package migrations

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

func init() {
	Register("20240310_create_logs_tables", &MigrationDefinition{
		Up: func(tx *gorm.DB) error {
			// Login logs table
			type LoginLog struct {
				ID        uint      `gorm:"primarykey"`
				UserID    uint      `gorm:"index"`
				Username  string    `gorm:"size:50"`
				IP        string    `gorm:"size:50"`
				UserAgent string    `gorm:"size:255"`
				Status    int       `gorm:"default:1;comment:'Status: 0-failed, 1-success'"`
				Message   string    `gorm:"size:255"`
				LoginTime time.Time `gorm:"not null"`
			}

			// Operation logs table
			type OperationLog struct {
				ID            uint      `gorm:"primarykey"`
				UserID        uint      `gorm:"index"`
				Username      string    `gorm:"size:50"`
				IP            string    `gorm:"size:50"`
				Method        string    `gorm:"size:20"`
				Path          string    `gorm:"size:255"`
				Action        string    `gorm:"size:100"`
				Module        string    `gorm:"size:100"`
				BusinessID    string    `gorm:"size:100"`
				BusinessType  string    `gorm:"size:100"`
				RequestParams string    `gorm:"type:text"`
				Status        int       `gorm:"default:1;comment:'Status: 0-failed, 1-success'"`
				ErrorMessage  string    `gorm:"size:255"`
				Duration      int64     `gorm:"comment:'Duration in milliseconds'"`
				OperationTime time.Time `gorm:"not null"`
			}

			// Create tables
			if err := tx.AutoMigrate(&LoginLog{}); err != nil {
				return err
			}
			if err := tx.AutoMigrate(&OperationLog{}); err != nil {
				return err
			}

			// Add indexes for login_logs (check if they exist first)
			var count int64

			// Check and create idx_login_logs_login_time
			tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'login_logs' AND index_name = 'idx_login_logs_login_time'").Scan(&count)
			if count == 0 {
				if err := tx.Exec("CREATE INDEX idx_login_logs_login_time ON login_logs(login_time)").Error; err != nil {
					return err
				}
			}

			// Check and create idx_login_logs_status
			tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'login_logs' AND index_name = 'idx_login_logs_status'").Scan(&count)
			if count == 0 {
				if err := tx.Exec("CREATE INDEX idx_login_logs_status ON login_logs(status)").Error; err != nil {
					return err
				}
			}

			// Add indexes for operation_logs (check if they exist first)
			indexQueries := []struct {
				name  string
				query string
			}{
				{"idx_operation_logs_operation_time", "CREATE INDEX idx_operation_logs_operation_time ON operation_logs(operation_time)"},
				{"idx_operation_logs_module", "CREATE INDEX idx_operation_logs_module ON operation_logs(module)"},
				{"idx_operation_logs_action", "CREATE INDEX idx_operation_logs_action ON operation_logs(action)"},
				{"idx_operation_logs_status", "CREATE INDEX idx_operation_logs_status ON operation_logs(status)"},
			}

			for _, idx := range indexQueries {
				tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'operation_logs' AND index_name = ?", idx.name).Scan(&count)
				if count == 0 {
					if err := tx.Exec(idx.query).Error; err != nil {
						return err
					}
				}
			}

			return nil
		},
		Down: func(tx *gorm.DB) error {
			// Drop indexes first (MySQL compatible syntax)
			indexes := []string{
				"idx_login_logs_login_time",
				"idx_login_logs_status",
				"idx_operation_logs_operation_time",
				"idx_operation_logs_module",
				"idx_operation_logs_action",
				"idx_operation_logs_status",
			}

			for _, idx := range indexes {
				tableName := tableNameFromIndex(idx)
				if tableName != "" {
					if err := tx.Exec("ALTER TABLE " + tableName + " DROP INDEX " + idx).Error; err != nil {
						// Ignore error if index doesn't exist
					}
				}
			}

			// Drop tables
			if err := tx.Migrator().DropTable("login_logs"); err != nil {
				return err
			}
			return tx.Migrator().DropTable("operation_logs")
		},
	})
}

// tableNameFromIndex returns the table name from an index name
func tableNameFromIndex(indexName string) string {
	if len(indexName) > 4 && indexName[:4] == "idx_" {
		parts := strings.Split(indexName[4:], "_")
		if len(parts) >= 2 {
			return parts[0] + "_" + parts[1]
		}
	}
	return ""
}
