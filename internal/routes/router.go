package routes

import (
	"app/internal/api/admin/handlers"
	"app/internal/api/admin/middleware"
	adminv1 "app/internal/api/admin/v1"
	openv1 "app/internal/api/open/v1"
	"app/internal/config"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	// Global middleware
	r.Use(middleware.I18n())                // Add i18n middleware globally
	r.Use(middleware.ServiceInjection(cfg)) // Add service injection middleware globally

	// Serve static files
	r.Static("/static", "./static")

	// Admin API routes (v1)
	adminV1 := r.Group("/api/admin/v1")
	{
		// Auth routes (no JWT protection needed)
		auth := adminV1.Group("/auth")
		{
			auth.POST("/login", adminv1.Login)
			auth.POST("/refresh", adminv1.RefreshToken)
		}

		// WebSocket routes (no JWT middleware needed, token passed via query params)
		wsHandler := handlers.NewWSHandler()
		adminV1.GET("/ws", wsHandler.HandleWebSocket)
		adminV1.POST("/ws/join", middleware.JWT(), wsHandler.JoinGroup)
		adminV1.POST("/ws/send", middleware.JWT(), wsHandler.SendMessage)
	}

	// Protected Admin API routes
	adminV1Protected := r.Group("/api/admin/v1")
	adminV1Protected.Use(middleware.JWT())          // Protect all admin routes with JWT auth
	adminV1Protected.Use(middleware.OperationLog()) // Add operation logging
	{
		// User routes
		users := adminV1Protected.Group("/users")
		{
			users.GET("", middleware.RBAC("user:view"), adminv1.ListUsers)
			users.POST("", middleware.RBAC("user:create"), adminv1.CreateUser)
			users.GET("/:id", middleware.RBAC("user:view"), adminv1.GetUser)
			users.PUT("/:id", middleware.RBAC("user:edit"), adminv1.UpdateUser)
			users.DELETE("/:id", middleware.RBAC("user:delete"), adminv1.DeleteUser)
			users.GET("/:id/logs", middleware.RBAC("log:view"), adminv1.GetUserLogs)
		}

		// Role routes
		roles := adminV1Protected.Group("/roles")
		{
			roles.GET("", middleware.RBAC("role:view"), adminv1.ListRoles)
			roles.POST("", middleware.RBAC("role:create"), adminv1.CreateRole)
			roles.GET("/:id", middleware.RBAC("role:view"), adminv1.GetRole)
			roles.PUT("/:id", middleware.RBAC("role:edit"), adminv1.UpdateRole)
			roles.DELETE("/:id", middleware.RBAC("role:delete"), adminv1.DeleteRole)
		}

		// Permission routes
		permissions := adminV1Protected.Group("/permissions")
		{
			permissions.GET("", middleware.RBAC("permission:view"), adminv1.ListPermissions)
			permissions.POST("", middleware.RBAC("permission:create"), adminv1.CreatePermission)
			permissions.GET("/:id", middleware.RBAC("permission:view"), adminv1.GetPermission)
			permissions.PUT("/:id", middleware.RBAC("permission:edit"), adminv1.UpdatePermission)
			permissions.DELETE("/:id", middleware.RBAC("permission:delete"), adminv1.DeletePermission)
			permissions.GET("/modules", middleware.RBAC("permission:view"), adminv1.GetPermissionsByModule)
		}

		// Log routes
		logs := adminV1Protected.Group("/logs")
		logs.Use(middleware.RBAC("log:view"))
		{
			logs.GET("/login", adminv1.ListLoginLogs)
			logs.GET("/operation", adminv1.ListOperationLogs)
		}

		// I18n routes
		i18n := adminV1Protected.Group("/i18n")
		{
			i18n.GET("/locales", adminv1.GetLocales)
			i18n.GET("/translations", adminv1.GetTranslations)
		}

		// Dashboard routes (accessible to all authenticated users)
		dashboard := adminV1Protected.Group("/dashboard")
		dashboard.Use(middleware.RBAC("dashboard:view"))
		{
			// Add dashboard endpoints here when needed
		}

		// Profile routes (accessible to all authenticated users)
		profile := adminV1Protected.Group("/profile")
		{
			profile.GET("", middleware.RBAC("profile:view"), adminv1.GetCurrentUser)
			profile.PUT("", middleware.RBAC("profile:edit"), adminv1.UpdateCurrentUser)
		}
	}

	// Open API routes (v1)
	openV1 := r.Group("/api/open/v1")
	{
		// Public routes
		public := openV1.Group("/public")
		{
			public.GET("/health", openv1.HealthCheck)
		}

		// OAuth routes
		oauth := openV1.Group("/oauth")
		{
			oauth.GET("/github", openv1.GithubOAuth)
			oauth.GET("/github/callback", openv1.GithubOAuthCallback)
		}
	}
}
