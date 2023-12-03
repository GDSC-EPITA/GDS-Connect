package handlers

import (
	"GDS-Connect/models"
	"GDS-Connect/utils"
	"cloud.google.com/go/firestore"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// curl {{base_url}}:{{server_port}}/api/users

// @BasePath /api

// GetUsers
// PingExample godoc
// @Summary Retrieves all users from the database
// @Schemes
// @Description Retrieves all users from the database
// @Tags Users
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func GetUsers(ctx *gin.Context) {

	client, dbContext := GetDatabase(ctx)

	users := utils.GetUsersFromDatabase(client, dbContext)

	ctx.IndentedJSON(http.StatusOK, users)
}

// GetDatabase Returns the client and the context of the database
func GetDatabase(ctx *gin.Context) (*firestore.Client, context.Context) {
	// Retrieves the client and dbContext from the gin context
	client, _ := ctx.MustGet("db").(*firestore.Client)
	dbContext, _ := ctx.MustGet("dbContext").(context.Context)
	return client, dbContext
}

// GetUserById Returns the user with the given <id> from the database
func GetUserById(ctx *gin.Context) {
	id := ctx.Param("id")

	client, dbContext := GetDatabase(ctx)

	user, err := client.Collection("users").Doc(id).Get(dbContext)
	if err != nil {
		log.Println(err)
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Couldn't find the requested user."})
		return
	}

	ctx.IndentedJSON(http.StatusOK, user.Data())
}

// CreateUser Adds a user from the body of the request to the database
func CreateUser(ctx *gin.Context) {
	var newUser models.User

	// Call BindJSON to bind the received JSON to the <userTmp> type
	if err := ctx.BindJSON(&newUser); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Incorrect POST request: Couldn't parse the body."})
		return
	}

	// Adds the new userTmp to the slice
	client, dbContext := GetDatabase(ctx)
	utils.InsertUserInDatabase(nil, client, dbContext, newUser)
	ctx.IndentedJSON(http.StatusCreated, newUser)
}
