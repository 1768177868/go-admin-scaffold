package models

import (
	"time"
)

type User struct {
	BaseModel
	Username    string    `gorm:"size:50;not null;unique" json:"username"`
	Password    string    `gorm:"size:255;not null" json:"-"`
	Nickname    string    `gorm:"size:50" json:"nickname"`
	Email       string    `gorm:"size:100;unique" json:"email"`
	Phone       string    `gorm:"size:20" json:"phone"`
	Avatar      string    `gorm:"size:255" json:"avatar"`
	Status      int       `gorm:"default:1" json:"status"` // 1: active, 0: inactive
	LastLoginAt time.Time `json:"last_login_at"`
	Roles       []Role    `gorm:"many2many:user_roles;" json:"roles"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "sys_users"
}

// BeforeSave hook is called before saving the user
func (u *User) BeforeSave() error {
	// Add any validation or data processing before save
	return nil
}

// ValidatePassword checks if the provided password matches the user's password
func (u *User) ValidatePassword(password string) bool {
	// TODO: Implement password validation using bcrypt
	return false
}
