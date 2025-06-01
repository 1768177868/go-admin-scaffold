package migrations

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

func init() {
	up := func(tx *gorm.DB) error {
		// Login logs table
		type LoginLog struct {
			ID        uint           `gorm:"primarykey"`
			UserID    uint           `gorm:"index;comment:'用户ID'"`
			Username  string         `gorm:"size:50;comment:'用户名'"`
			IP        string         `gorm:"size:50;comment:'登录IP'"`
			UserAgent string         `gorm:"size:255;comment:'用户代理'"`
			Status    int            `gorm:"default:1;comment:'状态：0-失败，1-成功'"`
			Message   string         `gorm:"size:255;comment:'消息'"`
			LoginTime time.Time      `gorm:"type:timestamp;not null;comment:'登录时间'"`
			CreatedAt time.Time      `gorm:"type:timestamp"`
			UpdatedAt time.Time      `gorm:"type:timestamp"`
			DeletedAt gorm.DeletedAt `gorm:"index;type:timestamp"`
		}

		// Operation logs table
		type OperationLog struct {
			ID            uint           `gorm:"primarykey"`
			UserID        uint           `gorm:"index;comment:'用户ID'"`
			Username      string         `gorm:"size:50;comment:'用户名'"`
			IP            string         `gorm:"size:50;comment:'操作IP'"`
			Method        string         `gorm:"size:20;comment:'请求方法'"`
			Path          string         `gorm:"size:255;comment:'请求路径'"`
			Action        string         `gorm:"size:100;comment:'操作类型'"`
			Module        string         `gorm:"size:100;comment:'模块名称'"`
			BusinessID    string         `gorm:"size:100;comment:'业务ID'"`
			BusinessType  string         `gorm:"size:100;comment:'业务类型'"`
			RequestParams string         `gorm:"type:text;comment:'请求参数'"`
			Status        int            `gorm:"default:1;comment:'状态：0-失败，1-成功'"`
			ErrorMessage  string         `gorm:"size:255;comment:'错误信息'"`
			Duration      int64          `gorm:"comment:'执行时长(毫秒)'"`
			OperationTime time.Time      `gorm:"type:timestamp;not null;comment:'操作时间'"`
			UserAgent     string         `gorm:"size:255;comment:'用户代理'"`
			ReqBody       string         `gorm:"type:text;comment:'请求体'"`
			RespBody      string         `gorm:"type:text;comment:'响应体'"`
			CreatedAt     time.Time      `gorm:"type:timestamp"`
			UpdatedAt     time.Time      `gorm:"type:timestamp"`
			DeletedAt     gorm.DeletedAt `gorm:"index;type:timestamp"`
		}

		// Set table names and create tables
		if err := tx.Set("gorm:table_options", "").Table("login_logs").AutoMigrate(&LoginLog{}); err != nil {
			return err
		}
		if err := tx.Set("gorm:table_options", "").Table("operation_logs").AutoMigrate(&OperationLog{}); err != nil {
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

		// Check and create idx_login_logs_user_id
		tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'login_logs' AND index_name = 'idx_login_logs_user_id'").Scan(&count)
		if count == 0 {
			if err := tx.Exec("CREATE INDEX idx_login_logs_user_id ON login_logs(user_id)").Error; err != nil {
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
			{"idx_operation_logs_user_id", "CREATE INDEX idx_operation_logs_user_id ON operation_logs(user_id)"},
			{"idx_operation_logs_business_type", "CREATE INDEX idx_operation_logs_business_type ON operation_logs(business_type)"},
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
	}

	down := func(tx *gorm.DB) error {
		// Drop indexes first (MySQL compatible syntax)
		indexes := []string{
			"idx_login_logs_login_time",
			"idx_login_logs_status",
			"idx_login_logs_user_id",
			"idx_operation_logs_operation_time",
			"idx_operation_logs_module",
			"idx_operation_logs_action",
			"idx_operation_logs_status",
			"idx_operation_logs_user_id",
			"idx_operation_logs_business_type",
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
	}

	Register("create_logs_tables", NewMigration("20240310_create_logs_tables.go", up, down))
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
