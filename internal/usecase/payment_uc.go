package usecase

import (
	"context"
	"errors"
	"payment-service/internal/domain"

	"github.com/google/uuid"
)

type PaymentUseCase struct {
	repo domain.PaymentRepository
}

func NewPaymentUseCase(r domain.PaymentRepository) *PaymentUseCase {
	return &PaymentUseCase{repo: r}
}

func (uc *PaymentUseCase) ProcessPayment(ctx context.Context, orderID string, amount int64) (*domain.Payment, error) {
	status := domain.StatusAuthorized
	if amount > 100000 {
		status = domain.StatusDeclined
	}

	payment := &domain.Payment{
		ID:            uuid.New().String(),
		OrderID:       orderID,
		TransactionID: uuid.New().String(),
		Amount:        amount,
		Status:        status,
	}

	if err := uc.repo.Create(ctx, payment); err != nil {
		return nil, err
	}

	if status == domain.StatusDeclined {
		return payment, errors.New("payment declined: amount exceeds limit")
	}

	return payment, nil
}

func (uc *PaymentUseCase) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	return uc.repo.GetByOrderID(ctx, orderID)
}
