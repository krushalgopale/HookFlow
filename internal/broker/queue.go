package broker

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

const DeliveryQueue = "hookflow.deliveries"

func DeclareDeliveryQueue(conn *amqp091.Connection) error {
	// Open a Channel on the existing RabbitMQ connection
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return err
	}

	// Close Channel when the function returns
	defer ch.Close()

	_, err = ch.QueueDeclare(
		DeliveryQueue, // queue name
		true,          // durability
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return err
	}

	return nil
}
