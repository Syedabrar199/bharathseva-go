package main

// import (
//     "fmt"
//     "log"

//     "bharat-seva-space/config"
//     "bharat-seva-space/database"
//     "bharat-seva-space/models"
//     "bharat-seva-space/utils"
// )

// func main() {
//     // Load environment variables
//     if err := config.LoadEnv(); err != nil {
//         log.Fatal("Failed to load environment variables:", err)
//     }

//     // Initialize database
//     if err := database.InitDB(); err != nil {
//         log.Fatal("Failed to initialize database:", err)
//     }

//     // Admin user details
//     adminEmail := "admin@bharatseva.com"
//     adminPassword := "Admin@2024"
//     adminName := "Admin User"
//     adminPhone := "9999999999"

//     // Hard delete any user with this email (including soft-deleted)
//     database.DB.Unscoped().Where("email = ?", adminEmail).Delete(&models.User{})

//     // Hash password
//     hashedPassword, err := utils.HashPassword(adminPassword)
//     if err != nil {
//         log.Fatal("Failed to hash password:", err)
//     }

//     // Create admin user
//     adminUser := models.User{
//         Email:    adminEmail,
//         Phone:    adminPhone,
//         Name:     adminName,
//         Password: hashedPassword,
//         Role:     "admin",
//         IsActive: true,
//     }

//     if err := database.DB.Create(&adminUser).Error; err != nil {
//         log.Fatal("Failed to create admin user:", err)
//     }

//     fmt.Printf("Admin user created successfully!\n")
//     fmt.Printf("Email: %s\n", adminEmail)
//     fmt.Printf("Password: %s\n", adminPassword)
//     fmt.Printf("Role: %s\n", adminUser.Role)
//     fmt.Printf("ID: %d\n", adminUser.ID)
//     fmt.Printf("\nYou can now log in with these credentials to access admin endpoints.\n")
// } 