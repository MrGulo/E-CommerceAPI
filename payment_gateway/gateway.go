package payment_gateway

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SessionRequest struct {
	Amount int64 `json:"amount"`
}

func GatewayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		var req SessionRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Wrong JSON", http.StatusBadRequest)
			return
		}

		if req.Amount < 0 {
			http.Error(w, "Wrong amount", http.StatusBadRequest)
			return
		}

		sessionID := fmt.Sprintf("fake_cs_%d", time.Now().UnixNano())
		paymentURL := fmt.Sprintf("http://localhost:8080/pay/%s", sessionID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"checkout_url": paymentURL,
		})
	}
}
