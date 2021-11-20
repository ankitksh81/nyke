package models

type Product struct {
	ID             string  `json:"product_id"`
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	ProductPicture string  `json:"picture"`
}
