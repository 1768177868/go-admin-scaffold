package migrations

import (
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

			// Add indexes
			return tx.Exec(`
				CREATE INDEX idx_login_logs_login_time ON login_logs(login_time);
				CREATE INDEX idx_login_logs_status ON login_logs(status);
				CREATE INDEX idx_operation_logs_operation_time ON operation_logs(operation_time);
				CREATE INDEX idx_operation_logs_module ON operation_logs(module);
				CREATE INDEX idx_operation_logs_action ON operation_logs(action);
				CREATE INDEX idx_operation_logs_status ON operation_logs(status);
			`).Error
		},
		Down: func(tx *gorm.DB) error {
			// Drop indexes first
			if err := tx.Exec(`
				DROP INDEX IF EXISTS idx_login_logs_login_time ON login_logs;
				DROP INDEX IF EXISTS idx_login_logs_status ON login_logs;
				DROP INDEX IF EXISTS idx_operation_logs_operation_time ON operation_logs;
				DROP INDEX IF EXISTS idx_operation_logs_module ON operation_logs;
				DROP INDEX IF EXISTS idx_operation_logs_action ON operation_logs;
				DROP INDEX IF EXISTS idx_operation_logs_status ON operation_logs;
			`).Error; err != nil {
				return err
			}

			// Drop tables
			if err := tx.Migrator().DropTable("login_logs"); err != nil {
				return err
			}
			return tx.Migrator().DropTable("operation_logs")
		},
	})
}
