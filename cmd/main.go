package main

import (
	"log"
	"payment-service/internal/repository"
	"payment-service/internal/transport/http"
	"payment-service/internal/usecase"

	"github.com/gin-gonic/gin"
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
	handler := http.NewPaymentHandler(uc)

	r := gin.Default()
	
	r.POST("/payments", handler.CreatePayment)   
	r.GET("/payments/:order_id", handler.GetPayment) 

	log.Println("Payment Service started on :8081")
	r.Run(":8081")
}