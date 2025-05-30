package migrations

import (
	"time"

	"gorm.io/gorm"
)

func init() {
	Register("20240310_create_seeder_histories_table", &MigrationDefinition{
		Up: func(tx *gorm.DB) error {
			type SeederHistory struct {
				ID         uint      `gorm:"primarykey"`
				Name       string    `gorm:"size:255;not null;unique"`
				ExecutedAt time.Time `gorm:"type:timestamp;not null"`
			}

			// 先检查表是否存在
			var count int64
			if err := tx.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'seeder_histories'").Scan(&count).Error; err != nil {
				return err
			}

			if count == 0 {
				// 如果表不存在，使用 GORM 创建表
				if err := tx.AutoMigrate(&SeederHistory{}); err != nil {
					return err
				}

				// 创建索引
				if err := tx.Exec("CREATE INDEX idx_seeder_histories_executed_at ON seeder_histories(executed_at)").Error; err != nil {
					return err
				}
			} else {
				// 如果表存在，只修改字段类型
				if err := tx.Exec("ALTER TABLE seeder_histories MODIFY COLUMN executed_at TIMESTAMP NOT NULL").Error; err != nil {
					return err
				}
			}

			return nil
		},
		Down: func(tx *gorm.DB) error {
			// Drop index first (MySQL compatible syntax)
			if err := tx.Exec("ALTER TABLE seeder_histories DROP INDEX idx_seeder_histories_executed_at").Error; err != nil {
				// Ignore error if index doesn't exist
			}

			// Drop table
			return tx.Migrator().DropTable("seeder_histories")
		},
	})
}
