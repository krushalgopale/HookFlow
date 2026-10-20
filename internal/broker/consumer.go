package broker

import (
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"
)

const deliveryWorkerCount = 5

type deliveryConsumer struct {
	id       int
	channel  *amqp091.Channel
	messages <-chan amqp091.Delivery
}

func ConsumeDeliveries(conn *amqp091.Connection, db *pgxpool.Pool) error {
	consumers := make([]deliveryConsumer, 0, deliveryWorkerCount)

	for i := 1; i <= deliveryWorkerCount; i++ {
		// Create a separate channel and consumer for each worker
		ch, err := conn.Channel()
		if err != nil {
			closeConsumers(consumers)
			return fmt.Errorf("Failed to open consumer channel for worker %d: %w", i, err)
		}

		// Allow one unacknowledged delivery per worker
		if err := ch.Qos(1, 0, false); err != nil {
			ch.Close()
			closeConsumers(consumers)
			return fmt.Errorf("Failed to configure consumer prefetch for worker %d: %w", i, err)
		}

		// Consumes the delivery message from the queue
		messages, err := ch.Consume(
			DeliveryQueue, // queue name
			"",            // consumer name
			false,         // ack
			false,         // exclusive
			false,         // no-local
			false,         // no-wait
			nil,           // arguments
		)
		if err != nil {
			ch.Close()
			closeConsumers(consumers)
			return fmt.Errorf("Failed to consume delivery message for worker %d: %w", i, err)
		}

		// Register consumers
		consumers = append(consumers, deliveryConsumer{
			id:       i,
			channel:  ch,
			messages: messages,
		})
	}

	// Process incoming delivery messages in the background
	for _, consumer := range consumers {
		go processDelivery(consumer.id, conn, consumer.channel, db, consumer.messages)
	}

	log.Printf("Started %d delivery workers", deliveryWorkerCount)

	return nil
}

func closeConsumers(consumers []deliveryConsumer) {
	for _, consumer := range consumers {
		_ = consumer.channel.Close()
	}
}
