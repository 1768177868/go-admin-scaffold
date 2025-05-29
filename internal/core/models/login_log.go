package models

import (
	"time"
)

// LoginLog represents a user login record
type LoginLog struct {
	BaseModel
	UserID    uint      `gorm:"index" json:"user_id"`
	Username  string    `gorm:"size:50" json:"username"`
	IP        string    `gorm:"size:50" json:"ip"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	Status    int       `gorm:"default:1" json:"status"` // 1: success, 0: failed
	Message   string    `gorm:"size:255" json:"message"`
	LoginTime time.Time `json:"login_time"`
}

// TableName specifies the table name for LoginLog model
func (LoginLog) TableName() string {
	return "sys_login_logs"
}
