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

// Seed initializes the database with default data
func Seed() error {
	db := database.DB()

	// Create default admin role with permissions
	adminRole := &models.Role{
		Name:        "admin",
		Code:        "admin",
		Description: "System administrator with full access",
		Status:      1,
		PermList:    []string{"users.manage", "roles.manage", "logs.view"},
	}

	if err := db.Where(models.Role{Name: "admin"}).FirstOrCreate(adminRole).Error; err != nil {
		return err
	}

	// Create default admin user
	adminUser := &models.User{
		Username: "admin",
		Password: "$2a$10$ThyIwBtFCrqtP8OzNTHxdOgkM8/zXJoZF.ZLEgRy4F6qv/SGvhNx2", // password: admin123
		Email:    "admin@example.com",
		Nickname: "Administrator",
		Status:   1,
		Roles:    []models.Role{*adminRole},
	}

	if err := db.Where(models.User{Username: "admin"}).FirstOrCreate(adminUser).Error; err != nil {
		return err
	}

	return nil
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
