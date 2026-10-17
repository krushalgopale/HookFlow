package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/worker"
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
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			ch.Close()
			closeConsumers(consumers)
			return fmt.Errorf("Failed to consume delivery message for worker %d: %w", i, err)
		}
		
		// Register consumers 
		consumers = append(consumers, deliveryConsumer{
			id: i,
			channel: ch,
			messages: messages,
		})
	}

	// Process incoming delivery messages in the background
	for _, consumer := range consumers {
		go processDelivery(consumer.id, consumer.channel, db, consumer.messages)
	}

	log.Printf("Started %d delivery workers", deliveryWorkerCount)

	return nil
}

func closeConsumers(consumers []deliveryConsumer) {
	for _, consumer := range consumers {
		_ = consumer.channel.Close()
	}
}

func processDelivery(workerID int, ch *amqp091.Channel, db *pgxpool.Pool, messages <-chan amqp091.Delivery) {
	defer ch.Close()

	for msg := range messages {
		var delivery DeliveryMessage

		// Don't Requeue delivery message if decoding from JSON fails
		err := json.Unmarshal(msg.Body, &delivery)
		if err != nil {
			log.Printf("Worker %d: failed to unmarshal delivery message: %v", workerID, err)
			rejectErr := msg.Reject(false)
			if rejectErr != nil {
				log.Printf("Worker %d: failed to reject delivery message: %v", workerID, rejectErr)
			}
			continue
		}

		// Don't Requeue delivery message if delivery ID is empty
		if delivery.DeliveryID == "" {
			log.Printf("Worker %d: Delivery ID is empty", workerID)
			rejectErr := msg.Reject(false)
			if rejectErr != nil {
				log.Printf("Worker %d: failed to reject delivery message: %v", workerID, rejectErr)
			}
			continue
		}

		log.Printf("Worker %d processing delivery job %s", workerID, delivery.DeliveryID)

		workErr := worker.ExecuteDelivery(
			db,
			context.Background(),
			delivery.DeliveryID,
			1,
		)

		// If successfully delivered the message acknowledge it
		if workErr == nil {
			ackErr := msg.Ack(false)
			if ackErr != nil {
				log.Printf("Worker %d: failed to acknowledge delivery message: %v", workerID, ackErr)
			}
			continue
		}

		// If delivery failed then also acknowledge it
		if worker.IsDeliveryFailed(workErr) {
			log.Printf("Worker %d: delivery %s failed permanently: %v", workerID, delivery.DeliveryID, workErr)

			ackErr := msg.Ack(false)
			if ackErr != nil {
				log.Printf("Worker %d: failed to acknowledge delivery message: %v", workerID, ackErr)
			}
			continue
		}

		log.Printf("Worker %d: processing delivery %s failed %v; requeueing", workerID, delivery.DeliveryID, workErr)

		// Requeue the delivery message if worker processing failed
		rejectErr := msg.Reject(true)
		if rejectErr != nil {
			log.Printf("Worker %d: failed to requeue delivery %s: %v", workerID, delivery.DeliveryID, rejectErr)
		}

	}
	log.Printf("Worker %d stopped", workerID)
}
