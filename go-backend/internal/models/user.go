package models

import (
	"time"
)

// User represents the user entity in database
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	DOB       time.Time `json:"dob"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// UserResponse represents the API response for a user (includes calculated age)
type UserResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Age  int    `json:"age,omitempty"`
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	DOB  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	DOB  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

// PaginationQuery represents pagination parameters
type PaginationQuery struct {
	Page     int `query:"page" validate:"min=1"`
	PageSize int `query:"page_size" validate:"min=1,max=100"`
}

// PaginatedResponse represents a paginated list response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalCount int64       `json:"total_count"`
	TotalPages int         `json:"total_pages"`
}

// CalculateAge calculates age from date of birth
func CalculateAge(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()

	// Adjust if birthday hasn't occurred yet this year
	if now.YearDay() < dob.YearDay() {
		years--
	}

	return years
}

// ToResponse converts User to UserResponse with calculated age
func (u *User) ToResponse(includeAge bool) UserResponse {
	response := UserResponse{
		ID:   u.ID,
		Name: u.Name,
		DOB:  u.DOB.Format("2006-01-02"),
	}

	if includeAge {
		response.Age = CalculateAge(u.DOB)
	}

	return response
}
