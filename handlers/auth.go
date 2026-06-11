package handlers

import (
	"E-CommerceAPI/database"
	"E-CommerceAPI/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret")

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Wrong JSON", http.StatusBadRequest)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Wrong password", http.StatusBadRequest)
			return
		}

		user.Password = string(hash)

		insertQuery := "INSERT INTO users(username, password) VALUES($1, $2) RETURNING id"
		err = database.DB.QueryRow(insertQuery, user.Username, user.Password).Scan(&user.ID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		user.Password = ""
		bytes, _ := json.Marshal(user)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
		return
	}
	http.Error(w, "Only POST", http.StatusMethodNotAllowed)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Wrong JSON", http.StatusBadRequest)
			return
		}

		selectQuery := "SELECT id,password,role FROM users WHERE username = $1"
		var dbID int
		var dbPasswordHash string
		var dbRole string
		err = database.DB.QueryRow(selectQuery, user.Username).Scan(&dbID, &dbPasswordHash, &dbRole)
		if err != nil {
			http.Error(w, "Incorrect login or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(dbPasswordHash), []byte(user.Password))
		if err != nil {
			http.Error(w, "Incorrect login or password", http.StatusUnauthorized)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":  dbID,
			"username": user.Username,
			"role":     dbRole,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"token": tokenString,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
	http.Error(w, "Only POST", http.StatusMethodNotAllowed)
}
