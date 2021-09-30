package models

import (
	"github.com/google/uuid"
)

type User struct {
	UserID      uuid.UUID `pg:"type:uuid,default:uuid_generate_v4()"`
	Email       string    `json:"email,omitempty"`
	PictureLink string    `json:"picture_link"`
}
