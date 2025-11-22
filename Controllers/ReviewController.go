package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/safanadhira/FOODREVIEW/initializers" 
	"github.com/safanadhira/FOODREVIEW/models"
	"gorm.io/gorm"
)

type ReviewInput struct {
	FoodID       uint    `json:"food_id" binding:"required"` 
	ReviewerName string  `json:"reviewer_name" binding:"required"` 
	Comment      *string `json:"comment"`
	Rating       int     `json:"rating" binding:"required,min=1,max=5"` 
}

// --- Controller Methods ---

func ReviewIndex(c *gin.Context) {
	foodIDStr := c.Param("foodId")
	
	// 1. Find the parent food 
	var food models.Food
	if err := initializers.DB.First(&food, foodIDStr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Food not found"})
		return
	}

	// 2. Fetch reviews 
	var reviews []models.Review
	if err := initializers.DB.
		Where("food_id = ?", food.ID).
		Order("created_at desc").
		Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"food_id": food.ID, "reviews": reviews})
}

func ReviewStore(c *gin.Context) {
	var input ReviewInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation Failed: " + err.Error()})
		return
	}

	var food models.Food
	if err := initializers.DB.First(&food, input.FoodID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Food item specified by food_id does not exist."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking food existence"})
		return
	}

	review := models.Review{
		FoodID:       input.FoodID,
		ReviewerName: input.ReviewerName,
		Comment:      input.Comment,
		Rating:       input.Rating,
	}

	initializers.DB.Create(&review)

	c.JSON(http.StatusCreated, gin.H{"message": "Review submitted!", "review": review})
}

func ReviewDestroy(c *gin.Context) {
	reviewID := c.Param("id")

	var review models.Review

	if err := initializers.DB.Preload("Food").First(&review, reviewID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	foodID := review.FoodID 

	result := initializers.DB.Delete(&review)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Review deleted!",
		"food_id": foodID,
	})
}