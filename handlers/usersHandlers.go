package handlers

import (
	"GDS-Connect/models"
	"GDS-Connect/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// curl {{base_url}}:{{server_port}}/api/users

// GetUsers godoc
// @Summary Retrieves all users from the database with their document IDs
// @Schemes
// @Description Retrieves all users from the database along with their Firestore document IDs
// @Tags Users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} string "Error: Error retrieving users"
// @Router /users [get]
func GetUsers(ctx *gin.Context) {

	client, dbContext := utils.GetDatabase(ctx)

	users := utils.GetUsersFromDatabase(client, dbContext)

	if users == nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving users"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, users)
}

// GetUserById godoc
// @Summary Get a user by ID
// @Description Returns the user with the given ID from the database
// @Tags Users
// @Produce  json
// @Param id path string true "User ID"
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
// @Failure 400 {object} string "Error: Interests cannot be null"
// @Router /users [post]
func CreateUser(ctx *gin.Context) {
	var newUser models.User

	// Call BindJSON to bind the received JSON to the <userTmp> type
	if err := ctx.BindJSON(&newUser); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Incorrect POST request: Couldn't parse the body."})
		return
	}

	// Validate that interests are not null
	if newUser.Interests == nil || len(newUser.Interests) == 0 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Interests cannot be null"})
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
// @Param id path string true "User ID"
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

type OtherUsers struct {
	IDs []string `json:"otherUsers"`
}

// MakeVisible godoc
// @Summary Make a user visible to other given users
// @Description Makes a user visible to other users based on the given user ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} string "User visibility updated"
// @Failure 400 {object} string "Error: Invalid user ID"
// @Failure 404 {object} string "Error: User not found"
// @Failure 500 {object} string "Error: Internal server error"
// @Router /users/{id}/visible [post]
func MakeVisible(ctx *gin.Context) {
	// Get other users from the body of the request
	var otherUsers OtherUsers
	if err := ctx.BindJSON(&otherUsers); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Incorrect POST request: Couldn't parse the body."})
		return
	}

	id := ctx.Param("id")
	client, dbContext := utils.GetDatabase(ctx)

	err := utils.MakeVisibleToUsers(client, dbContext, id, otherUsers.IDs)
	if err != nil {
		if err == strconv.ErrSyntax {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		} else {
			println(err.Error())
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error making user visible"})
		}
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "User visibility updated"})
}

// MakeInvisible godoc
// @Summary Make a user invisible to other given users
// @Description Makes a user invisible to other users based on the given user ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} string "User visibility updated"
// @Failure 400 {object} string "Error: Invalid user ID"
// @Failure 404 {object} string "Error: User not found"
// @Failure 500 {object} string "Error: Internal server error"
// @Router /users/{id}/invisible [post]
func MakeInvisible(ctx *gin.Context) {
	// Get other users from the body of the request
	var otherUsers OtherUsers
	if err := ctx.BindJSON(&otherUsers); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Incorrect POST request: Couldn't parse the body."})
		return
	}

	id := ctx.Param("id")
	client, dbContext := utils.GetDatabase(ctx)

	err := utils.MakeInvisibleToUsers(client, dbContext, id, otherUsers.IDs)
	if err != nil {
		if err == strconv.ErrSyntax {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		} else {
			println(err.Error())
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error making user invisible"})
		}
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "User visibility updated"})
}

// UpdateVisibilityToAll godoc
// @Summary Make a user visibility to all other users
// @Description Makes a user visible or invisible to all other users based on the given user ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} string "User visibility updated"
// @Failure 400 {object} string "Error: Invalid user ID"
// @Failure 404 {object} string "Error: User not found"
// @Failure 500 {object} string "Error: Internal server error"
// @Router /users/{id}/visibleToAll [post]
func UpdateVisibilityToAll(ctx *gin.Context) {
	id := ctx.Param("id")
	client, dbContext := utils.GetDatabase(ctx)

	res, err := utils.UpdateVisibility(client, dbContext, id)
	if err != nil {
		if err == strconv.ErrSyntax {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		} else {
			println(err.Error())
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error making user visible"})
		}
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"anonymous": *res})
}
