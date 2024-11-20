package main

import (
	"log"
	"net/http"
	_ "net/http/pprof" 

	"mngmnt/internal/database"
	"mngmnt/internal/handlers"
	"mngmnt/internal/rabbitmq"
	"mngmnt/internal/redisclient"
	"github.com/gin-gonic/gin"
)


// main function
func main() {
	go func() {
		log.Println("Starting pprof server on port 6060")
		log.Println(http.ListenAndServe("localhost:6060", nil)) 
	}()


	//redis part
	redisclient.InitializeRedis()

    log.Println("Starting application...")


	if err := database.ConnectDB(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}


	redisclient.InitializeRedis()

	//rabbit mq part
	if err := rabbitmq.InitializeRabbitMQ(); err != nil {
		log.Fatalf("Error initializing RabbitMQ: %v", err)
	}


	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

	// handling routes
	router.GET("/products/:id", handlers.GetProductByID)
	router.POST("/products", handlers.CreateProduct)   

	log.Println("Server is running on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
