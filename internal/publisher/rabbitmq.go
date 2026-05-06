package publisher

import (
	"context"
	"encoding/json"
	"log"
	"payment-service/internal/domain"

	"github.com/rabbitmq/amqp091-go"
)

type rabbitPublisher struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

func NewRabbitPublisher(url string) (domain.EventPublisher, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare dead letter exchange
	err = ch.ExchangeDeclare(
		"payment.dlx", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	// Declare dead letter queue
	_, err = ch.QueueDeclare(
		"payment.dlq", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	// Bind DLQ to DLX
	err = ch.QueueBind(
		"payment.dlq", // queue name
		"failed",      // routing key
		"payment.dlx", // exchange
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	// Declare main queue with DLX
	args := amqp091.Table{
		"x-dead-letter-exchange":    "payment.dlx",
		"x-dead-letter-routing-key": "failed",
	}
	_, err = ch.QueueDeclare(
		"payment.completed", // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		args,                // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &rabbitPublisher{conn: conn, ch: ch}, nil
}

func (p *rabbitPublisher) PublishPaymentEvent(ctx context.Context, event domain.PaymentEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.ch.PublishWithContext(ctx,
		"",                  // exchange
		"payment.completed", // routing key
		false,               // mandatory
		false,               // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Printf("Failed to publish event: %v", err)
		return err
	}

	log.Printf("Published payment event for order %s", event.OrderID)
	return nil
}

func (p *rabbitPublisher) Close() error {
	if p.ch != nil {
		p.ch.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
