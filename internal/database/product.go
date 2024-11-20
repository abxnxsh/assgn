package database

import (
	"encoding/json"
	"fmt"
	"mngmnt/internal/models"
	"github.com/lib/pq"      

)

func InsertProduct(product *models.Product) error {
    query := `
        INSERT INTO products (user_id, product_name, product_description, product_images, product_price)
        VALUES ($1, $2, $3, $4, $5) RETURNING id
    `


    err := DB.QueryRow(query, 
        product.UserID, 
        product.ProductName, 
        product.ProductDescription, 
        pq.Array(product.ProductImages),
        product.ProductPrice,
    ).Scan(&product.ID)

    if err != nil {
        return fmt.Errorf("failed to insert product: %v", err)
    }

    return nil
}
func GetProductByID(productID int) (*models.Product, error) {
	var product models.Product

	query := `
		SELECT id, user_id, product_name, product_description, product_images, 
		       compressed_product_images, product_price
		FROM products WHERE id = $1
	`
	var compressedImagesJSON []byte
	err := DB.QueryRow(query, productID).Scan(
		&product.ID,
		&product.UserID,
		&product.ProductName,
		&product.ProductDescription,
		pq.Array(&product.ProductImages),
		&compressedImagesJSON,
		&product.ProductPrice,
	)
	if err != nil {
		return nil, err
	}

	if len(compressedImagesJSON) > 0 { 
		err = json.Unmarshal(compressedImagesJSON, &product.CompressedProductImages)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal compressed_product_images: %v", err)
		}
	}

	return &product, nil
}
