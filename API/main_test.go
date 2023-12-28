package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_validImageURL(t *testing.T) {
	tests := []struct {
		name string
		urls []string
		want bool
	}{
		{
			name: "valid urls",
			urls: []string{
				"https://example.com/image1.jpg",
				"https://example.com/image2.jpg",
				"https://example.com/image3.jpg",
			},
			want: true,
		},
		{
			name: "invalid urls",
			urls: []string{
				"https://example.com/image1.jpg",
				"https://example.com/image2.jpg",
				"https://example.com/image3.jpg",
				"invalid-url",
			},
			want: false,
		},
		{
			name: "empty urls",
			urls: []string{},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validURLs(tt.urls); got != tt.want {
				t.Errorf("validImageURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddProducts(t *testing.T) {
	// This is a temporary fix, it need to be replaced with an mock database later
	if err := Connect(); err != nil {
		panic(err)
	}
	defer db.Close()

	// Create a test server with the AddProducts handler and mocked addProduct
	testServer := httptest.NewServer(http.HandlerFunc(AddProducts))
	defer testServer.Close()

	// Test cases for various scenarios
	testCases := []struct {
		name          string
		method        string
		path          string
		body          []byte
		expectedCode  int
		expectedError string
	}{
		{
			name:          "Successful product addition",
			method:        "POST",
			path:          "/",
			body:          []byte(`{"user_id":1,"product_name":"Cool Gadget","product_description":"This gadget is super cool and does amazing things!","product_images":["https://example.com/image1.jpg","https://example.com/image2.png"],"product_price":49.99}`),
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name:          "Invalid method",
			method:        "GET",
			path:          "/",
			body:          []byte(``),
			expectedCode:  http.StatusMethodNotAllowed,
			expectedError: "Invalid request method\n",
		},
		{
			name:          "Invalid JSON",
			method:        "POST",
			path:          "/",
			body:          []byte(``),
			expectedCode:  http.StatusBadRequest,
			expectedError: "Invalid JSON\n",
		},
		{
			name:          "Invalid user",
			method:        "POST",
			path:          "/",
			body:          []byte(`{"user_id":0,"product_name":"Cool Gadget", "product_images":[],"product_price":49.99}`),
			expectedCode:  http.StatusBadRequest,
			expectedError: "Invalid product: [user 0 does not exist]\n",
		},
		{
			name:          "Invalid URL",
			method:        "POST",
			path:          "/",
			body:          []byte(`{"user_id":1,"product_name":"Cool Gadget", "product_images":["invalid-url", "/path/to/something"],"product_price":49.99}`),
			expectedCode:  http.StatusBadRequest,
			expectedError: "Invalid product: [product images must be valid URLs]\n",
		},
		{
			name:          "Negative price",
			method:        "POST",
			path:          "/",
			body:          []byte(`{"user_id":1,"product_name":"Cool Gadget", "product_images":[],"product_price":-49.99}`),
			expectedCode:  http.StatusBadRequest,
			expectedError: "Invalid product: [product price cannot be negative]\n",
		},
		{
			name:          "Empty name",
			method:        "POST",
			path:          "/",
			body:          []byte(`{"user_id":1,"product_name":"", "product_images":[],"product_price":49.99}`),
			expectedCode:  http.StatusBadRequest,
			expectedError: "Invalid product: [product name cannot be empty]\n",
		},
		{
			name:          "Invalid Path",
			method:        "POST",
			path:          "/invalid-path",
			body:          []byte(`{"user_id":1,"product_name":"Cool Gadget", "product_images":[],"product_price":49.99}`),
			expectedCode:  http.StatusNotFound,
			expectedError: "404 not found.\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, testServer.URL+tc.path, bytes.NewBuffer(tc.body))
			if err != nil {
				t.Fatal(err)
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			if res.StatusCode != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedCode, res.StatusCode)
			}

			if tc.expectedError != "" {
				body, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatal(err)
				}
				
				if string(body) != tc.expectedError {
					t.Errorf("Expected error message %q, got %q", tc.expectedError, string(body))
				}
			}
		})
	}
}
