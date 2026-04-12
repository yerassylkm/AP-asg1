package main

import (
	"log"
	"net"
	"os"
	
	"payment-service/internal/repository"
	"payment-service/internal/usecase"
	grpcHandler "payment-service/internal/transport/grpc"
	
	pb "github.com/yerassylkm/AP-asg2_generated"
	"google.golang.org/grpc"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", "host=payment-db port=5432 user=user password=password dbname=payment_db sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	repo := repository.NewPostgresRepo(db)
	uc := usecase.NewPaymentUseCase(repo)

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051" 
	}

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