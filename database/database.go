package database

import (
	"fmt"
	"log"

	"bharat-seva-space/config"
	"bharat-seva-space/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection and runs migrations
func InitDB() error {
	dbConfig := config.GetDBConfig()
	
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig["host"],
		dbConfig["port"],
		dbConfig["user"],
		dbConfig["password"],
		dbConfig["dbname"],
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Println("Database connected and migrations completed successfully")
	return nil
}

// runMigrations runs database migrations
func runMigrations() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Query{},
		&models.Application{},
		&models.Document{},
	)
}

// CreateAdminUser creates a default admin user if it doesn't exist
func CreateAdminUser() error {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	
	if count == 0 {
		// Create default admin user
		adminUser := models.User{
			Email:    "admin@bharatseva.com",
			Phone:    "9999999999",
			Name:     "Admin User",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     "admin",
			IsActive: true,
		}
		
		if err := DB.Create(&adminUser).Error; err != nil {
			return fmt.Errorf("failed to create admin user: %v", err)
		}
		
		log.Println("Default admin user created: admin@bharatseva.com / password")
	}
	
	return nil
} 