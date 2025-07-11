package handlers

import (
	"net/http"
	"strconv"

	"bharat-seva-space/database"
	"bharat-seva-space/models"

	"github.com/gin-gonic/gin"
)

// CreateQuery handles public query submission
func CreateQuery(c *gin.Context) {
	var req models.QueryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create query
	query := models.Query{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Service: req.Service,
		Message: req.Message,
		Status:  "new",
	}

	if err := database.DB.Create(&query).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create query"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Query submitted successfully",
		"query_id": query.ID,
	})
}

// GetQueries returns all queries (admin only)
func GetQueries(c *gin.Context) {
	var queries []models.Query
	
	// Get query parameters
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query := database.DB.Preload("AssignedUser")

	// Filter by status if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	var total int64
	query.Model(&models.Query{}).Count(&total)

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&queries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch queries"})
		return
	}

	// Convert to response format
	var responses []models.QueryResponse
	for _, query := range queries {
		response := models.QueryResponse{
			ID:        query.ID,
			Name:      query.Name,
			Email:     query.Email,
			Phone:     query.Phone,
			Service:   query.Service,
			Message:   query.Message,
			Status:    query.Status,
			AssignedTo: query.AssignedTo,
			Notes:     query.Notes,
			CreatedAt: query.CreatedAt,
			UpdatedAt: query.UpdatedAt,
		}

		if query.AssignedUser != nil {
			response.AssignedUser = &models.UserResponse{
				ID:        query.AssignedUser.ID,
				Email:     query.AssignedUser.Email,
				Phone:     query.AssignedUser.Phone,
				Name:      query.AssignedUser.Name,
				Role:      query.AssignedUser.Role,
				IsActive:  query.AssignedUser.IsActive,
				CreatedAt: query.AssignedUser.CreatedAt,
			}
		}

		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"queries": responses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetQuery returns a specific query by ID (admin only)
func GetQuery(c *gin.Context) {
	id := c.Param("id")
	
	var query models.Query
	if err := database.DB.Preload("AssignedUser").First(&query, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Query not found"})
		return
	}

	response := models.QueryResponse{
		ID:        query.ID,
		Name:      query.Name,
		Email:     query.Email,
		Phone:     query.Phone,
		Service:   query.Service,
		Message:   query.Message,
		Status:    query.Status,
		AssignedTo: query.AssignedTo,
		Notes:     query.Notes,
		CreatedAt: query.CreatedAt,
		UpdatedAt: query.UpdatedAt,
	}

	if query.AssignedUser != nil {
		response.AssignedUser = &models.UserResponse{
			ID:        query.AssignedUser.ID,
			Email:     query.AssignedUser.Email,
			Phone:     query.AssignedUser.Phone,
			Name:      query.AssignedUser.Name,
			Role:      query.AssignedUser.Role,
			IsActive:  query.AssignedUser.IsActive,
			CreatedAt: query.AssignedUser.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"query": response})
}

// UpdateQuery updates a query (admin only)
func UpdateQuery(c *gin.Context) {
	id := c.Param("id")
	
	var req models.QueryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if query exists
	var query models.Query
	if err := database.DB.First(&query, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Query not found"})
		return
	}

	// Update query
	updates := make(map[string]interface{})
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.AssignedTo != nil {
		updates["assigned_to"] = req.AssignedTo
	}
	if req.Notes != "" {
		updates["notes"] = req.Notes
	}

	if err := database.DB.Model(&query).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update query"})
		return
	}

	// Get updated query with relationships
	database.DB.Preload("AssignedUser").First(&query, id)

	response := models.QueryResponse{
		ID:        query.ID,
		Name:      query.Name,
		Email:     query.Email,
		Phone:     query.Phone,
		Service:   query.Service,
		Message:   query.Message,
		Status:    query.Status,
		AssignedTo: query.AssignedTo,
		Notes:     query.Notes,
		CreatedAt: query.CreatedAt,
		UpdatedAt: query.UpdatedAt,
	}

	if query.AssignedUser != nil {
		response.AssignedUser = &models.UserResponse{
			ID:        query.AssignedUser.ID,
			Email:     query.AssignedUser.Email,
			Phone:     query.AssignedUser.Phone,
			Name:      query.AssignedUser.Name,
			Role:      query.AssignedUser.Role,
			IsActive:  query.AssignedUser.IsActive,
			CreatedAt: query.AssignedUser.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Query updated successfully",
		"query":   response,
	})
}

// GetQueryStats returns query statistics (admin only)
func GetQueryStats(c *gin.Context) {
	var stats struct {
		Total      int64 `json:"total"`
		New        int64 `json:"new"`
		Contacted  int64 `json:"contacted"`
		Converted  int64 `json:"converted"`
		Closed     int64 `json:"closed"`
	}

	database.DB.Model(&models.Query{}).Count(&stats.Total)
	database.DB.Model(&models.Query{}).Where("status = ?", "new").Count(&stats.New)
	database.DB.Model(&models.Query{}).Where("status = ?", "contacted").Count(&stats.Contacted)
	database.DB.Model(&models.Query{}).Where("status = ?", "converted").Count(&stats.Converted)
	database.DB.Model(&models.Query{}).Where("status = ?", "closed").Count(&stats.Closed)

	c.JSON(http.StatusOK, gin.H{"stats": stats})
} 