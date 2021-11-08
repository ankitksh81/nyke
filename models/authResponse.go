package models

type AuthResponse struct {
	Email       string `json:"email,omitempty"`
	Name        string `json:"name,omitempty"`
	UserPicture string `json:"picture"`
}
