package handlers

import (
	"E-CommerceAPI/database"
	"encoding/json"
	"net/http"
)

type CartRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CartItemResponse struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Quantity  int    `json:"quantity"`
}

func CartRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userIDInterface, ok := r.Context().Value(UserIDKey).(float64)
		if !ok {
			http.Error(w, "Could not determine the user", http.StatusInternalServerError)
			return
		}

		userID := int(userIDInterface)
		if userID == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req CartRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Wrong JSON", http.StatusBadRequest)
			return
		}

		if req.Quantity == 0 {
			req.Quantity = 1
		}

		_, err = database.DB.Exec(`
									INSERT INTO cart_items (user_id, product_id, quantity) 
									VALUES ($1, $2, $3) 
									ON CONFLICT (user_id, product_id) 
									DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity;`, userID, req.ProductID, req.Quantity)
		if err != nil {
			http.Error(w, "Error add cart", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("CartItems added successfully"))
		return
	}
	if r.Method == http.MethodGet {
		userIDInterface, ok := r.Context().Value(UserIDKey).(float64)
		if !ok {
			http.Error(w, "Could not determine the user", http.StatusInternalServerError)
			return
		}

		userID := int(userIDInterface)
		if userID == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := database.DB.Query(`
									SELECT c.id, c.product_id, p.name, p.price, c.quantity
									FROM cart_items c
									JOIN products p ON c.product_id = p.id
									WHERE c.user_id = $1`, userID)
		if err != nil {
			http.Error(w, "Error get cart", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var cart []CartItemResponse

		for rows.Next() {
			var item CartItemResponse
			rows.Scan(&item.ID, &item.ProductID, &item.Name, &item.Price, &item.Quantity)
			cart = append(cart, item)
		}

		if cart == nil {
			cart = []CartItemResponse{}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cart)
		return
	}
	
	http.Error(w, "Only GET/POST", http.StatusMethodNotAllowed)
}
