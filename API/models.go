package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
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

func Connect() (err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	return
}

func userExists(userID int64) bool {
	return true
}

func addProduct(product Product) error {
	return nil
}
