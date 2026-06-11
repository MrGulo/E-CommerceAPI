package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error

	connStr := "postgres://postgres:secret@localhost:5433/shop?sslmode=disable"

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error Open", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("The base is unavailable. Did you definitely start Docker container?. ", err)
	}
	fmt.Println("Connected to Docker container")

	createProductsQuery := `
				CREATE TABLE IF NOT EXISTS products (
				    id SERIAL PRIMARY KEY,
				    name TEXT NOT NULL,
				    price INT NOT NULL
				);`
	userstableQuery := `
				CREATE TABLE IF NOT EXISTS users (
				    id SERIAL PRIMARY KEY,
				    username TEXT UNIQUE NOT NULL,
				    password TEXT NOT NULL
    
	)`
	usersColumnQuery := `ALTER TABLE users ADD COLUMN IF NOT EXISTS role TEXT DEFAULT 'user';`

	cartItemsQuery := `
				CREATE TABLE IF NOT EXISTS cart_items (
					id SERIAL PRIMARY KEY,
					user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
					product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
					quantity INT NOT NULL DEFAULT 1,
					UNIQUE(user_id, product_id)
	);`

	_, err = DB.Exec(createProductsQuery)
	if err != nil {
		log.Fatal("Error create table products", err)
	}
	fmt.Println("Created table products")

	_, err = DB.Exec(userstableQuery)
	if err != nil {
		log.Fatal("Error create table users", err)
	}
	fmt.Println("Created table users")

	_, err = DB.Exec(usersColumnQuery)
	if err != nil {
		log.Fatal("Error add column role", err)
	}
	fmt.Println("Created column role")

	_, err = DB.Exec(cartItemsQuery)
	if err != nil {
		log.Fatal("Error create table cart_items", err)
	}
	fmt.Println("Created table cart_items")
}
