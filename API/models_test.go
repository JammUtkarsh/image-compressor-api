package main

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func Test_addProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	tests := []struct {
		name          string
		product       Product
		wantProductID int
	}{
		{
			name: "add product",
			product: Product{
				UserID:             1,
				ProductName:        "test product",
				ProductDescription: "test product description",
				ProductImages:      []string{"https://example.com/image1.png", "https://example.com/image2.png"},
				ProductPrice:       1.99,
			},
			wantProductID: 1,
		},
		{
			name: "add product with empty product images",
			product: Product{
				UserID:             1,
				ProductName:        "test product",
				ProductDescription: "test product description",
				ProductImages:      []string{},
				ProductPrice:       1.99,
			},
			wantProductID: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta(addProductQuery)).
				WithArgs(tt.product.UserID, tt.product.ProductName, tt.product.ProductDescription, pq.Array(tt.product.ProductImages), tt.product.ProductPrice).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.wantProductID))
			if _, err := AddProduct(db, tt.product); err != nil {
				t.Errorf("error was not expected while updating stats: %s", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_userExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	tests := []struct {
		name   string
		userID int64
		want   bool
	}{
		{
			name:   "user does not exist",
			userID: 1,
			want:   false, // This test will fail because the mock doesn't have previous data in it.
		},
		{
			name:   "user exists",
			userID: 1,
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec(regexp.QuoteMeta(FindUserByIDQuery)).
				WithArgs(tt.userID).
				WillReturnResult(sqlmock.NewResult(tt.userID, 1))
			if got := UserExists(db, tt.userID); got != tt.want {
				t.Errorf("userExists() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
