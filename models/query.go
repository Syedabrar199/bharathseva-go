package models

import (
	"time"

	"gorm.io/gorm"
)

// Query represents a public query from the website
type Query struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Email       string         `json:"email" gorm:"not null"`
	Phone       string         `json:"phone" gorm:"not null"`
	Service     string         `json:"service" gorm:"not null"`
	Message     string         `json:"message" gorm:"type:text"`
	Status      string         `json:"status" gorm:"default:'new'"` // new, contacted, converted, closed
	AssignedTo  *uint          `json:"assigned_to"` // Admin user ID
	Notes       string         `json:"notes" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	AssignedUser *User `json:"assigned_user,omitempty" gorm:"foreignKey:AssignedTo"`
}

// QueryCreateRequest represents query creation request
type QueryCreateRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Phone   string `json:"phone" binding:"required"`
	Service string `json:"service" binding:"required"`
	Message string `json:"message" binding:"required"`
}

// QueryUpdateRequest represents query update request (admin only)
type QueryUpdateRequest struct {
	Status     string `json:"status"`
	AssignedTo *uint  `json:"assigned_to"`
	Notes      string `json:"notes"`
}

// QueryResponse represents query data in API responses
type QueryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Service     string    `json:"service"`
	Message     string    `json:"message"`
	Status      string    `json:"status"`
	AssignedTo  *uint     `json:"assigned_to"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AssignedUser *UserResponse `json:"assigned_user,omitempty"`
} 