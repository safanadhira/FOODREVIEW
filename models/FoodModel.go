package models

import (
	"strings"

	"gorm.io/gorm"
)

// Asumsi kamu punya struct Review
type Food struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:255"`
	Price       string `gorm:"type:decimal(8,2)"`
	ImagePath   string
	Image       *string `gorm:"size:255" json:"image"`
	Description string  `gorm:"type:text"`

	// --- RELASI PARENT (ke Restaurant) ---
	RestaurantID uint
	Restaurant   Restaurant `gorm:"constraint:OnDelete:CASCADE;"`

	Reviews []Review `gorm:"constraint:OnDelete:CASCADE;"`
	// ---------------------------------------------
}

func (f *Food) GetCorrectedImagePath() string {
	// Mengganti semua backslash (\) dengan forward slash (/)
	return strings.ReplaceAll(f.ImagePath, "\\", "/")
}
