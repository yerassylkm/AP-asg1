package repository

import (
	"context"
	"payment-service/internal/domain"

	"github.com/jmoiron/sqlx"
)

type postgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) domain.PaymentRepository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) Create(ctx context.Context, p *domain.Payment) error {
	query := `INSERT INTO payments (id, order_id, transaction_id, amount, status, customer_email) 
              VALUES (:id, :order_id, :transaction_id, :amount, :status, :customer_email)`
	_, err := r.db.NamedExecContext(ctx, query, p)
	return err
}

func (r *postgresRepo) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	var p domain.Payment
	err := r.db.GetContext(ctx, &p, "SELECT * FROM payments WHERE order_id=$1", orderID)
	return &p, err
}
