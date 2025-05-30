package seeders

import (
	"app/internal/database/seeder"
	"time"

	"gorm.io/gorm"
)

func init() {
	Register("permissions", &seeder.Seeder{
		Name:        "permissions",
		Description: "Create default permissions",
		Run: func(tx *gorm.DB) error {
			permissions := []map[string]interface{}{
				// 用户管理权限
				{
					"name":         "user:view",
					"display_name": "查看用户",
					"description":  "查看用户列表和详情",
					"module":       "user",
					"action":       "view",
					"resource":     "user",
					"status":       1,
					"created_at":   time.Now(),
					"updated_at":   time.Now(),
				},
				{
					"name":         "user:create",
					"display_name": "创建用户",
					"description":  "创建新用户",
					"module":       "user",
					"action":       "create",
					"resource":     "user",
					"status":       1,
					"created_at":   time.Now(),
					"updated_at":   time.Now(),
				},
				{
					"name":         "user:edit",
					"display_name": "编辑用户",
					"description":  "编辑用户信息",
					"module":       "user",
					"action":       "edit",
					"resource":     "user",
					"status":       1,
					"created_at":   time.Now(),
					"updated_at":   time.Now(),
				},
				{
					"name":         "user:delete",
					"display_name": "删除用户",
					"description":  "删除用户",
					"module":       "user",
					"action":       "delete",
					"resource":     "user",
					"status":       1,
					"created_at":   time.Now(),
					"updated_at":   time.Now(),
				},

				// 角色管理权限
				{
					"name":         "role:view",
					"display_name": "查看角色",
					"description":  "查看角色列表和详情",
					"module":       "rbac",
					"action":       "view",
					"resource":     "role",
					"status":       1,
					"created_at":   time.Now(),
					"updated_at":   time.Now(),
				},
				{
					"name":         "role:create",
					"display_name": "创建角色",
					"description":  "创建新角色",
					"module":       "rbac",
					"action":       "create",
					"resource":     "role",
					"status":       1,
					"created_at":   time.Now(),
					"updated_at":   time.Now(),
				},
				{
					"name":         "role:edit",
					"display_name": "编辑角色",
					"description":  "编辑角色信息",
					"module":       "rbac",
					"action":       "edit",
					"resource":     "role",
					"status":       1,
					"created_at":   time.Now(),
					"updated_at":   time.Now(),
				},
				{
					"name":         "role:delete",
					"display_name": "删除角色",
					"description":  "删除角色",
					"module":       "rbac",
					"action":       "delete",
					"resource":     "role",
					"status":       1,
					"created_at":   time.Now(),
					"updated_at":   time.Now(),
				},

				// 权限管理权限
				{
					"name":         "permission:view",
					"display_name": "查看权限",
					"description":  "查看权限列表和详情",
					"module":       "rbac",
					"action":       "view",
					"resource":     "permission",
					"status":       1,
					"created_at":   time.Now(),
					"updated_at":   time.Now(),
				},
				{
					"name":         "permission:create",
					"display_name": "创建权限",
					"description":  "创建新权限",
					"module":       "rbac",
					"action":       "create",
					"resource":     "permission",
					"status":       1,
					"created_at":   time.Now(),
					"updated_at":   time.Now(),
				},
				{
					"name":         "permission:edit",
					"display_name": "编辑权限",
					"description":  "编辑权限信息",
					"module":       "rbac",
					"action":       "edit",
					"resource":     "permission",
					"status":       1,
					"created_at":   time.Now(),
					"updated_at":   time.Now(),
				},
				{
					"name":         "permission:delete",
					"display_name": "删除权限",
					"description":  "删除权限",
					"module":       "rbac",
					"action":       "delete",
					"resource":     "permission",
					"status":       1,
					"created_at":   time.Now(),
					"updated_at":   time.Now(),
				},

				// 系统管理权限
				{
					"name":         "dashboard:view",
					"display_name": "查看仪表盘",
					"description":  "查看系统仪表盘",
					"module":       "system",
					"action":       "view",
					"resource":     "dashboard",
					"status":       1,
					"created_at":   time.Now(),
					"updated_at":   time.Now(),
				},
				{
					"name":         "log:view",
					"display_name": "查看日志",
					"description":  "查看系统日志",
					"module":       "system",
					"action":       "view",
					"resource":     "log",
					"status":       1,
					"created_at":   time.Now(),
					"updated_at":   time.Now(),
				},

				// 个人资料权限
				{
					"name":         "profile:view",
					"display_name": "查看个人资料",
					"description":  "查看个人资料",
					"module":       "profile",
					"action":       "view",
					"resource":     "profile",
					"status":       1,
					"created_at":   time.Now(),
					"updated_at":   time.Now(),
				},
				{
					"name":         "profile:edit",
					"display_name": "编辑个人资料",
					"description":  "编辑个人资料",
					"module":       "profile",
					"action":       "edit",
					"resource":     "profile",
					"status":       1,
					"created_at":   time.Now(),
					"updated_at":   time.Now(),
				},
			}

			return tx.Table("permissions").Create(permissions).Error
		},
	})
}
