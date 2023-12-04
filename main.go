package main

import (
	docs "GDS-Connect/docs"
	"GDS-Connect/handlers"
	"GDS-Connect/middlewares"
	"GDS-Connect/utils"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	docs.SwaggerInfo.BasePath = "/"

	// All routers can access the DB
	router.Use(middlewares.DbMiddleware(client, dbContext))

	// [/api] group
	api := router.Group("/api")
	{
		api.GET("/users", handlers.GetUsers)
		api.GET("/users/:id", handlers.GetUserById)
		api.POST("/users", handlers.CreateUser)
		api.GET("/users/:id/matches", handlers.GetMatches)
		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		// ginSwagger.URL("http://localhost:3000/swagger/doc.json"),
		// ginSwagger.DefaultModelsExpandDepth(-1)))
	}

	// Starts the server
	err = router.Run("localhost:3000")
	if err != nil {
		println("Error launching the server")
		return
	}
}
