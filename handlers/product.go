package handlers

import (
	"E-CommerceAPI/database"
	"E-CommerceAPI/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		rows, err := database.DB.Query("SELECT id, name, price FROM products")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []models.Product
		for rows.Next() {
			var p models.Product
			rows.Scan(&p.ID, &p.Name, &p.Price)
			products = append(products, p)
		}

		if products == nil {
			products = []models.Product{}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(products)
		return
	}
	if r.Method == http.MethodPost {
		var newProduct models.Product

		err := json.NewDecoder(r.Body).Decode(&newProduct)
		if err != nil {
			http.Error(w, "Wrong JSON", http.StatusBadRequest)
			return
		}

		insertQuery := "INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id"
		err = database.DB.QueryRow(insertQuery, newProduct.Name, newProduct.Price).Scan(&newProduct.ID)
		if err != nil {
			http.Error(w, "Error saved in database", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newProduct)
		return
	}
	http.Error(w, "The method is not supported", http.StatusMethodNotAllowed)
}

func ProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/products/")
	targetID, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Wrong ID product", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		var product models.Product

		err := database.DB.QueryRow("SELECT id,name,price FROM products WHERE id = $1", targetID).Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(product)
		return
	}
	if r.Method == http.MethodDelete {

		if userRole := r.Context().Value(RoleKey); userRole != "admin" {
			http.Error(w, "Only an administrator can manage products", http.StatusForbidden)
			return
		}

		result, err := database.DB.Exec("DELETE FROM products WHERE id = $1", targetID)
		if err != nil {
			http.Error(w, "Error when deleting", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Deleted"))
		return
	}
	if r.Method == http.MethodPut {

		if userRole := r.Context().Value(RoleKey); userRole != "admin" {
			http.Error(w, "Only an administrator can manage products", http.StatusForbidden)
			return
		}

		var product models.Product

		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, "Wrong JSON", http.StatusBadRequest)
			return
		}

		result, err := database.DB.Exec("UPDATE products SET name = $1, price = $2 WHERE id = $3", product.Name, product.Price, targetID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		product.ID = targetID

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(product)
		return
	}

	http.Error(w, "Only GET/DELETE/PUT", http.StatusMethodNotAllowed)
}
