package publisher

import (
	"context"
	"encoding/json"
	"log"
	"payment-service/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type rabbitPublisher struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

func NewRabbitPublisher(url string) (domain.EventPublisher, error) {
	var conn *amqp091.Connection
	var err error
	for i := 0; i < 5; i++ {
		conn, err = amqp091.Dial(url)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ connect attempt %d failed: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"payment.dlx",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		"payment.dlq",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	err = ch.QueueBind(
		"payment.dlq",
		"failed",
		"payment.dlx",
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	args := amqp091.Table{
		"x-dead-letter-exchange":    "payment.dlx",
		"x-dead-letter-routing-key": "failed",
	}
	_, err = ch.QueueDeclare(
		"payment.completed",
		true,
		false,
		false,
		false,
		args,
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
		"",
		"payment.completed",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
			MessageId:   uuid.New().String(),
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
