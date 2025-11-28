package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/safanadhira/FOODREVIEW/initializers"
	"github.com/safanadhira/FOODREVIEW/models"
)

func FoodIndex(c *gin.Context) {
	var foods []models.Food
	initializers.DB.Preload("Restaurant").Find(&foods)

	c.HTML(http.StatusOK, "foods/index.html", gin.H{
		"foods": foods,
	})
}

func FoodCreate(c *gin.Context) {
	restaurantID := c.Param("id")

	var restaurant models.Restaurant
	if err := initializers.DB.First(&restaurant, restaurantID).Error; err != nil {
		c.String(http.StatusNotFound, "Restaurant not found")
		return
	}

	c.HTML(http.StatusOK, "foods/create.html", gin.H{
		"restaurant": restaurant,
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

	food := models.Food{
		Name:         name,
		Description:  description,
		Price:        priceStr,
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
		"food":        food,
		"restaurants": restaurants,
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
	priceStr := c.PostForm("price")                // Ambil harga sebagai string
	description := c.PostForm("description")       // Ambil deskripsi sebagai string
	restaurantIDStr := c.PostForm("restaurant_id") // Ambil ID Restoran sebagai string

	// --- Konversi ID Restoran (Wajib karena RestaurantID adalah uint) ---
	var restaurantIDUint uint = 0
	if i, err := strconv.Atoi(restaurantIDStr); err == nil {
		restaurantIDUint = uint(i)
	}
	// ------------------------------------------------------------------

	// --- HAPUS logic lama (ParseFloat) yang menyebabkan error kompilasi ---
	// if p, err := strconv.ParseFloat(price, 64); err == nil { ... }
	// -----------------------------------------------------------------------

	food.Name = name

	// FIX 1: food.Price (Assign string langsung, mengatasi error float64 vs string)
	food.Price = priceStr

	// FIX 2: food.Description (Assign string langsung, mengatasi error *string vs string)
	food.Description = description

	// FIX 3: Assign ID Restoran yang sudah dikonversi
	food.RestaurantID = restaurantIDUint

	initializers.DB.Save(&food)

	// FIX 4: Redirect ke halaman detail restoran (/restaurants/%d), bukan /foods (yang menyebabkan 404)
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
