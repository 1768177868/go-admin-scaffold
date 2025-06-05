package models

import (
	"gorm.io/gorm"
)

// Menu represents a menu item in the system
type Menu struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	Name      string `json:"name" gorm:"size:50;not null;comment:'菜单名称'"`
	Title     string `json:"title" gorm:"size:50;not null;comment:'菜单标题'"`
	Icon      string `json:"icon" gorm:"size:50;comment:'菜单图标'"`
	Path      string `json:"path" gorm:"size:200;comment:'菜单路径'"`
	Component string `json:"component" gorm:"size:200;comment:'组件路径'"`
	ParentID  *uint  `json:"parent_id" gorm:"comment:'父菜单ID'"`
	Sort      int    `json:"sort" gorm:"default:0;comment:'排序值'"`
	Type      int    `json:"type" gorm:"default:1;comment:'菜单类型：1-菜单，2-按钮'"`
	Visible   int    `json:"visible" gorm:"default:1;comment:'是否可见：0-隐藏，1-显示'"`
	Status    int    `json:"status" gorm:"default:1;comment:'状态：0-禁用，1-启用'"`
	KeepAlive bool   `json:"keep_alive" gorm:"default:false;comment:'是否缓存'"`
	External  bool   `json:"external" gorm:"default:false;comment:'是否外链'"`

	// 权限相关
	Permission string `json:"permission" gorm:"size:100;comment:'权限标识'"`

	// Meta 信息 (JSON)
	Meta string `json:"meta" gorm:"type:json;comment:'菜单元信息'"`

	// 关联
	Parent   *Menu   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []*Menu `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Roles    []Role  `json:"roles,omitempty" gorm:"many2many:role_menus;"`

	// 时间戳
	CreatedAt CustomTime     `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt CustomTime     `json:"updated_at" gorm:"type:timestamp"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;type:timestamp"`
}

// MenuMeta represents menu metadata
type MenuMeta struct {
	Title      string `json:"title"`
	Icon       string `json:"icon,omitempty"`
	Hidden     bool   `json:"hidden,omitempty"`
	AlwaysShow bool   `json:"alwaysShow,omitempty"`
	NoCache    bool   `json:"noCache,omitempty"`
	Affix      bool   `json:"affix,omitempty"`
	Breadcrumb bool   `json:"breadcrumb,omitempty"`
	ActiveMenu string `json:"activeMenu,omitempty"`
	KeepAlive  bool   `json:"keepAlive,omitempty"`
}

// TableName specifies the table name for Menu model
func (Menu) TableName() string {
	return "menus"
}

// IsMenu returns true if this is a menu (not a button)
func (m *Menu) IsMenu() bool {
	return m.Type == 1
}

// IsButton returns true if this is a button
func (m *Menu) IsButton() bool {
	return m.Type == 2
}

// IsVisible returns true if the menu is visible
func (m *Menu) IsVisible() bool {
	return m.Visible == 1
}

// IsActive returns true if the menu is active
func (m *Menu) IsActive() bool {
	return m.Status == 1
}
