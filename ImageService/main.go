package main

import (
	"log"

	"github.com/JammUtkarsh/zocket-intern-Image/services"
)

func main() {
	db, err := services.Connect()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to DB")
	defer db.Close()
	services.KConsumer(db)
}