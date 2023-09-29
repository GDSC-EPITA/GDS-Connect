package middlewares

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/gin-gonic/gin"
)

// DbMiddleware Middleware that makes the database accessible to the routes
func DbMiddleware(client *firestore.Client, dbContext context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", client)
		c.Set("dbContext", dbContext)
		c.Next()
	}
}
