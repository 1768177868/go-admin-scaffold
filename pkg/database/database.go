package database

import (
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

type DBConfig struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// Setup initializes the database connection
func Setup(config *DBConfig) error {
	var err error
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			config.Database,
		)

		instance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	})
	return err
}

// DB returns the database instance
func DB() *gorm.DB {
	return instance
}
