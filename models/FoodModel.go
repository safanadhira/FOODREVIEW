package models

// Asumsi kamu punya struct Review
type Food struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"size:255"`
	Price       string  `gorm:"type:decimal(8,2)"`
	Image       *string `gorm:"size:255"`
	Description string  `gorm:"type:text"`

	// --- RELASI PARENT (ke Restaurant) ---
	RestaurantID uint
	Restaurant   Restaurant `gorm:"constraint:OnDelete:CASCADE;"`

	Reviews []Review `gorm:"constraint:OnDelete:CASCADE;"`
	// ---------------------------------------------
}
