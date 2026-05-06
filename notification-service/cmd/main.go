package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"notification-service/internal/repository"
	"notification-service/internal/usecase"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	db, err := sqlx.Connect("postgres", "host=notification-db port=5432 user=user password=password dbname=notification_db sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"payment.completed",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	repo := repository.NewPostgresRepo(db)
	uc := usecase.NewNotificationUseCase(repo)

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for d := range msgs {
			if err := uc.ProcessPaymentEvent(ctx, d.MessageId, d.Body); err != nil {
				log.Printf("Failed to process message: %v", err)
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
	}()

	log.Println("Notification Service started, waiting for messages...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	cancel()
}
