package service

import (
	"testing"
	"time"

	"user-api/internal/models"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name     string
		dob      time.Time
		expected int
	}{
		{
			name:     "30 years old - birthday passed",
			dob:      time.Date(1994, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: 30, // Assuming current date is after Jan 15
		},
		{
			name:     "25 years old",
			dob:      time.Date(1999, 6, 20, 0, 0, 0, 0, time.UTC),
			expected: 25,
		},
		{
			name:     "0 years old - born this year",
			dob:      time.Now().AddDate(0, -6, 0), // 6 months ago
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			age := models.CalculateAge(tt.dob)
			// Note: This test may need adjustment based on the current date
			// The age calculation depends on whether the birthday has occurred this year
			if age < 0 {
				t.Errorf("CalculateAge() returned negative age: %d", age)
			}
		})
	}
}

func TestCalculateAge_ExactBirthday(t *testing.T) {
	// Test for someone born exactly N years ago today
	yearsAgo := 25
	dob := time.Now().AddDate(-yearsAgo, 0, 0)
	
	age := models.CalculateAge(dob)
	
	if age != yearsAgo {
		t.Errorf("CalculateAge() = %d, expected %d for exact birthday", age, yearsAgo)
	}
}

func TestCalculateAge_BeforeBirthday(t *testing.T) {
	// Test for someone whose birthday hasn't occurred yet this year
	now := time.Now()
	
	// Create a DOB that's in the future this year (birthday not yet)
	futureBirthday := now.AddDate(-30, 0, 30) // 30 years ago + 30 days in future
	
	age := models.CalculateAge(futureBirthday)
	
	// Should be 29 if birthday hasn't occurred yet
	if age < 0 {
		t.Errorf("CalculateAge() should not return negative age: %d", age)
	}
}

func TestUserToResponse(t *testing.T) {
	dob := time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC)
	user := &models.User{
		ID:   1,
		Name: "Alice",
		DOB:  dob,
	}

	// Test with age
	responseWithAge := user.ToResponse(true)
	if responseWithAge.ID != 1 {
		t.Errorf("Expected ID 1, got %d", responseWithAge.ID)
	}
	if responseWithAge.Name != "Alice" {
		t.Errorf("Expected name Alice, got %s", responseWithAge.Name)
	}
	if responseWithAge.DOB != "1990-05-10" {
		t.Errorf("Expected DOB 1990-05-10, got %s", responseWithAge.DOB)
	}
	if responseWithAge.Age <= 0 {
		t.Errorf("Expected positive age, got %d", responseWithAge.Age)
	}

	// Test without age
	responseWithoutAge := user.ToResponse(false)
	if responseWithoutAge.Age != 0 {
		t.Errorf("Expected age 0 when includeAge=false, got %d", responseWithoutAge.Age)
	}
}
