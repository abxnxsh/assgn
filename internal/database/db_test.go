package database

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"mngmnt/internal/models"
	"github.com/lib/pq"
)

func TestInsertProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening mock DB: %v", err)
	}
	DB = db
	product := &models.Product{
		UserID:          1,
		ProductName:     "Sample Product",
		ProductDescription: "Sample Description",
		ProductImages:   []string{"image1.jpg", "image2.jpg"},
		ProductPrice:    99.99,
	}
	mock.ExpectQuery("INSERT INTO products").
		WithArgs(product.UserID, product.ProductName, product.ProductDescription, pq.Array(product.ProductImages), product.ProductPrice).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))

	err = InsertProduct(product)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetProductByID(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening mock DB: %v", err)
	}
	DB = db 

	productID := 123
	product := &models.Product{
		ID:                  productID,
		UserID:              1,
		ProductName:         "Sample Product",
		ProductDescription:  "Sample Description",
		ProductImages:       []string{"image1.jpg", "image2.jpg"},
		ProductPrice:        99.99,
	}

	mock.ExpectQuery("SELECT id, user_id, product_name, product_description, product_images, compressed_product_images, product_price").
		WithArgs(productID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "product_name", "product_description", "product_images", "compressed_product_images", "product_price"}).
			AddRow(product.ID, product.UserID, product.ProductName, product.ProductDescription, pq.Array(product.ProductImages), nil, product.ProductPrice))

	result, err := GetProductByID(productID)
	assert.NoError(t, err)
	assert.Equal(t, product, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
