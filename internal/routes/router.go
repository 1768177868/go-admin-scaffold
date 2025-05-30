package routes

import (
	"app/internal/api/admin/handlers"
	"app/internal/api/admin/middleware"
	adminv1 "app/internal/api/admin/v1"
	openv1 "app/internal/api/open/v1"
	"app/internal/config"
	"app/pkg/response"

	"github.com/gin-gonic/gin"
)

// responseWriter wraps gin.ResponseWriter to track if response was written
type responseWriter struct {
	gin.ResponseWriter
	written bool
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(data)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.written = true
	w.ResponseWriter.WriteHeader(statusCode)
}

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
			auth.POST("/login", wrapHandler(adminv1.Login))
			auth.POST("/refresh", wrapHandler(adminv1.RefreshToken))
		}

		// WebSocket routes (no JWT middleware needed, token passed via query params)
		wsHandler := handlers.NewWSHandler()
		adminV1.GET("/ws", wrapHandler(wsHandler.HandleWebSocket))
		adminV1.POST("/ws/join", middleware.JWT(), wrapHandler(wsHandler.JoinGroup))
		adminV1.POST("/ws/leave", middleware.JWT(), wrapHandler(wsHandler.LeaveGroup))
		adminV1.POST("/ws/send", middleware.JWT(), wrapHandler(wsHandler.SendMessage))

		// SSE routes
		sseHandler := handlers.NewSSEHandler()
		adminV1.GET("/sse", wrapHandler(sseHandler.HandleSSE))
		adminV1.POST("/sse/notify", middleware.JWT(), wrapHandler(sseHandler.SendNotification))
		adminV1.POST("/sse/join", middleware.JWT(), wrapHandler(sseHandler.JoinGroup))
		adminV1.POST("/sse/leave", middleware.JWT(), wrapHandler(sseHandler.LeaveGroup))
	}

	// Protected Admin API routes
	adminV1Protected := r.Group("/api/admin/v1")
	adminV1Protected.Use(middleware.JWT())          // Protect all admin routes with JWT auth
	adminV1Protected.Use(middleware.OperationLog()) // Add operation logging
	{
		// User routes
		users := adminV1Protected.Group("/users")
		{
			users.GET("", middleware.RBAC("user:view"), wrapHandler(adminv1.ListUsers))
			users.POST("", middleware.RBAC("user:create"), wrapHandler(adminv1.CreateUser))
			users.GET("/:id", middleware.RBAC("user:view"), wrapHandler(adminv1.GetUser))
			users.PUT("/:id", middleware.RBAC("user:edit"), wrapHandler(adminv1.UpdateUser))
			users.DELETE("/:id", middleware.RBAC("user:delete"), wrapHandler(adminv1.DeleteUser))
			users.GET("/:id/logs", middleware.RBAC("log:view"), wrapHandler(adminv1.GetUserLogs))
		}

		// Role routes
		roles := adminV1Protected.Group("/roles")
		{
			roles.GET("", middleware.RBAC("role:view"), wrapHandler(adminv1.ListRoles))
			roles.POST("", middleware.RBAC("role:create"), wrapHandler(adminv1.CreateRole))
			roles.GET("/:id", middleware.RBAC("role:view"), wrapHandler(adminv1.GetRole))
			roles.PUT("/:id", middleware.RBAC("role:edit"), wrapHandler(adminv1.UpdateRole))
			roles.DELETE("/:id", middleware.RBAC("role:delete"), wrapHandler(adminv1.DeleteRole))
		}

		// Permission routes
		permissions := adminV1Protected.Group("/permissions")
		{
			permissions.GET("", middleware.RBAC("permission:view"), wrapHandler(adminv1.ListPermissions))
			permissions.POST("", middleware.RBAC("permission:create"), wrapHandler(adminv1.CreatePermission))
			permissions.GET("/:id", middleware.RBAC("permission:view"), wrapHandler(adminv1.GetPermission))
			permissions.PUT("/:id", middleware.RBAC("permission:edit"), wrapHandler(adminv1.UpdatePermission))
			permissions.DELETE("/:id", middleware.RBAC("permission:delete"), wrapHandler(adminv1.DeletePermission))
			permissions.GET("/modules", middleware.RBAC("permission:view"), wrapHandler(adminv1.GetPermissionsByModule))
		}

		// Log routes
		logs := adminV1Protected.Group("/logs")
		logs.Use(middleware.RBAC("log:view"))
		{
			logs.GET("/login", wrapHandler(adminv1.ListLoginLogs))
			logs.GET("/operation", wrapHandler(adminv1.ListOperationLogs))
		}

		// I18n routes
		i18n := adminV1Protected.Group("/i18n")
		{
			i18n.GET("/locales", wrapHandler(adminv1.GetLocales))
			i18n.GET("/translations", wrapHandler(adminv1.GetTranslations))
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
			profile.GET("", middleware.RBAC("profile:view"), wrapHandler(adminv1.GetCurrentUser))
			profile.PUT("", middleware.RBAC("profile:edit"), wrapHandler(adminv1.UpdateCurrentUser))
		}
	}

	// Open API routes (v1)
	openV1 := r.Group("/api/open/v1")
	{
		// Public routes
		public := openV1.Group("/public")
		{
			public.GET("/health", wrapHandler(openv1.HealthCheck))
		}

		// OAuth routes
		oauth := openV1.Group("/oauth")
		{
			oauth.GET("/github", wrapHandler(openv1.GithubOAuth))
			oauth.GET("/github/callback", wrapHandler(openv1.GithubOAuthCallback))
		}
	}
}

// wrapHandler wraps a gin.HandlerFunc to ensure consistent response handling
func wrapHandler(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Store the original writer
		originalWriter := c.Writer

		// Create a custom response writer that captures the response
		responseWriter := &responseWriter{ResponseWriter: originalWriter}
		c.Writer = responseWriter

		// Call the original handler
		handler(c)

		// If no response has been written and it's not a redirect
		if !responseWriter.written && c.Writer.Status() < 300 {
			// Ensure response includes trace_id
			response.Success(c, nil)
		}
	}
}
