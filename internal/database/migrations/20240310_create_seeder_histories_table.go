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
				ExecutedAt time.Time `gorm:"not null"`
			}

			// Create seeder_histories table
			if err := tx.AutoMigrate(&SeederHistory{}); err != nil {
				return err
			}

			// Add index (check if it exists first)
			var count int64
			tx.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'seeder_histories' AND index_name = 'idx_seeder_histories_executed_at'").Scan(&count)
			if count == 0 {
				if err := tx.Exec("CREATE INDEX idx_seeder_histories_executed_at ON seeder_histories(executed_at)").Error; err != nil {
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
