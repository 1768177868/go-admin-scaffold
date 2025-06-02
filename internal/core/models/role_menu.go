package models

// RoleMenu represents the many-to-many relationship between roles and menus
type RoleMenu struct {
	RoleID uint `gorm:"primaryKey;column:role_id"`
	MenuID uint `gorm:"primaryKey;column:menu_id"`
}

// TableName specifies the table name for RoleMenu
func (RoleMenu) TableName() string {
	return "role_menus"
}
