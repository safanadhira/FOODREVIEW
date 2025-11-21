package models

type Food struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"size:255"`
	Price       float64 `gorm:"type:decimal(8,2)"`
	Image       *string `gorm:"size:255"`
	Description *string `gorm:"type:text"`

	RestaurantID uint
	Restaurant   Restaurant `gorm:"constraint:OnDelete:CASCADE;"`
}
