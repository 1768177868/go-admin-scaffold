package database

import (
	"sync"

	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// Init initializes the database connection
func Init(dbConn *gorm.DB) {
	once.Do(func() {
		db = dbConn
	})
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}

// Close closes the database connection
func Close() error {
	return CloseMySQL()
}
