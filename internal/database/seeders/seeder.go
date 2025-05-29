package seeders

import (
	"app/internal/models"
	"app/pkg/database"

	"golang.org/x/crypto/bcrypt"
)

// Seed runs all database seeders
func Seed() error {
	if err := seedUsers(); err != nil {
		return err
	}
	if err := seedRoles(); err != nil {
		return err
	}
	if err := seedPermissions(); err != nil {
		return err
	}
	return nil
}

// seedUsers seeds the users table
func seedUsers() error {
	db := database.GetDB()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create admin user
	admin := &models.User{
		Username: "admin",
		Password: string(hashedPassword),
		Email:    "admin@example.com",
		Name:     "Administrator",
		Status:   1,
	}

	return db.Create(admin).Error
}

// seedRoles seeds the roles table
func seedRoles() error {
	db := database.GetDB()

	roles := []models.Role{
		{Name: "admin", DisplayName: "Administrator", Description: "System Administrator"},
		{Name: "user", DisplayName: "Normal User", Description: "Normal User"},
	}

	for _, role := range roles {
		if err := db.Create(&role).Error; err != nil {
			return err
		}
	}

	return nil
}

// seedPermissions seeds the permissions table
func seedPermissions() error {
	db := database.GetDB()

	permissions := []models.Permission{
		{Name: "users.manage", DisplayName: "Manage Users", Description: "Can manage users"},
		{Name: "roles.manage", DisplayName: "Manage Roles", Description: "Can manage roles"},
		{Name: "permissions.manage", DisplayName: "Manage Permissions", Description: "Can manage permissions"},
	}

	for _, permission := range permissions {
		if err := db.Create(&permission).Error; err != nil {
			return err
		}
	}

	return nil
}
