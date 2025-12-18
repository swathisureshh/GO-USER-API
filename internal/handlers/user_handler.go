package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-user-api/internal/models"
	"go-user-api/internal/service"
)

// POST /users
func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	createdUser := service.CreateUser(user)
	c.JSON(http.StatusCreated, createdUser)
}

// GET /users/:id
func GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, found := service.GetUserByID(id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
func GetAllUsers(c *gin.Context) {
	users := service.GetAllUsers()
	c.JSON(200, users)
}
