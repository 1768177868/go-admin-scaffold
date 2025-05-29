package models

import (
	"time"

	"gorm.io/gorm"
)

// Role represents a user role in the system
type Role struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"size:50;not null"`
	Code        string         `json:"code" gorm:"uniqueIndex;size:50;not null"`
	Description string         `json:"description" gorm:"size:200"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	PermList    []string       `json:"permissions" gorm:"type:json;column:permissions"`
	Users       []User         `json:"users,omitempty" gorm:"many2many:user_roles;"`
}

// TableName specifies the table name for Role model
func (Role) TableName() string {
	return "sys_roles"
}

// GetPermissions returns the list of permissions associated with the role
func (r *Role) GetPermissions() []string {
	if r.PermList == nil {
		return []string{}
	}
	return r.PermList
}
