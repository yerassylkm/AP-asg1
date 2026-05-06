package main

import (
	"log"
	"order-service/internal/repository"
	"order-service/internal/transport/http"
	"order-service/internal/usecase"

	"order-service/internal/transport/grpc" 
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", "host=order-db port=5432 user=user password=password dbname=order_db sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	payClient, err := grpc.NewPaymentClient()
	if err != nil {
		log.Fatalf("Failed to init gRPC client: %v", err)
	}
	
	repo := repository.NewPostgresRepo(db)
	uc := usecase.NewOrderUseCase(repo, payClient)
	handler := http.NewOrderHandler(uc)

	r := gin.Default()

	r.POST("/orders", handler.CreateOrder)
	r.GET("/orders/:id", handler.GetOrder)
	r.GET("/orders", handler.GetOrdersByCustomer)
	r.PATCH("/orders/:id/cancel", handler.CancelOrder)

	log.Println("Order Service started on :8080")
	r.Run(":8080")
}
