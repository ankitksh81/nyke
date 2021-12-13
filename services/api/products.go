package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ankitksh81/nyke/logger"
	"github.com/ankitksh81/nyke/middleware"
	"github.com/ankitksh81/nyke/models"
	"github.com/gorilla/mux"
)

// Function to send products as json
func GetProducts(w http.ResponseWriter, r *http.Request) {
	// Get slice of products
	products, err := GetProductsHandler()
	if err != nil {
		middleware.Error500(w)
		return
	}

	middleware.SetContentJSON(w)
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

// get product by product_id
// @ /product/{id}
func GetProductByID(w http.ResponseWriter, r *http.Request) {
	middleware.SetContentJSON(w)
	params := mux.Vars(r)
	product_id := params["id"]

	// call getProduct function with product_id
	prod, err := GetProductByIDHandler(product_id)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No product found with id: ", product_id)
			middleware.Error404(w)
		} else {
			middleware.Error500(w)
		}
		logger.Log.Error("Could not get product by id. " + err.Error())
		return
	}
	json.NewEncoder(w).Encode(prod)
}

func GetProductByIDHandler(product_id string) (models.Product, error) {
	sqlQuery := `SELECT * FROM products WHERE product_id = $1`
	row := middleware.DB.QueryRow(sqlQuery, product_id)

	// An product to hold data from returned row.
	var product models.Product

	// Scan row into product
	err := row.Scan(&product.ID, &product.Name, &product.Price,
		&product.ProductPicture)
	if err != nil {
		return product, err
	}
	return product, nil
}
