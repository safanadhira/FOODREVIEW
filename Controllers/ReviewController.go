package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/safanadhira/FOODREVIEW/initializers"
	"github.com/safanadhira/FOODREVIEW/models"
	"gorm.io/gorm"
)

type ReviewInput struct {
	FoodID       uint    `form:"food_id" binding:"required"`
	ReviewerName string  `form:"reviewer_name" binding:"required"`
	Comment      *string `form:"comment"`
	Rating       int     `form:"rating" binding:"required,min=1,max=5"`
}

func ReviewIndex(c *gin.Context) {

	foodIDStr := c.Param("id")

	var food models.Food

	if err := initializers.DB.Preload("Restaurant").First(&food, foodIDStr).Error; err != nil {
		c.String(http.StatusNotFound, "Food not found")
		return
	}

	var reviews []models.Review
	if err := initializers.DB.
		Where("food_id = ?", food.ID).
		Order("created_at desc").
		Find(&reviews).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to fetch reviews")
		return
	}

	c.HTML(http.StatusOK, "reviews/index.html", gin.H{
		"Food":    food,
		"Reviews": reviews,
	})
}

func ReviewCreate(c *gin.Context) {

	foodID := c.Query("food_id")

	var food models.Food

	if err := initializers.DB.Preload("Restaurant").First(&food, foodID).Error; err != nil {
		c.String(http.StatusNotFound, "Food ID not found")
		return
	}

	c.HTML(http.StatusOK, "reviews/create.html", gin.H{
		"Food": food,
	})
}

func ReviewStore(c *gin.Context) {
	var input ReviewInput
	if err := c.ShouldBind(&input); err != nil {
		c.String(http.StatusBadRequest, "Validation Failed: "+err.Error())
		return
	}

	var food models.Food
	if err := initializers.DB.First(&food, input.FoodID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.String(http.StatusBadRequest, "Food item specified by food_id does not exist.")
			return
		}
		c.String(http.StatusInternalServerError, "Database error checking food existence")
		return
	}

	review := models.Review{
		FoodID:       input.FoodID,
		ReviewerName: input.ReviewerName,
		Comment:      input.Comment,
		Rating:       input.Rating,
	}

	initializers.DB.Create(&review)

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/reviews/%d", input.FoodID))
}

func ReviewDestroy(c *gin.Context) {
	reviewID := c.Param("id")

	var review models.Review
	if err := initializers.DB.Preload("Food").First(&review, reviewID).Error; err != nil {
		c.String(http.StatusNotFound, "Review not found")
		return
	}

	foodID := review.FoodID

	if err := initializers.DB.Delete(&review).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to delete review")
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/reviews/%d", foodID))
}
