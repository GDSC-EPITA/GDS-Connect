package handlers

import (
	"GDS-Connect/models"
	"GDS-Connect/utils"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

// curl {{base_url}}:{{server_port}}/api/users

// GetUsers godoc
// @Summary Retrieves all users from the database
// @Schemes
// @Description Retrieves all users from the database
// @Tags Users
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func GetUsers(ctx *gin.Context) {

	client, dbContext := utils.GetDatabase(ctx)

	users := utils.GetUsersFromDatabase(client, dbContext)

	ctx.IndentedJSON(http.StatusOK, users)
}

// GetUserById godoc
// @Summary Get a user by ID
// @Description Returns the user with the given ID from the database
// @Tags Users
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} models.User "User data"
// @Failure 400 {object} string "Error: Invalid user ID"
// @Failure 404 {object} string "Error: Couldn't find the requested user"
// @Router /users/{id} [get]
func GetUserById(ctx *gin.Context) {
    id := ctx.Param("id")

    client, dbContext := utils.GetDatabase(ctx)

    user, err := utils.GetUserById(client, dbContext, id)
    if err != nil {
        if err == strconv.ErrSyntax {
            ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        } else {
            ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Couldn't find the requested user."})
        }
        return
    }

    ctx.IndentedJSON(http.StatusOK, user)
}


// CreateUser godoc
// @Summary Creates a new user
// @Schemes
// @Description Adds a user from the body of the request to the database
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User required "User info"
// @Success 201 {object} string "User created"
// @Failure 400 {object} string "Error: Incorrect POST request: Couldn't parse the body."
// @Router /users [post]
func CreateUser(ctx *gin.Context) {
	var newUser models.User

	// Call BindJSON to bind the received JSON to the <userTmp> type
	if err := ctx.BindJSON(&newUser); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Incorrect POST request: Couldn't parse the body."})
		return
	}

	// Adds the new userTmp to the slice
	client, dbContext := utils.GetDatabase(ctx)
	utils.InsertUserInDatabase(nil, client, dbContext, newUser)
	ctx.IndentedJSON(http.StatusCreated, newUser)
}

// GetMatches godoc
// @Summary Match users by shared interests
// @Description Finds users with matching interests based on the given user ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {array} models.User "List of users with matching interests"
// @Failure 400 {object} string "Error: Invalid user ID"
// @Failure 404 {object} string "Error: User not found"
// @Failure 500 {object} string "Error: Internal server error"
// @Router /users/{id}/matches [get]
func GetMatches(ctx *gin.Context) {
    id := ctx.Param("id")

    client, dbContext := utils.GetDatabase(ctx)

    // Utilize the FindMatchingUsers utility function
    matchedUsers, err := utils.FindMatchingUsers(client, dbContext, id)
    if err != nil {
        if err == strconv.ErrSyntax {
            ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
        } else {
            ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error finding matches"})
        }
        return
    }

    ctx.IndentedJSON(http.StatusOK, matchedUsers)
}

