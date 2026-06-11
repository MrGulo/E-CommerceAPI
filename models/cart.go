package models

type CartItem struct {
	ID        int `json:"id"`
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
