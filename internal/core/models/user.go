package models

import (
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	Username     string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password     string         `json:"-" gorm:"size:255;not null"`
	Email        string         `json:"email" gorm:"uniqueIndex;size:100"`
	Nickname     string         `json:"nickname" gorm:"size:50"`
	Avatar       string         `json:"avatar" gorm:"size:255"`
	Status       int            `json:"status" gorm:"default:1"`
	IsSuperAdmin bool           `json:"is_super_admin" gorm:"-"` // Virtual field, not stored in database
	Roles        []Role         `json:"roles" gorm:"many2many:user_roles"`
	CreatedAt    CustomTime     `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt    CustomTime     `json:"updated_at" gorm:"type:timestamp"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index;type:timestamp"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// BeforeSave hook is called before saving the user
func (u *User) BeforeSave(tx *gorm.DB) error {
	// Add any validation or data processing before save
	return nil
}

// ValidatePassword checks if the provided password matches the user's password
func (u *User) ValidatePassword(password string) bool {
	// TODO: Implement password validation using bcrypt
	return false
}
