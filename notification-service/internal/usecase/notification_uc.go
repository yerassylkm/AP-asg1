package usecase

import (
	"context"
	"encoding/json"
	"log"
	"notification-service/internal/domain"
)

type NotificationUseCase struct {
	repo domain.ProcessedMessageRepository
}

func NewNotificationUseCase(repo domain.ProcessedMessageRepository) *NotificationUseCase {
	return &NotificationUseCase{repo: repo}
}

func (uc *NotificationUseCase) ProcessPaymentEvent(ctx context.Context, messageID string, body []byte) error {
	processed, err := uc.repo.IsProcessed(ctx, messageID)
	if err != nil {
		return err
	}
	if processed {
		log.Printf("Message %s already processed, skipping", messageID)
		return nil
	}

	var event domain.PaymentEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}

	log.Printf("[Notification] Sent email to %s for Order #%s. Amount: $%.2f", event.CustomerEmail, event.OrderID, float64(event.Amount)/100)

	return uc.repo.MarkProcessed(ctx, messageID)
}
