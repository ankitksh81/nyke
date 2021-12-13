package models

type CartRequest struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CartItemCheck struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CartItemResponse struct {
	ID       string  `json:"id"`
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}
