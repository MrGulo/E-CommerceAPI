package handlers

import (
	"E-CommerceAPI/database"
	"bytes"
	"encoding/json"
	"net/http"
)

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {

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

		var totalAmount int
		err := database.DB.QueryRow(`
										SELECT COALESCE(SUM(p.price * c.quantity), 0)
										FROM cart_items c
										JOIN products p ON c.product_id = p.id
										WHERE c.user_id = $1`, userID).Scan(&totalAmount)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if totalAmount == 0 {
			http.Error(w, "The basket is empty", http.StatusBadRequest)
			return
		}

		amountInCents := int64(totalAmount * 100)

		requestBody, _ := json.Marshal(map[string]int64{"amount": amountInCents})

		resp, err := http.Post("http://localhost:8080/internal/stripe/create-session", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			http.Error(w, "Error connecting to the payment gateway", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			http.Error(w, "Payment gateway rejected the request", http.StatusInternalServerError)
			return
		}

		var gatewayResponse map[string]string
		json.NewDecoder(resp.Body).Decode(&gatewayResponse)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(gatewayResponse)
		return

	}
}
