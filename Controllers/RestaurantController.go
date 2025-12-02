package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/safanadhira/FOODREVIEW/initializers"
	"github.com/safanadhira/FOODREVIEW/models"
	"gorm.io/gorm"
)

type RestaurantInput struct {
	Name        string                `form:"name" binding:"required"`
	Description *string               `form:"description"`
	Image       *multipart.FileHeader `form:"images"`
}

func RestaurantIndex(c *gin.Context) {
	var restaurants []models.Restaurant
	initializers.DB.Find(&restaurants)

	c.HTML(http.StatusOK, "restaurants/index.html", gin.H{
		"restaurants": restaurants,
	})
}

func RestaurantCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "restaurants/create.html", nil)
}

func RestaurantStore(c *gin.Context) {

	name := c.PostForm("name")
	description := c.PostForm("description")

	var count int64
	initializers.DB.Model(&models.Restaurant{}).Where("name = ?", name).Count(&count)

	if count > 0 {
		c.String(http.StatusConflict, "Restaurant name already exists")
		return
	}

	file, err := c.FormFile("image_upload_field")
	var imagePathForDB string

	if err == nil {
		filename := filepath.Base(file.Filename)
		dst := filepath.Join("./static/images", filename)

		if c.SaveUploadedFile(file, dst) == nil {
			imagePathForDB = "/static/images/" + filename
		} else {

			imagePathForDB = "/static/images/upload_failed.png"
		}
	} else {
		imagePathForDB = "/static/images/default.png"
	}

	restaurant := models.Restaurant{
		Name: name,

		Description: &description,
		Image:       &imagePathForDB,
	}
	initializers.DB.Create(&restaurant)

	c.Redirect(http.StatusSeeOther, "/restaurants")
}

func RestaurantShow(c *gin.Context) {
	param := c.Param("id")

	var restaurant models.Restaurant

	query := initializers.DB.Preload("Foods")

	idInt, err := strconv.Atoi(param)

	if err == nil {

		query = query.Where("id = ?", idInt)
	} else {

		query = query.Where("name = ?", param)
	}

	if err := query.First(&restaurant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.String(http.StatusNotFound, "Restaurant not found")
			return
		}
		c.String(http.StatusInternalServerError, "Database Error")
		return
	}

	c.HTML(http.StatusOK, "restaurants/show.html", gin.H{
		"Restaurant": restaurant,
	})
}

func RestaurantEdit(c *gin.Context) {
	id := c.Param("id")

	var restaurant models.Restaurant
	if err := initializers.DB.First(&restaurant, id).Error; err != nil {
		c.String(http.StatusNotFound, "Restaurant not found")
		return
	}

	c.HTML(http.StatusOK, "restaurants/edit.html", gin.H{
		"restaurant": restaurant,
	})
}

// Update restoran
func RestaurantUpdate(c *gin.Context) {
	id := c.Param("id")

	var restaurant models.Restaurant
	if err := initializers.DB.First(&restaurant, id).Error; err != nil {
		c.String(http.StatusNotFound, "Restaurant not found")
		return
	}

	var input RestaurantInput
	if err := c.ShouldBind(&input); err != nil {
		c.String(http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	restaurant.Name = input.Name
	restaurant.Description = input.Description
	initializers.DB.Save(&restaurant)

	c.Redirect(http.StatusSeeOther, "/restaurants")
}

func RestaurantDestroy(c *gin.Context) {
	id := c.Param("id")

	var restaurant models.Restaurant
	if err := initializers.DB.First(&restaurant, id).Error; err != nil {
		c.String(http.StatusNotFound, "Restaurant not found")
		return
	}

	if restaurant.Image != nil && *restaurant.Image != "" {
		ImagePath := *restaurant.Image

		if ImagePath != "/static/images/default.png" && ImagePath != "/statc/images/default-restaurant.png" {
			systempath := "." + ImagePath

			err := os.Remove(systempath)
			if err != nil {
				fmt.Println("Gagal menghapus", err)
			} else {
				fmt.Println("Berhasil Hapus file:", systempath)
			}
		}
	}

	initializers.DB.Delete(&restaurant)
	c.Redirect(http.StatusSeeOther, "/restaurants")
}
