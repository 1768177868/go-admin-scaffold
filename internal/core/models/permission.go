package models

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"size:100;not null;unique;comment:'权限名称，如user:create'"`
	DisplayName string         `json:"display_name" gorm:"size:100;not null;comment:'权限显示名称'"`
	Description string         `json:"description" gorm:"size:255;comment:'权限描述'"`
	Module      string         `json:"module" gorm:"size:50;not null;comment:'所属模块'"`
	Action      string         `json:"action" gorm:"size:50;not null;comment:'操作类型：view,create,edit,delete'"`
	Resource    string         `json:"resource" gorm:"size:50;not null;comment:'资源类型：user,role,permission等'"`
	Status      int            `json:"status" gorm:"default:1;comment:'状态：0-禁用，1-启用'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Roles       []Role         `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
}

// TableName specifies the table name for Permission model
func (Permission) TableName() string {
	return "permissions"
}

// IsActive returns true if the permission is active
func (p *Permission) IsActive() bool {
	return p.Status == 1
}
