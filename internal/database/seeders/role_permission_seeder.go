package seeders

import (
	"app/internal/database/seeder"

	"gorm.io/gorm"
)

func init() {
	Register("role_permissions", &seeder.Seeder{
		Name:         "role_permissions",
		Description:  "Create role-permission associations",
		Dependencies: []string{"roles", "permissions"},
		Run: func(tx *gorm.DB) error {
			// 获取所有权限
			var permissions []struct {
				ID   uint
				Name string
			}
			if err := tx.Table("permissions").Select("id, name").Find(&permissions).Error; err != nil {
				return err
			}

			// 创建权限名称到ID的映射
			permissionMap := make(map[string]uint)
			for _, p := range permissions {
				permissionMap[p.Name] = p.ID
			}

			// 获取角色
			var adminRole struct{ ID uint }
			if err := tx.Table("roles").Where("code = ?", "admin").First(&adminRole).Error; err != nil {
				return err
			}

			var managerRole struct{ ID uint }
			if err := tx.Table("roles").Where("code = ?", "manager").First(&managerRole).Error; err != nil {
				return err
			}

			var userRole struct{ ID uint }
			if err := tx.Table("roles").Where("code = ?", "user").First(&userRole).Error; err != nil {
				return err
			}

			// 管理员角色 - 拥有所有权限
			var adminRolePermissions []map[string]interface{}
			for _, p := range permissions {
				adminRolePermissions = append(adminRolePermissions, map[string]interface{}{
					"role_id":       adminRole.ID,
					"permission_id": p.ID,
				})
			}

			// 管理者角色权限
			managerPermissionNames := []string{
				"dashboard:view",
				"user:view", "user:edit",
				"role:view",
				"log:view",
				"profile:view", "profile:edit",
			}
			var managerRolePermissions []map[string]interface{}
			for _, permName := range managerPermissionNames {
				if permID, ok := permissionMap[permName]; ok {
					managerRolePermissions = append(managerRolePermissions, map[string]interface{}{
						"role_id":       managerRole.ID,
						"permission_id": permID,
					})
				}
			}

			// 普通用户角色权限
			userPermissionNames := []string{
				"dashboard:view",
				"profile:view", "profile:edit",
			}
			var userRolePermissions []map[string]interface{}
			for _, permName := range userPermissionNames {
				if permID, ok := permissionMap[permName]; ok {
					userRolePermissions = append(userRolePermissions, map[string]interface{}{
						"role_id":       userRole.ID,
						"permission_id": permID,
					})
				}
			}

			// 插入角色权限关联
			if err := tx.Table("role_permissions").Create(adminRolePermissions).Error; err != nil {
				return err
			}
			if err := tx.Table("role_permissions").Create(managerRolePermissions).Error; err != nil {
				return err
			}
			if err := tx.Table("role_permissions").Create(userRolePermissions).Error; err != nil {
				return err
			}

			return nil
		},
	})
}
