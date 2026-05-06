package main

import (
	"log"
	"net"
	"os"

	"payment-service/internal/publisher"
	"payment-service/internal/repository"
	grpcHandler "payment-service/internal/transport/grpc"
	httpTransport "payment-service/internal/transport/http"
	"payment-service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	pb "github.com/yerassylkm/AP-asg2_generated"
	"google.golang.org/grpc"
)

func main() {
	db, err := sqlx.Connect("postgres", "host=payment-db port=5432 user=user password=password dbname=payment_db sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	pub, err := publisher.NewRabbitPublisher(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer pub.Close()

	repo := repository.NewPostgresRepo(db)
	uc := usecase.NewPaymentUseCase(repo, pub)

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8081"
	}

	go func() {
		r := gin.Default()
		handler := httpTransport.NewPaymentHandler(uc)

		r.POST("/payments", handler.CreatePayment)
		r.GET("/payments/:order_id", handler.GetPayment)

		log.Printf("HTTP Payment Service started on :%s", httpPort)
		if err := r.Run(":" + httpPort); err != nil {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s, grpcHandler.NewPaymentGRPCHandler(uc))

	log.Printf("gRPC Payment Service started on :%s", grpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
