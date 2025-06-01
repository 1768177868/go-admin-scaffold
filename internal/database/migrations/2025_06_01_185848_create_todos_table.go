package migrations

import (
	"time"

	"gorm.io/gorm"
)

func init() {
	up := func(tx *gorm.DB) error {
		type Todo struct {
			ID        uint           `gorm:"primarykey"`
			Title     string         `gorm:"size:255;not null;comment:'标题'"`
			Content   string         `gorm:"type:text;comment:'内容'"`
			Status    int            `gorm:"default:0;comment:'状态：0-未完成，1-已完成'"`
			Priority  int            `gorm:"default:0;comment:'优先级：0-低，1-中，2-高'"`
			DueDate   *time.Time     `gorm:"type:timestamp;comment:'截止日期'"`
			UserID    uint           `gorm:"not null;comment:'创建者ID'"`
			CreatedAt time.Time      `gorm:"type:timestamp"`
			UpdatedAt time.Time      `gorm:"type:timestamp"`
			DeletedAt gorm.DeletedAt `gorm:"index;type:timestamp"`
		}

		// Create todos table
		if err := tx.AutoMigrate(&Todo{}); err != nil {
			return err
		}

		// Add indexes (check if they exist first)
		var count int64

		// Check and create idx_todos_user_id
		tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'todos' AND index_name = 'idx_todos_user_id'").Scan(&count)
		if count == 0 {
			if err := tx.Exec("CREATE INDEX idx_todos_user_id ON todos(user_id)").Error; err != nil {
				return err
			}
		}

		// Check and create idx_todos_status
		tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'todos' AND index_name = 'idx_todos_status'").Scan(&count)
		if count == 0 {
			if err := tx.Exec("CREATE INDEX idx_todos_status ON todos(status)").Error; err != nil {
				return err
			}
		}

		// Check and create idx_todos_due_date
		tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'todos' AND index_name = 'idx_todos_due_date'").Scan(&count)
		if count == 0 {
			if err := tx.Exec("CREATE INDEX idx_todos_due_date ON todos(due_date)").Error; err != nil {
				return err
			}
		}

		// Add foreign key constraint (check if it exists first)
		tx.Raw("SELECT COUNT(*) FROM information_schema.key_column_usage WHERE table_schema = DATABASE() AND table_name = 'todos' AND constraint_name = 'fk_todos_user_id'").Scan(&count)
		if count == 0 {
			if err := tx.Exec("ALTER TABLE todos ADD CONSTRAINT fk_todos_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE").Error; err != nil {
				return err
			}
		}

		return nil
	}

	down := func(tx *gorm.DB) error {
		// Drop foreign key first (MySQL compatible syntax)
		if err := tx.Exec("ALTER TABLE todos DROP FOREIGN KEY fk_todos_user_id").Error; err != nil {
			// Ignore error if foreign key doesn't exist
		}

		// Drop indexes (MySQL compatible syntax)
		if err := tx.Exec("ALTER TABLE todos DROP INDEX idx_todos_user_id").Error; err != nil {
			// Ignore error if index doesn't exist
		}
		if err := tx.Exec("ALTER TABLE todos DROP INDEX idx_todos_status").Error; err != nil {
			// Ignore error if index doesn't exist
		}
		if err := tx.Exec("ALTER TABLE todos DROP INDEX idx_todos_due_date").Error; err != nil {
			// Ignore error if index doesn't exist
		}

		// Drop table
		return tx.Migrator().DropTable("todos")
	}

	Register("create_todos_table", NewMigration("2025_06_01_185848_create_todos_table.go", up, down))
}
