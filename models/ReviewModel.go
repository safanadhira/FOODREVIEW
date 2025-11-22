package models

import "time"

type Review struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FoodID       uint      `json:"food_id"`
	Food         Food      `json:"food"`
	ReviewerName string   `gorm:"size:255" json:"reviewer_name"`
	Comment      *string    `gorm:"type:text" json:"comment"`
	Rating       int     `json:"rating"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
