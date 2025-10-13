package model

// User struct for registration
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}