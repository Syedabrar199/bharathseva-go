package models

import (
	"time"

	"gorm.io/gorm"
)

// Document represents uploaded files for applications
type Document struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	ApplicationID uint           `json:"application_id" gorm:"not null"`
	FileName      string         `json:"file_name" gorm:"not null"`
	FilePath      string         `json:"file_path" gorm:"not null"`
	FileSize      int64          `json:"file_size"`
	FileType      string         `json:"file_type"`
	Description   string         `json:"description"`
	UploadedAt    time.Time      `json:"uploaded_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Application Application `json:"application,omitempty" gorm:"foreignKey:ApplicationID"`
}

// DocumentCreateRequest represents document creation request
type DocumentCreateRequest struct {
	Description string `json:"description"`
}

// DocumentResponse represents document data in API responses
type DocumentResponse struct {
	ID            uint      `json:"id"`
	ApplicationID uint      `json:"application_id"`
	FileName      string    `json:"file_name"`
	FilePath      string    `json:"file_path"`
	FileSize      int64     `json:"file_size"`
	FileType      string    `json:"file_type"`
	Description   string    `json:"description"`
	UploadedAt    time.Time `json:"uploaded_at"`
	CreatedAt     time.Time `json:"created_at"`
} 