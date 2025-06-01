package migrations

import (
	"time"
	"gorm.io/gorm"
)

func init() {
	Register("create_todos_table", &MigrationDefinition{
		Up: func(tx *gorm.DB) error {
			// Create todos table
			type Todos struct {
				ID        uint           `gorm:"primarykey"`
				// Add your columns here
				CreatedAt time.Time      `gorm:"type:timestamp"`
				UpdatedAt time.Time      `gorm:"type:timestamp"`
				DeletedAt gorm.DeletedAt `gorm:"index;type:timestamp"`
			}

			return tx.Table("todos").AutoMigrate(&Todos{})
		},
		Down: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("todos")
		},
	})
}
