package api

import (
	"encoding/json"
	"net/http"

	"github.com/ankitksh81/nyke/middleware"
	"github.com/ankitksh81/nyke/models"
)

// Function to send products as json
func GetProducts(w http.ResponseWriter, r *http.Request) {
	// Get slice of products
	products, err := GetProductsHandler()
	if err != nil {
		middleware.Error500(w)
		return
	}

	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		middleware.Error500(w)
	}
}

// Function to get all products
func GetProductsHandler() (prod []models.Product, err error) {
	sqlQuery := `SELECT * FROM products`
	rows, err := middleware.DB.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// An product slice to hold data from returned rows.
	var products []models.Product

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price,
			&product.ProductPicture); err != nil {
			return products, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return products, err
	}

	return products, nil
}
