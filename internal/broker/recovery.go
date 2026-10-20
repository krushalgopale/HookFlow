package broker

import (
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"
)

const workerRestartDelay = time.Second

func runDeliveryWorker(consumer deliveryConsumer, conn *amqp091.Connection, db *pgxpool.Pool) {
	workerID := consumer.id

	for {
		processDelivery(workerID, conn, consumer.channel, db, consumer.messages)

		if conn == nil || conn.IsClosed() {
			log.Printf("Worker %d stopped: RabbitMQ connection is closed", workerID)
			return
		}

		log.Printf("Worker %d stopped; restarting in %s", workerID, workerRestartDelay)

		// Wait for 1 second before restarting
		time.Sleep(workerRestartDelay)

		// Restart
		newConsumer, err := createDeliveryConsumer(conn, workerID)
		if err != nil {
			log.Printf("Worker %d restarted failed: %v", workerID, err)

			if conn.IsClosed() {
				return
			}
			continue
		}
		consumer = newConsumer
		log.Printf("Worker %d restarted successfully", workerID)
	}
}

func createDeliveryConsumer(conn *amqp091.Connection, workerID int) (deliveryConsumer, error) {
	if conn == nil || conn.IsClosed() {
		return deliveryConsumer{}, amqp091.ErrClosed
	}

	// Create a separate channel and consumer for each worker
	ch, err := conn.Channel()
	if err != nil {
		return deliveryConsumer{}, err
	}

	// Allow only one unacknowledged message per worker.
	if err := ch.Qos(1, 0, false); err != nil {
		_ = ch.Close()
		return deliveryConsumer{}, err
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
		return deliveryConsumer{}, err
	}

	return deliveryConsumer{
		id:       workerID,
		channel:  ch,
		messages: messages,
	}, nil
}
