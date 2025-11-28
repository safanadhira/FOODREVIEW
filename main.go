package main

import (
	"github.com/gin-gonic/gin"
	"github.com/safanadhira/FOODREVIEW/initializers"
	"github.com/safanadhira/FOODREVIEW/models"
	"github.com/safanadhira/FOODREVIEW/routes"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.DB.AutoMigrate(&models.Restaurant{}, &models.Food{}, &models.Review{})
}

func main() {
	r := gin.Default()

	r.Static("/static", "./static")

	r.LoadHTMLGlob("templates/**/*")

	routes.RegisterRoutes(r)

	r.Run()
}
