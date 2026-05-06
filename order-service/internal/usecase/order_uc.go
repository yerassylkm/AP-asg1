package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"order-service/internal/domain"
	"order-service/internal/transport/grpc"
	"time"

	"github.com/google/uuid"
)

type OrderUseCase struct {
	repo          domain.OrderRepository
	paymentClient *grpc.PaymentClient
}

func NewOrderUseCase(r domain.OrderRepository, pc *grpc.PaymentClient) *OrderUseCase {
	return &OrderUseCase{
		repo:          r,
		paymentClient: pc,
	}
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, customerID, customerEmail, itemName string, amount int64) (*domain.Order, error) {
	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}

	order := &domain.Order{
		ID:            uuid.New().String(),
		CustomerID:    customerID,
		CustomerEmail: customerEmail,
		ItemName:      itemName,
		Amount:        amount,
		Status:        domain.StatusPending,
		CreatedAt:     time.Now(),
	}

	if err := uc.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	go uc.processAsyncPayment(order.ID, amount, customerEmail)

	return order, nil
}

func (uc *OrderUseCase) processAsyncPayment(orderID string, amount int64, customerEmail string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status, err := uc.paymentClient.ProcessPayment(ctx, orderID, amount)

	newStatus := domain.StatusFailed
	if err == nil && status == "SUCCESS" {
		newStatus = domain.StatusPaid
	}

	if err := uc.repo.UpdateStatus(context.Background(), orderID, newStatus); err != nil {
		log.Printf("Failed to update order status: %v", err)
	}
}

func (uc *OrderUseCase) CancelOrder(ctx context.Context, id string) error {
	order, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status != domain.StatusPending {
		return fmt.Errorf("cannot cancel order in status: %s", order.Status)
	}

	return uc.repo.UpdateStatus(ctx, id, domain.StatusCancelled)
}

func (uc *OrderUseCase) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *OrderUseCase) GetCustomerOrders(customerId string) ([]domain.Order, error) {
	return uc.repo.GetCustomerOrders(customerId)
}
