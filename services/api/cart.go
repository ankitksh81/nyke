package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ankitksh81/nyke/logger"
	"github.com/ankitksh81/nyke/middleware"
	"github.com/ankitksh81/nyke/models"
)

// Handler to add item to cart
func AddItemToCart(w http.ResponseWriter, r *http.Request) {
	middleware.SetContentJSON(w)
	var cartItem models.CartRequest
	err := json.NewDecoder(r.Body).Decode(&cartItem)
	if err != nil {
		middleware.Error500(w)
		return
	}

	// check if user already have item in cart
	// if yes then increment the quantity
	item, err := CheckItemExists(cartItem.UserID, cartItem.ProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			// item does not exists in cart
			// add item to cart with quantity 1
			err := AddItemToCartHandler(cartItem.UserID, cartItem.ProductID, cartItem.Quantity)
			if err != nil {
				middleware.Error500(w) // change to could not add
				logger.Log.Error("Could not add item to cart. " + err.Error())
			}
		} else {
			middleware.Error500(w)
			logger.Log.Error("Could not check if item exists. " + err.Error())
		}
		return
	} else {
		quantity := item.Quantity
		quantity += cartItem.Quantity
		// update item with new quantity
		err := UpdateCartWithItemHandler(cartItem.UserID, cartItem.ProductID, quantity)
		if err != nil {
			middleware.Error500(w) // change to could not add
			logger.Log.Error("Could not update item. " + err.Error())
			return
		}
	}
}

func AddItemToCartHandler(user_id, product_id string, qnt int) error {
	sqlQuery := `INSERT INTO cart(user_id, product_id, quantity)
		VALUES ($1, $2, $3)`

	err := middleware.DB.QueryRow(sqlQuery, user_id, product_id, qnt).Err()
	if err != nil {
		return err
	}
	return nil
}

// Handler to update cart
// Updates the quantity of the item if added more than once to the cart
func UpdateCartWithItemHandler(user_id, product_id string, quantity int) error {
	sqlQuery := `UPDATE cart SET quantity = $1 WHERE user_id = $2 AND product_id = $3`

	err := middleware.DB.QueryRow(sqlQuery, quantity, user_id, product_id).Err()
	if err != nil {
		return err
	}
	return nil
}

// handler to check if the item exists
func CheckItemExists(user_id, product_id string) (models.CartItemCheck, error) {
	sqlQuery := `SELECT * FROM cart WHERE user_id = $1 AND product_id = $2`

	row := middleware.DB.QueryRow(sqlQuery, user_id, product_id)

	var item models.CartItemCheck
	err := row.Scan(&item.ID, &item.UserID, &item.ProductID, &item.Quantity)
	if err != nil {
		return item, err
	}
	return item, nil
}

// get cart items for a user
