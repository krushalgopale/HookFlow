package broker

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type DeliveryMessage struct {
	DeliveryID string `json:"delivery_id"`
}

func PublishDelivery(conn *amqp091.Connection, deliveryID string) error {
	if deliveryID == "" {
		err := errors.New("Delivery ID cannot be empty")
		log.Println(err)
		return err
	}

	// Open a Channel on the existing RabbitMQ connection
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open publisher channel: %v", err)
		return err
	}

	// Close Channel when the funnction returns
	defer ch.Close()

	// Encode the delivery message in JSON
	body, err := json.Marshal(DeliveryMessage{
		DeliveryID: deliveryID,
	})
	if err != nil {
		log.Printf("Failed to marshal delivery message: %v", err)
		return err
	}

	// Publish the delivery message to the queue
	err = ch.Publish(
		"",            // exchange default
		DeliveryQueue, // routing key or queue name
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp091.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish delivery message: %v", err)
		return err
	}

	return nil
}
