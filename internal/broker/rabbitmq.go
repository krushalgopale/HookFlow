package broker

import (
	"fmt"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

func Connect() (*amqp091.Connection, error) {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is not set in .env")
	}

	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err 
	}

	return conn, nil
}
