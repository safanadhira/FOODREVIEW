package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/safanadhira/FOODREVIEW/initializers"
	"github.com/safanadhira/FOODREVIEW/models"
)

func FoodIndex(c *gin.Context) {
	var foods []models.Food
	initializers.DB.Find(&foods)

	c.JSON(http.StatusOK, foods)
}

func FoodStore(c *gin.Context) {
	restaurantID := c.Param("ID")

	var restaurant models.Restaurant

	if err := initializers.DB.First(&restaurant, restaurantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
	}

	var body struct {
		Name        string  `json:"name"`
		Price       float64 `json:"price"`
		Description string  `json:"Description"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
		return
	}

	food := models.Food{
		Name:         body.Name,
		Price:        body.Price,
		Description:  &body.Description,
		RestaurantID: restaurant.ID,
	}

	initializers.DB.Create(&food)

	c.JSON(http.StatusOK, gin.H{"message": "Food added!", "food": food})

}

func FoodEdit(c *gin.Context) {
	id := c.Param("id")

	var food models.Food

	if err := initializers.DB.First(&food, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "food not found"})
		return
	}

	var restaurants []models.Restaurant
	initializers.DB.Find(&restaurants)

	c.JSON(http.StatusOK, gin.H{
		"food":        food,
		"restaurants": restaurants,
	})
}

func FoodUpdate(c *gin.Context) {
	id := c.Param("id")

	var food models.Food

	if err := initializers.DB.First(&food, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Food Not Found"})
		return
	}

	var body struct {
		Name           string  `json:"name"`
		Price          float64 `json:"price"`
		Description    string  `json:"description"`
		RestaurantName string  `json:"restaurant_name"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	var restaurant models.Restaurant

	initializers.DB.FirstOrCreate(&restaurant, models.Restaurant{Name: body.RestaurantName})

	desc := body.Description

	initializers.DB.Model(&food).Updates(models.Food{
		Name:         body.Name,
		Price:        body.Price,
		Description:  &desc,
		RestaurantID: restaurant.ID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Food Updated", "food": food})

}

func FoodDelete(c *gin.Context) {
	id := c.Param("id")

	var food models.Food

	if err := initializers.DB.Preload("Restaurant").First(&food, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Food not found"})
	}

	initializers.DB.Delete(&food)

	c.JSON(http.StatusOK, gin.H{"message": "Food Deleted"})
}
