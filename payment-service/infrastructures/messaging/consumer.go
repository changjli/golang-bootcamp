package messaging

import (
	"context"
	"encoding/json"
	"events"
	"log"
	"payment-service/domains/transaction"

	"github.com/rabbitmq/amqp091-go"
)

// PaymentConsumer listens for payment events from RabbitMQ and orchestrates their processing.
type PaymentConsumer struct {
	conn       *amqp091.Connection
	channel    *amqp091.Channel
	trxUsecase transaction.TransactionUsecase // The use case to handle the business logic
}

// NewPaymentConsumer creates a new consumer service.
func NewPaymentConsumer(trxUsecase transaction.TransactionUsecase) (*PaymentConsumer, error) {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare the queue to ensure it exists.
	_, err = ch.QueueDeclare(
		"payment.queue", // queue name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &PaymentConsumer{
		conn:       conn,
		channel:    ch,
		trxUsecase: trxUsecase,
	}, nil
}

// StartConsumer begins listening for messages on the payment.queue.
// This should be run as a goroutine from your main function.
func (c *PaymentConsumer) StartConsumer(ctx context.Context) {
	msgs, err := c.channel.Consume(
		"payment.queue", // queue
		"",              // consumer tag
		true,            // auto-ack: automatically acknowledge messages. For production, this should be false.
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	log.Println("RabbitMQ consumer started. Waiting for payment messages...")

	// Create a forever loop to process messages
	for {
		select {
		case <-ctx.Done():
			// If the context is cancelled (e.g., application shutdown), stop the consumer.
			log.Println("Shutting down consumer...")
			c.conn.Close()
			c.channel.Close()
			return
		case d := <-msgs:
			log.Printf("Received a message: %s", d.Body)
			var event events.PaymentRequestedEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error unmarshalling message: %s", err)
				// In a real system, you might move this to a dead-letter queue.
				continue
			}

			// Pass the event to the use case for processing.
			err := c.trxUsecase.ProcessPaymentRequest(ctx, &event)
			if err != nil {
				log.Printf("Error processing payment %s: %v", event.PaymentID, err)
			} else {
				log.Printf("Successfully processed payment %s.", event.PaymentID)
			}
		}
	}
}
