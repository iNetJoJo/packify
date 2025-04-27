package services

import (
	"packify/internal/models"
	"packify/pkg/calculator"

	"gorm.io/gorm"
)

// PackService handles pack calculation business logic
type PackService struct {
	DB *gorm.DB
}

// NewPackService creates a new pack service
func NewPackService(db *gorm.DB) *PackService {
	return &PackService{
		DB: db,
	}
}

// CalculatePacks calculates the optimal packs for an order
func (s *PackService) CalculatePacks(itemsOrdered uint64) (*calculator.PackResult, error) {
	// Get available pack sizes from the database
	packSizes, err := models.GetPackSizes(s.DB)
	if err != nil {
		return nil, err
	}

	// Calculate the optimal packs using the most efficient algorithm based on order size and pack sizes
	result, err := calculator.OptimalCalculatePacks(itemsOrdered, packSizes)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetPackSizes returns all available pack sizes
func (s *PackService) GetPackSizes() ([]models.PackSize, error) {
	var packSizes []models.PackSize
	if err := s.DB.Find(&packSizes).Error; err != nil {
		return nil, err
	}
	return packSizes, nil
}

// AddPackSize adds a new pack size
func (s *PackService) AddPackSize(size uint64) error {
	packSize := models.PackSize{
		Size: size,
	}
	return s.DB.Create(&packSize).Error
}

// UpdatePackSize updates a pack size availability
func (s *PackService) UpdatePackSize(id uint, isAvailable bool) error {
	return s.DB.Model(&models.PackSize{}).Where("id = ?", id).Update("is_available", isAvailable).Error
}

// DeletePackSize deletes a pack size
func (s *PackService) DeletePackSize(id uint) error {
	return s.DB.Unscoped().Delete(&models.PackSize{}, id).Error
}
