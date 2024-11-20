package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

type ImageProcessingMessage struct {
	ProductID int      `json:"product_id"`
	ImageURLs []string `json:"image_urls"`
}

const (
	retryLimit = 3
)

func InitializeRabbitMQ() error {
	var err error

	// Connect to RabbitMQ
	RabbitConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	RabbitChannel, err = RabbitConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open RabbitMQ channel: %v", err)
	}

	_, err = RabbitChannel.QueueDeclare(
		"image_processing_queue", 
		true,                     
		false,                    
		false,                    
		false,                    
		nil,                      
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}


	_, err = RabbitChannel.QueueDeclare(
		"dlq.queue",
		true,        
		false,       
		false,      
		false,       
		nil,         
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ: %v", err)
	}

	log.Println("RabbitMQ connection, queue, and DLQ initialized successfully")
	return nil
}


func PublishImageProcessingTask(message ImageProcessingMessage) error {
	
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}


	err = RabbitChannel.Publish(
		"",                        
		"image_processing_queue",  
		false,                     
		false,                     
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	log.Printf("Published image processing task for product ID %d", message.ProductID)
	return nil
}


func ProcessMessage(msg amqp.Delivery) {
	var message ImageProcessingMessage
	err := json.Unmarshal(msg.Body, &message)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		msg.Ack(false)
		return
	}

	log.Printf("Processing message for ProductID: %d", message.ProductID)


	if err := processImages(message); err != nil {
		log.Printf("Error processing message: %v", err)

		if retryCount := getRetryCount(msg); retryCount < retryLimit {
			log.Printf("Retrying message (Attempt %d/%d)", retryCount+1, retryLimit)
			retryMessage(msg, retryCount+1)
			return
		}
		log.Printf("Message failed after %d retries, sending to DLQ", retryLimit)
		sendToDeadLetterQueue(msg)
		return
	}

	msg.Ack(false)
	log.Printf("Message processed successfully for ProductID: %d", message.ProductID)
}

func processImages(message ImageProcessingMessage) error {
	return fmt.Errorf("simulated error in image processing")
}

func getRetryCount(msg amqp.Delivery) int {
	if deaths, ok := msg.Headers["x-death"]; ok {
		if deathArray, ok := deaths.([]interface{}); ok && len(deathArray) > 0 {
			if deathMap, ok := deathArray[0].(amqp.Table); ok {
				if count, ok := deathMap["count"].(int64); ok {
					return int(count)
				}
			}
		}
	}
	return 0
}

func retryMessage(msg amqp.Delivery, retryCount int) {
	msg.Nack(false, false)
	time.Sleep(5 * time.Second)
}

func sendToDeadLetterQueue(msg amqp.Delivery) {
	err := publishMessageToQueue("dlq.queue", msg.Body)
	if err != nil {
		log.Printf("Failed to publish message to DLQ: %v", err)
	}
	msg.Ack(false)
}

func publishMessageToQueue(queueName string, body []byte) error {
	channel, err := RabbitConn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return channel.Publish(
		"",       
		queueName, 
		false,     
		false,    
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
