package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a registered user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Phone     string         `json:"phone" gorm:"uniqueIndex;not null"`
	Name      string         `json:"name" gorm:"not null"`
	Password  string         `json:"-" gorm:"not null"` // "-" means don't include in JSON
	Role      string         `json:"role" gorm:"default:'user'"` // user, admin
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Applications []Application `json:"applications,omitempty" gorm:"foreignKey:UserID"`
}

// UserLoginRequest represents login request
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserRegisterRequest represents registration request
type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserUpdateRequest represents user profile update request
type UserUpdateRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// UserResponse represents user data in API responses
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
} 