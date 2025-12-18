package repository

import "go-user-api/internal/models"

// default user so browser always shows data
var users = []models.User{
	{
		ID:   1,
		Name: "Swathi",
		DOB:  "2002-01-15",
	},
}

var idCounter = 2

func CreateUser(user models.User) models.User {
	user.ID = idCounter
	idCounter++
	users = append(users, user)
	return user
}

func GetAllUsers() []models.User {
	return users
}

func GetUserByID(id int) (models.User, bool) {
	for _, user := range users {
		if user.ID == id {
			return user, true
		}
	}
	return models.User{}, false
}
