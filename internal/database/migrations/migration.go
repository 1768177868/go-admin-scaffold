package migrations

import (
	"app/internal/core/models"
	"app/pkg/database"
)

// Migrate runs all database migrations
func Migrate() error {
	db := database.DB()

	// Add your models here
	return db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.LoginLog{},
		&models.OperationLog{},
	)
}

// Rollback rolls back all migrations
func Rollback() error {
	db := database.DB()

	// Drop tables in reverse order
	if err := db.Migrator().DropTable(&models.OperationLog{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&models.LoginLog{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&models.Role{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&models.User{}); err != nil {
		return err
	}

	return nil
}
