package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/safanadhira/FOODREVIEW/initializers"
	"github.com/safanadhira/FOODREVIEW/models"
)

const UploadsDir = "static/images/foods"

// --- Image Upload Helper Function (Reused by FoodStore and FoodUpdate) ---
func handleImageUpload(c *gin.Context) (string, error) {
	file, handler, err := c.Request.FormFile("image")

	if err == http.ErrMissingFile {
		// No new file was uploaded, return empty path without error
		return "", nil 
	}
	if err != nil {
		// Other file retrieval error
		return "", fmt.Errorf("error retrieving file: %w", err)
	}
	defer file.Close()

	// Ensure UploadsDir Exists
	if _, statErr := os.Stat(UploadsDir); os.IsNotExist(statErr) {
		if mkdirErr := os.MkdirAll(UploadsDir, os.ModePerm); mkdirErr != nil {
			return "", fmt.Errorf("error creating upload directory: %w", mkdirErr)
		}
	}

	destFilePath := filepath.Join(UploadsDir, handler.Filename)

	dst, createErr := os.Create(destFilePath)
	if createErr != nil {
		return "", fmt.Errorf("error creating destination file: %w", createErr)
	}
	defer dst.Close()

	if _, copyErr := io.Copy(dst, file); copyErr != nil {
		return "", fmt.Errorf("error copying file contents: %w", copyErr)
	}

	// Return the web-accessible path
	return "/" + destFilePath, nil
}
// --------------------------------------------------------------------------

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

	// 1. Handle the Image Upload using the helper
	imagePath, err := handleImageUpload(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 2. Handle Text Form Fields
	name := c.PostForm("name")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")

	food := models.Food{
		Name: name,
		Description: description,
		Price: cleanPrice,
		RestaurantID: restaurant.ID,
		ImagePath: imagePath, // Store the new path (or empty string if none)
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
		"Food": food,
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

	// --- 1. Handle New Image Upload for Update ---
	newImagePath, err := handleImageUpload(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// If a new image was successfully uploaded
	if newImagePath != "" {
		// Delete old image file if it exists
		if food.ImagePath != "" {
			localFilePath := strings.TrimPrefix(food.ImagePath, "/")
			os.Remove(localFilePath)
		}
		// Update the ImagePath field
		food.ImagePath = newImagePath
	}
	// ---------------------------------------------


	// 2. Handle Text Form Fields
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
	food.RestaurantID = restaurantIDUint // Note: You might want to handle this assignment differently

	if restaurantIDStr != "" {
		if i, err := strconv.Atoi(restaurantIDStr); err == nil {
			food.RestaurantID = uint(i)
		}
	}
	
	if err := initializers.DB.Save(&food).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to update food: "+err.Error())
		return
	}

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
	if food.ImagePath != "" {
		// The file path is relative to the web root, so we strip the starting '/'
		localFilePath := strings.TrimPrefix(food.ImagePath, "/")
		// The `os` package is now imported, so this works:
		os.Remove(localFilePath)
	}

	initializers.DB.Where("food_id = ?", food.ID).Delete(&models.Review{})
	initializers.DB.Delete(&food)
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/restaurants/%d", restaurantID))
}