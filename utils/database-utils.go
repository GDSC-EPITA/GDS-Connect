package utils

import (
	"GDS-Connect/models"
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"os"
)

// InitDatabase Inits the database and returns the client and context generated
func InitDatabase() (error, context.Context, *firestore.Client) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	ctx := context.Background()
	firebaseConfig := &firebase.Config{
		// AuthOverride:  os.Getenv("FIREBASE_API_KEY"),
		DatabaseURL:   "https://" + os.Getenv("FIREBASE_PROJECT_ID") + ".firebaseio.com",
		ProjectID:     os.Getenv("FIREBASE_PROJECT_ID"),
		StorageBucket: os.Getenv("FIREBASE_STORAGE_BUCKET"),
	}
	opt := option.WithCredentialsFile(os.Getenv("SERVICE_ACCOUNT_KEY_PATH"))
	app, err := firebase.NewApp(context.Background(), firebaseConfig, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatal(err)
	}

	/*
		defer func(client *firestore.Client) {
			err := client.Close()
			if err != nil {
				log.Fatalf("error closing the database %v", err)
			}
		}(client)
	*/

	return err, ctx, client
}

// GetUsersFromDatabase Retrieves the users from the database as a list
func GetUsersFromDatabase(client *firestore.Client, ctx context.Context) []map[string]interface{} {
	users, err := client.Collection("users").Documents(ctx).GetAll()
	if err != nil {
		log.Println(err)
		return nil
	}

	// Converts the iterator to a list
	var userArray []map[string]interface{}
	for _, user := range users {
		userArray = append(userArray, user.Data())
	}

	// fmt.Printf("%#v\n", userArray)
	return userArray
}

// InsertUserInDatabase Inserts user <newUser> in the database
func InsertUserInDatabase(err error, client *firestore.Client, ctx context.Context, newUser models.User) {
	if newUser.Age <= 0 {
		log.Printf("Error: User <%s> has a non-valid age: <%d>", newUser.Name, newUser.Age)
		return
	}

	// Adds a userTmp to the <users> collection in the Firestore DB
	_, _, err = client.Collection("users").Add(ctx, map[string]interface{}{
		"name": newUser.Name,
		"age":  newUser.Age,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("OK: Inserted userTmp %s", newUser.Name)
}
