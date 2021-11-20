package models

type Error struct {
	Code int   `json:"code"`
	Err  error `json:"message"`
}
