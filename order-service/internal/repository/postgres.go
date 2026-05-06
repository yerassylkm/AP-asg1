package repository

import (
	"context"
	"order-service/internal/domain"

	"github.com/jmoiron/sqlx"
)

type postgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) domain.OrderRepository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) Create(ctx context.Context, o *domain.Order) error {
	query := `INSERT INTO orders (id, customer_id, customer_email, item_name, amount, status, created_at) 
			  VALUES (:id, :customer_id, :customer_email, :item_name, :amount, :status, :created_at)`
	_, err := r.db.NamedExecContext(ctx, query, o)
	return err
}

func (r *postgresRepo) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	var o domain.Order
	err := r.db.GetContext(ctx, &o, "SELECT * FROM orders WHERE id=$1", id)
	return &o, err
}

func (r *postgresRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE orders SET status=$1 WHERE id=$2", status, id)
	return err
}

func (r *postgresRepo) GetCustomerOrders(customerID string) ([]domain.Order, error) {
	orders := []domain.Order{}

	query := `SELECT id, customer_id, customer_email, item_name, amount, status, created_at 
              FROM orders WHERE customer_id = $1`

	err := r.db.Select(&orders, query, customerID)
	return orders, err
}
