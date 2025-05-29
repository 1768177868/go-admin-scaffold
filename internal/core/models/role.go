package models

import (
	"time"

	"gorm.io/gorm"
)

// Role represents a user role in the system
type Role struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"size:100;not null;comment:'角色名称'"`
	Code        string         `json:"code" gorm:"size:50;not null;unique;comment:'角色编码'"`
	Description string         `json:"description" gorm:"size:255;comment:'角色描述'"`
	Status      int            `json:"status" gorm:"default:1;comment:'状态：0-禁用，1-启用'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Users       []User         `json:"users,omitempty" gorm:"many2many:user_roles;"`
	Permissions []Permission   `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
}

// TableName specifies the table name for Role model
func (Role) TableName() string {
	return "roles"
}

// IsActive returns true if the role is active
func (r *Role) IsActive() bool {
	return r.Status == 1
}

// GetPermissionNames returns a list of permission names for this role
func (r *Role) GetPermissionNames() []string {
	var permissions []string
	for _, perm := range r.Permissions {
		if perm.IsActive() {
			permissions = append(permissions, perm.Name)
		}
	}
	return permissions
}
