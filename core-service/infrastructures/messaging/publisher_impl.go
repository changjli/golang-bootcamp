package messaging

import (
	"context"
	"encoding/json"
	"events"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

// RabbitMqPublisher is the implementation of the PaymentPublisher interface.
type RabbitMqPublisher struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

// NewRabbitMQPublisher creates a new publisher and connects to RabbitMQ.
func NewRabbitMQPublisher() (*RabbitMqPublisher, error) {
	conn, err := amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare the queue to ensure it exists. This is idempotent.
	_, err = ch.QueueDeclare(
		"payment.queue", // queue name
		true,            // durable: the queue will survive a broker restart
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMqPublisher{conn: conn, channel: ch}, nil
}

// PublishPaymentRequested serializes the event to JSON and sends it to the queue.
func (p *RabbitMqPublisher) PublishPaymentRequested(ctx context.Context, event *events.PaymentRequestedEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.channel.PublishWithContext(ctx,
		"",              // exchange: use the default exchange
		"payment.queue", // routing key: the queue name
		false,           // mandatory
		false,           // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err == nil {
		log.Printf("Successfully published payment request for PaymentID: %s", event.PaymentID)
	}
	return err
}
