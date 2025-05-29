package database

import (
	"sync"

	"app/internal/config"

	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

// Init initializes the database connection
func Init(cfg *config.Config) error {
	var err error
	once.Do(func() {
		err = InitMySQL(&Config{
			Host:     cfg.MySQL.Host,
			Port:     cfg.MySQL.Port,
			Username: cfg.MySQL.Username,
			Password: cfg.MySQL.Password,
			Database: cfg.MySQL.Database,
		})
		if err == nil {
			instance = GetMySQLDB()
		}
	})
	return err
}

// DB returns the database instance
func DB() *gorm.DB {
	return instance
}

// Close closes the database connection
func Close() error {
	return CloseMySQL()
}
