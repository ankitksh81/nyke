package models

type Response struct {
	ID          string `json:"userID"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserPicture string `json:"user_picture" json:"-"`
}
