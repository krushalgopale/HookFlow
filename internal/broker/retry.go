package broker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

const (
	maxDeliveryRetries = 3
	retryCountHeader   = "x-hookflow-retry-count"
	publishTimeout     = 5 * time.Second
)

var ErrRetryLimitReached = errors.New("retry limit reached")

// Retry Delivery publishes a failed message to the retry queue
func RetryDelivery(conn *amqp091.Connection, messages amqp091.Delivery) error {
	headers := copyRetryHeaders(messages.Headers)
	retries := getRetryCount(headers)

	// Indicates that the message has exhausted its retry attempts
	if retries >= maxDeliveryRetries {
		return ErrRetryLimitReached
	}

	headers[retryCountHeader] = int32(retries + 1)

	return publishSafely(
		conn,
		DeliveryRetryQueue,
		messages.Body,
		headers,
	)
}

// Creates a separate copy so the original message headers remain unchanged.
func copyRetryHeaders(original amqp091.Table) amqp091.Table {
	headers := make(amqp091.Table, len(original))

	for key, value := range original {
		headers[key] = value
	}

	return headers
}

// Reads the current retry count from message headers.
func getRetryCount(headers amqp091.Table) int {
	value, exists := headers[retryCountHeader]
	if !exists {
		return 0
	}

	switch count := value.(type) {
	case int:
		if count >= 0 {
			return count
		}

	case int8:
		if count >= 0 {
			return int(count)
		}

	case int32:
		if count >= 0 {
			return int(count)
		}

	case int64:
		if count >= 0 && int64(int(count)) == count {
			return int(count)
		}
	}

	return maxDeliveryRetries
}

// Publishes a message and verifies the publish result.
func publishSafely(
	conn *amqp091.Connection,
	queueName string,
	body []byte,
	headers amqp091.Table,
) error {
	if conn == nil {
		return errors.New("RabbitMQ connection is nil")
	}

	if conn.IsClosed() {
		return errors.New("RabbitMQ connection is closed")
	}

	if queueName == "" {
		return errors.New("Destination queue name is empty")
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open publisher channel: %w", err)
	}
	defer ch.Close()

	err = ch.Confirm(false)
	if err != nil {
		return fmt.Errorf("Failed to enable publisher confirms: %w", err)
	}

	confirmations := ch.NotifyPublish(
		make(chan amqp091.Confirmation, 1),
	)

	returns := ch.NotifyReturn(
		make(chan amqp091.Return, 1),
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		publishTimeout,
	)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"",        // default exchange
		queueName, // routing key or queue name
		true,      // return message if unroutable
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp091.Persistent,
			Headers:      copyRetryHeaders(headers),
			Body:         append([]byte(nil), body...),
		},
	)
	if err != nil {
		return fmt.Errorf("Failed publish to %s: %w", queueName, err)
	}

	var confirmation amqp091.Confirmation
	var ok bool

	select {
	case confirmation, ok = <-confirmations:
		if !ok {
			return fmt.Errorf("Publisher confirmation channel closed for %s", queueName)
		}
	case <-ctx.Done():
		return fmt.Errorf("Publisher confirmation timed out for queue %s", queueName)
	}

	select {
	case returned, ok := <-returns:
		if !ok {
			return fmt.Errorf("Return notification channel closed %s", queueName)
		}

		return fmt.Errorf(
			"Message unroutable: exchange=%s routing=%s reply_code=%d reply_text=%s",
			returned.Exchange,
			returned.RoutingKey,
			returned.ReplyCode,
			returned.ReplyText,
		)

	default:
	}
	if !confirmation.Ack {
		return fmt.Errorf("RabbitMQ NACK for queue %s", queueName)
	}
	return nil
}
