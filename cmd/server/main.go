package main

import (
	"github.com/gin-gonic/gin"
	"go-user-api/internal/routes"
)

func main() {
	r := gin.Default()

	// root test route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server running"})
	})

	// register user routes
	routes.UserRoutes(r)

	r.Run(":8080")
}
