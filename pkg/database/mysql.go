package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

var mysqlDB *gorm.DB

func InitMySQL(config *Config) error {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	mysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	sqlDB, err := mysqlDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return nil
}

// GetMySQLDB returns the database instance
func GetMySQLDB() *gorm.DB {
	return mysqlDB
}

// CloseMySQL closes the database connection
func CloseMySQL() error {
	if mysqlDB != nil {
		sqlDB, err := mysqlDB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
