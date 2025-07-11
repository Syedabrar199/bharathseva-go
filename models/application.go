package models

import (
	"time"

	"gorm.io/gorm"
)

// Application represents a service application from registered users
type Application struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"user_id" gorm:"not null"`
	ServiceType   string         `json:"service_type" gorm:"not null"`
	Status        string         `json:"status" gorm:"default:'pending'"` // pending, in_progress, completed, cancelled
	Progress      string         `json:"progress" gorm:"default:'0%'"`
	PaymentStatus string         `json:"payment_status" gorm:"default:'pending'"` // pending, paid, refunded
	Amount        float64        `json:"amount" gorm:"default:0"`
	Description   string         `json:"description" gorm:"type:text"`
	AssignedCA    *uint          `json:"assigned_ca"` // Admin user ID assigned as CA
	Notes         string         `json:"notes" gorm:"type:text"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User        User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	AssignedCAUser *User    `json:"assigned_ca_user,omitempty" gorm:"foreignKey:AssignedCA"`
	Documents   []Document  `json:"documents,omitempty" gorm:"foreignKey:ApplicationID"`
}

// ApplicationCreateRequest represents application creation request
type ApplicationCreateRequest struct {
	ServiceType string  `json:"service_type" binding:"required"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

// ApplicationUpdateRequest represents application update request (admin only)
type ApplicationUpdateRequest struct {
	Status        string  `json:"status"`
	Progress      string  `json:"progress"`
	PaymentStatus string  `json:"payment_status"`
	AssignedCA    *uint   `json:"assigned_ca"`
	Notes         string  `json:"notes"`
	Amount        float64 `json:"amount"`
}

// ApplicationResponse represents application data in API responses
type ApplicationResponse struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	ServiceType   string    `json:"service_type"`
	Status        string    `json:"status"`
	Progress      string    `json:"progress"`
	PaymentStatus string    `json:"payment_status"`
	Amount        float64   `json:"amount"`
	Description   string    `json:"description"`
	AssignedCA    *uint     `json:"assigned_ca"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	User          UserResponse `json:"user,omitempty"`
	AssignedCAUser *UserResponse `json:"assigned_ca_user,omitempty"`
	Documents     []DocumentResponse `json:"documents,omitempty"`
} 