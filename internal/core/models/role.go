package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string   `json:"name" gorm:"size:50;not null;unique"`
	Code        string   `json:"code" gorm:"size:50;not null;unique"`
	Description string   `json:"description" gorm:"size:200"`
	Status      int      `json:"status" gorm:"default:1"` // 1: active, 0: inactive
	PermList    []string `json:"permissions" gorm:"type:json;column:permissions"`
	Users       []User   `json:"users,omitempty" gorm:"many2many:user_roles;"`
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
