package service

import (
	"time"

	"go-user-api/internal/models"
	"go-user-api/internal/repository"
)

// helper function to calculate age
func calculateAge(dob string) int {
	layout := "2006-01-02"
	birthDate, err := time.Parse(layout, dob)
	if err != nil {
		return 0
	}

	now := time.Now()
	age := now.Year() - birthDate.Year()

	// adjust if birthday hasn't occurred yet this year
	if now.YearDay() < birthDate.YearDay() {
		age--
	}

	return age
}

func CreateUser(user models.User) models.User {
	user.Age = calculateAge(user.DOB)
	return repository.CreateUser(user)
}

func GetAllUsers() []models.User {
	users := repository.GetAllUsers()

	for i := range users {
		users[i].Age = calculateAge(users[i].DOB)
	}

	return users
}

func GetUserByID(id int) (models.User, bool) {
	user, found := repository.GetUserByID(id)
	if found {
		user.Age = calculateAge(user.DOB)
	}
	return user, found
}
