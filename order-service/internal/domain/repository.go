package domain

import "context"

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	GetCustomerOrders(customerID string) ([]Order, error)
}
