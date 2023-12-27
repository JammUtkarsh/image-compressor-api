package main

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        int
	Name      string  `json:"name"`
	Mobile    string  `json:"mobile"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Product struct {
	ProductID               int64    `json:"product_id"`
	UserID                  int64    `json:"user_id"`
	ProductName             string   `json:"product_name"`
	ProductDescription      string   `json:"product_description"`
	ProductImages           []string `json:"product_images"`
	ProductPrice            float64  `json:"product_price"`
	CompressedProductImages []string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

var db *gorm.DB

var (
	port     = 5432
	host     = os.Getenv("POSTGRES_HOST")
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DB")
)

func Connect() (err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return
}

func userExists(userID int64) bool {
	return true
}

func addProduct(product Product) error {
	return nil
}
