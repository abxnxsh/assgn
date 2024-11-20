package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

type ImageProcessingMessage struct {
	ProductID int      `json:"product_id"`
	ImageURLs []string `json:"image_urls"`
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"image_processing_queue", 
		"",                      
		true,                     
		false,                   
		false,                    
		false,                    
		nil,                      
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Postgre---- Connection
	db, err := sql.Open("postgres", "postgres://product_user:password@localhost:5432/product_db?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	log.Println("Waiting for messages. To exit press CTRL+C")
	for msg := range msgs {
		var task ImageProcessingMessage
		if err := json.Unmarshal(msg.Body, &task); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		log.Printf("Processing images for product ID: %d", task.ProductID)

		var compressedURLs []string
		for _, imageURL := range task.ImageURLs {
			compressedURL := strings.Replace(imageURL, "http", "compressed-http", 1)
			compressedURLs = append(compressedURLs, compressedURL)
			log.Printf("Compressed image URL: %s", compressedURL)
		}
		compressedURLsJSON, err := json.Marshal(compressedURLs)
		if err != nil {
			log.Printf("Failed to marshal compressed URLs: %v", err)
			continue
		}

		_, err = db.Exec(
			"UPDATE products SET compressed_product_images = $1 WHERE id = $2",
			compressedURLsJSON,
			task.ProductID,
		)
		if err != nil {
			log.Printf("Failed to update database for product ID %d: %v", task.ProductID, err)
		} else {
			log.Printf("Successfully updated database for product ID %d", task.ProductID)
		}
	}
}
