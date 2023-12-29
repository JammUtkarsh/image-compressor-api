package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	defer db.Close()
	log.Println("Server listening on: http://127.0.0.1:8080")
	http.HandleFunc("/", AddProducts)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func AddProducts(w http.ResponseWriter, r *http.Request) {
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
		if errs := basicProductChecks(product); errs != nil {
			errorMessages := make([]string, len(errs))
			for i, err := range errs {
				errorMessages[i] = err.Error()
			}
			http.Error(w, fmt.Sprintf("Invalid product: %s", errorMessages), http.StatusBadRequest)
			return
		}
		if _, err := addProduct(product); err != nil {
			http.Error(w, "Unable to add the product", http.StatusInternalServerError)
			return
		}
		// return a response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Product added successfully"))
		return
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

func basicProductChecks(product Product) []error {
	var validationErrors []error
	switch {
	case !userExists(product.UserID) || product.UserID == 0:
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
		_, err := url.ParseRequestURI(urlString)
		if err != nil {
			return false
		}
	}
	return true
}
