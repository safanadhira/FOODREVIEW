package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/safanadhira/FOODREVIEW/initializers"
	"github.com/safanadhira/FOODREVIEW/models"
	"gorm.io/gorm"
)

type RestaurantInput struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

// --- Controller Methods ---

func RestaurantIndex(c *gin.Context) {
	var restaurants []models.Restaurant
	
	initializers.DB.Find(&restaurants) 

	c.JSON(http.StatusOK, restaurants)
}

func RestaurantStore(c *gin.Context) {
	var input RestaurantInput
	
	// Bind input and check basic validation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data or Missing Fields: " + err.Error()})
		return
	}

	var existingRestaurant models.Restaurant
	result := initializers.DB.Where("name = ?", input.Name).First(&existingRestaurant)
	
	if result.Error == nil { // If name is not unique
		c.JSON(http.StatusConflict, gin.H{"error": "Restaurant name already exists."})
		return
	}
	if result.Error != gorm.ErrRecordNotFound { // Other database error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error while checking uniqueness."})
		return
	}

	restaurant := models.Restaurant{
		Name:        input.Name,
		Description: input.Description,
	}

	initializers.DB.Create(&restaurant)

	c.JSON(http.StatusCreated, gin.H{"message": "Restaurant added!", "restaurant": restaurant})
}

func RestaurantShow(c *gin.Context) {
	name := c.Param("name") 

	var restaurant models.Restaurant
	
	result := initializers.DB.Preload("Foods").Where("name = ?", name).First(&restaurant)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

func RestaurantUpdate(c *gin.Context) {
	idStr := c.Param("id")

	var food models.Restaurant
	
	// 1. Find the restaurant 
	if err := initializers.DB.First(&food, idStr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant Not Found"})
		return
	}

	var input RestaurantInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}
	
	// 2. Update the fields
	initializers.DB.Model(&food).Updates(models.Restaurant{
		Name: input.Name, 
		Description: input.Description,
	})

	// 3. Return updated resource
	c.JSON(http.StatusOK, gin.H{"message": "Restaurant Updated", "restaurant": food})
}

func RestaurantDestroy(c *gin.Context) {
	id := c.Param("id")
	
	// Check if the restaurant exists
	var restaurant models.Restaurant
	if err := initializers.DB.First(&restaurant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	result := initializers.DB.Delete(&restaurant)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete restaurant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restaurant Deleted"})
}