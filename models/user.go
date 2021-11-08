package models

type User struct {
	Email     string `json:"email,omitempty" schema:"email, required"`
	FirstName string `json:"first_name,omitempty" schema:"first_name, required"`
	LastName  string `json:"last_name,omitempty" schema:"last_name, required"`
	Picture   string `json:"user_picture" json:"-" schema:"user_picture"`
}
