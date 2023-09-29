package main

import (
	"GDS-Connect/handlers"
	"GDS-Connect/middlewares"
	"GDS-Connect/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	err, dbContext, client := utils.InitDatabase()
	if err != nil {
		println("Error initializing database connection")
		err := client.Close()
		if err != nil {
			println("Error closing the client")
			return
		}
		return
	}

	router := gin.Default()

	// All routers can access the DB
	router.Use(middlewares.DbMiddleware(client, dbContext))

	api := router.Group("/api")
	{
		api.GET("/users", handlers.GetUsers)
		api.GET("/users/:id", handlers.GetUserById)
		api.POST("/users", handlers.CreateUser)
	}

	router.Run("localhost:3000")
}
