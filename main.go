package main

import (
	"E-CommerceAPI/database"
	"E-CommerceAPI/handlers"
	"E-CommerceAPI/payment_gateway"
	"net/http"
)

func main() {
	database.InitDB()

	http.HandleFunc("/products", handlers.AuthMiddleware(handlers.ProductsHandler))
	http.HandleFunc("/products/", handlers.AuthMiddleware(handlers.ProductByIDHandler))
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/cart", handlers.AuthMiddleware(handlers.CartRouter))
	http.HandleFunc("/checkout", handlers.AuthMiddleware(handlers.CheckoutHandler))

	http.HandleFunc("/internal/stripe/create-session", payment_gateway.GatewayHandler)

	http.ListenAndServe(":8080", nil)
}
