package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/lib/pq"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Mobile    string    `json:"mobile" db:"mobile"`
	Latitude  float64   `json:"latitude" db:"latitude"`
	Longitude float64   `json:"longitude" db:"longitude"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Product struct {
	ProductID               int64     `json:"product_id" db:"product_id"`
	UserID                  int64     `json:"user_id" db:"user_id"`
	ProductName             string    `json:"product_name" db:"product_name"`
	ProductDescription      string    `json:"product_description" db:"product_description"`
	ProductImages           []string  `json:"product_images" db:"product_images"`
	ProductPrice            float64   `json:"product_price" db:"product_price"`
	CompressedProductImages []string  `json:"compressed_product_images" db:"compressed_product_images"`
	CreatedAt               time.Time `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time `json:"updated_at" db:"updated_at"`
}

var db *sql.DB

var (
	port     = 5432
	host     = os.Getenv("POSTGRES_HOST")
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DB")
)

const (
	FindUserByIDQuery = "SELECT * FROM users WHERE id=$1"
	addProductQuery   = "INSERT INTO products (user_id, product_name, product_description, product_images, product_price) VALUES ($1, $2, $3, $4, $5) RETURNING product_id"
)

func Connect() (err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	return db.Ping()
}

func userExists(userID int64) bool {
	_, err := db.Exec(FindUserByIDQuery, userID)
	return err == nil
}

func addProduct(product Product) (productID int, err error) {
	err = db.QueryRow(addProductQuery, product.UserID, product.ProductName, product.ProductDescription, pq.Array(product.ProductImages), product.ProductPrice).Scan(&productID)
	if err != nil {
		return 0, err
	}
	return productID, nil
}
