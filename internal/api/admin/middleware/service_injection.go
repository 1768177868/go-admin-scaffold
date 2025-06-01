package middleware

import (
	"app/internal/config"
	"app/internal/core/repositories"
	"app/internal/core/services"
	"app/pkg/database"

	"github.com/gin-gonic/gin"
)

// ServiceInjection injects all required services into the gin context
func ServiceInjection(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDB()

		// Initialize repositories
		userRepo := repositories.NewUserRepository(db)
		logRepo := repositories.NewLogRepository(db)
		roleRepo := repositories.NewRoleRepository(db)
		todoRepo := repositories.NewTodoRepository(db)

		// Initialize services
		logSvc := services.NewLogService(logRepo)
		authSvc := services.NewAuthService(userRepo, logSvc, cfg)
		userSvc := services.NewUserService(userRepo, logSvc)
		rbacSvc := services.NewRBACService(db)
		roleSvc := services.NewRoleService(roleRepo, db, logSvc)
		permissionSvc := services.NewPermissionService(db)
		todoService := services.NewTodoService(todoRepo)

		// Inject services into context
		c.Set("authService", authSvc)
		c.Set("userService", userSvc)
		c.Set("logService", logSvc)
		c.Set("rbacService", rbacSvc)
		c.Set("roleService", roleSvc)
		c.Set("permissionService", permissionSvc)
		c.Set("todoService", todoService)

		c.Next()
	}
}
