package domain

import "context"

type EventPublisher interface {
	PublishPaymentEvent(ctx context.Context, event PaymentEvent) error
	Close() error
}
