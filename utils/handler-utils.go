package utils

import (
	"GDS-Connect/models"
	"cloud.google.com/go/firestore"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

// GetDatabase Retrieves the database client and context from the gin context
func GetDatabase(ctx *gin.Context) (*firestore.Client, context.Context) {
	// Retrieves the client and dbContext from the gin context
	client, _ := ctx.MustGet("db").(*firestore.Client)
	dbContext, _ := ctx.MustGet("dbContext").(context.Context)
	return client, dbContext
}

// GetUserById retrieves a user from the database by their document ID
func GetUserById(client *firestore.Client, dbContext context.Context, docId string) (models.User, error) {
	var user models.User

	docSnap, err := client.Collection("users").Doc(docId).Get(dbContext)
	if err != nil {
		return user, err
	}

	if err := docSnap.DataTo(&user); err != nil {
		return user, err
	}

	return user, nil
}

// FindMatchingUsers finds users with at least one matching interest with the given user document ID
func FindMatchingUsers(client *firestore.Client, dbContext context.Context, userDocId string) ([]models.User, error) {
	user, err := GetUserById(client, dbContext, userDocId)
	if err != nil {
		return nil, err
	}

	userInterests := user.Interests
	if len(userInterests) == 0 {
		// No interests to match with, return empty slice
		return []models.User{}, nil
	}

	// Query for other users with overlapping interests
	query := client.Collection("users").Where("interests", "array-contains-any", userInterests)
	iter := query.Documents(dbContext)
	defer iter.Stop()

	var matchedUsers []models.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var otherUser models.User
		if err := doc.DataTo(&otherUser); err == nil {
			// Ensure not to match the same user
			if doc.Ref.ID != userDocId {
				matchedUsers = append(matchedUsers, otherUser)
			}
		}
	}

	return matchedUsers, nil
}

// MakeVisibleToUser makes the given user visible to another given user, based on their document IDs
func MakeVisibleToUser(client *firestore.Client, dbContext context.Context, userDocId string, otherUserDocId string) error {
	_, err := client.Collection("users").Doc(userDocId).Update(dbContext, []firestore.Update{
		{
			Path:  "usersVisibility",
			Value: firestore.ArrayUnion(otherUserDocId),
		},
	})

	return err
}

func MakeInvisibleToUser(client *firestore.Client, dbContext context.Context, userDocId string, otherUserDocId string) error {
	_, err := client.Collection("users").Doc(userDocId).Update(dbContext, []firestore.Update{
		{
			Path:  "usersVisibility",
			Value: firestore.ArrayRemove(otherUserDocId),
		},
	})

	return err
}

// MakeVisibleToAllUsers makes the given user visible to all other users
func MakeVisibleToAllUsers(client *firestore.Client, dbContext context.Context, userDocId string) error {
	_, err := client.Collection("users").Doc(userDocId).Update(dbContext, []firestore.Update{
		{
			Path:  "Anonymous",
			Value: true,
		},
	})

	return err
}

// MakeInvisibleToAllUsers makes the given user invisible to all other users
func MakeInvisibleToAllUsers(client *firestore.Client, dbContext context.Context, userDocId string) error {
	_, err := client.Collection("users").Doc(userDocId).Update(dbContext, []firestore.Update{
		{
			Path:  "Anonymous",
			Value: false,
		},
	})

	return err
}
