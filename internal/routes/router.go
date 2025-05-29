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
	r.Use(middleware.I18n()) // Add i18n middleware globally

	// Serve static files
	r.Static("/static", "./static")

	// Admin API routes (v1)
	adminV1 := r.Group("/api/admin/v1")
	adminV1.Use(middleware.JWT(cfg))       // Protect all admin routes with JWT auth
	adminV1.Use(middleware.OperationLog()) // Add operation logging
	{
		// Auth routes
		auth := adminV1.Group("/auth")
		{
			auth.POST("/login", adminv1.Login)
			auth.POST("/refresh", adminv1.RefreshToken)
		}

		// User routes
		users := adminV1.Group("/users")
		users.Use(middleware.RBAC()) // Add role-based access control
		{
			users.GET("", adminv1.ListUsers)
			users.POST("", adminv1.CreateUser)
			users.GET("/:id", adminv1.GetUser)
			users.PUT("/:id", adminv1.UpdateUser)
			users.DELETE("/:id", adminv1.DeleteUser)
			users.GET("/:id/logs", adminv1.GetUserLogs) // Add user logs endpoint
		}

		// Role routes
		roles := adminV1.Group("/roles")
		roles.Use(middleware.RBAC())
		{
			roles.GET("", adminv1.ListRoles)
			roles.POST("", adminv1.CreateRole)
			roles.GET("/:id", adminv1.GetRole)
			roles.PUT("/:id", adminv1.UpdateRole)
			roles.DELETE("/:id", adminv1.DeleteRole)
		}

		// WebSocket routes
		wsHandler := handlers.NewWSHandler()
		adminV1.GET("/ws", wsHandler.HandleWebSocket)
		adminV1.POST("/ws/join", wsHandler.JoinGroup)
		adminV1.POST("/ws/send", wsHandler.SendMessage)

		// Log routes
		logs := adminV1.Group("/logs")
		logs.Use(middleware.RBAC())
		{
			logs.GET("/login", adminv1.ListLoginLogs)
			logs.GET("/operation", adminv1.ListOperationLogs)
		}

		// I18n routes
		i18n := adminV1.Group("/i18n")
		{
			i18n.GET("/locales", adminv1.GetLocales)
			i18n.GET("/translations", adminv1.GetTranslations)
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
