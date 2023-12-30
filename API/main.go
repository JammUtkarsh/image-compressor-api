package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type api struct {
	db *sql.DB
}

func main() {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	log.Println("Connected to database")
	a := api{db: db}

	log.Println("Server listening on: http://127.0.0.1:8080")
	http.HandleFunc("/", a.AddProducts)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func (a *api) AddProducts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "POST":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to parse body", http.StatusInternalServerError)
			return
		}
		var product Product
		err = json.Unmarshal(body, &product)

		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if errs := a.basicProductChecks(product); errs != nil {
			errorMessages := make([]string, len(errs))
			for i, err := range errs {
				errorMessages[i] = err.Error()
			}
			http.Error(w, fmt.Sprintf("Invalid product: %s", errorMessages), http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, "Unable to connect to database", http.StatusInternalServerError)
			return
		}
		var productID int64
		if productID, err = AddProduct(a.db, product); err != nil {
			http.Error(w, "Unable to add the product", http.StatusInternalServerError)
			return
		}
		log.Println("Product added successfully")
		msg := struct {
			UserID    int64 `json:"u"`
			ProductID int64 `json:"p"`
		}{UserID: product.UserID, ProductID: productID}
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, "Unable to marshal message", http.StatusInternalServerError)
			return
		}

		KProducer(msgBytes)
		log.Println("Message sent to Kafka")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Product added successfully"))
		return

	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

func (a *api) basicProductChecks(product Product) (validationErrors []error) {
	switch {
	case !UserExists(a.db, product.UserID) || product.UserID == 0:
		validationErrors = append(validationErrors, fmt.Errorf("user %d does not exist", product.UserID))
	case product.ProductName == "":
		validationErrors = append(validationErrors, errors.New("product name cannot be empty"))
	case !validURLs(product.ProductImages):
		validationErrors = append(validationErrors, errors.New("product images must be valid URLs"))
	case product.ProductPrice < 0:
		validationErrors = append(validationErrors, errors.New("product price cannot be negative"))
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validURLs(urls []string) bool {
	for _, urlString := range urls {
		if !strings.HasPrefix(urlString, "http") {
			return false
		}
	}
	return true
}
