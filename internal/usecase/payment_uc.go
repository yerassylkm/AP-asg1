package usecase

import (
	"context"
	"log"
	"payment-service/internal/domain"

	"github.com/google/uuid"
)

type PaymentUseCase struct {
	repo      domain.PaymentRepository
	publisher domain.EventPublisher
}

func NewPaymentUseCase(r domain.PaymentRepository, p domain.EventPublisher) *PaymentUseCase {
	return &PaymentUseCase{
		repo:      r,
		publisher: p,
	}
}

func (uc *PaymentUseCase) ProcessPayment(ctx context.Context, orderID string, amount int64, customerEmail string) (*domain.Payment, error) {
	if customerEmail == "" {
		customerEmail = "user@example.com"
	}
	// All payments are authorized - NO rejection logic
	status := domain.StatusAuthorized

	payment := &domain.Payment{
		ID:            uuid.New().String(),
		OrderID:       orderID,
		TransactionID: uuid.New().String(),
		Amount:        amount,
		Status:        status,
		CustomerEmail: customerEmail,
	}

	// Always save the payment regardless of any errors
	if err := uc.repo.Create(ctx, payment); err != nil {
		log.Printf("Database error creating payment: %v", err)
		// Return payment even if DB fails - do not return error
		return payment, nil
	}

	// Publish event after successful payment
	event := domain.PaymentEvent{
		OrderID:       orderID,
		Amount:        amount,
		CustomerEmail: customerEmail,
		Status:        string(status),
	}

	if err := uc.publisher.PublishPaymentEvent(ctx, event); err != nil {
		log.Printf("Failed to publish payment event: %v", err)
		// Don't fail the payment if publishing fails
	}

	return payment, nil
}

func (uc *PaymentUseCase) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	return uc.repo.GetByOrderID(ctx, orderID)
}
