package models

import (
	"gorm.io/gorm"
)

// PackSize represents a pack size option
type PackSize struct {
	gorm.Model
	Size uint64 `gorm:"not null;uniqueIndex:idx_size_deleted_at"`
}

// SetupDatabase initializes the database with default pack sizes
func SetupDatabase(db *gorm.DB) error {
	// Auto migrate the schemas
	err := db.AutoMigrate(&PackSize{})
	if err != nil {
		return err
	}

	// Check if pack sizes already exist
	var count int64
	db.Model(&PackSize{}).Count(&count)
	if count == 0 {
		// Create default pack sizes
		defaultPackSizes := []PackSize{
			{Size: 250},
			{Size: 500},
			{Size: 1000},
			{Size: 2000},
			{Size: 5000},
		}

		// Insert default pack sizes
		for _, packSize := range defaultPackSizes {
			if err := db.Create(&packSize).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// GetPackSizes returns all available pack sizes in descending order
func GetPackSizes(db *gorm.DB) ([]uint64, error) {
	var packSizes []PackSize
	if err := db.Order("size DESC").Find(&packSizes).Error; err != nil {
		return nil, err
	}

	sizes := make([]uint64, len(packSizes))
	for i, pack := range packSizes {
		sizes[i] = pack.Size
	}

	return sizes, nil
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Error: message,
	}
}

// SuccessResponse represents an API success response
type SuccessResponse struct {
	Message string `json:"message"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string) SuccessResponse {
	return SuccessResponse{
		Message: message,
	}
}
