package seeders

import (
	"encoding/json"
	"time"

	"app/internal/database/seeder"

	"gorm.io/gorm"
)

func init() {
	Register("menus", &seeder.Seeder{
		Name:        "menus",
		Description: "Create default menu structure",
		Run: func(tx *gorm.DB) error {
			// Menu metadata structure
			type MenuMeta struct {
				Title      string `json:"title"`
				Icon       string `json:"icon"`
				Hidden     bool   `json:"hidden"`
				AlwaysShow bool   `json:"alwaysShow"`
				NoCache    bool   `json:"noCache"`
				Affix      bool   `json:"affix"`
				Breadcrumb bool   `json:"breadcrumb"`
				ActiveMenu string `json:"activeMenu"`
			}

			// Dashboard menu
			dashboardMeta, _ := json.Marshal(MenuMeta{
				Title:      "仪表盘",
				Icon:       "Odometer",
				Affix:      true,
				Breadcrumb: true,
			})

			// System management menu
			systemMeta, _ := json.Marshal(MenuMeta{
				Title:      "系统管理",
				Icon:       "Setting",
				AlwaysShow: true,
				Breadcrumb: true,
			})

			// User management menu
			userMeta, _ := json.Marshal(MenuMeta{
				Title:      "用户管理",
				Icon:       "User",
				Breadcrumb: true,
			})

			// Role management menu
			roleMeta, _ := json.Marshal(MenuMeta{
				Title:      "角色管理",
				Icon:       "UserFilled",
				Breadcrumb: true,
			})

			// Permission management menu
			permissionMeta, _ := json.Marshal(MenuMeta{
				Title:      "权限管理",
				Icon:       "Key",
				Breadcrumb: true,
			})

			// Menu management menu
			menuMeta, _ := json.Marshal(MenuMeta{
				Title:      "菜单管理",
				Icon:       "Menu",
				Breadcrumb: true,
			})

			// Log management menu
			logMeta, _ := json.Marshal(MenuMeta{
				Title:      "日志管理",
				Icon:       "Document",
				AlwaysShow: true,
				Breadcrumb: true,
			})

			// Login log menu
			loginLogMeta, _ := json.Marshal(MenuMeta{
				Title:      "登录日志",
				Icon:       "Key",
				Breadcrumb: true,
			})

			// Operation log menu
			operationLogMeta, _ := json.Marshal(MenuMeta{
				Title:      "操作日志",
				Icon:       "Document",
				Breadcrumb: true,
			})

			// Profile menu
			profileMeta, _ := json.Marshal(MenuMeta{
				Title:      "个人中心",
				Icon:       "User",
				Hidden:     true,
				Breadcrumb: true,
			})

			menus := []map[string]interface{}{
				// Dashboard
				{
					"name":       "Dashboard",
					"title":      "仪表盘",
					"icon":       "Odometer",
					"path":       "/dashboard",
					"component":  "@/views/dashboard/index.vue",
					"parent_id":  nil,
					"sort":       1,
					"type":       1,
					"visible":    1,
					"status":     1,
					"keep_alive": false,
					"external":   false,
					"permission": "dashboard:view",
					"meta":       string(dashboardMeta),
					"created_at": time.Now(),
					"updated_at": time.Now(),
				},

				// System Management - Parent Menu
				{
					"name":       "System",
					"title":      "系统管理",
					"icon":       "Setting",
					"path":       "/system",
					"component":  "Layout",
					"parent_id":  nil,
					"sort":       2,
					"type":       1,
					"visible":    1,
					"status":     1,
					"keep_alive": false,
					"external":   false,
					"permission": "",
					"meta":       string(systemMeta),
					"created_at": time.Now(),
					"updated_at": time.Now(),
				},

				// User Management
				{
					"name":       "User",
					"title":      "用户管理",
					"icon":       "User",
					"path":       "user",
					"component":  "@/views/system/user/index.vue",
					"parent_id":  2, // System parent ID will be 2
					"sort":       1,
					"type":       1,
					"visible":    1,
					"status":     1,
					"keep_alive": false,
					"external":   false,
					"permission": "user:view",
					"meta":       string(userMeta),
					"created_at": time.Now(),
					"updated_at": time.Now(),
				},

				// Role Management
				{
					"name":       "Role",
					"title":      "角色管理",
					"icon":       "UserFilled",
					"path":       "role",
					"component":  "@/views/system/role/index.vue",
					"parent_id":  2,
					"sort":       2,
					"type":       1,
					"visible":    1,
					"status":     1,
					"keep_alive": false,
					"external":   false,
					"permission": "role:view",
					"meta":       string(roleMeta),
					"created_at": time.Now(),
					"updated_at": time.Now(),
				},

				// Permission Management
				{
					"name":       "Permission",
					"title":      "权限管理",
					"icon":       "Key",
					"path":       "permission",
					"component":  "@/views/system/permission/index.vue",
					"parent_id":  2,
					"sort":       3,
					"type":       1,
					"visible":    1,
					"status":     1,
					"keep_alive": false,
					"external":   false,
					"permission": "permission:view",
					"meta":       string(permissionMeta),
					"created_at": time.Now(),
					"updated_at": time.Now(),
				},

				// Menu Management
				{
					"name":       "Menu",
					"title":      "菜单管理",
					"icon":       "Menu",
					"path":       "menu",
					"component":  "@/views/system/menu/index.vue",
					"parent_id":  2,
					"sort":       4,
					"type":       1,
					"visible":    1,
					"status":     1,
					"keep_alive": false,
					"external":   false,
					"permission": "menu:view",
					"meta":       string(menuMeta),
					"created_at": time.Now(),
					"updated_at": time.Now(),
				},

				// Log Management - Parent Menu
				{
					"name":       "Log",
					"title":      "日志管理",
					"icon":       "Document",
					"path":       "/log",
					"component":  "Layout",
					"parent_id":  nil,
					"sort":       3,
					"type":       1,
					"visible":    1,
					"status":     1,
					"keep_alive": false,
					"external":   false,
					"permission": "",
					"meta":       string(logMeta),
					"created_at": time.Now(),
					"updated_at": time.Now(),
				},

				// Login Log
				{
					"name":       "LoginLog",
					"title":      "登录日志",
					"icon":       "Key",
					"path":       "login",
					"component":  "@/views/log/login/index.vue",
					"parent_id":  7, // Log parent ID will be 7
					"sort":       1,
					"type":       1,
					"visible":    1,
					"status":     1,
					"keep_alive": false,
					"external":   false,
					"permission": "log:view",
					"meta":       string(loginLogMeta),
					"created_at": time.Now(),
					"updated_at": time.Now(),
				},

				// Operation Log
				{
					"name":       "OperationLog",
					"title":      "操作日志",
					"icon":       "Document",
					"path":       "operation",
					"component":  "@/views/log/operation/index.vue",
					"parent_id":  7,
					"sort":       2,
					"type":       1,
					"visible":    1,
					"status":     1,
					"keep_alive": false,
					"external":   false,
					"permission": "log:view",
					"meta":       string(operationLogMeta),
					"created_at": time.Now(),
					"updated_at": time.Now(),
				},

				// Profile
				{
					"name":       "Profile",
					"title":      "个人中心",
					"icon":       "User",
					"path":       "/profile",
					"component":  "@/views/profile/index.vue",
					"parent_id":  nil,
					"sort":       4,
					"type":       1,
					"visible":    0, // Hidden
					"status":     1,
					"keep_alive": false,
					"external":   false,
					"permission": "profile:view",
					"meta":       string(profileMeta),
					"created_at": time.Now(),
					"updated_at": time.Now(),
				},
			}

			return tx.Table("menus").Create(menus).Error
		},
	})
}
