package models

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	DOB   string `json:"dob"`
	Age  int     `json:"age"`
}
