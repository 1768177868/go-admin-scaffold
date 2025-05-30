package models

import (
	"gorm.io/gorm"
)

// OperationLog represents an operation log record
type OperationLog struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	UserID        uint           `gorm:"index" json:"user_id"`
	Username      string         `gorm:"size:50" json:"username"`
	IP            string         `gorm:"size:50" json:"ip"`
	Method        string         `gorm:"size:20" json:"method"`                         // HTTP method
	Path          string         `gorm:"size:255" json:"path"`                          // Request path
	Action        string         `gorm:"size:100" json:"action"`                        // Operation action
	Module        string         `gorm:"size:100" json:"module"`                        // System module
	BusinessID    string         `gorm:"size:100" json:"business_id"`                   // Related business ID
	BusinessType  string         `gorm:"size:100" json:"business_type"`                 // Business type
	RequestParams string         `gorm:"type:text" json:"request_params"`               // Request parameters
	Status        int            `gorm:"default:1" json:"status"`                       // 1: success, 0: failed
	ErrorMessage  string         `gorm:"size:255" json:"error_message"`                 // Error message if failed
	Duration      int64          `json:"duration"`                                      // Request duration in milliseconds
	OperationTime CustomTime     `gorm:"type:timestamp;not null" json:"operation_time"` // Operation time
	UserAgent     string         `gorm:"size:255" json:"user_agent"`                    // User agent
	ReqBody       string         `gorm:"type:text" json:"req_body"`                     // Request body
	RespBody      string         `gorm:"type:text" json:"resp_body"`                    // Response body
	CreatedAt     CustomTime     `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt     CustomTime     `gorm:"type:timestamp" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index;type:timestamp" json:"-"`
}

// TableName returns the table name
func (OperationLog) TableName() string {
	return "operation_logs"
}
