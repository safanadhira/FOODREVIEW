package models

type Restaurant struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"size:255;unique" json:"class_name"`
	Location    *string `gorm:"size:255" json:"class_location"` // nullable
	Rating      *int    `json:"class_rating"`                   // nullable
	Image       *string `json:"image"`                          // nullable
	Description *string `gorm:"type:text" json:"description"`   // nullable large text
	Foods       []Food  `gorm:"foreignKey:RestaurantID"`
}
