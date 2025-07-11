package routes

import (
	"bharat-seva-space/handlers"
	"bharat-seva-space/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// API routes
	api := router.Group("/api")
	{
		// Public routes (no authentication required)
		public := api.Group("")
		{
			// Query submission
			public.POST("/queries", handlers.CreateQuery)
			
			// Authentication
			public.POST("/auth/register", handlers.Register)
			public.POST("/auth/login", handlers.Login)
		}

		// User routes (authentication required)
		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			// User profile
			user.GET("/profile", handlers.GetProfile)
			user.PUT("/profile", handlers.UpdateProfile)
			
			// Dashboard
			user.GET("/dashboard", handlers.GetDashboardData)
			
			// Applications
			user.GET("/applications", handlers.GetUserApplications)
			user.POST("/applications", handlers.CreateApplication)
			user.GET("/applications/:id", handlers.GetApplication)
			user.GET("/applications/stats", handlers.GetApplicationStats)
		}

		// Admin routes (admin authentication required)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.AdminMiddleware())
		{
			// Query management
			admin.GET("/queries", handlers.GetQueries)
			admin.GET("/queries/:id", handlers.GetQuery)
			admin.PUT("/queries/:id", handlers.UpdateQuery)
			admin.GET("/queries/stats", handlers.GetQueryStats)
			
			// User management
			admin.GET("/users", handlers.GetAllUsers)
			admin.GET("/users/:id", handlers.GetUser)
			admin.PUT("/users/:id", handlers.UpdateUser)
			admin.GET("/users/stats", handlers.GetUserStats)
			
			// Application management
			admin.GET("/applications", handlers.GetAllApplications)
			admin.PUT("/applications/:id", handlers.UpdateApplication)
			admin.GET("/applications/stats", handlers.GetApplicationStats)
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Bharat Seva Space API is running",
		})
	})

	return router
} 