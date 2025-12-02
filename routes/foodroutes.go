package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/safanadhira/FOODREVIEW/controllers"
)

func RegisterRoutes(r *gin.Engine) {

	// root route langsung ke daftar restoran
	r.GET("/", controllers.RestaurantIndex)

	restaurant := r.Group("/restaurants")
	{
		restaurant.GET("", controllers.RestaurantIndex)
		restaurant.GET("/", controllers.RestaurantIndex)

		restaurant.GET("/create", controllers.RestaurantCreate)
		restaurant.POST("", controllers.RestaurantStore)
		restaurant.POST("/", controllers.RestaurantStore)
		restaurant.GET("/:id", controllers.RestaurantShow)
		restaurant.GET("/:id/foods/create", controllers.FoodCreate)

		restaurant.POST("/:id/foods", controllers.FoodStore)

		restaurant.GET("/:id/edit", controllers.RestaurantEdit)
		restaurant.POST("/:id", controllers.RestaurantUpdate)
		restaurant.POST("/:id/delete", controllers.RestaurantDestroy)
	}

	food := r.Group("/foods")
	{
		// removed undefined handler controllers.FoodIndex to fix compile error;
		// add a matching handler in controllers package (e.g., FoodIndex) if needed.
		// food.GET("/", controllers.FoodIndex)
		food.GET("/:id/edit", controllers.FoodEdit)

		food.POST("/:id", controllers.FoodUpdate)

		food.DELETE("/:id", controllers.FoodDelete)

		food.POST("/:id/delete", controllers.FoodDelete)

	}

	review := r.Group("/reviews")
	{
		review.GET("/create", controllers.ReviewCreate)

		review.GET("/:id", controllers.ReviewIndex)

		review.POST("", controllers.ReviewStore)
		review.POST("/delete/:id", controllers.ReviewDestroy)
	}
}
