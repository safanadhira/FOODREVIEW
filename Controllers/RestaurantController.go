package controllers

import (
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/safanadhira/FOODREVIEW/initializers"
	"github.com/safanadhira/FOODREVIEW/models"
	"gorm.io/gorm"
)

// --- Input Struct ---
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
	// 1. Ambil data teks (Name, Description) secara manual dari form
	//    Kita TIDAK menggunakan c.ShouldBind(&input) secara langsung karena ada file.
	name := c.PostForm("name")
	description := c.PostForm("description") // Asumsi Description adalah string non-pointer

	// 2. Cek Unik Nama
	var count int64
	initializers.DB.Model(&models.Restaurant{}).Where("name = ?", name).Count(&count)

	if count > 0 {
		c.String(http.StatusConflict, "Restaurant name already exists")
		return
	}

	// 3. --- LOGIKA FILE UPLOAD: Menggantikan c.FormFile & imagePathForDB ---
	file, err := c.FormFile("image_upload_field") // Nama field harus sama dengan di form HTML!
	var imagePathForDB string

	if err == nil {
		// Simpan file ke folder static
		filename := filepath.Base(file.Filename)
		dst := filepath.Join("./static/images", filename)

		if c.SaveUploadedFile(file, dst) == nil {
			imagePathForDB = "/static/images/" + filename
		} else {
			// Jika SaveUploadedFile gagal, kita bisa log error atau menggunakan default
			imagePathForDB = "/static/images/upload_failed.png"
		}
	} else {
		// Jika tidak ada file diupload (atau error form)
		imagePathForDB = "/static/images/default.png"
	}
	// ------------------------------------------------------------------------

	// 4. Deklarasikan dan Isi struct Restaurant
	restaurant := models.Restaurant{
		Name: name,
		// Asumsi Description di model adalah string non-pointer
		Description: &description,
		Image:       &imagePathForDB,
	}
	initializers.DB.Create(&restaurant)

	c.Redirect(http.StatusSeeOther, "/restaurants")
}

func RestaurantShow(c *gin.Context) {
	param := c.Param("id") // Menangkap parameter, bisa "1" atau "agus"

	var restaurant models.Restaurant

	// Siapkan base query dengan Preload Foods
	query := initializers.DB.Preload("Foods")

	// 1. Coba konversi parameter menjadi integer (membutuhkan import "strconv")
	idInt, err := strconv.Atoi(param)

	if err == nil {
		// 2. Jika konversi berhasil (adalah angka/ID), cari berdasarkan Primary Key
		query = query.Where("id = ?", idInt)
	} else {
		// 3. Jika konversi gagal (adalah string/Nama), cari berdasarkan kolom Name
		query = query.Where("name = ?", param)
	}

	// 4. Eksekusi query
	if err := query.First(&restaurant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.String(http.StatusNotFound, "Restaurant not found")
			return
		}
		c.String(http.StatusInternalServerError, "Database Error")
		return
	}

	// 5. Render
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

// Delete restoran
func RestaurantDestroy(c *gin.Context) {
	id := c.Param("id")

	var restaurant models.Restaurant
	if err := initializers.DB.First(&restaurant, id).Error; err != nil {
		c.String(http.StatusNotFound, "Restaurant not found")
		return
	}

	initializers.DB.Delete(&restaurant)
	c.Redirect(http.StatusSeeOther, "/restaurants")
}
