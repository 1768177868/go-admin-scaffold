package models

import (
	"time"

	"gorm.io/gorm"
)

// LoginLog represents a login log record
type LoginLog struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"index" json:"user_id"`
	Username  string         `gorm:"size:50" json:"username"`
	IP        string         `gorm:"size:50" json:"ip"`
	UserAgent string         `gorm:"size:255" json:"user_agent"`
	Status    int            `gorm:"default:1" json:"status"`
	Message   string         `gorm:"size:255" json:"message"`
	LoginTime time.Time      `gorm:"not null" json:"login_time"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name
func (LoginLog) TableName() string {
	return "login_logs"
}
