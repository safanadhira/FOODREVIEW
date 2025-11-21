package controllers

import (
	"gin/models"

	"github.com/gin-gonic/gin"
)

func FoodIndex(c *gin.Context) {
	var foods []models.Food
	intitializers.DB.Find(&foods)

	c.JSON(http.statusOK, foods)
}
