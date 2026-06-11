package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

var RoleKey ContextKey = "role"
var UserIDKey ContextKey = "userID"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if (r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/product")) || r.URL.Path == "/register" {
			next(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No Token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			ctx := context.WithValue(r.Context(), RoleKey, claims["role"])
			ctx = context.WithValue(ctx, UserIDKey, claims["user_id"])

			next(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid Token claims", http.StatusUnauthorized)
			return
		}
	}
}
