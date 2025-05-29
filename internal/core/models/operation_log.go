package models

import (
	"time"

	"gorm.io/gorm"
)

// OperationLog represents a user operation record
type OperationLog struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	UserID        uint           `gorm:"index" json:"user_id"`
	Username      string         `gorm:"size:50" json:"username"`
	IP            string         `gorm:"size:50" json:"ip"`
	Method        string         `gorm:"size:20" json:"method"`           // HTTP method
	Path          string         `gorm:"size:255" json:"path"`            // Request path
	Action        string         `gorm:"size:100" json:"action"`          // Operation action
	Module        string         `gorm:"size:50" json:"module"`           // System module
	BusinessID    string         `gorm:"size:50" json:"business_id"`      // Related business ID
	BusinessType  string         `gorm:"size:50" json:"business_type"`    // Business type
	RequestParams string         `gorm:"type:text" json:"request_params"` // Request parameters
	Status        int            `gorm:"default:1" json:"status"`         // 1: success, 0: failed
	ErrorMessage  string         `gorm:"size:500" json:"error_message"`   // Error message if failed
	OperationTime time.Time      `json:"operation_time"`                  // Operation time
	Duration      int64          `json:"duration"`                        // Request duration in milliseconds
	UserAgent     string         `gorm:"size:255" json:"user_agent"`      // User agent
	ReqBody       string         `gorm:"type:text" json:"req_body"`       // Request body
	RespBody      string         `gorm:"type:text" json:"resp_body"`      // Response body
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for OperationLog model
func (OperationLog) TableName() string {
	return "sys_operation_logs"
}
