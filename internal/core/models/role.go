package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// StringSlice 是一个自定义的字符串切片类型，用于处理 JSON 字段
type StringSlice []string

// Value 实现 driver.Valuer 接口
func (s StringSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

// Scan 实现 sql.Scanner 接口
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = StringSlice{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("unsupported type for StringSlice")
	}

	return json.Unmarshal(bytes, s)
}

// Role represents a user role in the system
type Role struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"size:50;not null;comment:'角色名称'"`
	Code        string         `json:"code" gorm:"size:50;not null;unique;comment:'角色编码'"`
	Description string         `json:"description" gorm:"size:255;comment:'角色描述'"`
	Status      int            `json:"status" gorm:"default:1;comment:'状态：0-禁用，1-启用'"`
	PermList    StringSlice    `json:"perm_list" gorm:"type:json"`
	CreatedAt   CustomTime     `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt   CustomTime     `json:"updated_at" gorm:"type:timestamp"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;type:timestamp"`
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
