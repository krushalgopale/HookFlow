package broker

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/worker"
	"github.com/rabbitmq/amqp091-go"
)

func processDelivery(workerID int, conn *amqp091.Connection, ch *amqp091.Channel, db *pgxpool.Pool, messages <-chan amqp091.Delivery) {
	defer ch.Close()

	for msg := range messages {
		var delivery DeliveryMessage

		// Invalid JSON: publish to dead-letter queue, then ACK.
		err := json.Unmarshal(msg.Body, &delivery)
		if err != nil {
			log.Printf("Worker %d: failed to unmarshal delivery message: %v", workerID, err)

			deadErr := DeadLetterDelivery(conn, msg, "invalid delivery message JSON")
			if deadErr != nil {
				log.Printf("Worker %d: dead-letter publishing failed: %v", workerID, deadErr)
				return
			}

			ackErr := msg.Ack(false)
			if ackErr != nil {
				log.Printf("Worker %d: failed to acknowledge delivery message: %v", workerID, ackErr)
				return
			}
			continue
		}

		// Empty delivery ID: publish to dead-letter queue, then ACK.
		if delivery.DeliveryID == "" {
			log.Printf("Worker %d: Delivery ID is empty", workerID)

			deadErr := DeadLetterDelivery(conn, msg, "delivery ID is empty")
			if deadErr != nil {
				log.Printf("Worker %d: dead-letter publishing failed: %v", workerID, deadErr)
				return
			}

			ackErr := msg.Ack(false)
			if ackErr != nil {
				log.Printf("Worker %d: failed to acknowledge delivery message: %v", workerID, ackErr)
				return
			}
			continue
		}

		log.Printf("Worker %d processing delivery job %s", workerID, delivery.DeliveryID)

		// Execute webhook delivery.
		workErr := worker.ExecuteDelivery(
			db,
			context.Background(),
			delivery.DeliveryID,
			1,
		)

		// Successful delivery: ACK.
		if workErr == nil {
			ackErr := msg.Ack(false)
			if ackErr != nil {
				log.Printf("Worker %d: failed to acknowledge delivery message: %v", workerID, ackErr)
				return
			}
			continue
		}

		// Permanent delivery failure: ACK, same as before.
		if worker.IsDeliveryFailed(workErr) {
			log.Printf("Worker %d: delivery %s failed permanently: %v", workerID, delivery.DeliveryID, workErr)

			ackErr := msg.Ack(false)
			if ackErr != nil {
				log.Printf("Worker %d: failed to acknowledge delivery message: %v", workerID, ackErr)
				return
			}
			continue
		}

		// Temporary failure: publish to retry queue.
		retryErr := RetryDelivery(conn, msg)

		if retryErr == nil {
			ackErr := msg.Ack(false)
			if ackErr != nil {
				log.Printf("Worker %d: failed to acknowledge delivery message: %v", workerID, ackErr)
				return
			}
			continue
		}

		// Retry limit reached: publish to dead-letter queue, then ACK.
		if errors.Is(retryErr, ErrRetryLimitReached) {
			deadErr := DeadLetterDelivery(conn, msg, workErr.Error())
			if deadErr != nil {
				log.Printf("Worker %d: dead-letter publishing failed: %v", workerID, deadErr)
				return
			}

			ackErr := msg.Ack(false)
			if ackErr != nil {
				log.Printf("Worker %d: dead-letter acknowledge failed: %v", workerID, ackErr)
				return
			}

			log.Printf("Worker %d: delivery %s sent to dead-letter queue", workerID, delivery.DeliveryID)
			continue
		}

		// Retry publishing failed: do not ACK the original message.
		log.Printf("Worker %d: retry publishing failed for %s: %v", workerID, delivery.DeliveryID, retryErr)
		return
	}	
}
