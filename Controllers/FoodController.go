package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/safanadhira/FOODREVIEW/initializers"
	"github.com/safanadhira/FOODREVIEW/models"
)

func FoodCreate(c *gin.Context) {
	restaurantID := c.Param("id")

	var restaurant models.Restaurant
	if err := initializers.DB.First(&restaurant, restaurantID).Error; err != nil {
		c.String(http.StatusNotFound, "Restaurant not found")
		return
	}

	c.HTML(http.StatusOK, "foods/create.html", gin.H{
		"Restaurant": restaurant,
	})
}

func FoodStore(c *gin.Context) {

	restaurantID := c.Param("id")

	var restaurant models.Restaurant
	if err := initializers.DB.First(&restaurant, restaurantID).Error; err != nil {
		c.String(http.StatusNotFound, "Restaurant not found")
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")

	cleanPrice := strings.ReplaceAll(priceStr, "Rp.", "")
	cleanPrice = strings.ReplaceAll(cleanPrice, "$", "")
	cleanPrice = strings.TrimSpace(cleanPrice)

	food := models.Food{
		Name:         name,
		Description:  description,
		Price:        cleanPrice,
		RestaurantID: restaurant.ID,
	}

	initializers.DB.Create(&food)

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/restaurants/%d", restaurant.ID))
}
func FoodEdit(c *gin.Context) {
	id := c.Param("id")

	var food models.Food
	if err := initializers.DB.Preload("Restaurant").First(&food, id).Error; err != nil {
		c.String(http.StatusNotFound, "Food not found")
		return
	}

	var restaurants []models.Restaurant
	initializers.DB.Find(&restaurants)

	c.HTML(http.StatusOK, "foods/edit.html", gin.H{
		"Food":        food,
		"Restaurants": restaurants,
	})
}

func FoodUpdate(c *gin.Context) {
	id := c.Param("id")

	var food models.Food
	if err := initializers.DB.First(&food, id).Error; err != nil {
		c.String(http.StatusNotFound, "Food Not Found")
		return
	}

	name := c.PostForm("name")
	priceStr := c.PostForm("price")
	description := c.PostForm("description")
	restaurantIDStr := c.PostForm("restaurant_id")

	var restaurantIDUint uint = 0

	food.Name = name

	food.Price = priceStr

	food.Description = description

	food.RestaurantID = restaurantIDUint

	if restaurantIDStr != "" {
		if i, err := strconv.Atoi(restaurantIDStr); err == nil {
			food.RestaurantID = uint(i)
		}
	}
	if err := initializers.DB.Save(&food).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to update food: "+err.Error())
		return
	}

	initializers.DB.Save(&food)

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/restaurants/%d", food.RestaurantID))
}

func FoodDelete(c *gin.Context) {
	id := c.Param("id")

	var food models.Food

	err := initializers.DB.First(&food, id).Error

	if err != nil {

		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/restaurants/%d", food.RestaurantID))
		return
	}

	restaurantID := food.RestaurantID

	initializers.DB.Where("food_id = ?", food.ID).Delete(&models.Review{})

	initializers.DB.Delete(&food)

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/restaurants/%d", restaurantID))
}
