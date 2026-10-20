package broker

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

const (
	DeliveryQueue      = "hookflow.deliveries"
	DeliveryRetryQueue = "hookflow.deliveries.retry"
	DeliveryDeadQueue  = "hookflow.deliveries.dead"

	retryDelayMilliseconds int32 = 5000
)

func DeclareDeliveryQueue(conn *amqp091.Connection) error {
	// Open a Channel on the existing RabbitMQ connection
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel: %w", err)
	}

	// Close Channel when the function returns
	defer ch.Close()

	// Main queue
	_, err = ch.QueueDeclare(
		DeliveryQueue, // queue name
		true,          // durability
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("Failed to declare a main queue: %w", err)
	}

	// Dead-letter queue
	_, err = ch.QueueDeclare(
		DeliveryDeadQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Failed to declare dead-letter queue: %w", err)
	}

	// Retry queue
	_, err = ch.QueueDeclare(
		DeliveryRetryQueue,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-message-ttl":             retryDelayMilliseconds,
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": DeliveryQueue,
		},
	)
	if err != nil {
		return fmt.Errorf("Failed to declare retry queue: %w", err)
	}

	return nil
}
