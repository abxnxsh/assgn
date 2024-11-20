package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"mngmnt/internal/database"
	"mngmnt/internal/models"
	"mngmnt/internal/redisclient"
	"mngmnt/internal/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
}
func CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		log.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := database.InsertProduct(&product); err != nil {
		log.WithError(err).Error("Database error while inserting product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	log.WithFields(logrus.Fields{
		"product_id":   product.ID,
		"user_id":      product.UserID,
		"product_name": product.ProductName,
	}).Info("Product created successfully")

	message := rabbitmq.ImageProcessingMessage{
		ProductID:  product.ID,        
		ImageURLs:  product.ProductImages, 
	}

	err := rabbitmq.PublishImageProcessingTask(message)
	if err != nil {
		log.WithError(err).Error("Failed to publish message to RabbitMQ")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send image processing task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":                  product.ID,
		"user_id":             product.UserID,
		"product_name":        product.ProductName,
		"product_description": product.ProductDescription,
		"product_images":      product.ProductImages,
		"product_price":       product.ProductPrice,
		"created_at":          product.CreatedAt,
	})
}

func GetProductByID(c *gin.Context) {

	id := c.Param("id")
	cacheKey := "product:" + id
	var product models.Product
	cachedData, err := redisclient.GetCache(cacheKey)
	if err == nil {
		
		log.WithField("product_id", id).Info("Cache hit for product")
		if err := json.Unmarshal([]byte(cachedData), &product); err != nil {
			log.WithError(err).Error("Error unmarshaling cached data")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error unmarshaling cached data"})
			return
		}
		c.JSON(http.StatusOK, product)
		return
	}

	query := `
		SELECT id, user_id, product_name, product_description, product_images, 
		       compressed_product_images, product_price, created_at
		FROM products WHERE id = $1
	`

	var compressedImagesJSON []byte 

	err = database.DB.QueryRow(query, id).Scan(
		&product.ID,
		&product.UserID,
		&product.ProductName,
		&product.ProductDescription,
		pq.Array(&product.ProductImages), 
		&compressedImagesJSON,          
		&product.ProductPrice,
		&product.CreatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			log.WithError(err).Error("Failed to fetch product from database")
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch product: %v", err)})
		}
		return
	}

	log.WithField("product_id", id).Info("Product retrieved from database")

	
	if len(compressedImagesJSON) > 0 { 
		err = json.Unmarshal(compressedImagesJSON, &product.CompressedProductImages)
		if err != nil {
			log.WithError(err).Error("Failed to unmarshal compressed_product_images")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process compressed images"})
			return
		}
	}
	productJSON, _ := json.Marshal(product)
	err = redisclient.SetCache(cacheKey, string(productJSON), 10*time.Minute)
	if err != nil {
		log.WithError(err).Error("Error caching product in Redis")
	}

	c.JSON(http.StatusOK, product)
}
