package handlers

import (
	"net/http"
	"strconv"

	"bharat-seva-space/database"
	"bharat-seva-space/models"

	"github.com/gin-gonic/gin"
)

// GetAllUsers returns all users (admin only)
func GetAllUsers(c *gin.Context) {
	var users []models.User
	
	// Get query parameters
	role := c.Query("role")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query := database.DB

	// Filter by role if provided
	if role != "" {
		query = query.Where("role = ?", role)
	}

	// Get total count
	var total int64
	query.Model(&models.User{}).Count(&total)

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Convert to response format
	var responses []models.UserResponse
	for _, user := range users {
		response := models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Phone:     user.Phone,
			Name:      user.Name,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		}
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"users": responses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetUser returns a specific user by ID (admin only)
func GetUser(c *gin.Context) {
	id := c.Param("id")
	
	var user models.User
	if err := database.DB.Preload("Applications").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userResponse := models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		Name:      user.Name,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	// Get user's applications
	var applications []models.ApplicationResponse
	for _, app := range user.Applications {
		appResponse := models.ApplicationResponse{
			ID:            app.ID,
			UserID:        app.UserID,
			ServiceType:   app.ServiceType,
			Status:        app.Status,
			Progress:      app.Progress,
			PaymentStatus: app.PaymentStatus,
			Amount:        app.Amount,
			Description:   app.Description,
			AssignedCA:    app.AssignedCA,
			Notes:         app.Notes,
			CreatedAt:     app.CreatedAt,
			UpdatedAt:     app.UpdatedAt,
		}
		applications = append(applications, appResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userResponse,
		"applications": applications,
	})
}

// UpdateUser updates a user (admin only)
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	
	var req struct {
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Role     string `json:"role"`
		IsActive bool   `json:"is_active"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update user
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	updates["is_active"] = req.IsActive

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Get updated user
	database.DB.First(&user, id)

	userResponse := models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		Name:      user.Name,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    userResponse,
	})
}

// GetUserStats returns user statistics (admin only)
func GetUserStats(c *gin.Context) {
	var stats struct {
		TotalUsers    int64 `json:"total_users"`
		ActiveUsers   int64 `json:"active_users"`
		InactiveUsers int64 `json:"inactive_users"`
		AdminUsers    int64 `json:"admin_users"`
		RegularUsers  int64 `json:"regular_users"`
	}

	database.DB.Model(&models.User{}).Count(&stats.TotalUsers)
	database.DB.Model(&models.User{}).Where("is_active = ?", true).Count(&stats.ActiveUsers)
	database.DB.Model(&models.User{}).Where("is_active = ?", false).Count(&stats.InactiveUsers)
	database.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&stats.AdminUsers)
	database.DB.Model(&models.User{}).Where("role = ?", "user").Count(&stats.RegularUsers)

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// GetDashboardData returns dashboard data for the current user
func GetDashboardData(c *gin.Context) {
	userInterface, _ := c.Get("user")
	currentUser := userInterface.(*models.User)

	// Get application stats
	var appStats struct {
		Total     int64 `json:"total"`
		Pending   int64 `json:"pending"`
		InProgress int64 `json:"in_progress"`
		Completed int64 `json:"completed"`
		Cancelled int64 `json:"cancelled"`
	}

	query := database.DB.Model(&models.Application{})
	
	// If user is not admin, only show their applications
	if currentUser.Role != "admin" {
		query = query.Where("user_id = ?", currentUser.ID)
	}

	query.Count(&appStats.Total)
	query.Where("status = ?", "pending").Count(&appStats.Pending)
	query.Where("status = ?", "in_progress").Count(&appStats.InProgress)
	query.Where("status = ?", "completed").Count(&appStats.Completed)
	query.Where("status = ?", "cancelled").Count(&appStats.Cancelled)

	// Get recent applications
	var recentApplications []models.Application
	recentQuery := database.DB.Preload("AssignedCAUser")
	if currentUser.Role != "admin" {
		recentQuery = recentQuery.Where("user_id = ?", currentUser.ID)
	}
	recentQuery.Order("created_at DESC").Limit(5).Find(&recentApplications)

	var recentResponses []models.ApplicationResponse
	for _, app := range recentApplications {
		response := models.ApplicationResponse{
			ID:            app.ID,
			UserID:        app.UserID,
			ServiceType:   app.ServiceType,
			Status:        app.Status,
			Progress:      app.Progress,
			PaymentStatus: app.PaymentStatus,
			Amount:        app.Amount,
			Description:   app.Description,
			AssignedCA:    app.AssignedCA,
			Notes:         app.Notes,
			CreatedAt:     app.CreatedAt,
			UpdatedAt:     app.UpdatedAt,
		}

		if app.AssignedCAUser != nil {
			response.AssignedCAUser = &models.UserResponse{
				ID:        app.AssignedCAUser.ID,
				Email:     app.AssignedCAUser.Email,
				Phone:     app.AssignedCAUser.Phone,
				Name:      app.AssignedCAUser.Name,
				Role:      app.AssignedCAUser.Role,
				IsActive:  app.AssignedCAUser.IsActive,
				CreatedAt: app.AssignedCAUser.CreatedAt,
			}
		}

		recentResponses = append(recentResponses, response)
	}

	// If admin, get additional stats
	var adminStats *gin.H
	if currentUser.Role == "admin" {
		var queryStats struct {
			Total      int64 `json:"total"`
			New        int64 `json:"new"`
			Contacted  int64 `json:"contacted"`
			Converted  int64 `json:"converted"`
			Closed     int64 `json:"closed"`
		}

		database.DB.Model(&models.Query{}).Count(&queryStats.Total)
		database.DB.Model(&models.Query{}).Where("status = ?", "new").Count(&queryStats.New)
		database.DB.Model(&models.Query{}).Where("status = ?", "contacted").Count(&queryStats.Contacted)
		database.DB.Model(&models.Query{}).Where("status = ?", "converted").Count(&queryStats.Converted)
		database.DB.Model(&models.Query{}).Where("status = ?", "closed").Count(&queryStats.Closed)

		adminStats = &gin.H{
			"queries": queryStats,
		}
	}

	response := gin.H{
		"user": currentUser,
		"application_stats": appStats,
		"recent_applications": recentResponses,
	}

	if adminStats != nil {
		response["admin_stats"] = adminStats
	}

	c.JSON(http.StatusOK, response)
} 