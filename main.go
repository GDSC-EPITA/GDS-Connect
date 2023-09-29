package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users = []user{
	{
		ID:   "0",
		Name: "Maxence",
		Age:  21,
	},
	{
		ID:   "1",
		Name: "Lucas",
		Age:  23,
	},
	{
		ID:   "3",
		Name: "Felipe",
		Age:  25,
	},
}

// getUsers returns the users as JSON
func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

func createUser(c *gin.Context) {
	var newUser user

	// Call BindJSON to bind the received JSON to the <user> type
	if err := c.BindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Incorrect POST request: Couldn't parse the body."})
		return
	}

	// Adds the new user to the slice
	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}

func getUserById(ctx *gin.Context) {
	id := ctx.Param("id")

	for _, elt := range users {
		if elt.ID == id {
			ctx.IndentedJSON(http.StatusOK, elt)
			return
		}
	}

	ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Couldn't find the requested user."})
}

func main() {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/users", getUsers)
		api.GET("/users/:id", getUserById)
		api.POST("/users", createUser)
	}

	router.Run("localhost:3000")
}
