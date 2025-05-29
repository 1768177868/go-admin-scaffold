package models

// RolePermission represents the association between roles and permissions
type RolePermission struct {
	ID           uint `json:"id" gorm:"primarykey"`
	RoleID       uint `json:"role_id" gorm:"not null;comment:'角色ID'"`
	PermissionID uint `json:"permission_id" gorm:"not null;comment:'权限ID'"`
}

// TableName specifies the table name for RolePermission model
func (RolePermission) TableName() string {
	return "role_permissions"
}
