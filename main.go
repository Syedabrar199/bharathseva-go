package main

import (
	"log"
	"os"

	"bharat-seva-space/config"
	"bharat-seva-space/database"
	"bharat-seva-space/routes"
)

func main() {
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatal("Error loading environment variables:", err)
	}

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Error initializing database:", err)
	}

	// Create default admin user
	if err := database.CreateAdminUser(); err != nil {
		log.Fatal("Error creating admin user:", err)
	}

	// Setup routes
	router := routes.SetupRoutes()

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error starting server:", err)
	}
} 