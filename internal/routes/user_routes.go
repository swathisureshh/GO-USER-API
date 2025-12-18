package routes

import (
	"github.com/gin-gonic/gin"
	"go-user-api/internal/handlers"
)

func UserRoutes(r *gin.Engine) {
	r.POST("/users", handlers.CreateUser)
	r.GET("/users/:id", handlers.GetUserByID)
	r.GET("/users", handlers.GetAllUsers)
}
