package models

import (
	"gorm.io/gorm"
)

// MenuPermission represents a permission associated with a menu item
type MenuPermission struct {
	ID         uint   `json:"id" gorm:"primarykey"`
	MenuID     uint   `json:"menu_id" gorm:"not null;comment:'菜单ID'"`
	Permission string `json:"permission" gorm:"size:100;not null;comment:'权限标识'"`
	Name       string `json:"name" gorm:"size:100;not null;comment:'权限名称'"`
	Action     string `json:"action" gorm:"size:50;not null;comment:'操作类型：view,create,edit,delete'"`
	Status     int    `json:"status" gorm:"default:1;comment:'状态：0-禁用，1-启用'"`

	// 关联
	Menu  Menu   `json:"menu,omitempty" gorm:"foreignKey:MenuID"`
	Roles []Role `json:"roles,omitempty" gorm:"many2many:role_menu_permissions;"`

	// 时间戳
	CreatedAt CustomTime     `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt CustomTime     `json:"updated_at" gorm:"type:timestamp"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;type:timestamp"`
}

// TableName specifies the table name for MenuPermission model
func (MenuPermission) TableName() string {
	return "menu_permissions"
}

// IsActive returns true if the menu permission is active
func (mp *MenuPermission) IsActive() bool {
	return mp.Status == 1
}

// RoleMenuPermission represents the association between roles and menu permissions
type RoleMenuPermission struct {
	RoleID           uint `json:"role_id" gorm:"primaryKey;column:role_id"`
	MenuPermissionID uint `json:"menu_permission_id" gorm:"primaryKey;column:menu_permission_id"`
}

// TableName specifies the table name for RoleMenuPermission model
func (RoleMenuPermission) TableName() string {
	return "role_menu_permissions"
}
