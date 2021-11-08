package models

type User struct {
	Email     string `json:"email,omitempty" schema:"email, required"`
	FirstName string `json:"first_name,omitempty" schema:"email, required"`
	LastName  string `json:"last_name,omitempty" schema:"email, required"`
	Picture   string `json:"user_picture" json:"-" schema:"user_picture"`
}
