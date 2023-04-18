package main

import (
	"github.com/streadway/amqp"
)

const (
	EXCHANGE_NAME = "messages"
	RABBIT_MQ_CONN = "amqp://user:pass@localhost:7001/"
	SERVER_PORT= ":8080"
)

func initRabbitMQ() (*amqp.Channel, func(), error) {
	conn, err := amqp.Dial(RABBIT_MQ_CONN)
	if err != nil {
		return nil, func() {
			conn.Close()
		}, err
	}

	// Create a new channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, func() {
			conn.Close()
			ch.Close()
		}, err
	}

	// Declare a new exchange
	if err := ch.ExchangeDeclare(
		EXCHANGE_NAME, // Exchange name
		"fanout",      // Exchange type
		true,          // Durable
		false,         // Auto-deleted
		false,         // Internal
		false,         // No-wait
		nil,           // Arguments
	); err != nil {
		return ch, func() {
			conn.Close()
			ch.Close()
		}, err
	}

	return ch, func() {
		conn.Close()
		ch.Close()
	}, nil
}
