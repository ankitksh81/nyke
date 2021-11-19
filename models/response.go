package models

type Response struct {
	User     ResponseUser `json:"user"`
	JwtToken string       `json:"token"`
}

type ResponseUser struct {
	ID          string `json:"userID"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserPicture string `json:"user_picture" json:"-"`
}
