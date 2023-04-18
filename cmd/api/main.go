package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"net/http"
	"fmt"
)

type Payload struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

func main() {
	r := gin.Default()

	ch, closeQueue, err := initRabbitMQ()
	if err != nil {
		fmt.Printf("Error initializing RabbitMQ %w", err)
		closeQueue()
		return
	}
	defer closeQueue()

	r.POST("/message", func(c *gin.Context) {
		var payload Payload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Convert the message to JSON
		message, err := json.Marshal(payload)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Publish the message to the exchange
		if err := ch.Publish(
			EXCHANGE_NAME, // Exchange name
			"",            // Routing key
			false,         // Mandatory
			false,         // Immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        message,
			},
		); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Return a success response
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Run(SERVER_PORT)
}
