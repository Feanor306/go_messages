package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/streadway/amqp"
)

const (
	RABBIT_MQ_CONN = "amqp://guest:guest@localhost:5672/"
	QUEUE_NAME   = "messages"
	REDIS_ADDRESS   = "localhost:6379"
)

type Message struct {
	First  string
	Second string
	Third  string
}

func main() {
	// Connect to Redis server
	client := redis.NewClient(&redis.Options{
		Addr: REDIS_ADDRESS,
	})

	// Create the RabbitMQ connection and channel
	conn, err := amqp.Dial(RABBIT_MQ_CONN)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ server: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a RabbitMQ channel: %s", err)
	}
	defer ch.Close()

	// Declare the queue to consume from
	q, err := ch.QueueDeclare(
		QUEUE_NAME,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a RabbitMQ queue: %s", err)
	}

	// Consume messages from the RabbitMQ queue
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to start consuming RabbitMQ queue: %s", err)
	}

	// Process incoming messages
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// Parse the incoming message
			var message Message
			if err := json.Unmarshal(d.Body, &message); err != nil {
				log.Printf("Failed to parse incoming message: %s", err)
				continue
			}

			// Publish message to Redis server
			if err := client.LPush("messages", fmt.Sprintf("%v", message)); err != nil {
				log.Printf("Failed to publish message to Redis: %s", err)
				continue
			}

			log.Printf("Processed message: %s", message)
		}
	}()

	log.Printf("Waiting for messages...")

	<-forever
}
