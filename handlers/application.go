package handlers

import (
	"net/http"
	"strconv"

	"bharat-seva-space/database"
	"bharat-seva-space/models"

	"github.com/gin-gonic/gin"
)

// CreateApplication handles application creation by users
func CreateApplication(c *gin.Context) {
	userInterface, _ := c.Get("user")
	currentUser := userInterface.(*models.User)

	var req models.ApplicationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create application
	application := models.Application{
		UserID:        currentUser.ID,
		ServiceType:   req.ServiceType,
		Description:   req.Description,
		Amount:        req.Amount,
		Status:        "pending",
		Progress:      "0%",
		PaymentStatus: "pending",
	}

	if err := database.DB.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Application created successfully",
		"application_id": application.ID,
	})
}

// GetUserApplications returns applications for the current user
func GetUserApplications(c *gin.Context) {
	userInterface, _ := c.Get("user")
	currentUser := userInterface.(*models.User)

	var applications []models.Application
	
	// Get query parameters
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query := database.DB.Where("user_id = ?", currentUser.ID).Preload("AssignedCAUser").Preload("Documents")

	// Filter by status if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	var total int64
	query.Model(&models.Application{}).Count(&total)

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})
		return
	}

	// Convert to response format
	var responses []models.ApplicationResponse
	for _, app := range applications {
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

		// Add documents
		for _, doc := range app.Documents {
			response.Documents = append(response.Documents, models.DocumentResponse{
				ID:            doc.ID,
				ApplicationID: doc.ApplicationID,
				FileName:      doc.FileName,
				FilePath:      doc.FilePath,
				FileSize:      doc.FileSize,
				FileType:      doc.FileType,
				Description:   doc.Description,
				UploadedAt:    doc.UploadedAt,
				CreatedAt:     doc.CreatedAt,
			})
		}

		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"applications": responses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetApplication returns a specific application by ID
func GetApplication(c *gin.Context) {
	id := c.Param("id")
	userInterface, _ := c.Get("user")
	currentUser := userInterface.(*models.User)

	var application models.Application
	query := database.DB.Preload("AssignedCAUser").Preload("Documents").Preload("User")

	// If user is not admin, only allow access to their own applications
	if currentUser.Role != "admin" {
		query = query.Where("user_id = ?", currentUser.ID)
	}

	if err := query.First(&application, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	response := models.ApplicationResponse{
		ID:            application.ID,
		UserID:        application.UserID,
		ServiceType:   application.ServiceType,
		Status:        application.Status,
		Progress:      application.Progress,
		PaymentStatus: application.PaymentStatus,
		Amount:        application.Amount,
		Description:   application.Description,
		AssignedCA:    application.AssignedCA,
		Notes:         application.Notes,
		CreatedAt:     application.CreatedAt,
		UpdatedAt:     application.UpdatedAt,
	}

	// Add user info
	response.User = models.UserResponse{
		ID:        application.User.ID,
		Email:     application.User.Email,
		Phone:     application.User.Phone,
		Name:      application.User.Name,
		Role:      application.User.Role,
		IsActive:  application.User.IsActive,
		CreatedAt: application.User.CreatedAt,
	}

	if application.AssignedCAUser != nil {
		response.AssignedCAUser = &models.UserResponse{
			ID:        application.AssignedCAUser.ID,
			Email:     application.AssignedCAUser.Email,
			Phone:     application.AssignedCAUser.Phone,
			Name:      application.AssignedCAUser.Name,
			Role:      application.AssignedCAUser.Role,
			IsActive:  application.AssignedCAUser.IsActive,
			CreatedAt: application.AssignedCAUser.CreatedAt,
		}
	}

	// Add documents
	for _, doc := range application.Documents {
		response.Documents = append(response.Documents, models.DocumentResponse{
			ID:            doc.ID,
			ApplicationID: doc.ApplicationID,
			FileName:      doc.FileName,
			FilePath:      doc.FilePath,
			FileSize:      doc.FileSize,
			FileType:      doc.FileType,
			Description:   doc.Description,
			UploadedAt:    doc.UploadedAt,
			CreatedAt:     doc.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"application": response})
}

// GetAllApplications returns all applications (admin only)
func GetAllApplications(c *gin.Context) {
	var applications []models.Application
	
	// Get query parameters
	status := c.Query("status")
	serviceType := c.Query("service_type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query := database.DB.Preload("User").Preload("AssignedCAUser").Preload("Documents")

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if serviceType != "" {
		query = query.Where("service_type = ?", serviceType)
	}

	// Get total count
	var total int64
	query.Model(&models.Application{}).Count(&total)

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})
		return
	}

	// Convert to response format
	var responses []models.ApplicationResponse
	for _, app := range applications {
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

		// Add user info
		response.User = models.UserResponse{
			ID:        app.User.ID,
			Email:     app.User.Email,
			Phone:     app.User.Phone,
			Name:      app.User.Name,
			Role:      app.User.Role,
			IsActive:  app.User.IsActive,
			CreatedAt: app.User.CreatedAt,
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

		// Add documents
		for _, doc := range app.Documents {
			response.Documents = append(response.Documents, models.DocumentResponse{
				ID:            doc.ID,
				ApplicationID: doc.ApplicationID,
				FileName:      doc.FileName,
				FilePath:      doc.FilePath,
				FileSize:      doc.FileSize,
				FileType:      doc.FileType,
				Description:   doc.Description,
				UploadedAt:    doc.UploadedAt,
				CreatedAt:     doc.CreatedAt,
			})
		}

		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"applications": responses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// UpdateApplication updates an application (admin only)
func UpdateApplication(c *gin.Context) {
	id := c.Param("id")
	
	var req models.ApplicationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if application exists
	var application models.Application
	if err := database.DB.First(&application, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// Update application
	updates := make(map[string]interface{})
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Progress != "" {
		updates["progress"] = req.Progress
	}
	if req.PaymentStatus != "" {
		updates["payment_status"] = req.PaymentStatus
	}
	if req.AssignedCA != nil {
		updates["assigned_ca"] = req.AssignedCA
	}
	if req.Notes != "" {
		updates["notes"] = req.Notes
	}
	if req.Amount > 0 {
		updates["amount"] = req.Amount
	}

	if err := database.DB.Model(&application).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application"})
		return
	}

	// Get updated application with relationships
	database.DB.Preload("User").Preload("AssignedCAUser").Preload("Documents").First(&application, id)

	response := models.ApplicationResponse{
		ID:            application.ID,
		UserID:        application.UserID,
		ServiceType:   application.ServiceType,
		Status:        application.Status,
		Progress:      application.Progress,
		PaymentStatus: application.PaymentStatus,
		Amount:        application.Amount,
		Description:   application.Description,
		AssignedCA:    application.AssignedCA,
		Notes:         application.Notes,
		CreatedAt:     application.CreatedAt,
		UpdatedAt:     application.UpdatedAt,
	}

	// Add user info
	response.User = models.UserResponse{
		ID:        application.User.ID,
		Email:     application.User.Email,
		Phone:     application.User.Phone,
		Name:      application.User.Name,
		Role:      application.User.Role,
		IsActive:  application.User.IsActive,
		CreatedAt: application.User.CreatedAt,
	}

	if application.AssignedCAUser != nil {
		response.AssignedCAUser = &models.UserResponse{
			ID:        application.AssignedCAUser.ID,
			Email:     application.AssignedCAUser.Email,
			Phone:     application.AssignedCAUser.Phone,
			Name:      application.AssignedCAUser.Name,
			Role:      application.AssignedCAUser.Role,
			IsActive:  application.AssignedCAUser.IsActive,
			CreatedAt: application.AssignedCAUser.CreatedAt,
		}
	}

	// Add documents
	for _, doc := range application.Documents {
		response.Documents = append(response.Documents, models.DocumentResponse{
			ID:            doc.ID,
			ApplicationID: doc.ApplicationID,
			FileName:      doc.FileName,
			FilePath:      doc.FilePath,
			FileSize:      doc.FileSize,
			FileType:      doc.FileType,
			Description:   doc.Description,
			UploadedAt:    doc.UploadedAt,
			CreatedAt:     doc.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Application updated successfully",
		"application": response,
	})
}

// GetApplicationStats returns application statistics
func GetApplicationStats(c *gin.Context) {
	userInterface, _ := c.Get("user")
	currentUser := userInterface.(*models.User)

	var stats struct {
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

	query.Count(&stats.Total)
	query.Where("status = ?", "pending").Count(&stats.Pending)
	query.Where("status = ?", "in_progress").Count(&stats.InProgress)
	query.Where("status = ?", "completed").Count(&stats.Completed)
	query.Where("status = ?", "cancelled").Count(&stats.Cancelled)

	c.JSON(http.StatusOK, gin.H{"stats": stats})
} 