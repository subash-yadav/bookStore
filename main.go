package main

import (
	"assignment/controllers"
	"assignment/middlewares"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set the router as the default one shipped with Gin

	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())

	api := router.Group("/books")
	{
		api.GET("getBookByISBN", controllers.GetBookDetailsByISBN)
		api.GET("getBookByAuthor", controllers.GetBookByAuthor)
	}

	user := router.Group("/user")
	{
		user.Use(middlewares.JwtAuthMiddleware())
		user.POST("postAuthorDetails", controllers.PostAuthor)
	}

	fmt.Println("Starting server on the port 8080...")

	router.Run(":8080")
}
