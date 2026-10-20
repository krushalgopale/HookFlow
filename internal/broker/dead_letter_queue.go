package broker

import (
	"strings"
	"unicode/utf8"

	"github.com/rabbitmq/amqp091-go"
)

const deadLetterReasonHeader = "x-hookflow-dead-letter-reason"

// Publishes a message to the dead-letter queue after retries are exhausted
func DeadLetterDelivery(conn *amqp091.Connection, messages amqp091.Delivery, reason string) error {
	headers := copyRetryHeaders(messages.Headers)

	reason = strings.ToValidUTF8(reason, "?")

	if len(reason) > 1024 {
		limit := 1024

		for !utf8.ValidString(reason[:limit]) {
			limit--
		}

		reason = reason[:limit]
	}

	headers[deadLetterReasonHeader] = reason

	return publishSafely(
		conn,
		DeliveryDeadQueue,
		messages.Body,
		headers,
	)
}
